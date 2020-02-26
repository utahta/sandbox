VERSION ?= latest
CONTAINER_REG ?= registry.local:5000

all: deps

deps:
	go mod download

cluster.create:
	k3d create --enable-registry -n istio-expt -a 6447 -x "--no-deploy=traefik"

cluster.delete:
	k3d delete -n istio-expt

ns.create:
	kubectl create ns echo-service
	kubectl create ns greet-service

istio.install:
	istioctl manifest apply \
		-s profile=default \
		-s telemetry.enabled=false \
		-s values.prometheus.enabled=false \
		-s values.sidecarInjectorWebhook.enableNamespacesByDefault=true \
		-s values.global.proxy.autoInject=disabled \
		-s values.global.hub=registry.local:5000 \
		-s values.global.tag=cbd7908e9d30b3ab1c02cbd45e4f68c1735e70bb

istio.uninstall:
	istioctl manifest apply \
		-s profile=default \
		-s telemetry.enabled=false \
		-s values.prometheus.enabled=false \
		-s values.sidecarInjectorWebhook.enableNamespacesByDefault=true
		-s values.proxy.autoInject=disabled \
		-s values.global.hub=registry.local:5000 \
		-s values.global.tag=cbd7908e9d30b3ab1c02cbd45e4f68c1735e70bb \
		| kcl delete -f -

.PHONY: echod
echod:
	go build -ldflags -X=main.version=${VERSION} -o echod ./echo/main.go

.PHONY: greetd
greetd:
	go build -ldflags -X=main.version=${VERSION} -o greetd ./greet/main.go

docker.echod:
	docker build --rm --pull -t echod:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-echod .
	docker tag echod:latest ${CONTAINER_REG}/echod:${VERSION}

docker.echod.push: docker.echod
	docker push ${CONTAINER_REG}/echod:${VERSION}

docker.greetd:
	docker build --rm --pull -t greetd:latest --build-arg VERSION=$(VERSION) -f ./Dockerfile-greetd .
	docker tag greetd:latest ${CONTAINER_REG}/greetd:${VERSION}

docker.greetd.push: docker.greetd
	docker push ${CONTAINER_REG}/greetd:${VERSION}
