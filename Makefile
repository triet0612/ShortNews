clean_bin:
	rm -rf ./bin/ && mkdir ./bin/

build_api:
	go build ./cmd/main.go && mv main ./bin

build_web:
	cd ./frontend && npm run build && mv build/ ../bin/

run_api:
	cd ./bin/ && ./main

debug_api:
	make build_api && make run_api

run_full:
	make clean_bin && make build_api && make build_web && make run_api

test_rss_source1:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://vnexpress.net/rss/thoi-su.rss", "language":"Vietnamese"}' localhost:8000/api/rss

test_rss_source2:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://infonet.vietnamnet.vn/rss/doi-song.rss", "language":"Vietnamese"}' localhost:8000/api/rss

test_rss_source3:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://abcnews.go.com/abcnews/topstories", "language":"English"}' localhost:8000/api/rss

test_rss_source4:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://thanhnien.vn/rss/home.rss", "language":"Vietnamese"}' localhost:8000/api/rss

test_rss_source5:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://vneconomy.vn/tin-moi.rss", "language":"Vietnamese"}' localhost:8000/api/rss
