build:
	go test ./... -tags="integration acceptance"
	go build

run: build run-mongo run-flyte
	./flyte-bitbucket

stop: stop-mongo stop-flyte
	killall flyte

run-mongo:
	docker run -d --name mongo mongo:latest

run-flyte:
	docker run -d -e FLYTE_MGO_HOST=mongo --link mongo:mongo --name flyte hotelsdotcom/flyte

stop-mongo:
	docker rm -f mongo

stop-flyte:
	docker rm -f flyte

docker-build:
	docker build --rm -t flyte-bitbucket:latest .

docker-run: docker-build run-mongo run-flyte
	docker run -d -e FLYTE_API_URL=http://flyte:8080 -d --name flyte-bitbucket --link flyte:flyte flyte-bitbucket:latest

docker-stop: stop-mongo stop-flyte
	docker rm -f flyte-bitbucket
