package util

import "github.com/gin-gonic/gin"

func GetWayParam(c *gin.Context) string {
	queryValue := c.Query("way")
	if queryValue != "" {
		return queryValue
	}

	formValue := c.PostForm("way")
	if formValue != "" {
		return formValue
	}

	type WayBody struct {
		Way string `json:"way"`
	}

	var jsonBody WayBody

	err := c.BindJSON(&jsonBody)
	if err == nil && jsonBody.Way != "" {
		return jsonBody.Way
	}
	pathValue := c.Param("way")
	if pathValue != "" {
		return pathValue
	}

	return ""
}
