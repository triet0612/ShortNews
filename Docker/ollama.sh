#!/usr/bin/env bash

ollama serve &
ollama list
ollama pull gemma:2b-instruct-v1.1-q4_0
