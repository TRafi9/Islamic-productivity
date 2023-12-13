package main

import (
	"fmt"
	"net/http"
)

func init() {
	callUpdatePt()
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func callUpdatePt() (int, error) {
	url := "https://living-sacred-skunk.ngrok-free.app/api/v1/updatePt"
	resp, err := http.Get(url)
	if err != nil {
		return 0, nil
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return http.StatusOK, nil
	} else {
		return resp.StatusCode, fmt.Errorf("unexpected status code in updatePt call from cloud function, %d", resp.StatusCode)
	}
}
