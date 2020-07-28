package docker

import (
	"testing"
	"util/logger"
)

func TestDockerClient_CreateExec(t *testing.T) {
	dockerClient := &DockerClient{
		Scheme: "https",
		Host:   "10.2.62.1",
		Port:   "2375",
	}
	id, err := dockerClient.CreateExec("34978f84a385", "/bin/bash", "root")
	if err != nil {
		t.Error(err)
	} else {
		logger.Info.Println(id)
	}
}

func TestDockerClient_ExecResize(t *testing.T) {
	dockerClient := &DockerClient{
		Scheme: "https",
		Host:   "10.2.62.1",
		Port:   "2375",
	}
	err := dockerClient.ExecResize("709e16dac4d9ff7b2482f35bcb718572ed3cc00f09be7d44fc59901accf5c124", 157, 24)
	if err != nil {
		t.Error(err)
	}
}
