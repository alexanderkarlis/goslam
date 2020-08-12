package main

import (
	"fmt"
	"image/color"
	"os"
	"reflect"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

var (
	window = gocv.NewWindow("goslam")
	orb    = gocv.NewORB()
	bf     = gocv.NewBFMatcherWithParams(gocv.NormHamming, false)
	last   = lastValues{}
	feats  = gocv.NewMat()
)

type lastValues struct {
	kps []gocv.KeyPoint
	des gocv.Mat
}

type keypoints struct {
	kp1 gocv.KeyPoint
	kp2 gocv.KeyPoint
}

func main() {
	defer window.Close()

	if len(os.Args) < 3 {
		fmt.Println("How to run:\n\tgoslam [imgfile]")
		return
	}

	filename := os.Args[1]
	framerateMs, err := strconv.Atoi(os.Args[2])
	checkErr(err)
	frames, err := gocv.VideoCaptureFile(filename)
	defer frames.Close()
	checkErr(err)

	ticker := time.NewTicker(time.Duration(framerateMs) * time.Millisecond)
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

func (x lastValues) IsStructureEmpty() bool {
	return reflect.DeepEqual(x, lastValues{})
}

func getMatType(img *gocv.Mat) gocv.MatType {
	fmt.Printf("Mat type: %+v\n", img.Type())
	return img.Type()
}

func extract(img *gocv.Mat) (ret []keypoints) {
	// gocv.Resize(*img, img, image.Point{}, 400, 600, 0)
	// detection
	fmt.Printf("SIZE: %+v\n", img.Size())
	cv8uci := gocv.NewMatWithSize(2160, 3840, gocv.MatTypeCV8UC1)
	gocv.GoodFeaturesToTrack(cv8uci, &feats, 3000, 0.01, 3.0)
	fmt.Printf("%+v\n", feats.Mean())

	// extraction
	corners := gocv.NewMat()
	kps, _ := orb.DetectAndCompute(*img, corners)
	gocv.DrawKeyPoints(*img, kps, img, color.RGBA{0, 0, 255, 0}, 0)
	window.IMShow(*img)

	// matching
	// if !last.IsStructureEmpty() {
	// 	matches := bf.KnnMatch(des, last.des, 2)
	// 	for i, m := range matches {
	// 		for _, n := range m {
	// 			if m[i].Distance < 0.75*n.Distance {
	// 				kp1 := kps[m[i].QueryIdx]
	// 				kp2 := last.kps[m[i].TrainIdx]
	// 				ret = append(ret, keypoints{
	// 					kp1,
	// 					kp2,
	// 				})
	// 			}

	// 		}
	// 	}
	// }
	// last = lastValues{kps: kps, des: des}
	return
}

func processFrame(frames *gocv.VideoCapture) {
	img := gocv.NewMat()
	ok := frames.Read(&img)
	if !ok || img.Empty() {
		return
	}
	extract(&img)
	// fmt.Printf("matches: %+v\n", matches)
	// fmt.Printf("last: %+v\n", last)
	// window.IMShow(img)
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
