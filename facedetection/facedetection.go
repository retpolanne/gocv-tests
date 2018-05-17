// Taken mostly from https://gocv.io/writing-code/face-detect/
package main

import (
	"flag"
	"image/color"

	"gocv.io/x/gocv"
)

func startImage(imageFile string) gocv.Mat {
	mat := gocv.IMRead(imageFile, gocv.IMReadAnyColor)
	return mat
}

func startWebcam(deviceId int) (*gocv.VideoCapture, gocv.Mat) {
	webcam, err := gocv.VideoCaptureDevice(deviceId)

	if err != nil {
		panic(err)
	}

	mat := gocv.NewMat()
	return webcam, mat
}

func drawRectangles(mat gocv.Mat, classifierFile string) {
	blue := color.RGBA{0, 0, 255, 255}
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load(classifierFile) {
		panic("Error opening classifier file")
	}

	for {
		rectangles := classifier.DetectMultiScale(mat)

		for _, r := range rectangles {
			gocv.Rectangle(&mat, r, blue, 3)
		}
	}
	defer classifier.Close()
}

func renderMatWindow(mat gocv.Mat, webcam *gocv.VideoCapture) {
	window := gocv.NewWindow("Face detection")
	for {
		if webcam != nil {
			webcam.Read(&mat)
		}
		window.IMShow(mat)
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
		mat = startImage(*imagePtr)
	} else {
		webcam, mat = startWebcam(*deviceIdPtr)
	}

	go renderMatWindow(mat, webcam)
	drawRectangles(mat, *classifierPtr)

	// TODO - learn about panic and defer

	defer webcam.Close()
	defer mat.Close()
}
