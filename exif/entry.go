package exif

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Entry struct {
	Tag  string            // 2
	DataFormat uint16      // 2 4
	NumOfComponents uint32 // 4 8
	Offset uint32		   // 4 12
	offsetBts []byte
	byteOrder binary.ByteOrder
}

type DataFormat struct {
	Name string
	BtsPerComponents uint32
	ID uint16
}

func (ent Entry) GetDataFormat() DataFormat {
	return dataFormats[ent.DataFormat]
}

func (ent Entry) DataSize() int {
	df := ent.GetDataFormat()
	return int(ent.NumOfComponents * df.BtsPerComponents)
}

func (ent Entry) DecodeData(rd io.ReaderAt, i interface{}) error {
	dataSize := ent.DataSize()
	if dataSize <= 4 { // in case the data is smaller than 4 bytes we can get it from the offset bytes
		return decode(ent.byteOrder, ent.offsetBts, i)
	}
	dataBts := make([]byte, dataSize)
	_, err := rd.ReadAt(dataBts, int64(ent.Offset) + 12 ) // the offset, as stated in the offset field 
	if err != nil {											   // plus 12 bytes for the header
		return err
	}
	return decode(ent.byteOrder, dataBts, i)
}


const (
	UByte uint16 =  1
	AsciiString uint16 = 2
	UShort uint16 = 3
	ULong uint16 = 4
	Urational uint16 = 5
	Byte uint16 = 6
	Undefined uint16 = 7
	Short uint16 = 8
	Long uint16 = 9
	Rational uint16 = 10
	Float uint16 = 11
	DoubleFloat uint16 = 12
)


var dataFormats = map[uint16]DataFormat{
	1: {ID: 1, Name: "unsigned byte", BtsPerComponents: 1},
	2: {ID: 2, Name: "ascii strings", BtsPerComponents: 1},
	3: {ID: 3, Name: "unsigned short", BtsPerComponents: 2},
	4: {ID: 4, Name: "unsigned long", BtsPerComponents: 4},
	5: {ID: 5, Name: "unsigned rational", BtsPerComponents: 8},
	6: {ID: 6, Name: "signed byte", BtsPerComponents: 1},
	7: {ID: 7, Name: "undefined", BtsPerComponents: 1},
	8: {ID: 8, Name: "signed short", BtsPerComponents: 2},
	9: {ID: 9, Name: "signed long", BtsPerComponents: 4},
	10: {ID: 10, Name: "signed rational", BtsPerComponents: 8},
	11: {ID: 11, Name: "single float", BtsPerComponents: 4},
	12: {ID: 12, Name: "double float", BtsPerComponents: 8},
}


// DecodeEntry accepts a 12 byte array and returns an Entry
func DecodeEntry(bo binary.ByteOrder, bts []byte) (Entry, error) {
	entry := Entry{
		byteOrder: bo,
	}
	if len(bts) != 12 {
		return entry, fmt.Errorf("Expected a 12 bytes long slice but got %d", len(bts))
	}
	if bo == binary.BigEndian {
		entry.Tag = fmt.Sprintf("0x%x",string(bts[:2]))
	} else {
		entry.Tag = fmt.Sprintf("0x%x",string([]byte{bts[1], bts[0]})) // need to flip it
	}
	dataFormatBts := bts[2:4]
	numOFComponentsBts := bts[4:8]
	entry.offsetBts = bts[8:12]
	errs := []error {
		decode(bo, dataFormatBts, &entry.DataFormat),
		decode(bo, numOFComponentsBts, &entry.NumOfComponents),
		decode(bo, entry.offsetBts, &entry.Offset),
	}
	for _, err := range errs {
		if err != nil {
			return entry, err
		}
	}

	return entry, nil
}
