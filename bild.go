package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func upload(url string, data map[string]string,
	paramname string, filename string,
) error {
	client := &http.Client{}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile(paramname, filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error open file")
		return err
	}

	_, err = io.Copy(fileWriter, f)
	if err != nil {
		return err
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
		return err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(string(respBody))
	return nil

}

func main() {
	data := map[string]string{
		"t":      "1",
		"C1":     "ON",
		"upload": "1",
	}
	url := "http://www.bild.me/index.php"
	upload(url, data, "F1", "up-download.jpg")
}
