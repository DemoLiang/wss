
landfwSRC = $(wildcard ./*.go)
landfwLinuxTarget = docker/bins/landfw
$(landfwLinuxTarget):$(landfwSRC)
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o docker/bins/landfw  github.com/DemoLiang/wss/src

linux:$(landfwLinuxTarget)

docker:linux
	cd docker;pwd;docker build -t landfw .
	-docker stop landfw
	-docker rm landfw
	docker run --name landfw -d -p7777:7777 landfw

all:$(docker)

clean:
	rm -rf docker/bins/landfw
