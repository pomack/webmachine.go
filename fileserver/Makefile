# Copyright 2012 Aalok Shah. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

all: install

GOPATH:=$(GOPATH):`pwd`
BIN_NAME=fileserver

clean:
	GOPATH=$(GOPATH) go clean -x

install:
	GOPATH=$(GOPATH) go build -v -o $(BIN_NAME)

nuke:
	GOPATH=$(GOPATH) go clean -i -x

test:

check:
	GOPATH=$(GOPATH) go build -v -o $(BIN_NAME)
	-/bin/rm -f $(BIN_NAME)

