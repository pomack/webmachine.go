
include $(GOROOT)/src/Make.inc

all: install

TARG=fileserver
GOFILES=\
	server.go

DIRS=\


-include ../Make.deps

include $(GOROOT)/src/Make.cmd

check:
	gomake test

