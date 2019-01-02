package main

import (
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
)

type Dwarf struct {
	Name    string
	Birth   string
	Death   string
	Culture string
}

func getDwarves() (string, error) {
	resp, err := http.Get("https://thedwarves.pusherplatform.io/api/dwarves")
	if err != nil {
		fmt.Printf("Error with GET: %s\n", err)
	}
	fmt.Println("get returned")
	fmt.Printf("Response: %s, Headers: %s\n", resp.Status, resp.Header)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}

func main() {
	fmt.Println("vim-go")
	body, _ := getDwarves()
	fmt.Printf("Body: %s", body)
}
