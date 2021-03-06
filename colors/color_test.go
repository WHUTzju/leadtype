// Copyright 2011-2012 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package colors

import (
	"testing"
)

func TestColor_RGB(t *testing.T) {
	c := Color(0x030507)
	r, g, b := c.RGB()
	expectNI(t, "red", 3, int(r))
	expectNI(t, "green", 5, int(g))
	expectNI(t, "blue", 7, int(b))
}

func TestColor_RGB64(t *testing.T) {
	c := Color(0xFF3300)
	r, g, b := c.RGB64()
	expectF(t, 1.00, r)
	expectF(t, 0.20, g)
	expectF(t, 0.00, b)
}

func TestNamedColor(t *testing.T) {
	c1, err1 := NamedColor("AliceBlue")
	check(t, err1 == nil, "Error looking up AliceBlue.")
	check(t, c1 == AliceBlue, "Expecting AliceBlue.")

	c2, err2 := NamedColor("FF0000")
	check(t, err2 == nil, "Error lookup up FF0000.")
	check(t, c2 == Red, "Expecting Red.")

	_, err3 := NamedColor("")
	check(t, err3 != nil, "Expecting error looking up empty string.")
	check(t, err3.Error() == "Expected name of color or numeric value in hex format.", "Unexpected error message looking up color.")
}

func TestColor_String(t *testing.T) {
	if AliceBlue.String() != "aliceblue" {
		t.Errorf("Expecting <%s>, got <%s>.", "aliceblue", AliceBlue.String())
	}
	if YellowGreen.String() != "yellowgreen" {
		t.Errorf("Expecting <%s>, got <%s>.", "yellowgreen", YellowGreen.String())
	}
	if Color(0xABCDEF).String() != "ABCDEF" {
		t.Errorf("Expecting <%s>, got <%s>.", "ABCDEF", Color(0xABCDEF).String())
	}
}

func expectNI(t *testing.T, name string, expected, actual int) {
	if expected != actual {
		t.Errorf("%s: expected %d, got %d", name, expected, actual)
	}
}

func expectF(t *testing.T, expected, actual float64) {
	if expected != actual {
		t.Errorf("Expected %f, got %f", expected, actual)
	}
}

func check(t *testing.T, condition bool, msg string) {
	if !condition {
		t.Error(msg)
	}
}
