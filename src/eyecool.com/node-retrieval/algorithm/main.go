package algorithm

/*
#include <stdio.h>
#include <stdlib.h>
#include "ChlFaceSdk.h"
#cgo CFLAGS: -I .
#cgo LDFLAGS: -L/usr/local/lib -lChlFaceSdk
*/
import "C"
import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"
	"unsafe"
)

func xmain() {
	fmt.Printf("main\n")

	var nMaxChannelNum C.int = 1
	var pLibPath *C.char = C.CString("E:\\temp\\SDK_CHL")
	var pTmpDir *C.char = C.CString("E:\\temp\\SDK_CHL")
	defer C.free(unsafe.Pointer(pLibPath))
	defer C.free(unsafe.Pointer(pTmpDir))
	ret := C.ChlFaceSdkInit(nMaxChannelNum, pLibPath, pTmpDir)
	fmt.Printf("init ret :%d\n ", ret)
	var version C.int = C.ChlFaceSdkVer()
	fmt.Printf("version %d\n : ", version)
	fs := C.ChlFaceSdkFeatureSize()
	fmt.Printf("max feature size: %d\n : ", fs)
	var imagef *C.char = C.CString("/opt/local/test/111.jpg")
	var width, height C.int = 0, 0
	ret = C.ReadImageFile(imagef, nil, 0, &width, &height, 24)
	fmt.Printf("ReadImageFile ret: %d ,width : %d ,height :  %d \n : ", ret, width, height)
	nSize := width * height * 3
	var bt []byte = make([]byte, nSize)
	pRgb24 := (*C.BYTE)(unsafe.Pointer(&bt[0]))
	ret = C.ReadImageFile(imagef, pRgb24, nSize, &width, &height, 24)
	fmt.Printf("ReadImageFile ret: %d ,width : %d ,height :  %d  \n: ", ret, width, height)
	fmt.Printf(" len : %d \n ", len(bt))
	/*for  i :=0;i< len(bt) ;i++ {
		s:=bt[i]

		if i>1000 && i<10000{
			continue;
		}
		fmt.Printf(" %d ",s)
	}*/
	var pFaceResult1 *C.struct_FACE_DETECT_RESULT = new(C.struct_FACE_DETECT_RESULT)

	ret = C.ChlFaceSdkDetectFace(C.int(0), (*C.BYTE)(pRgb24), width, height, pFaceResult1)
	var pFaceResult = *pFaceResult1
	fmt.Printf("ChlFaceSdkDetectFace ret : %d  \n  ", ret)
	fmt.Println(" : pFaceResult :", pFaceResult)
	var rect C.struct_tagRECT = pFaceResult.rcFace
	fmt.Println("rect ", rect)
	fmt.Println("rect all :  ", rect.left, rect.top, rect.right, rect.bottom)

	var ptLeftEye, ptRightEye, ptMouth, ptNose C.struct_tagPOINT = pFaceResult.ptLeftEye, pFaceResult.ptRightEye, pFaceResult.ptMouth, pFaceResult.ptNose
	fmt.Println("point :  ", ptLeftEye, ptRightEye, ptMouth, ptNose)

	t := time.Now().UnixNano()
	//open a image file
	file, err := os.Open("/opt/local/test/111.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Println(file)
	m, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("===image.Decode : ", err)
	}
	rgbImg := m.(*image.YCbCr)
	subImg := rgbImg.SubImage(image.Rect(int(rect.left), int(rect.top), int(rect.right), int(rect.bottom))).(*image.YCbCr)
	fmt.Println(subImg.Bounds())
	f, _ := os.Create("/opt/local/test/111_cut.jpg") //创建文件
	defer f.Close()                                  //关闭文件
	jpeg.Encode(f, subImg, nil)                      //写入文件
	t2 := time.Now().UnixNano()
	fmt.Println("cut image spend time  :  ", (t2 - t))
	//===================================================================

	var FaceData [512]C.char = pFaceResult.FaceData
	fmt.Println("FaceData :  ", FaceData)
	// C.free(unsafe.Pointer(pFaceResult1))
	//fmt.Println(" : pFaceResult :",pFaceResult)
	/*var detectResult []FACE_DETECT_RESULT
	dfSize := unsafe.Sizeof(*pFaceResult)
	for (*pFaceResult).rcFace != nil {
		rr:=FACE_DETECT_RESULT{
			c
		}
	}*/

	var feature []byte = make([]byte, fs)
	features := (*C.BYTE)(unsafe.Pointer(&feature[0]))
	//var feats []byte=make([]byte,nSize)
	//ofeat := (*C.uchar)(unsafe.Pointer(&feats[0]))
	ret = C.ChlFaceSdkFeatureGet(C.int(0), (*C.BYTE)(pRgb24), width, height, pFaceResult1, features)
	fmt.Println("feature : ", feature)

	var score C.BYTE = C.ChlFaceSdkFeatureCompare(C.int(0), features, features)

	fmt.Println("MATCH score  : ", int(score))

	var maxFace C.int = 1
	var faceResultBuf []C.struct_FACE_DETECT_RESULT = make([]C.struct_FACE_DETECT_RESULT, maxFace)
	pFaceResultBuf := (*C.struct_FACE_DETECT_RESULT)(unsafe.Pointer(&faceResultBuf[0]))
	var featuresBuf []byte = make([]byte, fs*maxFace)
	featuresBufs := (*C.BYTE)(unsafe.Pointer(&featuresBuf[0]))
	ret = C.ChlFaceSdkFaceFeature(C.int(0), (*C.BYTE)(pRgb24), width, height, maxFace, pFaceResultBuf, featuresBufs)
	fmt.Println("ChlFaceSdkFaceFeature face cnt  : ", ret)
	//fmt.Println("ChlFaceSdkFaceFeature feature bytes1  : ",featuresBuf[0:2600])
	//fmt.Println("ChlFaceSdkFaceFeature feature bytes2  : ",featuresBuf[2600:5200])
	//fmt.Println("ChlFaceSdkFaceFeature feature bytes3  : ",featuresBuf[5200:5205])

	fmt.Println("struct_FACE_DETECT_RESULT   : ", faceResultBuf)
	//===================================================================================================================

	var listHandle C.HANDLE = C.ChlFaceSdkListCreate(C.int(10))

	fmt.Println("ChlFaceSdkListCreate handle : ", listHandle)
	var pos C.int = -1
	fmt.Println("pos: ", pos)
	for i := 0; i < 10; i++ {
		pos = -1
		total := C.ChlFaceSdkListInsert(listHandle, &pos, features, C.int(1))
		fmt.Println("ChlFaceSdkListInsert tatol : ", i, total, pos)
	}
	var SCORE []byte = make([]byte, 10)
	SCORES := (*C.BYTE)(unsafe.Pointer(&SCORE[0]))
	total := C.ChlFaceSdkListCompare(C.int(0), listHandle, features, C.int(0), C.int(0), SCORES)
	fmt.Println("ChlFaceSdkListCompare tatol : ", total)
	for i := 0; i < len(SCORE); i++ {
		s := SCORE[i]
		fmt.Printf(" idx[%d] score :  %d \n", i, s)
	}
	total = C.ChlFaceSdkListRemove(listHandle, C.int(0), C.int(1))
	fmt.Println("ChlFaceSdkListRemove tatol : ", total)
	C.ChlFaceSdkListClearAll(listHandle)
	C.ChlFaceSdkListDestroy(listHandle)

}

/*
type  FACE_DETECT_RESULT struct{
 rcFace RECT
 ptLeftEye, ptRightEye, ptMouth, ptNose POINT
 nAngleYaw, nAnglePitch, nAngleRoll, nQuality int
 FaceData [512]byte;
}

type POINT struct
{
 x, y int
}


type  RECT struct
{
 left, top, right, bottom int
}*/
