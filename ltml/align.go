// Copyright 2016 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package ltml

type Align int

const (
	AlignDefault = Align(iota)
	AlignLeft
	AlignRight
	AlignTop
	AlignBottom
)

var alignStrings = []string{"default", "left", "right", "top", "bottom"}

func (a Align) String() string {
	if int(a) < len(alignStrings) {
		return alignStrings[a]
	}
	return "unknown"
}

type HAlign int

const (
	HAlignLeft = HAlign(iota)
	HAlignCenter
	HAlignRight
	HAlignJustify
)

var hAlignStrings = []string{"left", "center", "right", "justify"}

func (ha HAlign) String() string {
	if int(ha) < len(hAlignStrings) {
		return hAlignStrings[ha]
	}
	return "unknown"
}

type VAlign int

const (
	VAlignTop = VAlign(iota)
	VAlignMiddle
	VAlignBottom
	VAlignBaseline
)

var vAlignStrings = []string{"top", "middle", "bottom", "baseline"}

func (va VAlign) String() string {
	if int(va) < len(vAlignStrings) {
		return vAlignStrings[va]
	}
	return "unknown"
}
