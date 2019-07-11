package cgospy

// #cgo pkg-config: opencv
// #include "opencv-wrapper.hpp"
import "C"
import (
	"fmt"
	// "unsafe"
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
type CvMat struct {
	ptr C.CvMatrix
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
