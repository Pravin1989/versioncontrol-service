package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExecuteAndParseJSON(t *testing.T, router http.Handler, method, url, body string, expectedStatus int, v interface{}) (*httptest.ResponseRecorder, error) {
	response, statusErr := SendRequest(t, router, method, url, body, expectedStatus)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		return nil, err
	}
	if statusErr != nil {
		t.Errorf("response body : %s", responseBody)
	}
	if (response.Code == 201 || response.Code == 200) && len(responseBody) > 0 {
		marshalErr := json.Unmarshal(responseBody, &v)
		if marshalErr != nil {
			t.Error(marshalErr)
			return nil, marshalErr
		}
	}
	return response, nil
}

func SendRequest(t *testing.T, router http.Handler, method, url, body string, expectedStatus int) (*httptest.ResponseRecorder, error) {
	resRec := httptest.NewRecorder()
	reqBody := bytes.NewReader([]byte(body))
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	router.ServeHTTP(resRec, req)
	statusErr := checkStatus(t, expectedStatus, resRec)
	return resRec, statusErr
}

func checkStatus(t *testing.T, expectedStatus int, rec *httptest.ResponseRecorder) error {
	if rec.Code != expectedStatus {
		err := fmt.Errorf("Expected staus code %d but got %d", expectedStatus, rec.Code)
		t.Error(err)
		return err
	}
	return nil
}
