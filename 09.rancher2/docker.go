package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
)

func main() {
	cli, err := client.NewClient("tcp://192.168.172.2:2345", "v1.39", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range images {
		fmt.Println(image.ID, image.Labels)
	}
}
