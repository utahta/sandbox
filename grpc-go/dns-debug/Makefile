VERSION ?= latest

vendor:
	mkdir -p vendor
	cd vendor && git clone git@github.com:grpc/grpc-go.git

init:
	kubectl create namespace echo-server
	kubectl create namespace greet-client

clean:
	kubectl delete namespace echo-server
	kubectl delete namespace greet-client

.PHONY: echo-server
echo-server:
	go build -o echo-server ./echo/main.go

.PHONY: greet-client
greet-client:
	go build -ldflags -X=main.version=${VERSION} -o greet-client ./greet/main.go

docker.echo:
	docker build --rm --pull -t echo-server:latest -f ./Dockerfile-echo .
	docker tag echo-server:latest registry.local:5000/echo-server:${VERSION}

docker.echo.push:
	docker push registry.local:5000/echo-server:${VERSION}

docker.greet:
	docker build --rm --pull -t greet-client:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-greet .
	docker tag greet-client:latest registry.local:5000/greet-client:${VERSION}

docker.greet.push:
	docker push registry.local:5000/greet-client:${VERSION}

apply.echo:
	kubectl apply -f ./echo/deployments.yaml

apply.greet:
	kubectl apply -f ./greet/deployments.yaml

log.echo:
	stern -n echo-server echo

log.greet:
	stern -n greet-client greet

