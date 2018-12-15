package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

// Web service tests

func TestSelectUserHandler(t *testing.T) {
    testName := "handleUser"

    request, err := http.NewRequest("GET", "localhost:8080/u/", nil)
    if err != nil {
        t.Fatal(err)
    }

    request.Header.Add("Authentication", "6DA8F3AE1CE67DDEF5B55D062C246981")

    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(handleUser)
    handler.ServeHTTP(recorder, request)

    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("%v returned status of %v.", testName, status)
        return
    }

    expectedSubbody := "1234567890ABCDEF"
    responseBody := recorder.Body.String()

    if !strings.Contains(responseBody, expectedSubbody) {
        t.Errorf("%v provided an unexpected value: %v, without %v.", testName, responseBody, expectedSubbody)
        return
    }

    return
}

func TestSelectBranchHandler(t *testing.T) {
    testName := "handleBranch"

    request, err := http.NewRequest("GET", "localhost:8080/b/1234567890ABCDEF", nil)
    if err != nil {
        t.Fatal(err)
    }

    request.Header.Add("Authentication", "6DA8F3AE1CE67DDEF5B55D062C246981")

    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(handleBranch)
    handler.ServeHTTP(recorder, request)

    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("%v returned status of %v.", testName, status)
        return
    }

    expectedSubbody := "2018-12-09 14:21"
    responseBody := recorder.Body.String()

    if !strings.Contains(responseBody, expectedSubbody1) {
        t.Errorf("%v provided an unexpected value: %v, without %v.", testName, responseBody, expectedSubbody)
        return
    }

    return
}
