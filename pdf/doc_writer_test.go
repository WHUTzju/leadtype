// Copyright 2011-2012 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package pdf

import (
	"bytes"
	"testing"

	"github.com/rowland/leadtype/afm_fonts"
	"github.com/rowland/leadtype/codepage"
	"github.com/rowland/leadtype/colors"
	"github.com/rowland/leadtype/options"
	"github.com/rowland/leadtype/ttf_fonts"
)

func TestNewDocWriter(t *testing.T) {
	dw := NewDocWriter()
	check(t, dw.nextSeq != nil, "DocWriter should have nextSeq func")
	check(t, dw.file != nil, "DocWriter should have file")
	check(t, dw.catalog != nil, "DocWriter should have catalog")
	check(t, len(dw.file.body.list) == 4, "DocWriter file body should be initialized")
	check(t, dw.file.trailer.dict["Root"] != nil, "DocWriter file trailer root should be set")
	check(t, dw.resources != nil, "DocWriter resources should be initialized")
	check(t, dw.PagesAcross() == 1, "DocWriter PagesAcross should default to 1")
	check(t, dw.PagesDown() == 1, "DocWriter PagesDown should default to 1")
	check(t, dw.PagesUp() == 1, "DocWriter PagesUp should default to 1")
	check(t, dw.pages == nil, "DocWriter pages should be nil")
	check(t, dw.curPage == nil, "DocWriter curPage should be nil")
	check(t, !dw.inPage(), "DocWriter should not be in page yet")
}

func TestDocWriter_AddFontSource(t *testing.T) {
	dw := NewDocWriter()
	check(t, len(dw.fontSources) == 0, "No font sources should exist.")
	var fonts ttf_fonts.TtfFonts
	dw.AddFontSource(&fonts)
	check(t, dw.fontSources[0] == &fonts, "Font source should exist.")
}

func TestDocWriter_Close(t *testing.T) {
	var buf bytes.Buffer
	dw := NewDocWriter()
	dw.NewPage()
	dw.WriteTo(&buf)
	check(t, !dw.inPage(), "DocWriter should not be in page anymore")

	dw2 := NewDocWriter()
	dw2.WriteTo(&buf)
	check(t, len(dw.pages) == 1, "DocWriter pages should have minimum of 1 page after Close")
}

func TestDocWriter_fontKey(t *testing.T) {
	fc, err := ttf_fonts.New("/Library/Fonts/*.ttf")
	if err != nil {
		t.Fatal(err)
	}
	dw := NewDocWriter()
	dw.AddFontSource(fc)
	dw.NewPage()
	fonts, err := dw.AddFont("Arial", options.Options{})
	if err != nil {
		t.Fatal(err)
	}

	key1 := dw.fontKey(fonts[0], codepage.CodepageIndex(0))
	check(t, key1 == "F0", "1st fontKey should be F0.")
	key2 := dw.fontKey(fonts[0], codepage.CodepageIndex(0))
	check(t, key2 == "F0", "Same font and cpi should yield same key.")
	key3 := dw.fontKey(fonts[0], codepage.CodepageIndex(1))
	check(t, key3 == "F1", "2nd fontKey should be F1.")
}

func TestDocWriter_indexOfPage(t *testing.T) {
	dw := NewDocWriter()

	p1 := dw.NewPage()
	p2 := dw.NewPage()
	p3 := dw.NewPage()

	i1 := dw.indexOfPage(p1)
	check(t, i1 == 0, "1st PageWriter should have index 0")
	i2 := dw.indexOfPage(p2)
	check(t, i2 == 1, "2nd PageWriter should have index 1")
	i3 := dw.indexOfPage(p3)
	check(t, i3 == 2, "3rd PageWriter should have index 2")
}

func TestDocWriter_NewPage(t *testing.T) {
	dw := NewDocWriter()
	dw.NewPage()
	check(t, dw.curPage != nil, "DocWriter curPage should not be nil")
	check(t, dw.inPage(), "DocWriter should be in page now")
	check(t, len(dw.pages) == 1, "DocWriter should have 1 page now")
}

