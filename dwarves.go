package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

type ResponseArray struct {
	Dwarves []string `json:"dwarves"`
}

type ResponseDwarf struct {
	Dwarf `json:"dwarf"`
}

type ResponseError struct {
	Msg string `json:"error"`
}

const serviceUrl string = "https://thedwarves.pusherplatform.io/api/dwarves"
const retries int = 5

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

// getMap calls getJson and restructures the response to a
// *map[string]Dwarf. The parameters are the same as getJson
func getMap(url string, retries int) (*map[string]Dwarf, error) {
	// Get JSON information on the dwarves.
	dwarfJson, err := getJson(url, retries)
	if err != nil {
		return nil, err
	}

	if dwarfJson == nil {
		return nil, fmt.Errorf("Unknown error getting json")
	}

	//fmt.Printf("body: %s\n\n\n", string(*dwarfJson))

	var dwarves DwarfData
	err = json.Unmarshal(*dwarfJson, &dwarves)
	if err != nil {
		return nil, err
	}

	dwarfMap := make(map[string]Dwarf)
	for _, dwarf := range dwarves.Dwarves {
		//fmt.Printf("adding %s to map\n", dwarf.Name)
		dwarfMap[dwarf.Name] = dwarf
	}

	return &dwarfMap, nil
}

// getDwarves returns a json array of all the dwarf
// names in the map returned by dwarfMap().
// The json structure follows the ResponseArray struct
func getDwarves() (*[]byte, error) {
	dwarfMap, err := getMap(serviceUrl, retries)
	if err != nil {
		fmt.Println(err)
	}
	if dwarfMap == nil {
		fmt.Println("Unknown error getting dwarf map")
	}

	var dwarves []string
	for name, _ := range *dwarfMap {
		dwarves = append(dwarves, name)
		//fmt.Printf("Dwarf Name: %s, Value: %s\n", name, val)
	}

	response := ResponseArray{
		Dwarves: dwarves,
	}

	json, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}

	return &json, nil

}

func getDwarf(name string) (*[]byte, error) {
	dwarfMap, err := getMap(serviceUrl, retries)
	if err != nil {
		return nil, err
	}
	if dwarfMap == nil {
		return nil, fmt.Errorf("Unknown error getting dwarf map")
	}

	dwarf, present := (*dwarfMap)[name]

	var output []byte
	if present {
		response := ResponseDwarf{
			dwarf,
		}

		output, err = json.Marshal(response)
		if err != nil {
			return nil, err
		}
	} else {
		response := ResponseError{
			Msg: "dwarf not found",
		}

		output, err = json.Marshal(response)
		if err != nil {
			return nil, err
		}
	}

	return &output, nil
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Path[len("/api/"):]

	if len(request) == 0 {

		out, err := getDwarves()
		if err != nil {
			log.Print(err)
		}
		fmt.Fprintf(w, string(*out))
	} else {

		out, err := getDwarf(request)
		if err != nil {
			log.Print(err)
		}
		fmt.Fprintf(w, string(*out))
	}
}

func main() {
	http.HandleFunc("/api/", apiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
