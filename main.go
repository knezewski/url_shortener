package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urls = make(map[string]string)

func main() {
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/short", handleRedirect)

	fmt.Println("URL Shortener is running on: 8080")
	http.ListenAndServe(":8080", nil)
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/shorten", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
    	<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<form method="post" action="/shorten">
				<input type="url" name="url" placeholder="Enter a URL" required>
				<input type="submit" value="Shorten">
			</form>
		</body>
		</html>
    `)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")
	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	shortKey := generateShortKey()
	urls[shortKey] = originalURL

	shortenedURL := fmt.Sprintf("http://localhost:8080/short/%s", shortKey)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html>
		<head>
			<title>URL Shortener</title>
		</head>
		<body>
			<h2>URL Shortener</h2>
			<p>Original URL: `, originalURL, `</p>
			<p>Shortened URL: <a href="`, shortenedURL, `">`, shortenedURL, `</a></p>
		</body>
		</html>
	`)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimPrefix(r.URL.Path, "/short/")
	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	originalURL, found := urls[shortKey]
	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6

	time.Now().UnixNano()
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[rand.Intn(len(charset))]
	}

	return string(shortKey)
}
