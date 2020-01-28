VERSION ?= latest
CONTAINER_REG ?= registry.local:5000

all: deps

deps:
	go mod download

cluster.create:
	k3d create --enable-registry -n istio-hello -a 6444 -x "--no-deploy=traefik"

cluster.delete:
	k3d delete -n istio-hello

ns.create:
	kubectl create ns echo-service
	kubectl create ns echo2-service
	kubectl create ns greet-service

istio.install:
	istioctl manifest apply \
		-s profile=default \
		-s telemetry.enabled=false \
		-s values.prometheus.enabled=false \
		-s values.sidecarInjectorWebhook.enableNamespacesByDefault=true \
		-s values.global.proxy.autoInject=disabled

istio.uninstall:
	istioctl manifest apply \
		-s profile=default \
		-s telemetry.enabled=false \
		-s values.prometheus.enabled=false \
		-s values.sidecarInjectorWebhook.enableNamespacesByDefault=true
		-s values.proxy.autoInject=disabled \
		| kcl delete -f -

.PHONY: echod
echod:
	go build -ldflags -X=main.version=${VERSION} -o echod ./echo/main.go

.PHONY: echo2d
echo2d:
	go build -ldflags -X=main.version=${VERSION} -o echo2d ./echo2/main.go

.PHONY: greetd
greetd:
	go build -ldflags -X=main.version=${VERSION} -o greetd ./greet/main.go

.PHONY: webpaged
webpaged:
	go build -ldflags -X=main.version=${VERSION} -o webpaged ./webpage/main.go

docker.echod:
	docker build --rm --pull -t echod:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-echod .
	docker tag echod:latest ${CONTAINER_REG}/echod:${VERSION}

docker.echod.push: docker.echod
	docker push ${CONTAINER_REG}/echod:${VERSION}

docker.echo2d:
	docker build --rm --pull -t echo2d:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-echo2d .
	docker tag echo2d:latest ${CONTAINER_REG}/echo2d:${VERSION}

docker.echo2d.push: docker.echo2d
	docker push ${CONTAINER_REG}/echo2d:${VERSION}

docker.greetd:
	docker build --rm --pull -t greetd:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-greetd .
	docker tag greetd:latest ${CONTAINER_REG}/greetd:${VERSION}

docker.greetd.push: docker.greetd
	docker push ${CONTAINER_REG}/greetd:${VERSION}

docker.webpaged:
	docker build --rm --pull -t webpaged:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-webpaged .
	docker tag webpaged:latest ${CONTAINER_REG}/webpaged:${VERSION}

docker.webpaged.push: docker.webpaged
	docker push ${CONTAINER_REG}/webpaged:${VERSION}

