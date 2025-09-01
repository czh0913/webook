.PHONY: docker
docker:
	@rm -f webook || true
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o webook .
	@docker ps -a -q --filter ancestor=chhh31/webook:v0.0.1 | xargs -r docker rm -f
	@docker rmi -f chhh31/webook:v0.0.1 || true
	@docker build -t chhh31/webook:v0.0.1 .
