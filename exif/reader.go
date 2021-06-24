package exif

import (
	"encoding/binary"
	"fmt"
	"os"
)

const SizeDataLen = 2
const ExifStringDataLen = 4

var ExifMarker = []byte{0xFF,  0xE1}
var fileEndian binary.ByteOrder



func GetExifData(fname string) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	fd, err := os.Open(fname) // open file read only
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	wr := NewWindowReader(fd)
	found, _, err := wr.ScanUntil(ExifMarker)
	if err != nil { return nil, err }
	if !found {
		return nil, fmt.Errorf("could not read exif data from %s", fname)
	}
	headerBts, err := wr.NextWindow(16); if err != nil { return nil, err }
	exifHeader, err := DecodeHeader(headerBts); if err != nil { return nil, err}
	fileEndian = exifHeader.ByteOrder
	if int(exifHeader.Offset) > wr.curr {
		wr.NextWindow(int(exifHeader.Offset) - wr.curr)
	}

	// Start IFD
	entryNumber, err := wr.NextWindow(2); if err != nil { return nil, err }
	numOfEntries, err  := DecodeNumEntries(fileEndian, entryNumber); if err != nil { return nil, err }
	// For loop entries
	for i := 0; i < numOfEntries ; i++ {
		entryBts, err := wr.NextWindow(12); if err != nil { return nil, err}
		entry, err := DecodeEntry(fileEndian, entryBts); if err != nil { return nil, err}
		df := entry.GetDataFormat()
		// uint16, *uint32, *int64, *int, *int32, *float32, *float64
		tag, ok := Tags[entry.Tag] 
		if ok {
			var d interface{}
			switch df.ID {
			case UByte:
				var v uint8
				err = entry.DecodeData(fd, &v)
				d = v
			case AsciiString: 
				var v string
				err = entry.DecodeData(fd, &v)
				d = v
			case UShort:
				var v uint16
				err = entry.DecodeData(fd, &v)
				d = v
			case ULong:
				var v uint32
				err = entry.DecodeData(fd, &v)
				d = v
			case Urational:
				var v uint64
				err = entry.DecodeData(fd, &v)
				d = v
			case Byte:
				var v int8
				err = entry.DecodeData(fd, &v)
				d = v
			case Undefined:
				// leave it interface
			case Short:
				var v int16
				err = entry.DecodeData(fd, &v)
				d = v
			case Long:
				var v int32
				err = entry.DecodeData(fd, &v)
				d = v
			case Rational:
				var v int64
				err = entry.DecodeData(fd, &v)
				d = v
			case Float:
				var v float32
				err = entry.DecodeData(fd, &v)
				d = v
			case DoubleFloat:
				var v float64
				err = entry.DecodeData(fd, &v)
				d = v
			}
			if err != nil { return nil, err}
			res[tag.Name] = d
		}


	}
	
	// lastDirBytes, err := wr.NextWindow(4); if err != nil { return nil, err }
	// fmt.Printf("lastDirbytes %x\n", lastDirBytes) // should be 4 zeros	
	return res, nil
}



