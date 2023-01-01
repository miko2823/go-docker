package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	PhoneNumber int    `json:"phone_number"`
	City        string `json:"city"`
}
type Person struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

func main() {
	b := []byte(`{"name": "kaori", "address": {"phone_number": 9026220741, "city": "tokyo"}}`)
	var p Person
	if err := json.Unmarshal(b, &p); err != nil {
		fmt.Println(err)
	}
	fmt.Println(p)

	v, _ := json.Marshal(p)
	fmt.Println(string(v))
}
