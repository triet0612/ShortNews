build_web:
	cd ./frontend && npm run build && mv build/ ../bin/build

run_api:
	cd ./bin && ./short_news.bin

debug_api:
	make build_linux && make run_api

run_full:
	make clean_bin && make build_linux && make build_web && make run_api

build_linux:
	rm -rf ./bin/ && mkdir ./bin/ && cp ./install.sh ./bin/ && \
	go build -o ./bin/short_news.bin ./cmd/main.go && \
	make build_web

compose:
	docker compose up -d
