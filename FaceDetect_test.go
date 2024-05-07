package main

import (
	"gocv.io/x/gocv"
	"testing"
)

func TestDetectFaces(t *testing.T) {
	testData := []struct {
		imagePath string
		expected  int // Expected number of detected faces
	}{
		{"/Users/ankitapanchal/Desktop/Ankita.jpg", 1}, // Test case 1: Normal case
		{"", 0}, // Test case 2: Empty image case
		{"/Users/ankitapanchal/Desktop/ankita.jpg", -1}, // Test case 3: Error case
	}
	for _, data := range testData {
		var frame gocv.Mat
		var err error
		if data.imagePath != "" {
			frame, err = gocv.IMRead(data.imagePath, gocv.IMReadColor)
			if err != nil {
				t.Errorf("Failed to read image file: %v", err)
				continue
			}
		} else {
			// Create an empty frame for testing empty image case
			frame = gocv.NewMat()
		}

		result, err := detectFaces(APIKey("22cbd6fcdd2c43d98351b3b40285997f"), frame)
		if err != nil {
			// Error case
			if data.expected != -1 {
				t.Errorf("Expected no error, but got error: %v", err)
			}
		} else {
			// Normal case or Empty image case
			if len(result) != data.expected {
				t.Errorf("Expected %d detected faces, got %d", data.expected, len(result))
			}
		}
	}
}
