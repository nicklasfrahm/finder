GO		:= go
CC		?= gcc
SOURCES := $(shell find . -name "*.go")
GOOS	?= linux
GOARCH	?= amd64
TARGET	:= finder

.PHONY: all clean

all: bin/$(TARGET)

bin/$(TARGET): $(SOURCES)
	-mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 CC=$(CC) $(GO) build -o bin/$(TARGET)-$(GOOS)-$(GOARCH) cmd/$(TARGET)/*

clean:
	-rm -rvf bin
