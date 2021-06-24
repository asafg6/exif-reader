package exif

type Tag struct {
	ID  string
	Name string
}

// names are in big endian
var Tags = map[string]Tag{
	"0x0100": {ID: "0x0100", Name: "ImageWidth"}, // ImageWidth has double tag
	"0x0001": {ID: "0x0001", Name: "ImageWidth"},
	"0x0101": {ID: "0x0101", Name: "ImageLength"},
	"0x9003": {ID: "0x9003", Name:"DateTimeOriginal"},
	"0x0132": {ID: "0x0132", Name:"DateTime"},
	"0x010f": {ID: "0x010f", Name: "Make"},
	"0x0110": {ID: "0x0110", Name: "Model"},
	"0x0112": {ID: "0x0112", Name: "Orientation"},
	"0x0131": {ID: "0x0131", Name: "Software"},
}


