package main

import (
	"fmt"
	"os"
)

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	go startCaptions(*modelPath, *projectorPath, *promptText)

	fmt.Println("Capturing. Point your browser to", *host)

	startWebServer(*host, *promptText)
}
