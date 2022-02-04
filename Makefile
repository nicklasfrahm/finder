GO			:= go
CC			?= gcc
GO_SOURCES	:= $(shell find . -name "*.go")
WEB_SOURCES	:= $(shell find ./web/src -name "*") $(shell find ./web/static -name "*")
WEB_URL		?= file:///web/build
GOOS		?= linux
GOARCH		?= amd64

.PHONY: all clean

all: bin/gui-$(GOOS)-$(GOARCH) bin/cli-$(GOOS)-$(GOARCH)

web/build/: $(WEB_SOURCES)
	echo $(WEB_SOURCES)
	cd web && npm run build

bin/cli-$(GOOS)-$(GOARCH): $(GO_SOURCES)
	-mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 CC=$(CC) $(GO) build -o bin/$(TARGET)-$(GOOS)-$(GOARCH) cmd/$(TARGET)/*

bin/app-$(GOOS)-$(GOARCH): $(GO_SOURCES) web/build/
	-mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 CC=$(CC) \
	$(GO) build -o bin/app-$(GOOS)-$(GOARCH) -ldflags "-X main.url=$(WEB_URL)" cmd/app/*

# Rerun target every time a file change is detected in the current directory.
watch:
	while true; do \
		clear; \
        make --no-print-directory $(TARGET); \
        inotifywait -qre close_write .; \
    done

clean:
	-rm -rvf bin
