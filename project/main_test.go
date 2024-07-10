package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Setup(t *testing.T) *http.Client {
	t.Helper()

	go main()
	time.Sleep(100 * time.Millisecond)

	client := &http.Client{}
	return client
}

func TestApp(t *testing.T) {
	client := Setup(t)

	response, err := client.Get("http://localhost:8080/app")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	containsTitle := strings.Contains(string(body), "<h1>Welcome to Chirpy</h1>")
	if !containsTitle {
		t.Errorf("expected %t | got %t", true, containsTitle)
	}
}

func TestAssets(t *testing.T) {
	client := Setup(t)

	response, err := client.Get("http://localhost:8080/app/assets/")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	containsLink := strings.Contains(string(body), "<a href=\"logo.png\">logo.png</a>")
	if !containsLink {
		t.Errorf("expected %t | got %t", true, containsLink)
	}
}

func TestAssetsImage(t *testing.T) {
	client := Setup(t)

	response, err := client.Get("http://localhost:8080/app/assets/logo.png")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	expectedImageContentLength := "35672"
	actualImageContentLength := response.Header.Get("Content-Length")
	t.Log(actualImageContentLength)

	if actualImageContentLength != expectedImageContentLength {
		t.Errorf("got %s | want %s", actualImageContentLength, expectedImageContentLength)
	}
}

