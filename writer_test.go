// BSD 3-Clause License

// Copyright (c) 2018, Tushar Dwivedi
// All rights reserved.

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// An alternative to standard CSV implementation of Go (https://github.com/golang/go/tree/master/src/encoding/csv),
// and modified by Tushar Dwivedi (https://github.com/tushar2708/altcsv),
// to support some more delimiters, and a few non-standard(?) CSV formats.

package altcsv

import (
	"bytes"
	"errors"
	"testing"
)

var writeTests = []struct {
	Input     [][]string
	Output    string
	AllQuotes bool
	UseCRLF   bool
	Quote     rune
	Comma     rune
	ZeroQuote bool
}{
	{Input: [][]string{{"abc"}}, Output: "abc\n"},
	{Input: [][]string{{"abc"}}, Output: "abc\r\n", UseCRLF: true},
	{Input: [][]string{{`"abc"`}}, Output: `"""abc"""` + "\n"},
	{Input: [][]string{{`a"b`}}, Output: `"a""b"` + "\n"},
	{Input: [][]string{{`"a"b"`}}, Output: `"""a""b"""` + "\n"},
	{Input: [][]string{{" abc"}}, Output: `" abc"` + "\n"},
	{Input: [][]string{{"abc,def"}}, Output: `"abc,def"` + "\n"},
	{Input: [][]string{{"abc", "def"}}, Output: "abc,def\n"},
	{Input: [][]string{{"abc"}, {"def"}}, Output: "abc\ndef\n"},
	{Input: [][]string{{"abc\ndef"}}, Output: "\"abc\ndef\"\n"},
	{Input: [][]string{{"abc\ndef"}}, Output: "\"abc\r\ndef\"\r\n", UseCRLF: true},
	{Input: [][]string{{"abc\rdef"}}, Output: "\"abcdef\"\r\n", UseCRLF: true},
	{Input: [][]string{{"abc\rdef"}}, Output: "\"abc\rdef\"\n", UseCRLF: false},
	{Input: [][]string{{"abc,def"}}, Output: `|abc,def|` + "\n", Quote: '|'},
	{Input: [][]string{{"abc", "def"}}, Output: "abc;def\n", Comma: ';'},
	{Input: [][]string{{`a|b`}}, Output: `|a||b|` + "\n", Quote: '|'},
	{Input: [][]string{{`|a|b|`}}, Output: `|||a||b|||` + "\n", Quote: '|'},
	{Input: [][]string{{`a`, ``, `b`}}, Output: `a,,b` + "\n", ZeroQuote: true},
	{Input: [][]string{{"abc", "def"}}, Output: `"abc","def"` + "\n", AllQuotes: true},
	{Input: [][]string{{"abc", "def"}}, Output: "abc,def\n", AllQuotes: false},
	{Input: [][]string{{"a,bc", "de\nf"}}, Output: `"a,bc","de` + "\n" + `f"` + "\n", AllQuotes: true},
	{Input: [][]string{{"abc", "def"}, {"uvw", "xyz"}}, Output: `"abc","def"` + "\n" + `"uvw","xyz"` + "\n", AllQuotes: true},
}

func TestWrite(t *testing.T) {
	for n, tt := range writeTests {
		b := &bytes.Buffer{}
		f := NewWriter(b)
		f.UseCRLF = tt.UseCRLF
		f.AllQuotes = tt.AllQuotes

		if tt.Comma != 0 {
			f.Comma = tt.Comma
		}
		if tt.Quote != 0 {
			f.Quote = tt.Quote
		}
		if tt.ZeroQuote {
			f.Quote = '\000'
		}
		err := f.WriteAll(tt.Input)
		if err != nil {
			t.Errorf("Unexpected error: %s\n", err)
		}
		out := b.String()
		if out != tt.Output {
			t.Errorf("#%d: out=%q want %q", n, out, tt.Output)
		}
	}
}

type errorWriter struct{}

func (e errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("Test")
}

func TestError(t *testing.T) {
	b := &bytes.Buffer{}
	f := NewWriter(b)
	f.Write([]string{"abc"})
	f.Flush()
	err := f.Error()

	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}

	f = NewWriter(errorWriter{})
	f.Write([]string{"abc"})
	f.Flush()
	err = f.Error()

	if err == nil {
		t.Error("Error should not be nil")
	}
}
