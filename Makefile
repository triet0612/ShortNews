clean_bin:
	rm -rf ./bin/news.db && rm -rf ./bin/build/

build_web:
	cd ./frontend && npm run build && mv build/ ../bin/build

run_api:
	cd ./bin && ./short_news.bin

debug_api:
	make build_linux && make run_api

run_full:
	make clean_bin && make build_linux && make build_web && make run_api

build_windows:
	rm -rf ./bin/ && mkdir ./bin/ && cp ./install.ps1 ./bin/ && cp ./install.bat ./bin/ && \
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags="-H=windowsgui" -o ./bin/short_news.exe ./cmd/main.go && \
	make build_web

build_linux:
	rm -rf ./bin/ && mkdir ./bin/ && cp ./install.sh ./bin/ && \
	go build -o ./bin/short_news.bin ./cmd/main.go && \
	make build_web
