package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

const (
	cmp  string = "error"
	url  string = "http://localhost:9001" //one of the backend servers
	host string = "domain1.com"
)

func GetAndCmp(resp *http.Response) error {
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return errors.New(fmt.Sprintf("error reading body of request ,%v", err))
	}

	if bytes.Contains(data, []byte(cmp)) {
		return errors.New(fmt.Sprintf("error in request, %v", string(data))) //print logs of RESTFULL services 
	}

	return nil
}

// 1 case of POST and compare
func TestPost(t *testing.T) {

	jsonBody := []byte(`{"ID": 1, "name": "Varun", "add": "Kotoor village"}`)

	req, err := http.NewRequest("POST", url, io.NopCloser(bytes.NewReader(jsonBody)))
	req.Host = host

	if err != nil {
		t.Errorf("error creating POST request")
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("error sending POST request")
		return
	}

	err = GetAndCmp(res)

	if err != nil {
		t.Fatalf("%v",err)
		return
	}

	return
}

// 1 case of GET and compare
func TestGet(t *testing.T) {

	req, err := http.NewRequest("GET", url, nil)
	req.Host = host

	if err != nil {
		t.Errorf("error creating GET request")
		return
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Errorf("error sending GET request")
		return
	}

	data, err := ioutil.ReadAll(res.Body)

	if bytes.Contains(data, []byte(cmp)) {
		t.Errorf("error in GET request")
		return
	}

	fmt.Println(string(data))

	return
}

// 2 cases of PUT with valid ID and invalid ID and compare
func TestPut(t *testing.T) {

	jsonBody1 := []byte(`{"ID": 1, "name": "Mayank", "add": "Kotoor village"}`)
	jsonBody2 := []byte(`{"ID": 2, "name": "Mayank", "add": "Kotoor village"}`)

	req1, err := http.NewRequest("PUT", url, io.NopCloser(bytes.NewReader(jsonBody1)))

	if err != nil {
		t.Errorf("error creating a PUT request")
		return
	}

	resp, err := http.DefaultClient.Do(req1)

	if err != nil {
		t.Errorf("error sending a PUT request")
		return
	}

	err = GetAndCmp(resp)

	if err != nil {
		t.Fatalf("%v",err)
		return
	}

	req1, err = http.NewRequest("PUT", url, io.NopCloser(bytes.NewReader(jsonBody2)))

	if err != nil {
		t.Errorf("error creating a PUT request")
		return
	}

	resp, err = http.DefaultClient.Do(req1)

	if err != nil {
		t.Errorf("error sending a PUT request")
		return
	}

	err = GetAndCmp(resp)

	if err != nil {
		t.Fatalf("%v",err)
		return
	}

}

//2 cases of DELETE with valid ID and invalid ID 
func TestDelete(t *testing.T) {

	const id string = "/1"

	req , err := http.NewRequest("DELETE",url+id,nil)

	if err!=nil {
		t.Errorf("error in creating DELETE request")
		return 
	}

	resp ,err := http.DefaultClient.Do(req)

	if err !=nil {
		t.Errorf("error in sending DELETE request")
		return 
	}

	err = GetAndCmp(resp)

	if err !=nil {
		t.Errorf(fmt.Sprintf("error in DELETE request , %v" ,err))
		return 
	}

	return 
}
