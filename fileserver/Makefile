
include $(GOROOT_FINAL)/src/Make.inc

all: install

TARG=fileserver
GOFILES=\
	server.go

DIRS=\


-include ../Make.deps

include $(GOROOT_FINAL)/src/Make.cmd

check:
	gomake test

