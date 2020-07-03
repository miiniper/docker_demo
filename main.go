package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

	dir := " /home/han/ifeng/k8scheck/"

	buf := new(bytes.Buffer)
	if err := Tar(dir, buf); err != nil {
		fmt.Println("112233 == : ", err)
	}

	ctx := context.Background()

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		buf,
		types.ImageBuildOptions{
			PullParent: true,
			Dockerfile: "/home/han/ifeng/k8scheck/Dockerfile",
			Remove:     true,
		})
	if err != nil {
		fmt.Println("223344 == : ", err)
	}

	fmt.Println("334455 == ", imageBuildResponse)

}

func Tar(src string, writers ...io.Writer) error {

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	gzw := gzip.NewWriter(mw)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}
