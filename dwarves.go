package main

import (
	"fmt"
	"io/ioutil"
	//"log"
	"net/http"
	"time"
)

type Dwarf struct {
	Name    string
	Birth   string
	Death   string
	Culture string
}

// getDwarves wraps a GET request in a retry function. This is because the dwarves are
// sometimes busy. Url specifies where to find the dwarves, and retries how many times to
// retry before giving up.
// The body of the get response is returned as a string.
// There is a back off between retries.
func getDwarves(url string, retries int) (string, error) {
	// build get request and response outside of retry func
	resp := &http.Response{}
	req, err := http.NewRequest(
		"GET",
		url,
		nil)

	if err != nil {
		return "", fmt.Errorf("unable to create request", err)
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
	return string(body), err_retry
}

func main() {
	// Get JSON information on the dwarves.
	body, _ := getDwarves("https://thedwarves.pusherplatform.io/api/dwarves", 5)
	fmt.Printf("Body: %s", body)
}
