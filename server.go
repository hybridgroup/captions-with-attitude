package main

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	_ "embed"
)

// The HTML index page for the web server is embedded here.

//go:embed html/index.html
var index string

var (
	img   image.Image
	mutex sync.Mutex
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Limit upload size to 5MB
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read image data", http.StatusBadRequest)
		return
	}
	mutex.Lock()
	defer mutex.Unlock()

	img, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// startWebServer starts the web server. The server serves the MJPEG stream
// and the caption text.
// It also provides endpoints to set the tone and humor of the captions.
func startWebServer(host string, promptText string) {
	mux := http.NewServeMux()

	// Serve the index page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(index))
	})

	// Serve the MJPEG stream
	//mux.Handle("/video", stream)
	mux.HandleFunc("/upload", uploadHandler)

	// Serve the latest caption
	mux.HandleFunc("/caption", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(caption))
	})

	// Endpoint to set the tone
	mux.HandleFunc("/tone", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type ToneRequest struct {
				Tone string `json:"tone"`
			}
			var req ToneRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil || (req.Tone != "flattering" && req.Tone != "neutral" && req.Tone != "insulting") {
				http.Error(w, "Invalid tone", http.StatusBadRequest)
				return
			}
			tone = req.Tone

			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Endpoint to set the humor
	mux.HandleFunc("/humor", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			type HumorRequest struct {
				Humor string `json:"humor"`
			}
			var req HumorRequest
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil || (req.Humor != "funny" && req.Humor != "neutral" && req.Humor != "serious") {
				http.Error(w, "Invalid humor", http.StatusBadRequest)
				return
			}
			humor = req.Humor

			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	// Endpoint to get the original prompt
	mux.HandleFunc("/prompt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(promptText))
	})

	server := &http.Server{
		Addr:         host,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
