
include $(GOROOT)/src/Make.inc

all: Make.deps install

DIRS=\
	webmachine/\
	fileserver/\

TEST=\
	$(filter-out $(NOTEST),$(DIRS))


clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))

%.clean:
	+cd $* && gomake clean

%.install:
	+cd $* && gomake install

%.nuke:
	+cd $* && gomake nuke

%.test:
	+cd $* && gomake test

%.check:
	+cd $* && gomake check

clean: clean.dirs

install: install.dirs

test:   test.dirs

check:	check.dirs

#nuke: nuke.dirs
#   rm -rf "$(GOROOT)"/pkg/*

echo-dirs:
	@echo $(DIRS)

Make.deps:
	./deps.bash

deps:
	./deps.bash

#-include Make.deps
