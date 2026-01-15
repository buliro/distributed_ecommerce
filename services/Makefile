test: docker_build_test
	docker compose down
	docker compose up -d
	docker compose exec app go test -v ./...
	docker compose down

unit_test:
	go test -v ./...

docker_build:
	docker build . -t api

docker_build_test:
	docker build . -t api_test --target=test

docker_run:
	docker run -p 8080:8080 api
