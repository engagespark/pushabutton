BINDATA_FILE := bindata.go
TESTSETUP_DIR := testsetup

GIT_COMMIT=$(shell git rev-parse HEAD)


.PHONY: build
build: assets main

.PHONY: build
build-release: assets main-release

.PHONY: build-debug
build-debug: assets-debug main-debug

# Not phony
main:
	go build -ldflags "-X main.Build=dev-$(GIT_COMMIT)" cmd/pushabutton/main.go

.PHONY:main-release
main-release:
	GIT_TAG=$(shell git describe --exact-match --abbrev=0)
	go build -ldflags "-X main.Build=release-$(GIT_TAG)" cmd/pushabutton/main.go

# Not phony
main-debug:
	go build -ldflags "-X main.Build=debug-$(GIT_COMMIT)" cmd/pushabutton/main.go

.PHONY: assets
assets:
	go-bindata -o $(BINDATA_FILE) -pkg pushabutton assets

.PHONY: assets-debug
assets-debug:
	go-bindata -debug -o $(BINDATA_FILE) -pkg pushabutton assets

.PHONY: clean
clean:
	rm $(BINDATA_FILE) main || true
	rm -r --preserve-root $(TESTSETUP_DIR) || true

.PHONY: serve
serve: build-debug
	./main

.PHONY: release
release: clean build

# Not phony
$(TESTSETUP_DIR): build
	mkdir -p $(TESTSETUP_DIR)
	(cd $(TESTSETUP_DIR) && ../main setup)
