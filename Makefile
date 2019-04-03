NAME=xmlpaste
VERSION=1.0.0-dev

GOCMD=go
GOBUILD=$(GOCMD) build
DIST_DIR=dist
MACOS_DIR=macos
WINDOWS_DIR=windows
DIST_MACOS_DIR=$(NAME)-$(VERSION)-$(MACOS_DIR)
DIST_WINDOWS_DIR=$(NAME)-$(VERSION)-$(WINDOWS_DIR)

all: build

.PHONY: clean
clean:
	@rm -rf $(DIST_DIR)

build: build-macos build-windows

build-macos:
	mkdir -p $(DIST_DIR)/$(MACOS_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(DIST_DIR)/$(MACOS_DIR)/$(NAME)

build-windows:
	mkdir -p $(DIST_DIR)/$(WINDOWS_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) -ldflags "-X main.version=$(VERSION)" -o $(DIST_DIR)/$(WINDOWS_DIR)/$(NAME).exe

.PHONY: dist
dist: build
	cd $(DIST_DIR) && \
	mv $(MACOS_DIR) $(DIST_MACOS_DIR) && \
	cp -p ../LICENSE.txt $(DIST_MACOS_DIR)/ && \
	cp -p ../README.md $(DIST_MACOS_DIR)/ && \
	cp -p ../release-notes.txt $(DIST_MACOS_DIR)/ && \
	zip -r $(DIST_MACOS_DIR).zip $(DIST_MACOS_DIR) && \
	cd ..

	cd $(DIST_DIR) && \
	mv $(WINDOWS_DIR) $(DIST_WINDOWS_DIR) && \
	cp -p ../LICENSE.txt $(DIST_WINDOWS_DIR)/ && \
	cp -p ../README.md $(DIST_WINDOWS_DIR)/ && \
	cp -p ../release-notes.txt $(DIST_WINDOWS_DIR)/ && \
	zip -r $(DIST_WINDOWS_DIR).zip $(DIST_WINDOWS_DIR) && \
	cd ..
