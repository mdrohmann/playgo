package cgospy

// #cgo pkg-config: opencv
// #include "opencv-wrapper.hpp"
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

// CgoSpy manages a capture session
type CgoSpy struct {
	device int
}

// New returns a new MiniOpencv
func New(device int) *CgoSpy {
	return &CgoSpy{device: device}
}

// CvMat wraps an OpenCV matrix structure
// I also implements the image Image interface
type CvMat struct {
	ptr C.CvMatrix
}

// ColorModel currently always returns GrayModel
func (m *CvMat) ColorModel() color.Model {
	return color.GrayModel
}

// Size returns the size of a CV array
func (m *CvMat) Size() []int {
	sizeLen := int32(0)
	intArray := C.cvMatrixSize(m.ptr, (*_Ctype_int)(&sizeLen))
	slice := (*[1 << 28]C.CInt)(unsafe.Pointer(intArray))[:sizeLen:sizeLen]
	res := make([]int, sizeLen)
	for i := int32(0); i < sizeLen; i++ {
		res[i] = int(slice[i])
	}
	// C.free(unsafe.Pointer(intArray))
	return res
}

// Bounds returns the bounds of the image
func (m *CvMat) Bounds() image.Rectangle {

	s := m.Size()

	// check if matrix has two dimensions
	if len(s) != 2 {
		// We have to panic here in order to fulfill the image.Image()
		// interface.  An error check should be made earlier if you are
		// really unsure if the matrix has more than two dimensions.
		panic(fmt.Sprintf("Expected matrix dimension to be 2, was: %d", len(s)))
	}

	return image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{s[1], s[0]}}
}

// At returns the color value at a specific position
func (m *CvMat) At(x, y int) color.Color {
	colorAtPoint := C.cvMatAt(m.ptr, C.int(x), C.int(y))
	return color.Gray{uint8(colorAtPoint)}
}

// CaptureCvMat captures a matrix representation of an image.
func (spy *CgoSpy) CaptureCvMat() (m *CvMat, err error) {
	m = &CvMat{ptr: C.newCvMat()}
	if r := C.captureImage(C.int(spy.device), m.ptr); r != 0 {
		err = fmt.Errorf("Error opening the default video capture")
	}
	return m, err
}

// Close frees the memory of a CVMat structure
func (m *CvMat) Close() error {
	C.freeCvMat(m.ptr)
	return nil
}

// ImShow shows the CvMat in a window
func (m *CvMat) ImShow() {
	C.imShow(m.ptr)
}
