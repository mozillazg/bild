package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mozillazg/request"
)

func upload(url string, data map[string]string,
	paramname string, filename string,
) (s []string, err error) {
	client := &http.Client{}
	a := request.NewArgs(client)

	a.Data = data
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		fmt.Println("error open file")
		return
	}
	a.Files = []request.FileField{
		request.FileField{paramname, filename, f},
	}
	a.Headers = map[string]string{
		"User-Agent": "go-bild/0.1.0",
	}

	resp, err := request.Post(url, a)

	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := resp.Text()
	if err != nil {
		return
	}
	if !resp.OK() {
		fmt.Println("Response Status:", resp.Status)
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
