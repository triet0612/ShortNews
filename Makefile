build_web:
	cd ./frontend && npm install && npm run build

build_debug:
	rm -rf ./bin/ && mkdir ./bin/ && \
	make build_web && \
	go build -o ./bin/short_news.bin ./cmd/main.go \

run_api:
	cd ./bin && \
	ollama_api=http://localhost:11434 \
	voice_api=http://localhost:8000 \
	./short_news.bin

debug_api:
	docker compose -f ./Docker/compose-test.yaml up -d && make build_debug && make run_api

docker_build:
	docker image build . -t shortnews

docker_publish:
	docker tag docker/welcome-to-docker YOUR-USERNAME/welcome-to-docker
