<div align=center>
<img src="https://s2.loli.net/2024/03/17/o3R6NUMbxwEd1zQ.jpg" style="width:100px;"/>
<h2>云阶</h2>
</div>

### 一、功能简述
- 安全的、高性能的、可独立部署的、代理（反向代理）的、自助代理的、负载均衡的、可持久化的服务地址管理工具。
- 通过返回集合中设定的地址后由开发者再次自行调用返回的地址实现的自助负载均衡。通过代理请求直接返回向目标地址请求的结果。功能相同或不同、相对路径相同或不同均可代理，即“代理”实现“找平”。
- 支持根路径代理、支持代理公网地址或域名、支持自动屏蔽失效地址、地址失效告警。
- 通过请求加"盐"，实现代理出口变化限制与自嵌套加密。

### 二、部署方式
1. 安装gcc环境或其他sqlite3支持环境
2. 解压程序压缩包，启动二进制主程序文件（启动时可携带一个作为启动端口号的参数，程序名后空格后写就行）

### 三、使用教程
1. 按照指定的端口号的访问地址下的/web路径进入管理界面
    ```
    默认用户名：admin96
    默认密码：admin96
    可以在“系统管理”中修改默认密码。
    ```
2. 添加映射集
   ```
   在"映射管理"中，上方为映射集设置，您可以在此添加新的映射集，下方可以对已经存在的映射集进行管理。
   ```
3. 添加代理模式和自助模式代号
    ```
    在“自助模式”和“代理模式”中添代理路径（way）和代理点（point、即映射集），以及负载均衡模式。
    ```
4. 调用程序
    ```
    对自主模式路径（/self）发送任意形式带参请求，携带way参数指定到自己的自助模式代号上从而获得一个好用的请求地址（服务器ping）
    对代理模式路径（/proxy）发送任意形式带参请求，携带way参数指定到自己的自助模式代号上从而获得请求结果。
    若开启加盐功能，则需按照加盐规则进行请求。
    自检失效的服务器或请求地址将被暂时停用至10分钟后再次检查状态，若检查三次无效则发送告警且长期停用（可在后台手动解除）。
    ```

### 四、接口文档

### 五、专注的点
- 使得多服务器或同服务器多服务启动的不同后端接口或多相对路径的服务程序能够统一请求模式和地址。
- 使得负载均衡策略囊括前端和后端两方面领域。例如四六级查分的时候，大家都去访问查分网站（调用逻辑和数据集相同的接口）的时候导致网站卡死，这时使用前端调用负载均衡接口即可降低对单服务器请求量，但如果实现负载均衡接口的单台服务器带宽有限（即受到网速限制）从而导致代理能力有限、则此时选用自助模式直接返回原本由代理服务器代理的地址由客户端自助请求即可解决（代理服务器性能同理）。
- 自动排除失效的地址从而为调用端提供健壮性保障。
- 实现后端请求地址和功能热更新（无需修改前端代码即可更新后端功能和服务地址）。
- 实现类似k8s那种用户无感的迭代后端接口服务程序（迭代部分地址的服务时，先停用部分地址上的服务/服务程序，更新完成后再逐步启用更新其他地址即可）
- 使得本地性能极低服务器设备（或容器）能够通过代理模式实现多服务代理，从而提高服务器性能利用率。

### 六、开发进度
- [X] 自带界面（/web）
- [X] 自助模式（/self）
- [ ] 代理模式（/proxy）
- [ ] 空路径默认模式
- [ ] 一键启动包
- [ ] 地址失效告警
- [ ] 加盐模式
- [ ] 手动生效或失效地址
- [X] TCP网络模式
