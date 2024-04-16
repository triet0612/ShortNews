## AI Enhanced RSS News App

1. Dependencies:
- [Ollama](https://ollama.com/) required to install
- Installation file will auto download gemma:2b-instruct-v1.1-q4_0

2. How does it works?
- Get news from multiple source hourly.
- Extract only main content using gemma:2b-instruct-v1.1-q4_0
- Create audio from piper-TTS model with summarized text.
- Store in the database.
- Host an HTTP server to get news content
- User open http://localhost:8000 to use the app. When open the binary, auto launch in web browser.

4. Language support:
- English, Vietnamese.

5. Installation:
- First, must install [Ollama](https://ollama.com/)
- Windows: run "install.bat" script and wait until the installation is done. Run the short_news.exe binary.
- Linux: use console to run "install.sh" script, wait until the installation is done. Run the short_news.bin binary.
