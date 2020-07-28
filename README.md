[<h1>中文版</h1>](https://github.com/zhanghaohao/k8s-web-terminal/blob/master/README_zh.md)
# What is k8s-web-terminal?
if you are using docker or kubernetes for your service, you will need to login the container and execute some commands.
you have some ways to achieve this.
- **login by `docker exec`**                   
Annoying, because you have to find the node, login to the node, find the container ...
- **login by management tool like Rancher**       
Enjoyable but not Flexible, because you can not develop by your own requirement.
- **login by implement your own tool**            
Enjoyable and Flexible.
Now you can login to the container by the third way. Cool! right?
# How does k8s-web-terminal work?
k8s-web-terminal builds websocket connection between docker REST API and frontend xterm.js plugin.  
Commands from the frontend are send to docker container, and the output of container is send back.  
If the container has any output, it can be sent to frontend immediately.
## docker REST API
There are three APIs for k8s-web-terminal.
- [Create an exec instance](https://docs.docker.com/engine/api/v1.30/#operation/ContainerExec)
- [Start an exec instance](https://docs.docker.com/engine/api/v1.30/#operation/ExecStart)
- [Resize an exec instance](https://docs.docker.com/engine/api/v1.30/#operation/ExecResize)
## xterm.js plugin
k8s-web-terminal uses xterm.js as frontend plugin of web terminal.  
xterm.js plugin has been included in this project, it is located at k8s-web-terminal/src/static/xterm.js.  
k8s-web-terminal is compatible with xterm.js with version of lower than 1.1.0, because xterm.js with higher version use typescript.  
If you want to replace xterm.js with higher version, you have to change the code accordingly.   
# How to run k8s-web-terminal?
## enable docker REST API
docker can be configured to expose REST API or not, so first you have to enable it in the configuration file.
## change certificate and private key path
certificate and private key path for docker are set in `k8s-web-terminal/src/plugin/docker/terminal.go`
```
defaultCertPath = "/root/ssl-docker/cert.pem"
defaultKeyPath = "/root/ssl-docker/key.pem"
```
## run 
```
export GOPATH=~/go:$YOURPATH
go build -o terminal k8s-web-terminal/src/start.go
./terminal
```
The service listen on :8080
You can open it in browser http://localhost:8080
## fill node IP and container ID
Here is the picture you will see after you start the service.

![image](https://raw.githubusercontent.com/zhanghaohao/pictures/master/terminal-1.png)

After Websocket connection is build, you can interact with docker container freely.

![image](https://raw.githubusercontent.com/zhanghaohao/pictures/master/terminal-2.png)

Replace the node IP and container ID with your own.
- node IP: the server IP hosts the container
- container ID: the container ID
# How you can use k8s-web-terminal in your own code?


