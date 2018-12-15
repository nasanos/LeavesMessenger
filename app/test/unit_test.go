package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// Web service tests

func TestSelectUserHandler(t *testing.T) {
    testName := "handleUser"

    request, err := http.NewRequest("GET", "localhost:8080/u/", nil)
    if err != nil {
        t.Fatal(err)
    }

    request.Header.Add("Authentication", "138CB39F47E8342D88293B86CA37F9E8")

    recorder := httptest.NewRecorder()
    handler := http.HandlerFunc(handleUser)
    handler.ServeHTTP(recorder, request)

    if status := recorder.Code; status != http.StatusOK {
        t.Errorf("%v returned status of %v.", testName, status)
        return
    }

    expectedBody := "{\"LastBranchKey\":\"1234567890ABCDEF\",\"SubscribedBranchKeys\":[\"\"],\"SubscribedBranchNames\":[\"\"], \"Error\": {\"Description\":\"\"}"
    responseBody := recorder.Body.String()

    if responseBody != expectedBody {
        t.Errorf("%v provided an unexpected value:\n    %v\n  instead of\n    %v.", testName, responseBody, expectedBody)
        return
    }

    return
}
