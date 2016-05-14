BINDATA_FILE := bindata.go
TESTSETUP_DIR := testsetup


.PHONY: build
build: assets main

.PHONY: build-debug
build-debug: assets-debug main

# Not phony
main:
	go build cmd/pushabutton/main.go

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

# Not phony
$(TESTSETUP_DIR): build
	mkdir -p $(TESTSETUP_DIR)
	(cd $(TESTSETUP_DIR) && ../main setup)
