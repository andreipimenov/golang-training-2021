package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	req, err := NewRequest().Method(http.MethodGet).URL("http://127.0.0.1:8000/ping").Build()
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(b))
}
