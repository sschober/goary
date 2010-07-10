# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=goary
GOFMT=gofmt -spaces=true -tabindent=false -tabwidth=4

GOFILES=\
	goary.go\

include $(GOROOT)/src/Make.cmd

format:
	for src in ${GOFILES}; do \
		${GOFMT} -w $$src; \
	done
