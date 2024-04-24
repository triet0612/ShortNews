wget https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/vi/vi_VN/vais1000/medium/vi_VN-vais1000-medium.onnx?download=true -OutFile vi_VN.onnx

wget https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/vi/vi_VN/vais1000/medium/vi_VN-vais1000-medium.onnx.json?download=true.json -OutFile vi_VN.onnx.json

wget https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/en/en_US/ryan/low/en_US-ryan-low.onnx?download=true -Outfile en_US.onnx

wget https://huggingface.co/rhasspy/piper-voices/resolve/v1.0.0/en/en_US/ryan/low/en_US-ryan-low.onnx.json?download=true.json -Outfile en_US.onnx.json

wget https://github.com/rhasspy/piper/releases/download/2023.11.14-2/piper_windows_amd64.zip -OutFile piper_windows_amd64.zip

Expand-Archive -Path "./piper_windows_amd64.zip" -DestinationPath "./"

Remove-Item "./piper_windows_amd64.zip"

ollama pull gemma:2b-instruct-v1.1-q4_0
