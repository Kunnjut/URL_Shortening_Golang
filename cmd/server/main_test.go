package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMainPageHandler(t *testing.T) {
	tests := []struct {
		name                string
		method              string
		body                string
		expectedCode        int
		expectedBody        string
		expectedContentType string
	}{
		{
			name:                "Main1",
			method:              "GET",
			body:                "",
			expectedCode:        405,
			expectedBody:        "Only POST method is allowed\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name:                "Main2",
			method:              "DELETE",
			body:                "",
			expectedCode:        405,
			expectedBody:        "Only POST method is allowed\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
		{
			name:                "Main3",
			method:              "POST",
			body:                "https://dzen.ru/",
			expectedCode:        201,
			expectedBody:        "http://example.com/8d26b0\n",
			expectedContentType: "text/plain; charset=utf-8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/main", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			mainPage(w, r)
			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			assert.Equal(t, tt.expectedContentType, res.Header.Get("Content-Type"))
			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedBody, string(resBody))

		})
	}
}

func TestIdPageHandler(t *testing.T) {
	mapURL["https://dzen.ru/"] = "8d26b0"

	tests := []struct {
		name             string
		method           string
		target           string
		expectedCode     int
		expectedLocation string
	}{
		{
			name:             "Id1",
			method:           "POST",
			target:           "http://localhost:8080/8d2er",
			expectedCode:     405,
			expectedLocation: "",
		},
		{
			name:             "Id2",
			method:           "DELETE",
			target:           "http://localhost:8080/8d2er",
			expectedCode:     405,
			expectedLocation: "",
		},
		//{
		//	name:             "Id3",
		//	method:           "GET",
		//	target:           "/8d26b0",
		//	expectedCode:     307,
		//	expectedLocation: "https://dzen.ru/", TODO test this case with location
		//},
		{
			name:             "Id4",
			method:           "GET",
			target:           "http://localhost:8080/8d223",
			expectedCode:     404,
			expectedLocation: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.target, nil)
			w := httptest.NewRecorder()
			idPage(w, r)
			assert.Equal(t, tt.expectedLocation, w.Header().Get("Location"))
			assert.Equal(t, tt.expectedCode, w.Result().StatusCode)
		})
	}
}
