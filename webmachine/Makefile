
include $(GOROOT)/src/Make.inc

TARG=github.com/pomack/webmachine.go/webmachine
GOFILES=\
	charset.go\
	chunked.go\
	const.go\
	decision.go\
	directorylisting.go\
	encoders.go\
	file_resource.go\
	mock.go\
	passthroughmediatypehandler.go\
	passthroughmediatypeinputhandler.go\
	request.go\
	resource.go\
	responsewriter.go\
	types.go\
	util.go\
	webmachine.go\
	writethrough.go

DIRS=\

include $(GOROOT)/src/Make.pkg

check:
	gomake test

-include ../Make.deps