func TestHealthz(t *testing.T) {
	client := Setup(t)

	response, err := client.Get("http://localhost:8080/api/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	expectedContentType := "text/plain; charset=utf-8"
	actualContentType := response.Header.Get("Content-Type")

	if actualContentType != expectedContentType {
		t.Errorf("got %s | want %s", actualContentType, expectedContentType)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	if string(body) != http.StatusText(http.StatusOK) {
		t.Errorf("expected %s | got %s", http.StatusText(http.StatusOK), string(body))
	}
}

func TestMetrics(t *testing.T) {
	client := Setup(t)

	visitCount := 5

	for i := 0; i < visitCount; i++ {
		response, err := client.Get("http://localhost:8080/app")
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
	}

	response, err := client.Get("http://localhost:8080/api/admin/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	containsVisitCount := strings.Contains(string(body), fmt.Sprintf("<p>Chirpy has been visited %d times!</p>", visitCount))

	if !containsVisitCount {
		t.Errorf("expected %t | got %t", true, containsVisitCount)
	}
}

func TestReset(t *testing.T) {
	client := Setup(t)

	for i := 0; i < 5; i++ {
		response, err := client.Get("http://localhost:8080/app")
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
	}

	response, err := client.Get("http://localhost:8080/api/admin/reset")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	expectedBody := "Hits reset to 0"

	if string(body) != expectedBody {
		t.Errorf("expected %s | got %s", expectedBody, string(body))
	}
}

func TestMethodRestriction(t *testing.T) {
	client := Setup(t)

	urls := []string{
		"http://localhost:8080/api/healthz",
		"http://localhost:8080/api/admin/metrics",
		"http://localhost:8080/api/admin/reset",
	}

	for _, url := range urls {
		response, err := client.Post(url, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("expected status code %d | got %d", http.StatusMethodNotAllowed, response.StatusCode)
		}
	}
}

func TestValidateChirp(t *testing.T) {
	client := Setup(t)

	testCases := []struct {
		name               string
		requestBody        Chirp
		expectedStatusCode int
		expectedBody       Response
	}{
		{
			name:               "valid chirp",
			requestBody:        Chirp{Body: "I had something interesting for breakfast"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       Response{Valid: true},
		},
		{
			name:               "invalid chirp",
			requestBody:        Chirp{Body: "lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       Response{Error: "Chirp is too long"},
		},
		{
			name:               "profane chirp 1",
			requestBody:        Chirp{Body: "This is a kerfuffle opinion I need to share with the world"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       Response{CleanedBody: "This is a **** opinion I need to share with the world"},
		},
		{
			name:               "profane chirp 2",
			requestBody:        Chirp{Body: "I hear Mastodon is better than Chirpy. sharbert I need to migrate"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       Response{CleanedBody: "I hear Mastodon is better than Chirpy. **** I need to migrate"},
		},
		{
			name:               "profane chirp 3",
			requestBody:        Chirp{Body: "I really need a kerfuffle to go to bed sooner, Fornax !"},
			expectedStatusCode: http.StatusOK,
			expectedBody:       Response{CleanedBody: "I really need a **** to go to bed sooner, **** !"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBodyJson, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			// we use `bytes.NewReader` on the request body JSON for three main reasons:
			// 1. Interface compatibility: It converts the byte slice to an `io.Reader`, which is required by `http.NewRequest`.
			// 2. Efficiency: It creates a reader without copying the data, using less memory than alternatives like `bytes.Buffer`.
			// 3. Simplicity: It provides a read-only view of the data, which is sufficient for sending an HTTP request.
			// this approach is memory-efficient, simple, and aligns with Go's idiomatic practices for handling byte slices in HTTP requests.
			response, err := client.Post("http://localhost:8080/api/validate_chirp", "application/json", bytes.NewReader(requestBodyJson))
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.expectedStatusCode {
				t.Errorf("expected status code %d | got %d", tc.expectedStatusCode, response.StatusCode)
			}

			if response.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type %s | got %s", "application/json", response.Header.Get("Content-Type"))
			}

			byteData, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatal(err)
			}

			var responseJSON Response
			err = json.Unmarshal(byteData, &responseJSON)
			if err != nil {
				t.Fatal(err)
			}

			if responseJSON != tc.expectedBody {
				t.Errorf("expected %v | got %v", tc.expectedBody, responseJSON)
			}
		})
	}
}

func TestPostChirps(t *testing.T) {
	DeleteDBFile(t)

	client := Setup(t)

	requestBody := Chirp{Body: "I had something interesting for breakfast"}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	response, err := client.Post("http://localhost:8080/api/chirps", "application/json", bytes.NewReader(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		t.Errorf("expected status code %d | got %d", http.StatusCreated, response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type %s | got %s", "application/json", response.Header.Get("Content-Type"))
	}

	byteData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	var responseJSON Response
	err = json.Unmarshal(byteData, &responseJSON)
	if err != nil {
		t.Fatal(err)
	}

	if responseJSON.Body != requestBody.Body {
		t.Errorf("expected %v | got %v", requestBody.Body, responseJSON.Body)
	}
}

func TestGetChirps(t *testing.T) {
	// TODO (!) in order to pass the above POST function has to run
	client := Setup(t)

	response, err := client.Get("http://localhost:8080/api/chirps")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status code %d | got %d", http.StatusOK, response.StatusCode)
	}

	if response.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type %s | got %s", "application/json", response.Header.Get("Content-Type"))
	}

	byteData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	var chirps []Chirp
	err = json.Unmarshal(byteData, &chirps)
	if err != nil {
		t.Fatal(err)
	}

	expectedChirps := []Chirp{
		{
			ID:   1,
			Body: "I had something interesting for breakfast",
		},
	}

	if !reflect.DeepEqual(chirps, expectedChirps) {
		t.Errorf("expected %v | got %v", expectedChirps, chirps)
	}
}

func TestDB(t *testing.T) {
	DeleteDBFile(t)

	dbDisk, err := NewDB("database.json")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		_, err := dbDisk.CreateChirp(fmt.Sprintf("body-value-%d", i+1))
		if err != nil {
			t.Fatal(err)
		}
	}

	chirps, err := dbDisk.GetChirps()
	if err != nil {
		t.Fatal(err)
	}

	expectedChirps := []Chirp{
		{ID: 1, Body: "body-value-1"},
		{ID: 2, Body: "body-value-2"},
		{ID: 3, Body: "body-value-3"},
	}

	if !reflect.DeepEqual(expectedChirps, chirps) {
		t.Errorf("expected %v | got %v", expectedChirps, chirps)
	}
}

func DeleteDBFile(t *testing.T) {
	t.Helper()
	_, err := os.Stat("database.json")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return
		}
	}
	err = os.Remove("database.json")
	if err != nil {
		t.Fatal(err)
	}
}
