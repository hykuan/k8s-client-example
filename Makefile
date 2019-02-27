BUILD_DIR = build
SERVICES = models k8s-client
DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
CGO_ENABLED ?= 0
GOOS ?= linux
# GOOS ?= darwin

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ${BUILD_DIR}/quaistudio-$(1) cmd/$(1)/main.go
endef

define make_docker
	docker build --no-cache --build-arg SVC_NAME=$(subst docker_,,$(1)) --tag=quaistudio/$(subst docker_,,$(1)) -f docker/Dockerfile .
endef

define make_docker_dev
	docker build --build-arg SVC_NAME=$(subst docker_dev_,,$(1)) --tag=quaistudio/$(subst docker_dev_,,$(1)) -f docker/Dockerfile.dev ./build
endef

all: $(SERVICES)

.PHONY: all $(SERVICES) dockers dockers_dev

cleandocker: cleanghost
	# Stop all containers (if running)
	docker-compose -f docker/docker-compose.yaml stop
	# Remove quaistudio containers
	docker ps -f name=quaistudio -aq | xargs -r docker rm
	# Remove old quaistudio images
	docker images -q quaistudio\/* | xargs -r docker rmi

# Clean ghost docker images
cleanghost:
	# Remove exited containers
	docker ps -f status=dead -f status=exited -aq | xargs -r docker rm -v
	# Remove unused images
	docker images -f dangling=true -q | xargs -r docker rmi
	# Remove unused volumes
	docker volume ls -f dangling=true -q | xargs -r docker volume rm

install:
	cp ${BUILD_DIR}/* $(GOBIN)

test:
	GOCACHE=off go test -v -race -tags test $(shell go list ./... | grep -v 'vendor\|cmd')

proto:
	protoc --gofast_out=plugins=grpc:. *.proto

$(SERVICES):
	$(call compile_service,$(@))

$(DOCKERS):
	$(call make_docker,$(@))

dockers: $(DOCKERS)

$(DOCKERS_DEV):
	$(call make_docker_dev,$(@))

dockers_dev: $(DOCKERS_DEV)

run:
	docker-compose -f docker/docker-compose.yaml up

down:
	docker-compose -f docker/docker-compose.yaml down