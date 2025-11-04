package main

import (
	"fmt"
	"sync"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

var (
	webcam *gocv.VideoCapture
	img    gocv.Mat
	mutex  sync.Mutex
)

func startCapture(deviceID string, stream *mjpeg.Stream) {
	var err error
	webcam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	img = gocv.NewMat()
	defer img.Close()

	for {
		captureFrame(deviceID, stream)
	}
}

func captureFrame(deviceID string, stream *mjpeg.Stream) {
	mutex.Lock()
	defer mutex.Unlock()

	if ok := webcam.Read(&img); !ok {
		fmt.Printf("Device closed: %v\n", deviceID)
		return
	}
	if img.Empty() {
		return
	}

	buf, _ := gocv.IMEncode(".jpg", img)
	stream.UpdateJPEG(buf.GetBytes())
	buf.Close()
}
