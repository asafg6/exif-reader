package exif

import (
	"encoding/binary"
	"fmt"
)

type CompanyID string

const Intel CompanyID = "II"
const Motorolla CompanyID = "MM"

const EXIF = "Exif"

var BigEndianOrder = []byte{0x00, 0x2a}
var LittleEndianOrder = []byte{0x2a, 0x00}

type ExifHeader struct {
	Size uint16                   // 2 b
	ExifString string             // 4 b 6
	TwoZeros []byte				  // 2 b 8
	CompanyID CompanyID           // 2 b 10
	ByteOrder binary.ByteOrder    // 2 b 12
	Offset  uint32                // 4 b 16
}

// DecodeHeader accepts 16 bytes []byte and returns ExifHeader
func DecodeHeader(bts []byte) (ExifHeader, error) {
	header := ExifHeader{}
	if len(bts) != 16 {
		return header, fmt.Errorf("expecting header length of 18 but got %d", len(bts))
	}
	sizeBts := bts[:2]
	exifStrBts := bts[2:6]
	twoZeros := bts[6:8]
	companyIDBts := bts[8:10]
	bomBts := bts[10:12]
	offsetBts := bts[12:16]
	
	var decodesResults []error

	// Start with validations
	err := validateExifString(exifStrBts); if err != nil { goto onError }
	err = validateTwoZeros(twoZeros); if err != nil {  goto onError }
	// Continue on to determine byte order
	header.CompanyID =  CompanyID(forceASCII(companyIDBts))
	header.ByteOrder, err = determineByteOrder(header.CompanyID, bomBts); if err != nil {  goto onError }

	// Now we can decode all those numbers
	decodesResults = []error{
		decode(header.ByteOrder, sizeBts, &header.Size), 
		decode(header.ByteOrder, offsetBts, &header.Offset), 
	}
	for _, r := range decodesResults {
		if r != nil { 
			err = r
			goto onError 
		}
	}
	onError:
		if err != nil{
			return header, err
		}

	return header, nil

}


func determineByteOrder(cID CompanyID, bom []byte) (binary.ByteOrder, error) {

	if cID == Motorolla && compBytes(bom, BigEndianOrder){
		return binary.BigEndian, nil
	} else if cID == Intel && compBytes(bom, LittleEndianOrder) {
		return binary.LittleEndian, nil
	} 
	return nil, fmt.Errorf("Could not figure out file byte order")
}

func validateExifString(bts []byte) error {
	exifStr := forceASCII(bts)
	if exifStr != EXIF {
		return fmt.Errorf("Unexpected Exif identifier %s", exifStr)
	}
	return nil
}

func validateTwoZeros(bts []byte) error {
	if !compBytes(bts, []byte{0, 0}) {
		return fmt.Errorf("Expected 0x0000 but got %x", bts)
	}
	return nil
}
