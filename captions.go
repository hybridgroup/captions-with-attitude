package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hybridgroup/yzma/pkg/llama"
	"github.com/hybridgroup/yzma/pkg/mtmd"
	"github.com/hybridgroup/yzma/pkg/vlm"
)

var (
	caption string
	tone    string
	humor   string
)

// startCaptions starts the caption generation loop.
// It initializes the VLM and continuously generates captions
// based on the current video frame.
func startCaptions(modelFile, projectorFile, prompt string) {
	if err := llama.Load(*libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}
	if err := mtmd.Load(*libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}

	if !*verbose {
		llama.LogSet(llama.LogSilent())
		mtmd.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.BackendFree()

	vlm := vlm.NewVLM(modelFile, projectorFile)
	if err := vlm.Init(); err != nil {
		fmt.Println("unable to initialize VLM:", err)
		os.Exit(1)
	}
	defer vlm.Close()

	for {
		caption = nextCaption(vlm, prompt)
		if caption != "" {
			fmt.Println("Caption:", caption)
		}

		time.Sleep(3 * time.Second)
	}
}

// nextCaption generates the next caption using the VLM
// based on the current video frame and prompt.
func nextCaption(vlm *vlm.VLM, prompt string) string {
	bitmap, err := imgToBitmap(img)
	if err != nil {
		switch err.Error() {
		case "empty image":
			fmt.Println("Open your browser to", *host, "and activate your camera to start generating captions.")
		default:
			fmt.Println("Error converting image to bitmap:", err)
		}
		return ""
	}
	defer mtmd.BitmapFree(bitmap)

	newPrompt := prompt + promptStyle()
	fmt.Println(newPrompt)

	messages := []llama.ChatMessage{llama.NewChatMessage("user", newPrompt+mtmd.DefaultMarker())}
	input := mtmd.NewInputText(vlm.ChatTemplate(messages, true), true, true)

	output := mtmd.InputChunksInit()
	defer mtmd.InputChunksFree(output)

	if err := vlm.Tokenize(input, []mtmd.Bitmap{bitmap}, output); err != nil {
		fmt.Println("Error tokenizing input:", err)
		return ""
	}

	results, err := vlm.Results(output)
	if err != nil {
		fmt.Println("Error obtaining VLM results:", err)
		return ""
	}

	return results
}

const keepShort = " Keep the response to 30 words or less."

// promptStyle builds the style instructions for the prompt
// based on the current tone and humor settings.
func promptStyle() string {
	switch {
	case tone == "" && humor == "":
		return keepShort
	case tone != "" && humor != "":
		return " Be both " + tone + " and " + humor + " in your response." + keepShort
	case tone == "" && humor != "":
		return " Be somewhat " + humor + " in your response." + keepShort
	case tone != "" && humor == "":
		return " Be somewhat " + tone + " in your response." + keepShort
	}

	return keepShort
}
