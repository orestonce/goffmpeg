package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	for _, item := range gDownloadList {
		fileName := item.MustGetFileName()
		if item.Sha256Hash != "" && getFileSha256(fileName) == item.Sha256Hash {
			fmt.Println("skip", fileName)
			continue
		}
		fmt.Println("begin download", fileName)

		resp, err := http.Get(item.Url)
		if err != nil {
			panic(err)
		}
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			panic(err)
		}
		resp.Body.Close()
		var value string
		{
			tmp := sha256.Sum256(content)
			value = hex.EncodeToString(tmp[:])
		}
		fmt.Println(item.Url, value)
		if item.Sha256Hash != "" && item.Sha256Hash != value {
			panic("Invalid download data: " + item.Url)
		}
		err = os.WriteFile(fileName, content, 0666)
		if err != nil {
			panic(err)
		}
	}
}

type DownloadItem struct {
	Url        string
	Sha256Hash string
}

func (item DownloadItem) MustGetFileName() string {
	tmp := strings.LastIndexByte(item.Url, '/')
	if tmp < 0 {
		panic("Invalid Url: " + item.Url)
	}
	return "w_" + item.Url[tmp+1:]
}

func getFileSha256(fileName string) string {
	fin, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	defer fin.Close()

	obj := sha256.New()
	_, err = io.Copy(obj, fin)
	if err != nil {
		return ""
	}
	tmp := obj.Sum(nil)
	v := hex.EncodeToString(tmp[:])
	return v
}

var gDownloadList = []DownloadItem{
	{
		Url:        "https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0/win32-ia32.gz",
		Sha256Hash: "e606b5a2b8aa4e6165b9765c6b5a3dabb6ac9e50cdcc2f7db8e57662ac45fde4",
	},
	{
		Url:        "https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0/linux-ia32.gz",
		Sha256Hash: "82996afde27459bba7c8204161110ffe74790ec1fdf5902f7ed8593fd050bd33",
	},
	{
		Url:        "https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0/linux-arm.gz",
		Sha256Hash: "1764dab0b427b5a8b80722e841e4202d652570532bd97754a0b81779e721b2d5",
	},
	{
		Url:        "https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0/darwin-x64.gz",
		Sha256Hash: "4c613760fa98ee60ef7177a9946ec2e92691d2a53f4e5e160a387bdc9792e573",
	},
	{
		Url:        "https://github.com/eugeneware/ffmpeg-static/releases/download/b5.0/darwin-arm64.gz",
		Sha256Hash: "ecac76e3fca84a5a04fbb3a9cba51e7f399ac290b8b6d6e8410145fda019ccb2",
	},
}
