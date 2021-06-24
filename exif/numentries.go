package exif

import "encoding/binary"


func DecodeNumEntries(byteOrder binary.ByteOrder, bts []byte) (int, error) {
	var entries uint16
	return int(entries), decode(byteOrder, bts, &entries)
}


