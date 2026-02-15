package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hybridgroup/yzma/pkg/download"
)

func checkInstall() {
	if _, err := os.Stat(*libPath); os.IsNotExist(err) {
		fmt.Println("no llama.cpp library directory for yzma.")
		fmt.Print("Do you want to create '" + *libPath + "'? (y/n): ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "Y" {
			if err := os.Mkdir(*libPath, 0755); err != nil {
				fmt.Println("failed to create llama directory:", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Println("Exiting.")
			os.Exit(0)
		}
	}

	if !download.AlreadyInstalled(*libPath) {
		fmt.Println("yzma is not installed.")
		fmt.Print("Do you want to install yzma now? (y/n): ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "Y" {
			version, err := download.LlamaLatestVersion()
			if err != nil {
				fmt.Println("could not obtain latest version:", err.Error())
				os.Exit(1)
			}

			if *processor == "" {
				*processor = "cpu"
				if cudaInstalled, cudaVersion := download.HasCUDA(); cudaInstalled {
					fmt.Printf("CUDA detected (version %s), using CUDA build\n", cudaVersion)
					*processor = "cuda"
				}
			}

			fmt.Println("installing llama.cpp version", version, "to", *libPath)
			if err := download.Get(runtime.GOARCH, runtime.GOOS, *processor, version, *libPath); err != nil {
				fmt.Println("failed to download llama.cpp:", err.Error())
				return
			}
			fmt.Println("yzma installed successfully.")
		} else {
			fmt.Println("Exiting.")
			os.Exit(0)
		}
	}
}

const (
	defaultModelURL = "https://huggingface.co/bartowski/Qwen_Qwen3-VL-2B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-2B-Instruct-Q4_K_M.gguf"
	defaultProjURL  = "https://huggingface.co/bartowski/Qwen_Qwen3-VL-2B-Instruct-GGUF/resolve/main/mmproj-Qwen_Qwen3-VL-2B-Instruct-f16.gguf"
)

func checkModels() {
	if len(*modelPath) == 0 || len(*projectorPath) == 0 {
		// use default models if not provided
		fmt.Println("No model or projector specified, using default models (Qwen3-VL-2B-Instruct-Q4_K_M)")
		*modelPath = filepath.Join(download.DefaultModelsDir(), "Qwen_Qwen3-VL-2B-Instruct-Q4_K_M.gguf")
		*projectorPath = filepath.Join(download.DefaultModelsDir(), "mmproj-Qwen_Qwen3-VL-2B-Instruct-f16.gguf")
	}

	if _, err := os.Stat(*modelPath); os.IsNotExist(err) {
		fmt.Println("model is not downloaded. Do you want to download it now? (y/n): ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "Y" {
			if err := download.GetModel(defaultModelURL, download.DefaultModelsDir()); err != nil {
				fmt.Println("failed to download model:", err.Error())
				os.Exit(1)
			}

			if err := download.GetModel(defaultProjURL, download.DefaultModelsDir()); err != nil {
				fmt.Println("failed to download projector:", err.Error())
				os.Exit(1)
			}
			fmt.Println("model and projector downloaded successfully.")
		} else {
			fmt.Println("Exiting.")
			os.Exit(0)
		}
	}
}
