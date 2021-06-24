package main

import (
	"fmt"
	"./exif"
)






func main () {
	
	fname := "<path to file>"
	props, err := exif.GetExifData(fname)	
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	for k, v := range props {
		fmt.Println(k,":   ", v)
	}


}


