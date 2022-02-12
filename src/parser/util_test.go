package parser

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestSafeXMLReader(t *testing.T) {
	var f io.Reader
	want := []byte("привет мир")
	f = bytes.NewReader(want)
	f = NewSafeXMLReader(f)

	have, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Fatalf("invalid output\nwant: %v\nhave: %v", want, have)
	}
}

func TestSafeXMLReaderRemoveUnwantedRunes(t *testing.T) {
	var f io.Reader
	input := []byte("\aпривет \x0cмир\ufffe\uffff")
	want := []byte("привет мир")
	f = bytes.NewReader(input)
	f = NewSafeXMLReader(f)

	have, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(want, have) {
		t.Fatalf("invalid output\nwant: %v\nhave: %v", want, have)
	}
}

func TestSafeXMLReaderPartial1(t *testing.T) {
	var f io.Reader
	input := []byte("\aпривет \x0cмир\ufffe\uffff")
	want := []byte("привет мир")
	f = bytes.NewReader(input)
	f = NewSafeXMLReader(f)

	buf := make([]byte, 1)
	for i := 0; i < len(want); i++ {
		n, err := f.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		if n != 1 {
			t.Fatalf("expected 1 byte, got %d", n)
		}
		if buf[0] != want[i] {
			t.Fatalf("invalid char at pos %d\nwant: %v\nhave: %v", i, want[i], buf[0])
		}
	}
	if x, err := f.Read(buf); err != io.EOF {
		t.Fatalf("expected EOF, %v, %v %v", buf, x, err)
	}
}

func TestSafeXMLReaderPartial2(t *testing.T) {
	var f io.Reader
	input := []byte("привет\a\a\a\a\a")
	f = bytes.NewReader(input)
	f = NewSafeXMLReader(f)

	buf := make([]byte, 12)
	n, err := f.Read(buf)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != 12 {
		t.Fatalf("expected 12 bytes")
	}

	n, err = f.Read(buf)
	if n != 0 {
		t.Fatalf("expected 0")
	}
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
}
