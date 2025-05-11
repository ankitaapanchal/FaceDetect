package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gocv.io/x/gocv"
	"golang.org/x/image/colornames"
	"image"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// APIKey represents an API key used for authentication.
type APIKey string

// ImagePath represents the path to an image file.
type ImagePath string

// FaceResponse represents the response structure from the Luxand Cloud API.
type FaceResponse struct {
	Faces []struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"faces"`
}

// String converts an APIKey to a string.
func (apiKey APIKey) String() string { return string(apiKey) }

// String converts an ImagePath to a string.
func (imagePath ImagePath) String() string { return string(imagePath) }

func main() {
	// Open a webcam for video capture
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatal(err)
	}
	defer webcam.Close()

	// Create a window to display the video feed
	window := gocv.NewWindow("Looking for a face....")
	defer window.Close()

	// Define the Luxand Cloud API key
	apiKey := APIKey("Your_API_Key")

	// Start face recognition
	findFaces(webcam, window, apiKey)
}

// findAndLabelFaces detects faces in the video stream from the webcam and labels them.
func findFaces(camera *gocv.VideoCapture, window *gocv.Window, apiKey APIKey) {
	for {
		// Read a frame from the camera
		frame := gocv.NewMat()

		camera.Read(&frame)

		if frame.Empty() {
			log.Println("cannot read from camera!")
			continue
		}

		// Detect faces in the current frame
		faces, err := detectFaces(apiKey, frame)
		if err != nil {
			log.Println(err)
			continue
		}

		// Draw rectangles around detected faces and label them
		for _, rectangle := range faces {
			gocv.Rectangle(&frame, rectangle, colornames.Blue, 3)
			textsize := gocv.GetTextSize("I KNOW YOU!", gocv.FontHersheyDuplex, 1.5, 2)
			textXloc := rectangle.Min.X + (rectangle.Max.X-rectangle.Min.X)/2 - textsize.X/2
			textYLoc := rectangle.Min.Y + (rectangle.Max.Y-rectangle.Min.Y)/2 - textsize.Y/2
			textLoc := image.Pt(textXloc, textYLoc)
			gocv.PutText(&frame, "I KNOW YOU!", textLoc, gocv.FontHersheyDuplex, 1.5,
				colornames.Red, 2)
		}

		// Show the processed image with detected faces in the window
		window.IMShow(frame)
		if window.WaitKey(1) >= 0 {
			break
		}

		// Release the frame memory
		frame.Close()
	}
}

// detectFaces detects faces in the given image using the Luxand Cloud API.
func detectFaces(apiKey APIKey, img gocv.Mat) ([]image.Rectangle, error) {
	// Write the image to a temporary file
	tempFile, err := ioutil.TempFile("", "frame_*.jpg")
	if err != nil {
		return nil, err
	}
	//defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	//Path to image
	imagePath := "/path/image.jpg"

	//Load image from path
	img = gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, errors.New("Failed to read image file")
	}

	//Write the loaded image to the temporary file
	if ok := gocv.IMWrite(tempFile.Name(), img); !ok {
		return nil, errors.New("failed to write image to file")
	}

	// Create a multipart request for the Luxand Cloud API
	url := "https://api.luxand.cloud/photo/search/v2"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	part, err := writer.CreateFormFile("photo", filepath.Base(imagePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("collections", "")
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	// Send the request to the Luxand Cloud API
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("token", string(apiKey))
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response from api: ", string(body))
	// Decode the response from the Luxand Cloud API
	var faceResponse FaceResponse
	err = json.NewDecoder(resp.Body).Decode(&faceResponse)
	if err != nil {
		return nil, err
	}

	// Convert the detected faces to image rectangles
	var rectangles []image.Rectangle
	for _, face := range faceResponse.Faces {
		rect := image.Rect(face.X, face.Y, face.X+face.Width, face.Y+face.Height)
		rectangles = append(rectangles, rect)
	}
	return rectangles, nil
}


