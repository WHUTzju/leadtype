// Copyright 2011-2012 Brent Rowland.
// Use of this source code is governed the Apache License, Version 2.0, as described in the LICENSE file.

package ttf

import (
	"encoding/binary"
	"io"
	"unicode/utf16"
)

type Fixed struct {
	base int16
	frac uint16
}

func (f *Fixed) Read(file io.Reader) error {
	return readValues(file, &f.base, &f.frac)
}

func (f *Fixed) Tof64() float64 {
	return float64(f.base) + float64(f.frac)/65536
}

type FWord int16
type uFWord uint16

type longDateTime int64

type PANOSE [10]byte

func readValues(r io.Reader, values ...interface{}) (err error) {
	for _, v := range values {
		if err = binary.Read(r, binary.BigEndian, v); err != nil {
			return
		}
	}
	return
}

func countUniqueUint16Values(ary []uint16) int {
	m := make(map[uint16]bool)
	for _, v := range ary {
		m[v] = true
	}
	return len(m)
}

func utf16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[0:i]
			break
		}
	}
	return string(utf16.Decode(s))
}
