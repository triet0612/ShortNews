clean_bin:
	rm -rf ./bin/news.db && rm -rf ./bin/build/

build_api:
	go build -o ./bin/main.bin ./cmd/main.go

build_web:
	cd ./frontend && npm run build && mv build/ ../bin/

run_api:
	cd ./bin/ && ./main.bin

debug_api:
	make build_api && make run_api

run_full:
	make clean_bin && make build_api && make build_web && make run_api

docker_build:
	docker build --no-cache . -t short_news_app