func TestDocWriter_NewPageAfter(t *testing.T) {
	dw := NewDocWriter()
	p1 := dw.NewPage()

	checkFatal(t, p1 != nil, "NewPageAfter should return a valid reference p1")
	checkFatal(t, len(dw.pages) == 1, "pages should have len 1")
	check(t, p1 == dw.pages[0], "p1 should be in 1st slot")

	check(t, p1.LineCapStyle() != ProjectingSquareCap, "LineCapStyle shouldn't default to ProjectingSquareCap")
	check(t, p1.LineColor() != colors.AliceBlue, "LineColor shouldn't default to AliceBlue")
	check(t, p1.LineDashPattern() != "dotted", "LineDashPattern shouldn't default to 'dotted'")
	check(t, p1.LineWidth("pt") != 42, "LineWidth shouldn't default to 42")

	p1.SetLineCapStyle(ProjectingSquareCap)
	p1.SetLineColor(colors.AliceBlue)
	p1.SetLineDashPattern("dotted")
	p1.SetLineWidth(42, "pt")

	p3 := dw.NewPageAfter(p1)

	checkFatal(t, p3 != nil, "NewPageAfter should return a valid reference p3")
	checkFatal(t, len(dw.pages) == 2, "pages should have len 2")
	check(t, p3 == dw.pages[1], "p3 should be in 2nd slot for now")

	check(t, p3.LineCapStyle() == ProjectingSquareCap, "LineCapStyle should still be ProjectingSquareCap")
	check(t, p3.LineColor() == colors.AliceBlue, "LineColor should still be AliceBlue")
	check(t, p3.LineDashPattern() == "dotted", "LineDashPattern should still be 'dotted'")
	check(t, p3.LineWidth("pt") == 42, "LineWidth should still be 42")

	p3.SetLineCapStyle(RoundCap)
	p3.SetLineColor(colors.AntiqueWhite)
	p3.SetLineDashPattern("dashed")
	p3.SetLineWidth(36, "pt")

	p2 := dw.NewPageAfter(p1)
	checkFatal(t, p2 != nil, "NewPageAfter should return a valid reference p2")
	checkFatal(t, len(dw.pages) == 3, "pages should have len 3")
	check(t, p2 == dw.pages[1], "p2 should now be in 2nd slot")
	check(t, p3 == dw.pages[2], "p3 should be in 3rd slot now")

	check(t, p2.LineCapStyle() == ProjectingSquareCap, "LineCapStyle should be ProjectingSquareCap")
	check(t, p2.LineColor() == colors.AliceBlue, "LineColor should be AliceBlue")
	check(t, p2.LineDashPattern() == "dotted", "LineDashPattern should be 'dotted'")
	check(t, p2.LineWidth("pt") == 42, "LineWidth should be 42")
}

func TestDocWriter_Open(t *testing.T) {
	dw := NewDocWriter()
	check(t, len(dw.options) == 0, "Default page options should be empty")
}

func TestDocWriter_PageHeight(t *testing.T) {
	dw := NewDocWriter()
	// height in (default) points
	expectF(t, 792, dw.PageHeight())
	dw.SetUnits("in")
	expectF(t, 11, dw.PageHeight())
	dw.SetUnits("cm")
	expectFdelta(t, 27.93, dw.PageHeight(), 0.01)
	// custom: "Dave points"
	UnitConversions.Add("dp", 0.072)
	dw.SetUnits("dp")
	expectF(t, 11000, dw.PageHeight())
}

func TestDocWriter_PageWidth(t *testing.T) {
	dw := NewDocWriter()
	// width in (default) points
	expectF(t, 612, dw.PageWidth())
	dw.SetUnits("in")
	expectF(t, 8.5, dw.PageWidth())
	dw.SetUnits("cm")
	expectFdelta(t, 21.58, dw.PageWidth(), 0.01)
	// custom: "Dave points"
	UnitConversions.Add("dp", 0.072)
	dw.SetUnits("dp")
	expectF(t, 8500, dw.PageWidth())
}

func TestDocWriter_SetFont_TrueType(t *testing.T) {
	dw := NewDocWriter()

	fc, err := ttf_fonts.New("/Library/Fonts/*.ttf")
	if err != nil {
		t.Fatal(err)
	}
	dw.AddFontSource(fc)

	dw.NewPage()

	st := SuperTest{t}
	st.True(dw.Fonts() == nil, "Document font list should be empty by default.")

	fonts, _ := dw.SetFont("Courier New", 10, options.Options{"weight": "Bold", "style": "Italic", "color": "AliceBlue"})
	st.Must(len(fonts) == 1, "length of fonts should be 1")
	st.Equal("Courier New", fonts[0].Family())
	st.True(10 == dw.FontSize())
	st.True(1.0 == fonts[0].RelativeSize)
	st.Equal("Bold", fonts[0].Weight)
	st.Equal("Italic", fonts[0].Style())
	st.True(fonts[0] == dw.Fonts()[0], "SetFont result should match new font list.")
	st.True(fonts[0].SubType() == "TrueType", "Font subType should be TrueType.")
}

func TestDocWriter_SetFont_Type1(t *testing.T) {
	dw := NewDocWriter()

	fc, err := afm_fonts.New("../afm/data/fonts/*.afm")
	if err != nil {
		t.Fatal(err)
	}
	dw.AddFontSource(fc)

	dw.NewPage()

	check(t, dw.Fonts() == nil, "Document font list should be empty by default.")

	fonts, _ := dw.SetFont("Courier", 10, options.Options{"weight": "Bold", "style": "Italic", "color": "AliceBlue"})
	checkFatal(t, len(fonts) == 1, "length of fonts should be 1")
	expectNS(t, "family", "Courier", fonts[0].Family())
	expectF(t, 10, dw.FontSize())
	expectF(t, 1.0, fonts[0].RelativeSize)
	expectNS(t, "weight", "Bold", fonts[0].Weight)
	expectNS(t, "style", "Italic", fonts[0].Style())
	check(t, fonts[0] == dw.Fonts()[0], "SetFont result should match new font list.")
	check(t, fonts[0].SubType() == "Type1", "Font subType should be Type1.")
}

func TestDocWriter_SetOptions(t *testing.T) {
	dw := NewDocWriter()
	dw.SetOptions(options.Options{"units": "in"})
	check(t, len(dw.options) == 1, "Default page options should have 1 option")
	check(t, dw.options["units"] == "in", "Default units should be in")
}

// TODO: TestPagesAcross
// TODO: TestPagesDown
// TODO: TestPagesUp
