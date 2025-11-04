# video-description

![video description](./images/video-description.png)

This a Go application that creates a web server that shows live captions using a Vision Language Model (VLM) from your local webcam in your browser running entirely on your local machine! 

It uses [yzma](https://github.com/hybridgroup/yzma) to perform local inference using [`llama.cpp`](https://github.com/ggml-org/llama.cpp) and [GoCV](https://github.com/hybridgroup/gocv) for the video processing.

## Installation

### yzma

You must install yzma and llama.cpp to run this program.

See https://github.com/hybridgroup/yzma/blob/main/INSTALL.md

### GoCV

You must also install OpenCV and GoCV, which unlike yzma require CGo.

See https://gocv.io/getting-started/

Although yzma does not use CGo itself, you can also use it in Go applications that use CGo.

### Models

Download the model and projector files from Hugging Face in `.gguf` format.

https://huggingface.co/ggml-org/Qwen3-VL-2B-Instruct-GGUF

## Running

```shell
go run . 0 localhost:8080 ~/models/Qwen3-VL-2B-Instruct-Q8_0.gguf ~/models/mmproj-Qwen3-VL-2B-Instruct-Q8_0.gguf "Give a very brief description of what is going on."
```

Now open your web browser pointed to http://localhost:8080/

