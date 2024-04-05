## AI Enhanced RSS News App

1. Dependencies:
- [Ollama](https://ollama.com/) install with [gemma:2b](https://ollama.com/library/gemma) model
- Require python3.11 install, with [mimic3](https://github.com/MycroftAI/mimic3) package.

2. Program will get news from rss source.
- Get news from multiple source hourly.
- Extract only main content using LLM
- Create audio from Text-to-Speech model with summarized text.
- Store in the database.

3. Host an HTTP server to get news content
- Can add news source of choice.
- Get news by source, with pagination.

4. Summarize articles into short news.
- A short news should:
    - shows a title, thumbnail with summarize text.
    - play voice audio of the summarize text.
    - autoplay next news.

4. Language support:
- English, Vietnamese first.
