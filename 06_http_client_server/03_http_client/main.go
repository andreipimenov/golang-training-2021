package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type User struct {
	Name string
	Age  int
}

type Error struct {
	Error string
}

func main() {
	httpClient := &http.Client{
		Timeout: time.Duration(time.Minute * 5),
	}

	// POST

	u := &User{Name: "Rudolf", Age: 25}

	b, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/user", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		u := User{}
		err := json.Unmarshal(raw, &u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(u)
	} else {
		e := Error{}
		err := json.Unmarshal(raw, &e)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(e)
	}

	// GET

	// req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/user?name=Rudolf", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// res, err := httpClient.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer res.Body.Close()

	// raw, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if res.StatusCode == http.StatusOK {
	// 	u := User{}
	// 	err := json.Unmarshal(raw, &u)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(u)
	// } else {
	// 	e := Error{}
	// 	err := json.Unmarshal(raw, &e)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(e)
	// }

}
