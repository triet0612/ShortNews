build_app:
	rm -rf ./bin/ && mkdir ./bin/ && go build ./cmd/main.go && mv main ./bin

run_app:
	cd ./bin/ && ./main

run_debug:
	make build_app && make run_app

test_rss_source1:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://vnexpress.net/rss/tin-moi-nhat.rss"}' localhost:8000/api/rss

test_rss_source2:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://infonet.vietnamnet.vn/rss/doi-song.rss"}' localhost:8000/api/rss

test_rss_source3:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://abcnews.go.com/abcnews/topstories"}' localhost:8000/api/rss

test_rss_source4:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://thanhnien.vn/rss/home.rss"}' localhost:8000/api/rss

test_rss_source5:
	curl -i -X POST -H "Content-Type: application/json" \
	-d '{"link":"https://vneconomy.vn/tin-moi.rss"}' localhost:8000/api/rss
