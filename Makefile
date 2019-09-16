#  The MIT License
# 
#  Copyright (c) 2019 Bravo Cognos, Inc.
# 
#  Permission is hereby granted, free of charge, to any person obtaining a copy
#  of this software and associated documentation files (the "Software"), to deal
#  in the Software without restriction, including without limitation the rights
#  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#  copies of the Software, and to permit persons to whom the Software is
#  furnished to do so, subject to the following conditions:
# 
#  The above copyright notice and this permission notice shall be included in
#  all copies or substantial portions of the Software.
# 
#  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
#  THE SOFTWARE.

ifneq (,)
	This makefile requires GNU Make.
endif

APP:=hydra

GOOS_WINDOWS:=windows
GOOS_LINUX:=linux
GOOS_MAC:=darwin

GOARCH_WINDOWS:=amd64
GOARCH_LINUX:=amd64
GOARCH_MAC:=amd64

BIN_FOLDER:=bin

BIN_WINDOWS_NAME:=$(APP)-$(GOOS_WINDOWS)-$(GOARCH_WINDOWS)
BIN_MAC_NAME:=$(APP)-$(GOOS_MAC)-$(GOARCH_MAC)
BIN_LINUX_NAME:=$(APP)-$(GOOS_LINUX)-$(GOARCH_LINUX)

ANALYZE_THRESHOLD:=20
ANALYZE_MAX_LENGTH:=80

default: releaseAll
prod: $(APP)

buildWindows:
	@GOOS=$(GOOS_WINDOWS) GOARCH=$(GOARCH_WINDOWS) go build -o $(BIN_FOLDER)/$(BIN_WINDOWS_NAME);
	@chmod a+x $(BIN_FOLDER)/$(BIN_WINDOWS_NAME);

buildMac:
	@GOOS=$(GOOS_MAC) GOARCH=$(GOARCH_MAC) go build -o $(BIN_FOLDER)/$(BIN_MAC_NAME);
	@chmod a+x $(BIN_FOLDER)/$(BIN_MAC_NAME);

buildLinux:
	@GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_LINUX) go build -o $(BIN_FOLDER)/$(BIN_LINUX_NAME);
	@chmod a+x $(BIN_FOLDER)/$(BIN_LINUX_NAME);

buildAll: buildMac buildLinux buildWindows

releaseAll: format analyze buildAll
	@# Create archives
	@tar -C $(BIN_FOLDER)/ -cvzf $(BIN_FOLDER)/$(BIN_LINUX_NAME).tar.gz $(BIN_LINUX_NAME)
	@tar -C $(BIN_FOLDER)/ -cvzf $(BIN_FOLDER)/$(BIN_MAC_NAME).tar.gz $(BIN_MAC_NAME)
	@tar -C $(BIN_FOLDER)/ -cvzf $(BIN_FOLDER)/$(BIN_WINDOWS_NAME).tar.gz $(BIN_WINDOWS_NAME)

	@# Remove binaries
	@rm $(BIN_FOLDER)/$(BIN_LINUX_NAME) $(BIN_FOLDER)/$(BIN_MAC_NAME) $(BIN_FOLDER)/$(BIN_WINDOWS_NAME);

format:
	@go fmt ./...

analyze:
	@# Find exported variables (const, var, func, struct) in Go that could be unexported.
	@echo "Running usedexports"
	@usedexports ./...

	@#Check against Golang CI Lint
	@golangci-lint run --enable-all --disable gochecknoglobals --disable gochecknoinits

test: format analyze
	@go test -cover -race -v ./...

detect-licenses:
	@echo "Detecting licenses in the vendor folder"
	@find vendor/ -mindepth 3 -maxdepth 3|xargs license-detector > resources/detected-licenses.txt

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY:	prod
	buildWindows
	buildMac
	buildLinux
	buildAll
	releaseAll
	format
	analyze
	test
	detect-licenses
	list
