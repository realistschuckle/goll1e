# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=goll1e
GOFILES=\
	goll1e.go\
	scanner.go\
	set.go\
	compute.go\
	devprint.go\
	parse.go\

include $(GOROOT)/src/Make.pkg
