package exif

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)





func compBytes(a []byte, b []byte) bool {
	lena := len(a)
	if lena != len(b) { // if they are not equal in length they cannot be equal
		return false 
	}
	// check if any byte is different
	for i := 0; i < lena; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true 
}

func forceASCII(s []byte) string {
	rs := make([]byte, 0, len(s))
	for _, r := range s {
	  if r <= 127 {
		rs = append(rs, r)
	  }
	}
	return string(rs)
  }


func decode(bo binary.ByteOrder,b []byte, i interface{}) error {
	switch i.(type) {
	case *uint, *uint8, *uint16, *uint32, *uint64, *int, *int8, *int32, *int64, *float32, *float64:
		return binary.Read(bytes.NewReader(b), bo, i)
	case *string:
		asc := forceASCII(b)	
		conv := i.(*string)
		*conv = asc
	default:
		return fmt.Errorf("Unsupported type %T", i)
	}
	return nil
}




type WindowReader struct {
	br *bufio.Reader
	curr int
}


func (wr *WindowReader) NextWindow(winLength int) ([]byte, error) {
	bts, err := wr.br.Peek(winLength)
	if err != nil {
		return nil, err
	}
	_, err = wr.br.Discard(winLength)
	if err != nil {
		return nil, err
	}
	wr.curr += winLength
	return bts, nil
}

func (wr *WindowReader) ReadUntil(when int) ([]byte, error) {
	readTo := make([]byte, when - wr.curr)
	n, err := io.ReadFull(wr.br, readTo); if err != nil { return nil, err}
	wr.curr += n
	return readTo, nil
}

func (wr *WindowReader) ScanUntil(toFind []byte) (bool, int, error) {
	scanned := 0
	for {
		bts, err := wr.br.Peek(len(toFind))
		if err != nil { return false, scanned, err } // Check for EOF?
		if compBytes(bts, toFind) {
			wr.br.Discard(len(toFind))
			wr.curr += len(toFind)
			return true, scanned, nil
		}
		_, err = wr.br.Discard(1)
		if err != nil {
			return false, scanned, err
		}
		scanned++
		wr.curr++
	}
}


func NewWindowReader(reader io.Reader) *WindowReader {
	br := bufio.NewReader(reader)
	return &WindowReader{br: br, curr: 0}
}

