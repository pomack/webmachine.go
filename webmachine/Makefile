
include $(GOROOT)/src/Make.inc

TARG=webmachine
GOFILES=\
	charset.go\
	chunked.go\
	const.go\
	decision.go\
	encoders.go\
	file_resource.go\
	passthroughmediatypehandler.go\
	request.go\
	resource.go\
	types.go\
	util.go\
	webmachine.go\
	writethrough.go

DIRS=\

include $(GOROOT)/src/Make.pkg

check:
	gomake test

-include ../Make.deps

