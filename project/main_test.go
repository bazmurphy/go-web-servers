package main

import (
	"io"
	"net/http"
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
		t.Errorf("expected status %s | got %s", http.StatusText(http.StatusOK), response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	containsTitle := strings.Contains(string(body), "<h1>Welcome to Chirpy</h1>")
	if !containsTitle {
		t.Errorf("expected %v | got %v", true, containsTitle)
	}
}

func TestAssets(t *testing.T) {
	client := Setup(t)

	response, err := client.Get("http://localhost:8080/app/assets/logo.png")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %s | got %s", http.StatusText(http.StatusOK), response.Status)
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

	response, err := client.Get("http://localhost:8080/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %s | got %s", http.StatusText(http.StatusOK), response.Status)
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

	for i := 0; i < 5; i++ {
		response, err := client.Get("http://localhost:8080/app")
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
	}

	response, err := client.Get("http://localhost:8080/metrics")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %s | got %s", http.StatusText(http.StatusOK), response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))

	expectedBody := "Hits: 5"

	if string(body) != expectedBody {
		t.Errorf("expected %s | got %s", expectedBody, string(body))
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

	response, err := client.Get("http://localhost:8080/reset")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %s | got %s", http.StatusText(http.StatusOK), response.Status)
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
