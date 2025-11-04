# video-description

![video description](./images/video-description.png)

Uses yzma and GoCV to create a web server that shows live captions from your local webcam.

## Running

```shell
go run . 0 localhost:8080 ~/models/Qwen3-VL-2B-Instruct-Q8_0.gguf ~/models/mmproj-Qwen3-VL-2B-Instruct-Q8_0.gguf "Give a very brief description of what is going on."
```