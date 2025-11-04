// What it does:
//
// This example opens a video capture device, then streams MJPEG from it.
// Once running point your browser to the hostname/port you passed in the
// command line (for example http://localhost:8080) and you should see
// the live video stream.
//
// How to run:
//
// mjpeg-streamer [camera ID] [host:port]
//
//		go get -u github.com/hybridgroup/mjpeg
// 		go run ./cmd/mjpeg-streamer/main.go 1 0.0.0.0:8080
//

package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/mjpeg"
)

func main() {
	if len(os.Args) < 6 {
		fmt.Println("How to run:\n\tvideo-description [camera ID] [host:port] [model path] [projector path] [prompt text]")
		return
	}

	// parse args
	deviceID := os.Args[1]
	host := os.Args[2]
	modelPath := os.Args[3]
	projectorPath := os.Args[4]
	promptText := os.Args[5]

	// create the mjpeg stream
	stream := mjpeg.NewStream()

	go startCapture(deviceID, stream)
	go startVLM(modelPath, projectorPath, promptText)

	fmt.Println("Capturing. Point your browser to " + host)

	// start http server
	startServer(host, stream)
}
