PLUGIN_NAME=ishant8/sdip

all: clean build-image build-plugin create-plugin

clean:
	rm -rf ./plugin ./bin
	docker plugin disable ${PLUGIN_NAME} || true
	docker plugin rm ${PLUGIN_NAME} || true
	docker rm -vf tmp || true
	docker rmi sdip-build-image || true
	docker rmi ${PLUGIN_NAME}:rootfs || true

build-image:
	docker build -t sdip-build-image -f Dockerfile.build .
	docker create --name tmp sdip-build-image
	docker cp tmp:/go/bin/docker-ipam-plugin .
	docker rm -vf tmp
	docker rmi sdip-build-image
	docker build -t ${PLUGIN_NAME}:rootfs .

build-plugin:
	mkdir -p ./plugin/rootfs
	docker create --name tmp ${PLUGIN_NAME}:rootfs
	docker export tmp | tar -x -C ./plugin/rootfs
	cp config.json ./plugin/
	docker rm -vf tmp

create-plugin:
	docker plugin create ${PLUGIN_NAME} ./plugin
	docker plugin enable ${PLUGIN_NAME}

push:  clean build-image build-plugin create-plugin
	docker plugin push ${PLUGIN_NAME}
