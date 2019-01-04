package main

import (
	"fmt"
	"io/ioutil"
	//"log"
	"encoding/json"
	"net/http"
	"time"
)

type Dwarf struct {
	Name    string `json:"name"`
	Birth   string `json:"birth"`
	Death   string `json:"death"`
	Culture string `json:"culture"`
}

type DwarfData struct {
	Dwarves []Dwarf `json:"dwarves"`
}

// getJson wraps a GET request in a retry function. This is because the dwarves are
// sometimes busy. Url specifies where to find the dwarves,  retries how many times to
// knock before giving up.
// The body of the get response is returned as a slice of bytes.
// There is a back off between retries.
func getJson(url string, retries int) (*[]byte, error) {
	// Build get request and response outside of retry func
	resp := &http.Response{}
	req, err := http.NewRequest(
		"GET",
		url,
		nil)

	if err != nil {
		return nil, fmt.Errorf("unable to create request", err)
	}

	err_retry := retry(retries, time.Second, func() error {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		s := resp.StatusCode

		switch {
		case s >= 500:
			// Retry (they were busy)
			return fmt.Errorf("server error (busy dwarves?): %v", s)
		default:
			// Dwarves happy
			return nil
		}
	})
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return &body, err_retry
}

//func getMap(url string, retries int) (*map[string]Dwarf, error) {
//}

func main() {
	// Get JSON information on the dwarves.
	dwarfJson, _ := getJson("https://thedwarves.pusherplatform.io/api/dwarves", 5)
	fmt.Printf("body: %s\n\n\n", string(*dwarfJson))

	var dwarves DwarfData
	err := json.Unmarshal(*dwarfJson, &dwarves)
	if err != nil {
		fmt.Println(err)
	}

	dwarfMap := make(map[string]Dwarf)
	for _, dwarf := range dwarves.Dwarves {
		fmt.Printf("adding %s to map\n", dwarf.Name)
		dwarfMap[dwarf.Name] = dwarf
	}

	for key, val := range dwarfMap {
		fmt.Printf("Dwarf Name: %s, Value: %s\n", key, val)
	}

}
