package scripts

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"log"
)

type Empty struct{}

func StartGoAPI() {
	cli := getClient()
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
		//PortBindings:map[nat.Port][]nat.PortBinding{   //不需要绑定端口
		//	"80/tcp":[]nat.PortBinding{
		//		nat.PortBinding{HostPort:"80"},//宿主机的端口
		//	},
		//},
		Binds: []string{"/home/custer/myweb:/app"},
	}
	//networkConfig:=&network.NetworkingConfig{
	//	EndpointsConfig: map[string]*network.EndpointSettings{
	//		"自定义网络":{
	//			IPAMConfig: &network.EndpointIPAMConfig{
	//				IPv4Address: "172.18.0.3",
	//			},
	//		},
	//	},
	//}
	ret, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "myweb")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.ContainerStart(ctx, ret.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GO API 容器启动成功,ID是:", ret.ID)
}
