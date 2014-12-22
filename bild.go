package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func upload(url string, data map[string]string,
	paramname string, filename string,
) (s []string, err error) {
	client := &http.Client{}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile(paramname, filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error open file")
		return
	}

	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return
	}

	for k, v := range data {
		bodyWriter.WriteField(k, v)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest("POST", url, bodyBuf)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "go-bild/0.1.0")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	body := string(respBody)
	if resp.StatusCode >= 400 {
		return
	}
	// fmt.Println(body)
	urls := strings.Split(body, "\n")
	// fmt.Println(urls[0])
	// fmt.Println(urls[len(urls)-1])
	// fmt.Println("\n")
	s = []string{
		urls[0],
		urls[len(urls)-1],
	}

	return
}

func main() {
	data := map[string]string{
		"t":      "1",
		"C1":     "ON",
		"upload": "1",
	}
	url := "http://www.bild.me/index.php"
	fileSlice := []string{}
	files := os.Args[1:]

	// 支持通配符
	for _, file := range files {
		matches, err := filepath.Glob(file)
		if err == nil {
			fileSlice = append(fileSlice, matches...)
		}
	}
	if len(fileSlice) == 0 {
		fmt.Println("need files: bild FILE [FILE ...]")
		os.Exit(1)
	}
	var wg sync.WaitGroup

	for _, f := range fileSlice {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			s, err := upload(url, data, "F1", f)
			if err == nil {
				fmt.Println(f+":", s[1])
			}
		}(f)
	}
	wg.Wait()
}
