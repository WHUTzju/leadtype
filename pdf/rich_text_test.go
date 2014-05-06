// Copyright 2012, 2013 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package pdf

import (
	"fmt"
	"leadtype/wordbreaking"
	"reflect"
	"testing"
	"unicode"
)

// import "fmt"

func TestNewRichText_English(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial")
	rt, err := NewRichText("abc", fonts, 10, Options{"color": Green, "underline": true, "line_through": true, "nobreak": true})
	if err != nil {
		t.Fatal(err)
	}
	st.Equal("abc", rt.Text)
	st.Equal(10.0, rt.FontSize)
	st.Equal(Green, rt.Color)
	st.True(rt.Underline)
	st.True(rt.LineThrough)
	st.True(rt.NoBreak)
	st.Equal(fonts[0], rt.Font, "Should be tagged with Arial font.")
}

func TestNewRichText_EnglishAndChinese_Fail(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial")
	_, err := NewRichText("abc所有测", fonts, 10, Options{})
	st.False(err == nil, "NewRichText should fail with Chinese text and only Arial.")
	st.Equal("No font found for 所有测.", err.Error())
}

func TestNewRichText_EnglishAndChinese_Pass(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial", "STSong")
	rt, err := NewRichText("abc所有测def", fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	st.Equal(3, len(rt.pieces))

	st.Equal("abc", rt.pieces[0].Text)
	st.Equal(10.0, rt.pieces[0].FontSize)
	st.Equal(Black, rt.pieces[0].Color, "Color should be (default) black.")
	st.False(rt.pieces[0].Underline, "Underline should be (default) false.")
	st.False(rt.pieces[0].LineThrough, "LineThrough should be (default) false.")

	st.Equal("所有测", rt.pieces[1].Text)
	st.Equal(10.0, rt.pieces[1].FontSize)
	st.Equal(Black, rt.pieces[1].Color, "Color should be (default) black.")
	st.False(rt.pieces[1].Underline, "Underline should be (default) false.")
	st.False(rt.pieces[1].LineThrough, "LineThrough should be (default) false.")

	st.Equal("def", rt.pieces[2].Text)
	st.Equal(10.0, rt.pieces[2].FontSize)
	st.Equal(Black, rt.pieces[2].Color, "Color should be (default) black.")
	st.False(rt.pieces[2].Underline, "Underline should be (default) false.")
	st.False(rt.pieces[2].LineThrough, "LineThrough should be (default) false.")

	st.Equal(fonts[0], rt.pieces[0].Font, "abc should be tagged with Arial font.")
	st.Equal(fonts[1], rt.pieces[1].Font, "Chinese should be tagged with STSong font.")
	st.Equal(fonts[0], rt.pieces[2].Font, "def should be tagged with Arial font.")
}

func TestNewRichText_ChineseAndEnglish(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial", "STSong")
	rt, err := NewRichText("所有测abc", fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	st.Must(len(rt.pieces) == 2)
	st.Equal("所有测", rt.pieces[0].Text)
	st.Equal(10.0, rt.pieces[0].FontSize)
	st.Equal("abc", rt.pieces[1].Text)
	st.Equal(10.0, rt.pieces[1].FontSize)
	st.Equal(fonts[1], rt.pieces[0].Font, "Should be tagged with Arial font.")
	st.Equal(fonts[0], rt.pieces[1].Font, "Should be tagged with STSong font.")
}

const englishRussianChinese = "Here is some Russian, Неприкосновенность, and some Chinese, 表明你已明确同意你的回答接受评估."

func TestNewRichText_EnglishRussianAndChineseLanguages(t *testing.T) {
	st := SuperTest{t}
	afmFonts := testAfmFonts("Helvetica")
	ttfFonts := testTtfFonts("Arial", "STSong")
	fonts := append(afmFonts, ttfFonts...)
	rt, err := NewRichText(englishRussianChinese, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	st.Must(len(rt.pieces) == 5)
	st.Equal("Here is some Russian, ", rt.pieces[0].Text)
	st.Equal(10.0, rt.pieces[0].FontSize)
	st.Equal("Неприкосновенность", rt.pieces[1].Text)
	st.Equal(10.0, rt.pieces[1].FontSize)
	st.Equal(", and some Chinese, ", rt.pieces[2].Text)
	st.Equal(10.0, rt.pieces[2].FontSize)
	st.Equal("表明你已明确同意你的回答接受评估", rt.pieces[3].Text)
	st.Equal(10.0, rt.pieces[3].FontSize)
	st.Equal(".", rt.pieces[4].Text)
	st.Equal(10.0, rt.pieces[4].FontSize)
	st.Equal(fonts[0], rt.pieces[0].Font, "Should be tagged with Helvetica font.")
	st.Equal(fonts[1], rt.pieces[1].Font, "Should be tagged with Arial font.")
	st.Equal(fonts[0], rt.pieces[2].Font, "Should be tagged with Helvetica font.")
	st.Equal(fonts[2], rt.pieces[3].Font, "Should be tagged with STSong font.")
	st.Equal(fonts[0], rt.pieces[4].Font, "Should be tagged with Helvetica font.")
}

// With Chinese font first in list, Arial is not called upon for English.
func TestNewRichText_ChineseAndEnglish_Reversed(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("STSong", "Arial")
	rt, err := NewRichText("所有测abc", fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	st.Must(len(rt.pieces) == 0)
	st.Equal("所有测abc", rt.Text)
	st.Equal(10.0, rt.FontSize)
	st.Equal(fonts[0], rt.Font, "Should be tagged with Arial font.")
}

func richTextMixedText() *RichText {
	fonts := testTtfFonts("Arial", "STSong")
	rt, err := NewRichText("abc所有测def", fonts, 10, Options{})
	if err != nil {
		panic(err)
	}
	return rt
}

func TestRichText_Ascent(t *testing.T) {
	st := SuperTest{t}
	p := new(RichText)
	st.Equal(0.0, p.Ascent())
	p = richTextMixedText()
	st.AlmostEqual(9.052734, p.Ascent(), 0.001)
}

func TestRichText_Chars(t *testing.T) {
	st := SuperTest{t}
	p := new(RichText)
	st.Equal(0, p.Chars())
	p = richTextMixedText()
	st.Equal(9, p.Chars())
}

func TestRichText_Descent(t *testing.T) {
	st := SuperTest{t}
	p := new(RichText)
	st.Equal(0.0, p.Descent())
	p = richTextMixedText()
	st.AlmostEqual(-2.119141, p.Descent(), 0.001)
}

func TestRichText_Height(t *testing.T) {
	st := SuperTest{t}
	p := new(RichText)
	st.Equal(0.0, p.Height())
	p = richTextMixedText()
	st.AlmostEqual(11.171875, p.Height(), 0.001)
}

func TestRichText_InsertStringAtOffsets(t *testing.T) {
	st := SuperTest{t}
	afmFonts := testAfmFonts("Helvetica")
	ttfFonts := testTtfFonts("Arial", "STSong")
	fonts := append(afmFonts, ttfFonts...)
	text := "Automatic "
	text1 := "hyphenation "
	text2 := "aids word wrapping."
	rt, err := NewRichText(text, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	rt, err = rt.Add(text1, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	rt, err = rt.Add(text2, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}

	offsets := []int{4, 16, 36}
	rt0 := rt.InsertStringAtOffsets("-", offsets)
	st.Equal("Auto-matic hyphen-ation aids word wrap-ping.", rt0.String())
}

func TestRichText_InsertStringAtOffsets_simple(t *testing.T) {
	st := SuperTest{t}
	p := &RichText{Text: "Automatic hyphenation aids word wrapping."}
	offsets := []int{4, 16, 36}
	p0 := p.InsertStringAtOffsets("-", offsets)
	st.Equal("Auto-matic hyphen-ation aids word wrap-ping.", p0.String())
}

func TestRichText_InsertStringAtOffsets_simple2(t *testing.T) {
	st := SuperTest{t}
	rt1 := &RichText{Text: "Automatic hyphenation "}
	rt2 := &RichText{Text: "aids word wrapping."}
	rt := &RichText{pieces: []*RichText{rt1, rt2}}
	offsets := []int{4, 16, 36}
	rt0 := rt.InsertStringAtOffsets("-", offsets)
	st.Equal("Auto-matic hyphen-ation aids word wrap-ping.", rt0.String())
}

func TestRichText_IsNewLine(t *testing.T) {
	st := SuperTest{t}
	font := testTtfFonts("Arial")[0]

	empty := RichText{
		Text:     "",
		Font:     font,
		FontSize: 10,
	}
	st.False(empty.IsNewLine(), "An empty string is not a newline.")

	newline := RichText{
		Text:     "\n",
		Font:     font,
		FontSize: 10,
	}
	st.True(newline.IsNewLine(), "It really is a newline.")

	nonNewline := RichText{
		Text:     "Lorem",
		Font:     font,
		FontSize: 10,
	}
	st.False(nonNewline.IsNewLine(), "This isn't a newline.")
}

func TestRichText_IsWhiteSpace(t *testing.T) {
	st := SuperTest{t}
	font := testTtfFonts("Arial")[0]

	empty := RichText{
		Text:     "",
		Font:     font,
		FontSize: 10,
	}
	st.False(empty.IsWhiteSpace(), "Empty string should not be considered whitespace.")

	singleWhite := RichText{
		Text:     " ",
		Font:     font,
		FontSize: 10,
	}
	st.True(singleWhite.IsWhiteSpace(), "A single space should be considered whitespace.")

	multiWhite := RichText{
		Text:     "  \t\n\v\f\r",
		Font:     font,
		FontSize: 10,
	}
	st.True(multiWhite.IsWhiteSpace(), "Multiple spaces should be considered whitespace.")

	startWhite := RichText{
		Text:     "  Lorem",
		Font:     font,
		FontSize: 10,
	}
	st.False(startWhite.IsWhiteSpace(), "A piece that only starts with spaces should not be considered whitespace.")

	nonWhite := RichText{
		Text:     "Lorem",
		Font:     font,
		FontSize: 10,
	}
	st.False(nonWhite.IsWhiteSpace(), "Piece contains no whitespace.")
}

func TestRichText_Len(t *testing.T) {
	st := SuperTest{t}
	rt := &RichText{
		pieces: []*RichText{
			{Text: "Hello"},
			{Text: " "},
			{Text: "World"},
		},
	}
	st.Equal(11, rt.Len())
}

func TestRichText_Len_complex(t *testing.T) {
	st := SuperTest{t}
	rt := &RichText{
		pieces: []*RichText{
			{Text: "Goodbye"},
			{pieces: []*RichText{
				{Text: ","},
				{Text: " "},
			}},
			{Text: "World"},
		},
	}
	st.Equal(14, rt.Len())
}

func TestRichText_Len_simple(t *testing.T) {
	st := SuperTest{t}
	rt := &RichText{Text: "Lorem"}
	st.Equal(5, rt.Len())
}

func arialText(s string) *RichText {
	fonts := testTtfFonts("Arial")
	rt, err := NewRichText(s, fonts, 10, Options{})
	if err != nil {
		panic(err)
	}
	return rt
}

func helveticaText(s string) *RichText {
	fonts := testAfmFonts("Helvetica")
	rt, err := NewRichText(s, fonts, 10, Options{})
	if err != nil {
		panic(err)
	}
	return rt
}

func courierNewText(s string) *RichText {
	fonts := testTtfFonts("Courier New")
	rt, err := NewRichText(s, fonts, 10, Options{})
	if err != nil {
		panic(err)
	}
	return rt
}

func TestRichText_MatchesAttributes(t *testing.T) {
	st := SuperTest{t}

	font := testTtfFonts("Arial")[0]
	p1 := RichText{
		Text:     "Lorem",
		Font:     font,
		FontSize: 10,
	}
	p2 := p1
	st.True(p1.MatchesAttributes(&p2), "Attributes should match.")

	p2.Font = testTtfFonts("Arial")[0]
	st.True(p1.MatchesAttributes(&p2), "Attributes should match.")

	p2.FontSize = 12
	st.False(p1.MatchesAttributes(&p2), "Attributes should not match.")

	p2 = p1
	p2.Color = Azure
	st.False(p1.MatchesAttributes(&p2), "Attributes should not match.")

	p2 = p1
	p2.Underline = true
	st.False(p1.MatchesAttributes(&p2), "Attributes should not match.")

	p2 = p1
	p2.LineThrough = true
	st.False(p1.MatchesAttributes(&p2), "Attributes should not match.")

	p2 = p1
	p2.CharSpacing = 1
	st.False(p1.MatchesAttributes(&p2), "Attributes should not match.")

	p2 = p1
	p2.WordSpacing = 1
	st.False(p1.MatchesAttributes(&p2), "Attributes should not match.")
}

func TestRichText_measure(t *testing.T) {
	piece := arialText("Lorem")
	piece.measure()
	expectNFdelta(t, "ascent", 9.052734, piece.ascent, 0.001)
	expectNFdelta(t, "descent", -2.119141, piece.descent, 0.001)
	expectNFdelta(t, "height", 11.171875, piece.height, 0.001)
	expectNFdelta(t, "UnderlinePosition", -1.059570, piece.UnderlinePosition, 0.001)
	expectNFdelta(t, "UnderlineThickness", 0.732422, piece.UnderlineThickness, 0.001)
	expectNFdelta(t, "width", 28.344727, piece.width, 0.001)
	expectNI(t, "chars", 5, piece.chars)
}

func TestRichText_Merge(t *testing.T) {
	st := SuperTest{t}
	afmFonts := testAfmFonts("Helvetica")
	ttfFonts := testTtfFonts("Arial", "STSong")
	fonts := append(afmFonts, ttfFonts...)
	text := "Here is some "
	text1 := "Russian, Неприкосновенность, "
	text2 := "and some Chinese, 表明你已明确同意你的回答接受评估."
	original, err := NewRichText(text, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	original, err = original.Add(text1, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	original, err = original.Add(text2, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println("Original pieces: ", original.Count())
	// i := 0
	// original.VisitAll(func(p *RichText) {
	// 	i++
	// 	fmt.Println(i, p.Text, len(p.Pieces))
	// })

	// Node
	//     Leaf: Here is some  0
	//     Node
	// 	       Leaf: Russian,
	//         Leaf: Неприкосновенность
	//         Leaf: ,
	//     Node
	//         Leaf: and some Chinese,
	//         Leaf: 表明你已明确同意你的回答接受评估
	//         Leaf: .

	piece0 := *original.pieces[0]
	merged := original.Merge()

	// fmt.Println("Merged pieces: ", merged.Count())
	// i = 0
	// merged.VisitAll(func(p *RichText) {
	// 	i++
	// 	fmt.Println(i, p.Text, len(p.pieces))
	// })

	// Node
	// Leaf: Here is some Russian,
	// Leaf: Неприкосновенность
	// Leaf: , and some Chinese,
	// Leaf: 表明你已明确同意你的回答接受评估
	// Leaf: .

	st.Equal(text+text1+text2, merged.String())
	st.Must(len(merged.pieces) == 5)

	st.Equal("Here is some Russian, ", merged.pieces[0].Text)
	st.Equal(10.0, merged.pieces[0].FontSize)
	st.Equal(fonts[0], merged.pieces[0].Font, "Should be tagged with Helvetica font.")

	st.Equal("Неприкосновенность", merged.pieces[1].Text)
	st.Equal(10.0, merged.pieces[1].FontSize)
	st.Equal(fonts[1], merged.pieces[1].Font, "Should be tagged with Arial font.")

	st.Equal(", and some Chinese, ", merged.pieces[2].Text)
	st.Equal(10.0, merged.pieces[2].FontSize)
	st.Equal(fonts[0], merged.pieces[2].Font, "Should be tagged with Helvetica font.")

	st.Equal("表明你已明确同意你的回答接受评估", merged.pieces[3].Text)
	st.Equal(10.0, merged.pieces[3].FontSize)
	st.Equal(fonts[2], merged.pieces[3].Font, "Should be tagged with STSong font.")

	st.Equal(".", merged.pieces[4].Text)
	st.Equal(10.0, merged.pieces[4].FontSize)
	st.Equal(fonts[0], merged.pieces[4].Font, "Should be tagged with Helvetica font.")

	st.Equal(piece0.Text, original.pieces[0].Text, "Original should be unchanged.")
}

func mixedText() *RichText {
	afmFonts := testAfmFonts("Helvetica")
	ttfFonts := testTtfFonts("Arial", "STSong")
	fonts := append(afmFonts, ttfFonts...)
	rt, err := NewRichText(englishRussianChinese, fonts, 10, Options{})
	if err != nil {
		panic(err)
	}
	return rt
}

func TestRichText_Split(t *testing.T) {
	st := SuperTest{t}
	rt := mixedText()
	byteOffsets := []int{0, 5, 8, 13, 22, 60, 64, 69, 78}
	words := []string{"Here ", "is ", "some ", "Russian, ", "Неприкосновенность, ", "and ", "some ", "Chinese, ", "表明你已明确同意你的回答接受评估."}
	for i := len(byteOffsets) - 1; i >= 0; i-- {
		var word *RichText
		rt, word = rt.Split(byteOffsets[i])
		st.Equal(words[i], word.String())
	}
}

func TestRichText_Split_simple(t *testing.T) {
	st := SuperTest{t}
	rt := &RichText{Text: "Hello, World!"}
	left, right := rt.Split(7)
	st.Equal("Hello, ", left.String())
	st.Equal("World!", right.String())
}

func TestRichText_Split_simple2(t *testing.T) {
	st := SuperTest{t}
	rt1 := &RichText{Text: "Hello, "}
	rt2 := &RichText{Text: "World!"}
	rt := &RichText{pieces: []*RichText{rt1, rt2}}
	left, right := rt.Split(7)
	st.Equal("Hello, ", left.String())
	st.Equal("World!", right.String())
}

func TestRichText_String(t *testing.T) {
	st := SuperTest{t}
	rt := mixedText()
	st.Equal(englishRussianChinese, rt.String())
}

const (
	leadingWhiteSpaceText                   = "\n\t  Here is some text with leading whitespace."
	leadingWhiteSpaceText_trimmed           = "Here is some text with leading whitespace."
	trailingwhiteSpaceText_simple           = "Here is some text with trailing whitespace \t\n"
	trailingwhiteSpaceText_simple_trimmed   = "Here is some text with trailing whitespace"
	whiteSpaceText_complex1                 = "Here is some text "
	whiteSpaceText_complex2                 = "with trailing whitespace \t\n"
	leadingAndTrailingWhitespaceText        = leadingWhiteSpaceText + " " + trailingwhiteSpaceText_simple
	leadingAndTrailingWhitespaceTextTrimmed = leadingWhiteSpaceText_trimmed + " " + trailingwhiteSpaceText_simple_trimmed
)

func TestRichText_TrimLeftSpace_simple(t *testing.T) {
	st := SuperTest{t}
	t1 := &RichText{Text: leadingWhiteSpaceText}
	st.Equal(leadingWhiteSpaceText, t1.String())
	t2 := t1.TrimLeftSpace()
	st.Equal(leadingWhiteSpaceText_trimmed, t2.String())
}

func TestRichText_TrimRightSpace_simple(t *testing.T) {
	st := SuperTest{t}
	t1 := &RichText{Text: trailingwhiteSpaceText_simple}
	st.Equal(trailingwhiteSpaceText_simple, t1.String())
	t2 := t1.TrimRightSpace()
	st.Equal(trailingwhiteSpaceText_simple_trimmed, t2.String())
}

func TestRichText_TrimSpace(t *testing.T) {
	st := SuperTest{t}
	t1 := &RichText{Text: leadingAndTrailingWhitespaceText}
	st.Equal(leadingAndTrailingWhitespaceText, t1.String())
	t2 := t1.TrimSpace()
	st.Equal(leadingAndTrailingWhitespaceTextTrimmed, t2.String())
}

func TestRichText_TrimRightFunc_complex(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial")
	t1, err := NewRichText(whiteSpaceText_complex1, fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	t2, err2 := t1.Add(whiteSpaceText_complex2, fonts, 10, Options{})
	if err2 != nil {
		t.Fatal(err2)
	}
	st.Equal(trailingwhiteSpaceText_simple, t2.String())
	t3 := t2.TrimRightFunc(unicode.IsSpace)
	st.Equal(trailingwhiteSpaceText_simple_trimmed, t3.String())
}

func TestRichText_Width(t *testing.T) {
	st := SuperTest{t}
	p := new(RichText)
	st.Equal(0.0, p.Width())
	p = richTextMixedText()
	st.AlmostEqual(60.024414, p.Width(), 0.001)
}

func TestRichText_WordsToWidth_empty(t *testing.T) {
	st := SuperTest{t}
	p := new(RichText)
	flags := make([]wordbreaking.Flags, p.Len())
	line, remainder, lineFlags, remainderFlags, err := p.WordsToWidth(10, flags, false)
	st.Equal(p, line, "Should return original piece.")
	if remainder != nil {
		t.Error("There should be nothing left over.")
	}
	st.True(reflect.DeepEqual(flags, lineFlags), "Should return original flags.")
	if remainderFlags != nil {
		t.Error("There should be no flags remaining.")
	}
	st.True(err == nil, "No error is expected.")
}

func TestRichText_WordsToWidth_short(t *testing.T) {
	st := SuperTest{t}
	p := arialText("Lorem")
	flags := make([]wordbreaking.Flags, p.Len())
	wordbreaking.MarkRuneAttributes(p.String(), flags)
	line, remainder, lineFlags, remainderFlags, err := p.WordsToWidth(30, flags, false)
	st.Equal(p, line, "Should return original piece.")
	if remainder != nil {
		t.Error("There should be nothing left over.")
	}
	st.True(reflect.DeepEqual(flags, lineFlags), "Should return original flags.")
	if remainderFlags != nil {
		t.Error("There should be no flags remaining.")
	}
	st.True(err == nil, "No error is expected.")
}

func TestRichText_WordsToWidth_mixed(t *testing.T) {
	st := SuperTest{t}
	p := mixedText()
	flags := make([]wordbreaking.Flags, p.Len())
	wordbreaking.MarkRuneAttributes(p.String(), flags)
	line, remainder, lineFlags, remainderFlags, err := p.WordsToWidth(60, flags, false)
	st.MustNot(line == nil, "Line must not be nil.")
	st.Equal("Here is some", line.String())
	st.MustNot(remainder == nil, "There should be text left over.")
	st.Equal(" Russian, Неприкосновенность, and some Chinese, 表明你已明确同意你的回答接受评估.", remainder.String())
	st.Equal(12, len(lineFlags))
	st.Equal(115, len(remainderFlags))
	if remainderFlags == nil {
		t.Error("There should be flags remaining.")
	}
	st.Equal(nil, err, "No error is expected.")
}

func TestRichText_WordsToWidth_hardbreak(t *testing.T) {
	st := SuperTest{t}
	p := arialText("Supercalifragilisticexpialidocious")
	flags := make([]wordbreaking.Flags, p.Len())
	wordbreaking.MarkRuneAttributes(p.String(), flags)
	line, remainder, lineFlags, remainderFlags, err := p.WordsToWidth(60, flags, true)
	st.MustNot(line == nil, "Line must not be nil.")
	st.Equal("Supercalifrag", line.String())
	st.MustNot(remainder == nil, "There should be text left over.")
	st.Equal("ilisticexpialidocious", remainder.String())
	st.Equal(13, len(lineFlags))
	st.Equal(21, len(remainderFlags))
	st.Equal(nil, err, "No error is expected.")
}

func TestRichText_WrapToWidth_hyphenated(t *testing.T) {
	const hyphenatedText = "Super-cali-fragil-istic-expi-ali-do-cious"
	expected := []string{
		"Super-cali-",
		"fragil-istic-",
		"expi-ali-do-",
		"cious",
	}
	rt := arialText(hyphenatedText)
	flags := make([]wordbreaking.Flags, rt.Len())
	wordbreaking.MarkRuneAttributes(rt.String(), flags)
	lines, err := rt.WrapToWidth(60, flags, false)
	if err != nil {
		t.Fatal(err)
	}
	st := SuperTest{t}
	st.Must(len(lines) == len(expected), "Unexpected number of lines.")
	for i := range lines {
		st.Equal(expected[i], lines[i].String())
	}
}

func TestRichText_WrapToWidth_soft_hyphenated(t *testing.T) {
	const hyphenatedText = "Super\u00ADcali\u00ADfragil\u00ADistic\u00ADexpi\u00ADali\u00ADdo\u00ADcious"
	expected := []string{
		"Super\u00ADcali\u00AD-",
		"fragil\u00AD-",
		"istic\u00ADexpi\u00AD-",
		"ali\u00ADdo\u00ADcious",
	}
	rt := courierNewText(hyphenatedText)
	flags := make([]wordbreaking.Flags, rt.Len())
	wordbreaking.MarkRuneAttributes(rt.String(), flags)
	lines, err := rt.WrapToWidth(62, flags, false)
	if err != nil {
		t.Fatal(err)
	}
	st := SuperTest{t}
	st.Must(len(lines) == len(expected), "Unexpected number of lines.")
	for i := range lines {
		st.Equal(expected[i], lines[i].String())
		st.True(lines[i].Width() <= 62, "Line too long.")
	}
}

// Should not add hyphen-minus to end of string that does not end with a soft hyphen.
func TestRichText_WrapToWidth_soft_hyphenated2(t *testing.T) {
	const hyphenatedText = "Super\u00ADbad example"
	expected := []string{
		"Super\u00ADbad",
		"example",
	}
	rt := courierNewText(hyphenatedText)
	flags := make([]wordbreaking.Flags, rt.Len())
	wordbreaking.MarkRuneAttributes(rt.String(), flags)
	lines, err := rt.WrapToWidth(62, flags, false)
	if err != nil {
		t.Fatal(err)
	}
	st := SuperTest{t}
	st.Must(len(lines) == len(expected), "Unexpected number of lines.")
	for i := range lines {
		st.Equal(expected[i], lines[i].String())
		st.True(lines[i].Width() <= 62, "Line too long.")
	}
}

func TestRichText_WrapToWidth_short(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial")
	p, err := NewRichText("Lorem ipsum.", fonts, 10, Options{"nobreak": true})
	if err != nil {
		t.Fatal(err)
	}
	flags := make([]wordbreaking.Flags, p.Len())
	wordbreaking.MarkRuneAttributes(p.String(), flags)
	lines, err := p.WrapToWidth(100, flags, false)
	if err != nil {
		t.Fatal(err)
	}
	st.Must(len(lines) == 1, "Should return one piece.")
	st.Equal(p.Text, lines[0].Text, "Text should match.")
	st.True(p.MatchesAttributes(lines[0]))
}

func TestRichText_WrapToWidth_mixed(t *testing.T) {
	st := SuperTest{t}
	p := mixedText()
	flags := make([]wordbreaking.Flags, p.Len())
	wordbreaking.MarkRuneAttributes(p.String(), flags)
	lines, err := p.WrapToWidth(60, flags, false)
	if err != nil {
		t.Fatal(err)
	}
	st.Must(len(lines) > 1, "Should return at least one piece.")
	expected := []string{
		"Here is some",
		"Russian,",
		"Неприкосновенность,",
		"and some",
		"Chinese,",
		"表明你已明确同意你的回答接受评估.",
	}
	for i := 0; i < len(expected); i++ {
		st.Equal(expected[i], lines[i].String(), "Unexpected text.")
	}
}

func TestRichText_WrapToWidth_hardBreak(t *testing.T) {
	st := SuperTest{t}
	p := arialText("Supercalifragilisticexpialidocious")
	flags := make([]wordbreaking.Flags, p.Len())
	wordbreaking.MarkRuneAttributes(p.String(), flags)
	lines, err := p.WrapToWidth(60, flags, true)
	if err != nil {
		t.Fatal(err)
	}
	st.Must(len(lines) == 3, "Should return 3 pieces.")
	expected := []string{
		"Supercalifrag",
		"ilisticexpialid",
		"ocious",
	}
	for i := 0; i < len(expected); i++ {
		st.Equal(expected[i], lines[i].String(), "Unexpected text.")
	}
}

func TestRichText_WrapToWidth_nobreak_simple(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial")
	rt, err := NewRichText("Here is a long sentence with mostly small words.", fonts, 10, Options{"nobreak": true})
	if err != nil {
		t.Fatal(err)
	}
	flags := make([]wordbreaking.Flags, rt.Len())
	wordbreaking.MarkRuneAttributes(rt.String(), flags)
	rt.MarkNoBreak(flags)
	lines, err2 := rt.WrapToWidth(60, flags, false)
	if err2 != nil {
		t.Fatal(err2)
	}
	st.Equal(1, len(lines), "Should return a single line.")
}

func dump(rt *RichText) {
	fmt.Println("---")
	rt.VisitAll(func(p *RichText) {
		fmt.Printf("Text: %s, FontSize: %v, NoBreak: %v\n", p.Text, p.FontSize, p.NoBreak)
	})
}

func TestRichText_WrapToWidth_nobreak_complex(t *testing.T) {
	st := SuperTest{t}
	fonts := testTtfFonts("Arial")
	var rt *RichText
	var err error
	rt, err = NewRichText("Here is a ", fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	rt, err = rt.Add("long sentence with mostly", fonts, 10, Options{"nobreak": true})
	if err != nil {
		t.Fatal(err)
	}
	rt, err = rt.Add(" small words.", fonts, 10, Options{})
	if err != nil {
		t.Fatal(err)
	}
	flags := make([]wordbreaking.Flags, rt.Len())
	wordbreaking.MarkRuneAttributes(rt.String(), flags)
	rt.MarkNoBreak(flags)
	lines, err := rt.WrapToWidth(60, flags, false)
	if err != nil {
		t.Fatal(err)
	}
	st.Equal(3, len(lines), "Should return 3 lines.")
}

// 14,930 ns go1.1.2
// 14,685 ns go1.2.1
func BenchmarkNewRichText(b *testing.B) {
	b.StopTimer()
	afmFonts := testAfmFonts("Helvetica")
	ttfFonts := testTtfFonts("Arial", "STSong")
	fonts := append(afmFonts, ttfFonts...)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := NewRichText(englishRussianChinese, fonts, 10, Options{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 75.5 ns
// 76.9 ns go1.1.1
// 70.8 ns go1.1.2
// 74.1 ns go1.2.1
func BenchmarkRichText_IsWhiteSpace(b *testing.B) {
	b.StopTimer()
	font := testTtfFonts("Arial")[0]
	piece := &RichText{
		Text:     "  \t\n\v\f\r",
		Font:     font,
		FontSize: 10,
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		piece.IsWhiteSpace()
	}
}

// 345 ns
// 301 ns go1.1.1
// 296 ns go1.1.2
// 299 ns go1.2.1
func BenchmarkRichText_measure(b *testing.B) {
	b.StopTimer()
	piece := arialText("Lorem")
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		piece.measure()
	}
}

// 4488 ns go1.1.2
// 4591 ns go1.2.1
func BenchmarkRichText_Ascent(b *testing.B) {
	b.StopTimer()
	text := mixedText()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		text.ascent = 0
		text.pieces[0].ascent = 0
		text.pieces[1].ascent = 0
		text.pieces[2].ascent = 0
		text.pieces[3].ascent = 0
		text.pieces[4].ascent = 0
		text.Ascent()
	}
}

// 4547 ns go1.1.2
// 4591 ns go1.2.1
func BenchmarkRichText_Chars(b *testing.B) {
	b.StopTimer()
	text := mixedText()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		text.chars = 0
		text.pieces[0].chars = 0
		text.pieces[1].chars = 0
		text.pieces[2].chars = 0
		text.pieces[3].chars = 0
		text.pieces[4].chars = 0
		text.Chars()
	}
}

// 4492 ns go1.1.2
// 4568 ns go1.2.1
func BenchmarkRichText_Descent(b *testing.B) {
	b.StopTimer()
	text := mixedText()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		text.descent = 0
		text.pieces[0].descent = 0
		text.pieces[1].descent = 0
		text.pieces[2].descent = 0
		text.pieces[3].descent = 0
		text.pieces[4].descent = 0
		text.Descent()
	}
}

// 4559 ns go1.1.2
// 4593 ns go1.2.1
func BenchmarkRichText_Height(b *testing.B) {
	b.StopTimer()
	text := mixedText()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		text.height = 0
		text.pieces[0].height = 0
		text.pieces[1].height = 0
		text.pieces[2].height = 0
		text.pieces[3].height = 0
		text.pieces[4].height = 0
		text.Height()
	}
}

// 30.2 ns go1.1.2
// 30.4 ns go1.2.1
func BenchmarkRichText_Len(b *testing.B) {
	b.StopTimer()
	text := mixedText()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		text.Len()
	}
}

// 4539 ns go1.1.2
// 4636 ns go1.2.1
func BenchmarkRichText_Width(b *testing.B) {
	b.StopTimer()
	text := mixedText()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		text.width = 0
		text.pieces[0].width = 0
		text.pieces[1].width = 0
		text.pieces[2].width = 0
		text.pieces[3].width = 0
		text.pieces[4].width = 0
		text.Width()
	}
}
