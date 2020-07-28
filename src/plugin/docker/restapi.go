package docker

import (
	"net/http"
	"util/logger"
	"util/ctx"
	"github.com/gorilla/websocket"
	"time"
	"strconv"
	"html/template"
)

const (
	DefaultDockerAPIScheme	string = "https"
	DefaultDockerAPIPort 	string	= "2375"
)

type DockerClient struct {
	Scheme 				string
	Host 				string
	Port 				string
}

func PermissionVerify(w http.ResponseWriter, r *http.Request)  {
	// todo:
	return
}

func GetDockerClient(scheme, nodeIP, port string) (dockerClient DockerClient) {
	if len(scheme) == 0 {
		scheme = DefaultDockerAPIScheme
	}
	if len(port) == 0 {
		port = DefaultDockerAPIPort
	}
	dockerClient = DockerClient{
		Scheme: scheme,
		Host: nodeIP,
		Port: port,
	}
	return
}

func CreateContainer(w http.ResponseWriter, r *http.Request) {
	containerID := r.URL.Query().Get("containerID")
	nodeIP := r.URL.Query().Get("nodeIP")
	command := r.URL.Query().Get("command")
	user := r.URL.Query().Get("user")
	if len(containerID) == 0 || len(nodeIP) == 0 || len(command) == 0 {
		ctx.WriteJSON(w, http.StatusBadRequest, "参数错误，containerID、nodeIP和command不能为空", "")
		return
	}
	client := GetDockerClient("", nodeIP, "")
	id, err := client.CreateExec(containerID, command, user)
	if err != nil {
		logger.Error.Println(err)
		ctx.WriteJSON(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	logger.Info.Println("created id: ", id)
	ctx.WriteJSON(w, http.StatusOK, "", map[string]string{"id": id})
	return
}

func ShellContainer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	nodeIP := r.URL.Query().Get("nodeIP")
	if len(id) == 0 || len(nodeIP) == 0 {
		ctx.WriteJSON(w, http.StatusBadRequest, "参数错误，id和nodeIP不能为空", "")
		return
	}
	client := GetDockerClient("", nodeIP, "")
	input := make(chan []byte)
	var upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error.Println(err)
		ctx.WriteJSON(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	output, err := client.ExecStart(id, input)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		return
	}
	go func() {
		defer ws.Close()
		for {
			if data, ok := <-output; ok {
				//logger.Info.Println(string(data))
				ws.WriteMessage(websocket.TextMessage, data)
			} else {
				break
			}
		}
	}()
	go func() {
		for {
			_, data, err := ws.ReadMessage()
			if err != nil {
				logger.Error.Println(err)
				//send EOF to close chan
				input <- []byte("EOF")
				close(input)
				return
			}
			input <- data
		}
	}()
	return
}

func ResizeContainer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	nodeIP := r.URL.Query().Get("nodeIP")
	if len(id) == 0 || len(nodeIP) == 0 {
		ctx.WriteJSON(w, http.StatusBadRequest, "参数错误,id和nodeIP不能为空", "")
		return
	}
	cols, err := strconv.Atoi(r.URL.Query().Get("cols"))
	if err != nil {
		logger.Error.Println(err)
		ctx.WriteJSON(w, http.StatusBadRequest, "参数错误:"+err.Error(), "")
		return
	}
	rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
	if err != nil {
		logger.Error.Println(err)
		ctx.WriteJSON(w, http.StatusBadRequest, "参数错误:"+err.Error(), "")
		return
	}
	if cols <= 0 || rows <= 0 {
		ctx.WriteJSON(w, http.StatusBadRequest, "参数错误,cols和rows必须大于零", "")
		return
	}
	//logger.Info.Println("resize: ", id)
	client := GetDockerClient("", nodeIP, "")
	err = client.ExecResize(id, cols, rows)
	if err != nil {
		ctx.WriteJSON(w, http.StatusInternalServerError, err.Error(), "")
		return
	}
	ctx.WriteJSON(w, http.StatusOK, "success", "")
	return
}

func Index(w http.ResponseWriter, r *http.Request)  {
	index, err := template.ParseFiles("src/template/index.html")
	if err != nil {
		logger.Error.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	err = index.Execute(w, nil)
	if err != nil {
		logger.Error.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func Terminal(w http.ResponseWriter, r *http.Request) {
	index, err := template.ParseFiles("src/template/terminal.html")
	if err != nil {
		logger.Error.Println(err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	err = index.Execute(w, nil)
	if err != nil {
		logger.Error.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}