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
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A Writer writes records to a CSV encoded file.
//
// As returned by NewWriter, a Writer writes records terminated by a
// newline and uses ',' as the field delimiter.  The exported fields can be
// changed to customize the details before the first call to Write or WriteAll.
//
// Comma is the field delimiter.
//
// If UseCRLF is true, the Writer ends each record with \r\n instead of \n.
type Writer struct {
	Comma     rune // Field delimiter (set to ',' by NewWriter)
	Quote     rune // Quote field (set to '"' by NewWriter)
	AllQuotes bool // True to always quote csv fields
	UseCRLF   bool // True to use \r\n as the line terminator
	w         *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Comma: ',',
		Quote: '"',
		w:     bufio.NewWriter(w),
	}
}

// writeFieldWithQuote writes a single CSV field (cell) with quotes
func (w *Writer) writeFieldWithQuote(field string) (err error) {

	if _, err = w.w.WriteRune(w.Quote); err != nil {
		return err
	}

	for _, r1 := range field {
		switch r1 {
		case w.Quote:
			_, err = w.w.WriteString(string([]rune{w.Quote, w.Quote}))
		case '\r':
			if !w.UseCRLF {
				err = w.w.WriteByte('\r')
			}
		case '\n':
			if w.UseCRLF {
				_, err = w.w.WriteString("\r\n")
			} else {
				err = w.w.WriteByte('\n')
			}
		default:
			_, err = w.w.WriteRune(r1)
		}
		if err != nil {
			return err
		}
	}

	if _, err = w.w.WriteRune(w.Quote); err != nil {
		return err
	}

	return nil
}

// Writer writes a single CSV record to w along with any necessary quoting.
// A record is a slice of strings with each string being one field.
func (w *Writer) Write(record []string) (err error) {
	for n, field := range record {
		if n > 0 {
			if _, err = w.w.WriteRune(w.Comma); err != nil {
				return err
			}
		}
		// If we don't have to have a quoted field then just
		// write out the field and continue to the next field.
		if !w.fieldNeedsQuotes(field) {
			if _, err = w.w.WriteString(field); err != nil {
				return err
			}
			continue
		}
		if err = w.writeFieldWithQuote(field); err != nil {
			return err
		}

	}
	if w.UseCRLF {
		_, err = w.w.WriteString("\r\n")
	} else {
		err = w.w.WriteByte('\n')
	}
	return
}

// Flush writes any buffered data to the underlying io.Writer.
// To check if an error occurred during the Flush, call Error.
func (w *Writer) Flush() {
	w.w.Flush()
}

// Error reports any error that has occurred during a previous Write or Flush.
func (w *Writer) Error() error {
	_, err := w.w.Write(nil)
	return err
}

// WriteAll writes multiple CSV records to w using Write and then calls Flush.
func (w *Writer) WriteAll(records [][]string) (err error) {
	for _, record := range records {
		err = w.Write(record)
		if err != nil {
			return err
		}
	}
	return w.w.Flush()
}

// fieldNeedsQuotes returns true if our field must be enclosed in quotes.
// Empty fields, files with a Comma, fields with a quote or newline, and
// fields which start with a space must be enclosed in quotes.
func (w *Writer) fieldNeedsQuotes(field string) bool {

	// If quotes are enforced by configuration, always return true
	if w.AllQuotes {
		return true
	}

	if w.Quote == 0 {
		return false
	}

	// If the field is empty, or has either CSV's Comma or Quote character(or rune), or the field has new-line character in it, return true
	if len(field) == 0 || strings.IndexRune(field, w.Comma) >= 0 || strings.IndexRune(field, w.Quote) >= 0 || strings.IndexAny(field, "\r\n") >= 0 {
		return true
	}

	// If the field is one of the valid "white-space" characters
	r1, _ := utf8.DecodeRuneInString(field)
	return unicode.IsSpace(r1)
}
