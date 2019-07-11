package cgospy_test

import (
	"github.com/mdrohmann/playgo/cgospy"
	"testing"
)

func TestCapture(t *testing.T) {
	mini := cgospy.New(0)
	m, err := mini.CaptureCvMat()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer m.Close()

	t.Run("imshow shows image", func(t *testing.T) {
		m.ImShow()

	})
}
