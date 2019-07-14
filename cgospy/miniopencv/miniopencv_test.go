package miniopencv_test

import (
	"github.com/mdrohmann/playgo/cgospy/miniopencv"
	"image/color"
	"testing"
)

func TestCapture(t *testing.T) {
	mini := miniopencv.New(0)
	m, err := mini.CaptureCvMat()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer m.Close()

	t.Run("Size() should return the 2 dimensional image size", func(t *testing.T) {
		size := m.Size()
		if len(size) != 2 {
			t.Fatalf("Expected image to be two dimensional.")
		}
	})

	t.Run("At() should return color values", func(t *testing.T) {
		s := m.Bounds()
		anyUnequalZero := false
	OuterLoop:
		for x := s.Min.X; x < s.Max.X; x++ {
			for y := s.Min.Y; y < s.Max.Y; y++ {
				if m.At(x, y) != (color.Gray{0}) {
					anyUnequalZero = true
					break OuterLoop
				}
			}
		}
		if !anyUnequalZero {
			t.Errorf("Expected to have at least one value unequal zero")
		}
	})

	/*
		t.Run("imshow shows image", func(t *testing.T) {
			m.ImShow()

		}) */
}
