// cgospy command line interface
// stores output of the camera in a png file
// edge filter is applied to output

package main

import (
	"flag"
	"fmt"
	"github.com/mdrohmann/playgo/cgospy/miniopencv"
	"image/png"
	"os"
)

func main() {

	outputFile := flag.String("output", "test.png", "the PNG file name to write the file to.")

	flag.Parse()

	mini := miniopencv.New(0)
	m, err := mini.CaptureCvMat()
	if err != nil {
		fmt.Printf("Unexpected error: %s", err)
		os.Exit(-1)
	}
	defer m.Close()

	size := m.Size()
	if len(size) != 2 {
		fmt.Printf("Expected image to be two dimensional.")
		os.Exit(-1)
	}

	f, err := os.Create(*outputFile)
	defer f.Close()

	if err != nil {
		fmt.Printf("Expected to be able to create the output png file, but received error: %s", err)
		os.Exit(-2)
	}
	if err = png.Encode(f, m); err != nil {
		fmt.Printf("Expected to be able to encode CvMatrix as png, but got error: %s", err)
		os.Exit(-2)
	}

	m.ImShow()
}
