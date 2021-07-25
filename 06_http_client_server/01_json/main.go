package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type User struct {
	Name string
	Age  int
	info string
}

func marshal() {
	u1 := User{
		Name: "Jane",
		Age:  25,
		info: "Sensitive data",
	}

	b, err := json.Marshal(u1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func unmarshal() {
	raw := []byte(`{"name": "John", "age": 33}`)
	u2 := User{}

	json.Unmarshal(raw, &u2)

	fmt.Println(u2)
}

func marshalWithTags() {
	type device struct {
		ID      string `json:"device_id"`
		Enabled bool   `json:"enabled,omitempty"`
	}

	d := device{
		ID:      "123",
		Enabled: true, // set to false
	}

	b, _ := json.Marshal(d)
	fmt.Println(string(b))
}

func unstructuredJSON() {
	raw := `
	[
		{
			"id": 123,
			"x": 66.66,
			"y": 10,
			"info": {"name": "Camera"},
			"ip": "192.168.1.4"
		}
	]
	`

	var j interface{}
	json.Unmarshal([]byte(raw), &j)

	fmt.Println(j)

	objects := j.([]interface{})
	firstObject := objects[0].(map[string]interface{})

	fmt.Printf("Type of y is %T and value %v\n", firstObject["y"], firstObject["y"])
	y := firstObject["y"].(float64)

	fmt.Println(y)

	info := firstObject["info"].(map[string]interface{})
	fmt.Println(info["name"])

	ipAddr := net.ParseIP(firstObject["ip"].(string))
	fmt.Printf("Type of ip address is %T and value is %v\n", ipAddr, ipAddr)
}

type Device struct {
	ID string
}

func (d *Device) UnmarshalJSON(data []byte) error {
	type incomingDeviceFormat struct {
		FirstID  string
		SecondID string
	}

	v := incomingDeviceFormat{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	d.ID = fmt.Sprintf("%s:%s", v.FirstID, v.SecondID)
	return nil
}

func customUnmarshal() {
	data := []byte(`{"firstID": "ABC", "secondID": "123"}`)
	d := Device{}

	json.Unmarshal(data, &d)

	fmt.Println(d.ID)
}

func main() {
	marshal()
	// unmarshal()
	// marshalWithTags()
	// unstructuredJSON()
	// customUnmarshal()
}
