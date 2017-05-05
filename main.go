package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

const (
	uploadURL = "https://i.nuuls.com/upload"
)

func main() {
	if len(os.Args) > 1 {
		url, err := upload(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(url)
		os.Exit(0)
	}
	fmt.Println("usage: ni dank_pepe.png")
}

func upload(p string) (string, error) {
	file, err := os.Open(p)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	mp := multipart.NewWriter(buf)
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", getMimeType(file))
	header.Set("Content-Disposition", "form-data; name=xddd; filename=xddd.png")
	w, err := mp.CreatePart(header)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return "", err
	}
	mp.Close()
	req, err := http.NewRequest(http.MethodPost, uploadURL, buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", mp.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode > 201 {
		return "", fmt.Errorf("unexpected status code: %s", res.Status)
	}
	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func getMimeType(file *os.File) string {
	file.Seek(0, 0)
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return ""
	}
	buf = buf[:n]
	mimeType := http.DetectContentType(buf)
	file.Seek(0, 0)
	return strings.Split(mimeType, ";")[0]
}
