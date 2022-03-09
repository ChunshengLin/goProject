.PHONY: build

SERVICE := elastic
MAIN := cmd/
AUTHOR := $(shell git log --pretty=format:"%an"|head -n 1)
VERSION := $(shell git rev-list HEAD | head -1)
BUILD_INFO := $(shell git log --pretty=format:"%s" | head -1)
BUILD_DATE := $(shell date +%Y-%m-%d\ %H:%M:%S)
CUR_PWD := $(shell pwd)
export GO111MODULE=on

LD_FLAGS='-X "$(SERVICE)/version.TAG=$(TAG)" -X "$(SERVICE)/version.VERSION=$(VERSION)" -X "$(SERVICE)/version.AUTHOR=$(AUTHOR)" -X "$(SERVICE)/version.BUILD_INFO=$(BUILD_INFO)" -X "$(SERVICE)/version.BUILD_DATE=$(BUILD_DATE)"'

default: build 

build:
		go build -ldflags $(LD_FLAGS) -gcflags "-N" -o ./bin/$(SERVICE) ./$(MAIN)
clean:
		rm bin/*