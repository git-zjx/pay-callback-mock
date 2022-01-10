package xhttp

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(url string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func Post(url string, data string, contentType string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, contentType, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), err
}
