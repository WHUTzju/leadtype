// Copyright 2011-2012 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package pdf

import (
	"testing"

	"github.com/rowland/leadtype/options"
)

func TestMakeSizeRectangle(t *testing.T) {
	r := makeSizeRectangle("letter", "portrait")
	expectF(t, 0, r.x1)
	expectF(t, 0, r.y1)
	expectF(t, 612, r.x2)
	expectF(t, 792, r.y2)
}

func TestNewPageStyle(t *testing.T) {
	opt1 := options.Options{}
	ps1 := newPageStyle(opt1)
	expectF(t, 0, ps1.pageSize.x1)
	expectF(t, 0, ps1.pageSize.y1)
	expectF(t, 612, ps1.pageSize.x2)
	expectF(t, 792, ps1.pageSize.y2)
	expectF(t, 0, ps1.cropSize.x1)
	expectF(t, 0, ps1.cropSize.y1)
	expectF(t, 612, ps1.cropSize.x2)
	expectF(t, 792, ps1.cropSize.y2)

	opt2 := options.Options{"orientation": "landscape"}
	ps2 := newPageStyle(opt2)
	expectF(t, 0, ps2.pageSize.x1)
	expectF(t, 0, ps2.pageSize.y1)
	expectF(t, 792, ps2.pageSize.x2)
	expectF(t, 612, ps2.pageSize.y2)
	expectF(t, 0, ps2.cropSize.x1)
	expectF(t, 0, ps2.cropSize.y1)
	expectF(t, 792, ps2.cropSize.x2)
	expectF(t, 612, ps2.cropSize.y2)

	opt3 := options.Options{"page_size": "A4"}
	ps3 := newPageStyle(opt3)
	expectF(t, 0, ps3.pageSize.x1)
	expectF(t, 0, ps3.pageSize.y1)
	expectF(t, 595, ps3.pageSize.x2)
	expectF(t, 842, ps3.pageSize.y2)
	expectF(t, 0, ps3.cropSize.x1)
	expectF(t, 0, ps3.cropSize.y1)
	expectF(t, 595, ps3.cropSize.x2)
	expectF(t, 842, ps3.cropSize.y2)

	opt4 := options.Options{"page_size": "A4", "orientation": "landscape"}
	ps4 := newPageStyle(opt4)
	expectF(t, 0, ps4.pageSize.x1)
	expectF(t, 0, ps4.pageSize.y1)
	expectF(t, 842, ps4.pageSize.x2)
	expectF(t, 595, ps4.pageSize.y2)
	expectF(t, 0, ps4.cropSize.x1)
	expectF(t, 0, ps4.cropSize.y1)
	expectF(t, 842, ps4.cropSize.x2)
	expectF(t, 595, ps4.cropSize.y2)

	// TODO: Handle crop_size differently to allow non-zero x1 and y1.

	opt5 := options.Options{"rotate": "portrait"}
	ps5 := newPageStyle(opt5)
	expectNI(t, "rotate", 0, ps5.rotate)

	opt6 := options.Options{"rotate": "landscape"}
	ps6 := newPageStyle(opt6)
	expectNI(t, "rotate", 270, ps6.rotate)
}
