# Copyright 2012 Aalok Shah. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: install

GOPATH:=$(GOPATH):`pwd`
PACKAGE_NAME=github.com/pomack/webmachine.go/webmachine

clean:
	GOPATH=$(GOPATH) go clean $(PACKAGE_NAME)

install:
	GOPATH=$(GOPATH) go install $(PACKAGE_NAME)

nuke:
	GOPATH=$(GOPATH) go clean -i $(PACKAGE_NAME)

test:
	GOPATH=$(GOPATH) go test $(PACKAGE_NAME)

check:
	GOPATH=$(GOPATH) go build $(PACKAGE_NAME)

