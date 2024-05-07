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
)

func main() {

	//{
	//	"uuid": "484185c5-0bec-11ef-b26a-0242ac160003",
	//	"url": "https://faces.nyc3.digitaloceanspaces.com/484185c5-0bec-11ef-b26a-0242ac160003.jpg",
	//	"galleries": []
	//}
	url := "https://api.luxand.cloud/v2/person/288a0691-0bed-11ef-b26a-0242ac160003"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, errFile1 := os.Open("/Users/ankitapanchal/Desktop/Ankita.jpg")
	defer file.Close()
	part1,
		errFile1 := writer.CreateFormFile("photos", filepath.Base("/Users/ankitapanchal/Desktop/Ankita.jpg"))
	_, errFile1 = io.Copy(part1, file)
	if errFile1 != nil {
		fmt.Println(errFile1)
		return
	}
	_ = writer.WriteField("store", "1")
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("token", "22cbd6fcdd2c43d98351b3b40285997f")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
