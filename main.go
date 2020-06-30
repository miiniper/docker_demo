package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	fmt.Println("start server ... ")

	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	//ListImages(cli)
	BuildImage(cli)

}

func ListImages(cli *client.Client) {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {

		fmt.Println(image.RepoTags)
	}

}

func BuildImage(cli *client.Client) {

	dockerBuildContext, err := os.Open("/home/han/ifeng/k8scheck/Dockerfile")
	if err != nil {
		fmt.Println("11121: ", err)
	}
	defer dockerBuildContext.Close()

	_, err = cli.ImageBuild(context.Background(), dockerBuildContext, types.ImageBuildOptions{

		Dockerfile: "Dockerfile",
	})
	if err != nil {
		fmt.Println("11111: ", err)
	}
}
