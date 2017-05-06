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
	"regexp"
	"strings"
)

var (
	uploadURL = "https://i.nuuls.com/upload"
	urlRe     = regexp.MustCompile(`(https?:\/\/)?i\.nuuls\.com\/.+`)
)

func main() {
	if len(os.Args) > 1 {
		if match := urlRe.FindString(os.Args[1]); match != "" {
			err := download(match)
			if err != nil {
				exit("%v", err)
			} else {
				exit("got em")
			}
		}
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
	if len(os.Args) > 2 {
		return os.Args[2]
	}
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

func download(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %s", res.Status)
	}
	filename := ""
	if len(os.Args) > 2 {
		filename = os.Args[2]
	} else {
		re := regexp.MustCompile(`nuuls\.com\/([\w\-]+\.\w+)`)
		matches := re.FindStringSubmatch(url)
		if len(matches) > 1 {
			filename = matches[1]
		}
	}
	if filename == "" {
		fmt.Println("invalid filename")
		os.Exit(1)
	}
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, res.Body)
	return err
}

func exit(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(0)
}
