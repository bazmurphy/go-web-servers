package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	go func() {
		main()
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("/app/index.html", func(t *testing.T) {
		client := &http.Client{}

		response, err := client.Get("http://localhost:8080/app")
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("expected status OK | got %v", response.Status)
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}

		containsTitle := strings.Contains(string(body), "<h1>Welcome to Chirpy</h1>")
		if !containsTitle {
			t.Errorf("expected %v | got %v", true, containsTitle)
		}
	})

	t.Run("/app/assets/logo.png", func(t *testing.T) {
		client := &http.Client{}

		response, err := client.Get("http://localhost:8080/app/assets/logo.png")
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("expected status %s | got %v", http.StatusText(http.StatusOK), response.Status)
		}

		expectedImageContentLength := "35672"
		actualImageContentLength := response.Header.Get("Content-Length")

		if actualImageContentLength != expectedImageContentLength {
			t.Errorf("got %s | want %s", actualImageContentLength, expectedImageContentLength)
		}
	})

	t.Run("/healthz", func(t *testing.T) {
		client := &http.Client{}

		response, err := client.Get("http://localhost:8080/healthz")
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			t.Errorf("expected status %s | got %v", http.StatusText(http.StatusOK), response.Status)
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

		if string(body) != http.StatusText(http.StatusOK) {
			t.Errorf("expected %s | got %s", http.StatusText(http.StatusOK), string(body))
		}
	})
}
