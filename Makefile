include .env
PROJECTNAME=$(shell basename "$(PWD)")

# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
SERVICE_NAME = ${PROJECTNAME}
STATIC_CHECK = golangci-lint run
GITCOMMITID=`git rev-parse HEAD`
GITTAG=`git describe --tag`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-X ${PROJECTNAME}/tools.Version=${VERSION} -X ${PROJECTNAME}/tools.GitTag=${GITTAG} -X ${PROJECTNAME}/tools.GitCommitId=${GITCOMMITID} -X ${PROJECTNAME}/tools.BuildTime=${BUILD_TIME}"

export GOPROXY=https://goproxy.cn

all: build test
build:
	rm -rf target/
	# 配置文件目录
	mkdir -p target/conf
	cp restart.sh target/
	cp conf/* target/
	$(GOBUILD) $(LDFLAGS) -o target/$(SERVICE_NAME) *.go

test:
	$(GOTEST) -v ./...

static_check:
	$(STATIC_CHECK) ./...

clean:
	rm -rf target/

run:
	nohup target/$(SERVICE_NAME) &

stop:
	pkill -f target/$(SERVICE_NAME)