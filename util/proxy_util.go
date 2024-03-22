package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Proxy(targetURL string, c *gin.Context) {
	proxyUrl, _ := url.Parse(targetURL)

	target, _ := url.Parse(proxyUrl.Scheme + "://" + proxyUrl.Host)
	if c.Param("name") != "" {
		c.Request.URL.Path = proxyUrl.Path + c.Param("name")
	} else {
		c.Request.URL.Path = proxyUrl.Path
	}
	c.Request.URL.Path = strings.TrimRight(c.Request.URL.Path, "/")

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.ServeHTTP(c.Writer, c.Request)
}

func RootProxy(targetURL string, c *gin.Context) error {
	// 获取原来请求的所有get参数
	queryParams := c.Request.URL.Query()
	// 若指定为get模式 则将原本post参数也赋予到get参数中

	// 若指定为post模式 则将原本get参数也赋予到post参数中

	// 全部赋予到新的请求中
	targetURL += "?"
	for key, values := range queryParams {
		for _, value := range values {
			targetURL += key + "=" + value + "&"
		}
	}
	targetURL = strings.TrimRight(targetURL, "&")

	// 获取用户请求的请求头
	userHeaders := c.Request.Header

	// 解析请求参数
	requestBody := c.PostForm("body")

	// 创建新的请求
	req, err := http.NewRequest(c.Request.Method, targetURL, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request: "+err.Error())
		return fmt.Errorf("创建代理请求失败")
	}

	// 将用户请求的请求头添加到新请求的请求头中
	for key, values := range userHeaders {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// 设置请求的 Body（如果是 POST 请求）
	if c.Request.Method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(strings.NewReader(requestBody))
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending request: "+err.Error())
		return fmt.Errorf("发送代理请求失败")
	}
	defer resp.Body.Close()

	// 返回响应给客户端
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading response: "+err.Error())
		return fmt.Errorf("返回响应失败")
	}
	c.String(http.StatusOK, string(body))
	return nil
}
