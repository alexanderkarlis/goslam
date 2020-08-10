package main

import (
	"fmt"
	"os"
	"time"

	"gocv.io/x/gocv"
)

var (
	window = gocv.NewWindow("goslam")
	orb    = gocv.NewORB()
)

func main() {
	defer window.Close()

	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tshowimage [imgfile]")
		return
	}

	filename := os.Args[1]

	frames, err := gocv.VideoCaptureFile(filename)
	defer frames.Close()
	checkErr(err)

	ticker := time.NewTicker(100 * time.Millisecond)
	done := make(chan bool)

	for {
		select {
		case <-done:
			fmt.Println("done")
			return
		case <-ticker.C:
			processFrame(frames)
		}
	}
}

func getMatType(img *gocv.Mat) gocv.MatType {
	fmt.Printf("Mat type: %+v\n", img.Type())
	return img.Type()
}

func extract(img gocv.Mat) {
	// var corners *gocv.Mat
	corners := gocv.NewMat()
	grey := gocv.NewMat()
	gocv.CvtColor(img, &grey, gocv.ColorBGRToGray)
	fmt.Printf("%v\n", grey.Rows())

	gocv.GoodFeaturesToTrack(grey, &corners, 3000, 0.01, 3)
	fmt.Println(corners)
}

func processFrame(frames *gocv.VideoCapture) {
	img := gocv.NewMat()
	ok := frames.Read(&img)
	if !ok || img.Empty() {
		return
	}
	extract(img)

	// var greyImg gocv.Mat
	// gocv.CvtColor(img, &greyImg, gocv.ColorBGRToGray)
	// gocv.GoodFeaturesToTrack(img, &greyImg, 3000, 0.01, 3)

	fmt.Printf("mean: %+v\n", img.Mean())
	window.IMShow(img)
	if window.WaitKey(1) >= 0 {
		print("exiting display")
		os.Exit(0)
	}
}

func print(value string) {
	fmt.Printf("%s\n", value)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
