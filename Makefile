all: build-pip build

build-pip:
	go build -o pip_package/tfacon_pip/tfacon_binary/tfacon ./main.go

build:
	go build .