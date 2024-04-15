clean_bin:
	rm -rf ./bin/news.db && rm -rf ./bin/build/

build_web:
	cd ./frontend && npm run build && mv build/ ../bin/

run_api:
	cd ./bin/ && ./short_news.bin

debug_api:
	make build_linux && make run_api

run_full:
	make clean_bin && make build_linux && make build_web && make run_api

build_windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -o ./bin/short_news.exe ./cmd/main.go

build_linux:
	go build -o ./bin/short_news.bin ./cmd/main.go
