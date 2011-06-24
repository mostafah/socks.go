# Copyright 2009 The Go Authors. All rights reserved.
# Copyright 2011 Mostafa Hajizadeh.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=socks
GOFILES=main.go

include $(GOROOT)/src/Make.cmd
