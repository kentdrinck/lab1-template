app:
	go build -o bin/app cmd/server/main.go  

run: app
	./bin/app 

tidy:
	go mod tidy

test:
	go test ./...

docker-test:
	docker build -t go-unit-tests-temp -f Dockerfile-test . && \
	docker run --rm go-unit-tests-temp && \
	docker rmi go-unit-tests-temp

docker-run:
	docker compose down
	docker compose pull
	docker compose up -d --build
	docker image prune -f