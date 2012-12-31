// Copyright 2012 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package pdf

import "testing"

func textPieceTestText() *TextPiece {
	font := testTtfFonts("Arial")[0]
	return &TextPiece{
		Text:     "Lorem",
		Font:     font,
		FontSize: 10,
	}
}

func TestTextPiece_IsWhiteSpace(t *testing.T) {
	font := testTtfFonts("Arial")[0]

	empty := TextPiece{
		Text:     "",
		Font:     font,
		FontSize: 10,
	}
	check(t, !empty.IsWhiteSpace(), "Empty string should not be considered whitespace.")

	singleWhite := TextPiece{
		Text:     " ",
		Font:     font,
		FontSize: 10,
	}
	check(t, singleWhite.IsWhiteSpace(), "A single space should be considered whitespace.")

	multiWhite := TextPiece{
		Text:     "  \t\n\v\f\r",
		Font:     font,
		FontSize: 10,
	}
	check(t, multiWhite.IsWhiteSpace(), "Multiple spaces should be considered whitespace.")

	startWhite := TextPiece{
		Text:     "  Lorem",
		Font:     font,
		FontSize: 10,
	}
	check(t, !startWhite.IsWhiteSpace(), "A piece that only starts with spaces should not be considered whitespace.")

	nonWhite := TextPiece{
		Text:     "Lorem",
		Font:     font,
		FontSize: 10,
	}
	check(t, !nonWhite.IsWhiteSpace(), "Piece contains no whitespace.")
}

func TestTextPiece_measure(t *testing.T) {
	piece := textPieceTestText()
	piece.measure(0, 0)
	expectNFdelta(t, "Ascent", 9.052734, piece.Ascent, 0.001)
	expectNFdelta(t, "Descent", -2.119141, piece.Descent, 0.001)
	expectNFdelta(t, "Height", 11.171875, piece.Height, 0.001)
	expectNFdelta(t, "UnderlinePosition", -1.059570, piece.UnderlinePosition, 0.001)
	expectNFdelta(t, "UnderlineThickness", 0.732422, piece.UnderlineThickness, 0.001)
	expectNFdelta(t, "Width", 28.344727, piece.Width, 0.001)
	expectNI(t, "Chars", 5, piece.Chars)
	expectNI(t, "Tokens", 1, piece.Tokens)
}

// 345 ns
func BenchmarkTextPiece_measure(b *testing.B) {
	b.StopTimer()
	font := testTtfFonts("Arial")[0]
	piece := &TextPiece{
		Text:     "Lorem",
		Font:     font,
		FontSize: 10,
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		piece.measure(0, 0)
	}
}
