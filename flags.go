package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	libPath       *string
	processor     *string
	modelPath     *string
	projectorPath *string
	promptText    *string
	verbose       *bool
	deviceID      *string
	host          *string
)

// showUsage displays the usage information for the program.
func showUsage() {
	fmt.Println(`
Usage:
captions-with-attitudes`)
	flag.PrintDefaults()
}

// handleFlags processes the command-line flags and validates them.
func handleFlags() error {
	libPath = flag.String("lib", "", "path to llama.cpp compiled library files")
	modelPath = flag.String("model", "", "model file to use")
	projectorPath = flag.String("projector", "", "projector file to use")
	promptText = flag.String("p", "Give a very brief description of what is going on.", "prompt")
	verbose = flag.Bool("v", false, "verbose logging")
	deviceID = flag.String("device", "0", "camera device ID")
	host = flag.String("host", "localhost:8080", "web server host:port")
	processor = flag.String("processor", "cpu", "processor to use (cpu, cuda, metal, vulkan)")

	flag.Parse()

	if len(*libPath) == 0 && os.Getenv("YZMA_LIB") != "" {
		*libPath = os.Getenv("YZMA_LIB")
	}

	if len(*libPath) == 0 {
		return errors.New("missing lib flag or YZMA_LIB env var")
	}

	return nil
}
