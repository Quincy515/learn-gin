package scripts

import (
	"github.com/docker/docker/client"
	"log"
)

func getClient() *client.Client  {
	cli, err := client.NewClient("tcp://192.168.172.2:2345", "v1.39", nil, nil)
	if err!=nil{
		log.Fatal(err)
	}
	return cli
}