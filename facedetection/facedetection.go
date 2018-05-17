// Taken mostly from https://gocv.io/writing-code/face-detect/
package main

import (
	"flag"
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
)

func startImage(imageFile string) gocv.Mat {
	mat := gocv.IMRead(imageFile, gocv.IMReadAnyColor)
	fmt.Printf("[DEBUG] startImage - loaded image on mat %v\n", mat)
	return mat
}

func startWebcam(deviceId int) (*gocv.VideoCapture, gocv.Mat) {
	webcam, err := gocv.VideoCaptureDevice(deviceId)

	if err != nil {
		panic(err)
	}

	mat := gocv.NewMat()
	fmt.Printf("[DEBUG] startWebcam - started device %v and mat %v\n", webcam, mat)
	return webcam, mat
}

func drawRectangles(mat gocv.Mat, classifierFile string) {
	blue := color.RGBA{0, 0, 255, 255}
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load(classifierFile) {
		panic("Error opening classifier file")
	}

	fmt.Printf("[DEBUG] drawRectangles - using mat %v\n", mat)
	for {
		rectangles := classifier.DetectMultiScale(mat)
		fmt.Printf("[DEBUG:1] drawRectangles - currently rendering mat %v and rectangles classifier %v\n", mat, rectangles)

		for _, r := range rectangles {
			fmt.Printf("[DEBUG:1] drawRectangles - currently rendering mat %v and rectangle %v\n", mat, r)
			gocv.Rectangle(&mat, r, blue, 3)
		}
	}
	defer classifier.Close()
}

func renderMatWindow(mat gocv.Mat, webcam *gocv.VideoCapture) {
	window := gocv.NewWindow("Face detection")
	fmt.Printf("[DEBUG] renderMatWindow - starting window with mat %v and webcam %v\n", mat, webcam)
	for {
		if (*webcam != gocv.VideoCapture{}) {
			fmt.Printf("[DEBUG] renderMatWindow - webcam exists and is located at %v\n", webcam)
			webcam.Read(&mat)
		}
		window.IMShow(mat)
		fmt.Printf("[DEBUG:1] renderMatWindow - currently rendering mat %v\n", mat)
		window.WaitKey(1)
	}
	defer window.Close()
}

func main() {
	deviceIdPtr := flag.Int(
		"deviceId", 0, "The device id for your camera",
	)
	classifierPtr := flag.String(
		"classifier", "", "The classifier .xml file",
	)
	imagePtr := flag.String(
		"image", "", "Uses image file instead of webcam feed",
	)
	flag.Parse()

	if *classifierPtr == "" {
		panic("Please provide a classifier file")
	}

	webcam, mat := &gocv.VideoCapture{}, gocv.Mat{}

	// TODO - i believe I should share the mat address, maybe
	if *imagePtr != "" {
		fmt.Println("[INFO] main - using picture")
		mat = startImage(*imagePtr)
	} else {
		fmt.Println("[INFO] main - using webcam")
		webcam, mat = startWebcam(*deviceIdPtr)
	}

	go drawRectangles(mat, *classifierPtr)
	renderMatWindow(mat, webcam)

	// TODO - learn about panic and defer

	defer webcam.Close()
	defer mat.Close()
}
