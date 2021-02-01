package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

type Empty struct{}

func main() {
	cli, err := client.NewClient("tcp://192.168.172.2:2345", "v1.39", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	config := &container.Config{
		WorkingDir: "/app",
		ExposedPorts: map[nat.Port]struct{}{
			"80/tcp": Empty{},
		},
		Image: "alpine:3.13",
		Cmd:   []string{"./myserver"},
	}
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			"80/tcp": []nat.PortBinding{
				nat.PortBinding{HostPort: "8080"}, //宿主机的端口
			},
		},
		Binds: []string{"/home/custer/myweb:/app"},
	}
	ret, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil,"myweb")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.ContainerStart(ctx, ret.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("容器启动成功,ID是:", ret.ID)
}
