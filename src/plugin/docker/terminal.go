package docker

import (
	"strings"
	"io/ioutil"
	"net/http"
	"util/logger"
	"encoding/json"
	"net/http/httputil"
	"time"
	"fmt"
	"crypto/tls"
	"k8s.io/client-go/transport"
)

const (
	defaultCertPath = "/root/ssl-docker/cert.pem"
	defaultKeyPath = "/root/ssl-docker/key.pem"
)

type ExecResponse struct {
	ID 					string
}

func (client *DockerClient) String() (apiURL string) {
	apiURL = client.Scheme + "://" + client.Host + ":" + client.Port
	return
}

func getTLSConfig(certPath string, keyPath string) (tlsConfig *tls.Config, err error) {
	config := &transport.Config{
		TLS: transport.TLSConfig{
			CertFile: certPath,
			KeyFile: keyPath,
			Insecure: true,    // used to skip cert verification
		},
	}
	tlsConfig, err = transport.TLSConfigFor(config)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	return
}

func clientWithTLS(certPath string, keyPath string) (client *http.Client, err error) {
	client = new(http.Client)
	tlsConfig, err := getTLSConfig(certPath, keyPath)
	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	return
}

func (client *DockerClient) CreateExec(containerID string, cmd string, user string) (id string, err error) {
	var jsonBody = strings.NewReader(`{
		"AttachStdin": true,
		"AttachStdout": true,
		"AttachStderr": true,
		"DetachKeys": "ctrl-p,ctrl-q",
		"Tty": true,
		"User": "` + user + `",
		"Cmd": [
		"` + cmd + `"
		]
	}`)
	clientWithTLS, err := clientWithTLS(defaultCertPath, defaultKeyPath)
	if err != nil {
		return
	}
	res, err := clientWithTLS.Post(client.String() + "/containers/"+containerID+"/exec", "application/json;charset=utf-8", jsonBody)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	//logger.Info.Printf("%+v", res)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	var result ExecResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	id = result.ID
	if len(id) == 0 {
		err = fmt.Errorf(string(body))
		return
	}
	return
}

func (client *DockerClient) ExecStart(id string, input chan []byte) (output chan []byte, err error) {
	output = make(chan []byte)
	url := client.String() + "/exec/" + id + "/start"
	//logger.Info.Println(url)
	req, _ := http.NewRequest("POST", url, strings.NewReader(
		`{
			"Detach": false,
			"Tty": true
		}`))
	tlsConfig, err := getTLSConfig(defaultCertPath, defaultKeyPath)
	dial, err := tls.Dial("tcp", client.Host + ":" + client.Port, tlsConfig)
	if err != nil {
		logger.Info.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	clientconn := httputil.NewClientConn(dial, nil)
	clientconn.Do(req)
	rwc, br := clientconn.Hijack()
	/*
	the following goroutines will continue run after the function return
	 */
	// receive input, and write to connection
	go func() {
		defer clientconn.Close()
		for {
			if data, ok := <-input; ok {
				rwc.Write(data)
			} else {
				break
			}
		}
	}()
	// read buffer as output
	go func() {
		defer rwc.Close()
		for {
			buf := make([]byte, 1024)
			_, err := br.Read(buf)
			if err != nil {
				logger.Error.Println(err)
				if err.Error() == "EOF" {
					output <- []byte("EOF")
					break
				}
				logger.Error.Println("Read Error: " + err.Error())
				break
			}
			output <- buf
			//Equal EOF
			if buf[0] == 69 && buf[1] == 79 && buf[2] == 70 {
				close(output)
				break
			}
			time.Sleep(500)
			buf = nil
		}
	}()
	return
}

func (client *DockerClient) ExecResize(id string, width int, height int) (err error) {
	clientWithTLS, err := clientWithTLS(defaultCertPath, defaultKeyPath)
	if err != nil {
		return
	}
	url := fmt.Sprintf(client.String()+"/exec/%s/resize?h=%d&w=%d", id, height, width)
	resp, err := clientWithTLS.Post(url, "application/json;charset=utf-8", nil)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error.Println(err)
		return
	}
	if len(body) != 0 {
		err = fmt.Errorf(string(body))
		logger.Error.Println(err)
		//logger.Error.Printf("%+v", resp)
		return
	}
	return
}