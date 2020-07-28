# k8s-web-terminal是什么
如果你在使用docker或者kubernetes来运行你的服务，你会需要登录到容器里面来执行一些命令。   
你有以下三种方法：   
- 用`docker exec`命令来登录
这种方法很麻烦，因为你需要先找到容器所在节点IP，然后登录到节点上，然后再找出容器ID。。。
- 用三方管理工具来登录，比如Rancher
这种方法很方便，但是不够灵活，因为你可能需要根据你自己的需求来开发工具。
- 自己开发工具来登录
这种方法既方便有灵活。
现在你可以使用第三种方法来实现你的需求了。感觉很酷！对不对？
# k8s-web-terminal的工作原理
k8s-web-terminal在docker rest api和前端之间简历Websocket连接。   
前端如果有输入就会被发送到容器，容器如果有输出就会被发送给前端。   
前端内嵌了一个三方插件，叫xterm.js，这个插件负责渲染控制台。  
## docker REST API
k8s-web-terminal使用了三个API接口，
- [Create an exec instance](https://docs.docker.com/engine/api/v1.30/#operation/ContainerExec)
- [Start an exec instance](https://docs.docker.com/engine/api/v1.30/#operation/ExecStart)
- [Resize an exec instance](https://docs.docker.com/engine/api/v1.30/#operation/ExecResize)
## xterm.js插件
xterm.js插件代码已经被包含在了项目代码里面，位置在`k8s-web-terminal/src/static/xterm.js`。    
k8s-web-terminal兼容xterm.js的版本小于1.1.0，因为高版本的xterm.js使用了typescript。    
如果你想替换xterm.js成高版本，你需要自行修改代码。   
# 怎样运行k8s-web-terminal
## 启用docker rest api相关配置
docker默认是没有启用暴露rest api监听接口的，所以你需要先修改配置启用。
## 修改证书和私钥文件路径
证书和私钥文件路径在代码中，位置是`k8s-web-terminal/src/plugin/docker/terminal.go`。
```
defaultCertPath = "/root/ssl-docker/cert.pem"
defaultKeyPath = "/root/ssl-docker/key.pem"
```
## 运行
```
export GOPATH=~/go:$YOURPATH
go build -o terminal k8s-web-terminal/src/start.go
./terminal
```
服务起来后会监听在8080端口，你可以通过 http://localhost:8080 来访问
## 修改node IP和container ID
服务起来后，你可以在浏览器看到如下所示，  
![image](https://raw.githubusercontent.com/zhanghaohao/pictures/master/terminal-1.png)

把node IP和container ID替换为相应的。  
- node IP: 容器所在节点的IP地址
- container ID: 容器ID      

点击“执行命令行”按钮，在正确建立Websocket连接后会看到如下所示，  

![image](https://raw.githubusercontent.com/zhanghaohao/pictures/master/terminal-2.png)

# 如何把k8s-web-terminal嵌入到你自己的代码中
