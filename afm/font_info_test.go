// Copyright 2011-2014 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package afm

import "testing"

func TestLoadFontInfo(t *testing.T) {
	fi, err := LoadFontInfo("data/fonts/Helvetica.afm")
	if err != nil {
		t.Fatalf("Error loading font info: %v", err)
	}

	expectS(t, "Filename", "data/fonts/Helvetica.afm", fi.Filename())
	expectS(t, "Family", "Helvetica", fi.Family())
	expectS(t, "Style", "Medium", fi.Style())

	expectI(t, "Ascent", 718, fi.Ascent())
	// expectI(t, "AvgWidth", 904, fi.AvgWidth())
	expectI(t, "CapHeight", 718, fi.CapHeight())
	expectS(t, "Copyright", "(c) 1985, 1987, 1989, 1990, 1997 Adobe Systems Incorporated.  All Rights Reserved.Helvetica is a trademark of Linotype-Hell AG and/or its subsidiaries.", fi.Copyright())
	expectI(t, "Descent", -207, fi.Descent())
	// expectS(t, "Designer", "Monotype Type Drawing Office - Robin Nicholas, Patricia Saunders 1982", fi.Designer())
	// expect(t, "Embeddable", fi.Embeddable())
	expectS(t, "FullName", "Helvetica", fi.FullName())
	expect(t, "IsFixedPitch", !fi.IsFixedPitch())
	expectF(t, "ItalicAngle", 0, fi.ItalicAngle())
	// expectS(t, "License", "You may use this font to display and print content as permitted by the license terms for the product in which this font is included. You may only (i) embed this font in content as permitted by the embedding restrictions included in this font; and (ii) temporarily download this font to a printer or other output device to help print content.", fi.License())
	expectI(t, "Leading", (718 - -207)*120/100, fi.Leading())
	// expectS(t, "Manufacturer", "The Monotype Corporation", fi.Manufacturer())
	expectS(t, "PostScriptName", "Helvetica", fi.PostScriptName())
	expectI(t, "StemV", 88, fi.StemV())
	// expectS(t, "Trademark", "Arial is a trademark of The Monotype Corporation in the United States and/or other countries.", fi.Trademark())
	expectI(t, "UnderlinePosition", -100, fi.UnderlinePosition())
	expectI(t, "UnderlineThickness", 50, fi.UnderlineThickness())
	// expectS(t, "UniqueName", "Monotype:Arial Regular:Version 5.01 (Microsoft)", fi.UniqueName())
	expectS(t, "Version", "002.000", fi.Version())
	expectS(t, "Weight", "Medium", fi.Weight())
	expectI(t, "XHeight", 523, fi.XHeight())

	expectI(t, "BoundingBox[0]", -166, fi.BoundingBox()[0])
	expectI(t, "BoundingBox[1]", -225, fi.BoundingBox()[1])
	expectI(t, "BoundingBox[2]", 1000, fi.BoundingBox()[2])
	expectI(t, "BoundingBox[3]", 931, fi.BoundingBox()[3])
}

// 323,331 ns
// 169,941 ns
// 169,467 ns
// 161,887 ns weekly.2012-02-22
// 165,638 ns go1
// 153,228 ns go1.1.1
// 148,584 ns go1.2.1
// 166,227 ns go1.4.2
//  66,428 ns go1.6.2 mbp
//  56,820 ns go1.7.3
func BenchmarkLoadFontInfo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LoadFontInfo("data/fonts/Helvetica.afm")
	}
}
