package algorithm

/*
#include <stdio.h>
#include <stdlib.h>
#include "ChlFaceSdk.h"
#include "typedef.h"
#cgo CFLAGS: -I .  -ltypedef
#cgo LDFLAGS: -L/usr/local/lib -lChlFaceSdk
*/
import "C"
import (
	"fmt"
	"unsafe"
	"sync/atomic"
)
const (
	DefaultChannelNo = -1
)
type FACE_DETECT_RESULTX C.struct_FACE_DETECT_RESULT

func (data *FACE_DETECT_RESULTX) SetRECT(left, top, right, bottom int) {
	fmt.Println(" SetRECT ", left, top, right, bottom)
	var rcFace C.struct_tagRECT
	rcFace.left = C.int(left)
	rcFace.top = C.int(top)
	rcFace.right = C.int(right)
	rcFace.bottom = C.int(bottom)
	data.rcFace = rcFace
	fmt.Println(" rcFace ", data.rcFace)
}
func (data *FACE_DETECT_RESULTX) GetRectFace() C.struct_tagRECT {
	return data.rcFace
}
func (data *FACE_DETECT_RESULTX) SetRectFace(rcFace C.struct_tagRECT) {
	data.rcFace = rcFace
}

func (data *FACE_DETECT_RESULTX) GetRECT() RECTX {
	var rcFace C.struct_tagRECT = C.struct_tagRECT(data.rcFace)
	return RECTX{Left: int(rcFace.left), Top: int(rcFace.top), Right: int(rcFace.right), Bottom: int(rcFace.bottom)}
}
func (data *FACE_DETECT_RESULTX) SetLeftEye(p POINTX) {
	data.ptLeftEye = C.struct_tagPOINT{C.int(p.X), C.int(p.Y)}
}

func (data *FACE_DETECT_RESULTX) GetLeftEye() C.struct_tagPOINT {
	var p C.struct_tagPOINT = C.struct_tagPOINT(data.ptLeftEye)
	return p
}
func (data *FACE_DETECT_RESULTX) GetLeftEyeX() POINTX {
	var p C.struct_tagPOINT = data.ptLeftEye
	return POINTX{int(p.x), int(p.y)}
}

func (data *FACE_DETECT_RESULTX) SetRightEye(p POINTX) {
	data.ptRightEye = C.struct_tagPOINT{C.int(p.X), C.int(p.Y)}
}

func (data *FACE_DETECT_RESULTX) GetRightEyeX() POINTX {
	var p C.struct_tagPOINT = data.ptRightEye
	return POINTX{int(p.x), int(p.y)}
}
func (data *FACE_DETECT_RESULTX) GetRightEye() C.struct_tagPOINT {
	var p C.struct_tagPOINT = C.struct_tagPOINT(data.ptRightEye)
	return p
}

func (data *FACE_DETECT_RESULTX) SetMouth(p POINTX) {
	data.ptMouth = C.struct_tagPOINT{C.int(p.X), C.int(p.Y)}
}
func (data *FACE_DETECT_RESULTX) GetMouth() C.struct_tagPOINT {
	var p C.struct_tagPOINT = C.struct_tagPOINT(data.ptMouth)
	return p
}

func (data *FACE_DETECT_RESULTX) GetMouthX() POINTX {
	var p C.struct_tagPOINT = C.struct_tagPOINT(data.ptMouth)
	return  POINTX{int(p.x), int(p.y)}
}

func (data *FACE_DETECT_RESULTX) SetNose(p POINTX) {
	data.ptNose = C.struct_tagPOINT{C.int(p.X), C.int(p.Y)}
}
func (data *FACE_DETECT_RESULTX) GetNose() C.struct_tagPOINT {
	var p C.struct_tagPOINT = C.struct_tagPOINT(data.ptNose)
	return p
}
func (data *FACE_DETECT_RESULTX) GetNoseX() POINTX {
	var p C.struct_tagPOINT = C.struct_tagPOINT(data.ptNose)
	return POINTX{int(p.x), int(p.y)}
}


func (data *FACE_DETECT_RESULTX) SetAttrs(nAngleYaw, nAnglePitch, nAngleRoll, nQuality int) {
	data.nAngleYaw, data.nAnglePitch, data.nAngleRoll, data.nQuality = C.int(nAngleYaw), C.int(nAnglePitch), C.int(nAngleRoll), C.int(nQuality)
}
func (data *FACE_DETECT_RESULTX) GetAttrs() (C.int, C.int, C.int, C.int) {
	return data.nAngleYaw, data.nAnglePitch, data.nAngleRoll, data.nQuality
}
func (data *FACE_DETECT_RESULTX) GetIntAttrs() (int, int, int, int) {
	return int(data.nAngleYaw), int(data.nAnglePitch), int(data.nAngleRoll), int(data.nQuality)
}

func (data *FACE_DETECT_RESULTX) SetFaceData(FaceData [512]C.char) {
	data.FaceData = FaceData
}
func (data *FACE_DETECT_RESULTX) GetFaceData() [512]C.char {
	return data.FaceData
}

type POINTX struct {
	X, Y int
}

type RECTX struct {
	Left, Top, Right, Bottom int
}

type ChlFaceX struct {
	FaceMaxFeatureSize   int
	MaxChannelNum        int
	CallNum				 int32
	LibPath              string
	TmpDir               string
	Handles              map[string]C.HANDLE
	HandlesCachedSize    map[string]int
	HandlesMaxCachedSize map[string]int
}

func NewChlFaceX() *ChlFaceX {
	this := &ChlFaceX{
		FaceMaxFeatureSize:   2600,
		MaxChannelNum:        1,
		CallNum       :       0,
		LibPath:              "/usr/local/lib",
		TmpDir:               "/tmp",
		Handles:              make(map[string]C.HANDLE),
		HandlesCachedSize:    make(map[string]int),
		HandlesMaxCachedSize: make(map[string]int),
	}
	return this
}

func Free() {

}
func (this *ChlFaceX) ChlFaceSdkInit() int {
	var pLibPath *C.char = C.CString(this.LibPath)
	var pTmpDir *C.char = C.CString(this.TmpDir)
	defer C.free(unsafe.Pointer(pLibPath))
	defer C.free(unsafe.Pointer(pTmpDir))
	ret := C.ChlFaceSdkInit(C.int(this.MaxChannelNum), pLibPath, pTmpDir)
	if ret == 0 {
		this.ChlFaceSdkFeatureSize()
	}
	return int(ret)
}

func (this *ChlFaceX) ChlFaceSdkFeatureSize() int {
	fs := C.ChlFaceSdkFeatureSize()
	this.FaceMaxFeatureSize = int(fs)
	return this.FaceMaxFeatureSize
}

func (this *ChlFaceX) ReadImageFile(path string, width, height int) (int, int, int, []byte) {
	var cpath *C.char = C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	var cwidth, cheight C.int = C.int(width), C.int(height)
	var ret C.int
	if cwidth <= 0 || cheight <= 0 {
		ret = C.ReadImageFile(cpath, nil, 0, &cwidth, &cheight, 24)
		if ret < 0 {
			return int(ret), int(cwidth), int(cheight), nil
		}
	}
	nSize := cwidth * cheight * 3
	var rgb24 []byte = make([]byte, nSize)
	pRgb24 := (*C.BYTE)(unsafe.Pointer(&rgb24[0]))
	ret = C.ReadImageFile(cpath, pRgb24, nSize, &cwidth, &cheight, 24)
	return int(ret), int(cwidth), int(cheight), rgb24
}

func (this *ChlFaceX) ChlFaceSdkDetectFace(channelNo int, pRgb24 []byte, width, height int, isNoOnlyDetect bool) (hasFace int, pFaceResultx *FACE_DETECT_RESULTX) {
	this.incr()
	if channelNo >= this.MaxChannelNum || channelNo == DefaultChannelNo{
		channelNo = int(this.CallNum) % this.MaxChannelNum
	}
	var pFaceResult *C.struct_FACE_DETECT_RESULT = nil
	if isNoOnlyDetect {
		pFaceResultx = new(FACE_DETECT_RESULTX)
		pFaceResult = new(C.struct_FACE_DETECT_RESULT)
	}
	ret := C.ChlFaceSdkDetectFace(C.int(channelNo), (*C.BYTE)(unsafe.Pointer(&pRgb24[0])), C.int(width), C.int(height), pFaceResult)
	fmt.Println(" struct_FACE_DETECT_RESULT ", pFaceResult)
	if pFaceResultx != nil {
		makeGO_FACE_DETECT_RESULT(pFaceResult, pFaceResultx)
	}
	return int(ret), pFaceResultx
}
func (this *ChlFaceX) incr() {
	for o, n := this.CallNum, atomic.AddInt32(&this.CallNum, 1); o == n; {
		n = atomic.AddInt32(&this.CallNum, 1)
	}
}

func (this *ChlFaceX) ChlFaceSdkFeatureGet(channelNo int, pRgb24 []byte, width int, height int, pFaceResultX *FACE_DETECT_RESULTX) (success int, pFeature []byte) {
	if pFaceResultX == nil {
		return -9, nil
	}
	this.incr()
	if channelNo >= this.MaxChannelNum || channelNo == DefaultChannelNo{
		channelNo = int(this.CallNum) % this.MaxChannelNum
	}
	pFaceResult := makeC_FACE_DETECT_RESULT(pFaceResultX)
	pFeature = make([]byte, this.FaceMaxFeatureSize)
	ret := C.ChlFaceSdkFeatureGet(C.int(channelNo), (*C.BYTE)(unsafe.Pointer(&pRgb24[0])), C.int(width), C.int(height), pFaceResult, (*C.BYTE)(unsafe.Pointer(&pFeature[0])))
	return int(ret), pFeature
}

func (this *ChlFaceX) ChlFaceSdkFeatureCompare(channelNo int, pFeature1 []byte, pFeature2 []byte) int {
	if pFeature1 == nil || pFeature2 == nil {
		return -9
	}
	this.incr()
	if channelNo >= this.MaxChannelNum || channelNo == DefaultChannelNo{
		channelNo = int(this.CallNum) % this.MaxChannelNum
	}
	var score C.BYTE = C.ChlFaceSdkFeatureCompare(C.int(0), (*C.BYTE)(unsafe.Pointer(&pFeature1[0])), (*C.BYTE)(unsafe.Pointer(&pFeature2[0])))
	return int(score)
}

func (this *ChlFaceX) ChlFaceSdkDetectFaceExtractFeature(channelNo int, pRgb24 []byte, width int, height int, maxFace int, isDetectAndExtract int) (faceNum int, pFaceResultBufX []FACE_DETECT_RESULTX, pFeaturesBuf []byte) {
	this.incr()
	if channelNo >= this.MaxChannelNum || channelNo == DefaultChannelNo{
		channelNo = int(this.CallNum) % this.MaxChannelNum
	}
	var faceResultBuf []C.struct_FACE_DETECT_RESULT = nil
	var pFaceResultBuf *C.struct_FACE_DETECT_RESULT = nil
	var pFeaturesBufIn (*C.BYTE) = nil
	if isDetectAndExtract == 1 || isDetectAndExtract == 2 {
		faceResultBuf = make([]C.struct_FACE_DETECT_RESULT, maxFace)
		pFaceResultBuf = (*C.struct_FACE_DETECT_RESULT)(unsafe.Pointer(&faceResultBuf[0]))
	}
	if isDetectAndExtract == 2 {
		pFeaturesBuf = make([]byte, this.FaceMaxFeatureSize*maxFace)
		pFeaturesBufIn = (*C.BYTE)(unsafe.Pointer(&pFeaturesBuf[0]))
	}
	ret := C.ChlFaceSdkFaceFeature(C.int(channelNo), (*C.BYTE)(unsafe.Pointer(&pRgb24[0])), C.int(width), C.int(height), C.int(maxFace), pFaceResultBuf, pFeaturesBufIn)
	faceNum = int(ret)
	if isDetectAndExtract == 1 || isDetectAndExtract == 2 {
		pFaceResultBufX = make([]FACE_DETECT_RESULTX, faceNum)
		for i := 0; i < faceNum; i++ {
			var gfr *FACE_DETECT_RESULTX = new(FACE_DETECT_RESULTX)
			makeGO_FACE_DETECT_RESULT(&faceResultBuf[i], gfr)
			pFaceResultBufX[i] = *gfr
		}
	}
	if isDetectAndExtract == 2 {
		pFeaturesBuf = pFeaturesBuf[0 : this.FaceMaxFeatureSize*faceNum]
	}
	return faceNum, pFaceResultBufX, pFeaturesBuf
}

func (this *ChlFaceX) ChlFaceSdkListCreate(cacheName string, maxFeatureNum int) bool {
	if  this.Handles[cacheName] == nil {
		this.Handles[cacheName] = C.ChlFaceSdkListCreate(C.int(maxFeatureNum))
		this.HandlesCachedSize[cacheName] = 0
		this.HandlesMaxCachedSize[cacheName] = maxFeatureNum
	}
	return this.Handles[cacheName] != nil
}

func (this *ChlFaceX) ChlFaceSdkListCachedIsFull(cacheName string) bool {
	if maxCachedSize, ok := this.HandlesMaxCachedSize[cacheName]; ok {
		if cachedSize, ok := this.HandlesCachedSize[cacheName]; ok {
			if cachedSize >= maxCachedSize {
				return true
			}
		}
	}
	return false
}
func (this *ChlFaceX) ChlFaceSdkListCachedRemainCap(cacheName string) int {
	if maxCachedSize, ok := this.HandlesMaxCachedSize[cacheName]; ok {
		if cachedSize, ok := this.HandlesCachedSize[cacheName]; ok {
			return maxCachedSize - cachedSize
		}
	}
	return 0
}
func (this *ChlFaceX) ChlFaceSdkListMaxCachedSize(cacheName string) int {
	if maxCachedSize, ok := this.HandlesMaxCachedSize[cacheName]; ok {
		return maxCachedSize
	}
	return -9
}

func (this *ChlFaceX) ChlFaceSdkListInsert(cacheName string, pos int, pFeatures []byte, nFeatureNum int) (int, int) {
	if v, ok := this.Handles[cacheName]; ok {
		if pFeatures == nil || len(pFeatures) != nFeatureNum*this.FaceMaxFeatureSize {
			return -9, -1
		}
		if maxCachedSize, ok := this.HandlesMaxCachedSize[cacheName]; ok {
			if cachedSize, ok := this.HandlesCachedSize[cacheName]; ok {
				var allowCachedSize int = cachedSize + nFeatureNum
				fmt.Println("cachedSize  :", cachedSize, allowCachedSize)
				if allowCachedSize >= maxCachedSize {
					this.HandlesCachedSize[cacheName] = allowCachedSize
					return -10, -1
				}
			}
		}
		var pos C.int = C.int(pos)
		total := C.ChlFaceSdkListInsert(v, &pos, (*C.BYTE)(unsafe.Pointer(&pFeatures[0])), C.int(nFeatureNum))
		this.HandlesCachedSize[cacheName] = int(total)
		return int(total), int(pos)
	}
	return -10, -1
}
func (this *ChlFaceX) ChlFaceSdkListRemove(cacheName string, nPos int, nFeatureNum int) int {
	if v, ok := this.Handles[cacheName]; ok {
		total := C.ChlFaceSdkListRemove(v, C.int(nPos), C.int(nFeatureNum))
		this.HandlesCachedSize[cacheName] = int(total)
		return int(total)
	}
	return -10
}
func (this *ChlFaceX) ChlFaceSdkListClearAll(cacheName string) {
	if v, ok := this.Handles[cacheName]; ok && v != nil {
		C.ChlFaceSdkListClearAll(v)
		this.HandlesCachedSize[cacheName] = 0
	}
}

func (this *ChlFaceX) ChlFaceSdkListCompare(cacheName string, channelNo int, pFeature []byte, nPosBegin int, nCompareNum int) (compareNum int, scores []byte) {

	if pFeature == nil {
		return -9, nil
	}
	this.incr()
	if channelNo >= this.MaxChannelNum || channelNo == DefaultChannelNo{
		channelNo = int(this.CallNum) % this.MaxChannelNum
	}
	if v, ok := this.Handles[cacheName]; ok && v != nil {
		if(this.HandlesCachedSize[cacheName]==0){
			return 0,nil
		}
		var SCORE []byte = make([]byte, this.HandlesCachedSize[cacheName])
		SCORES := (*C.BYTE)(unsafe.Pointer(&SCORE[0]))
		total := C.ChlFaceSdkListCompare(C.int(channelNo), v, (*C.BYTE)(unsafe.Pointer(&pFeature[0])), C.int(nPosBegin), C.int(nCompareNum), SCORES)
		return int(total), SCORE
	}
	return -10, nil
}

func (this *ChlFaceX) ChlFaceSdkListDestroy(cacheName string) {
	if v, ok := this.Handles[cacheName]; ok && v != nil {
		C.ChlFaceSdkListDestroy(v)
		delete(this.Handles, cacheName)
		delete(this.HandlesCachedSize, cacheName)
		delete(this.HandlesMaxCachedSize, cacheName)
	}
}

func makeC_FACE_DETECT_RESULT(pFaceResultX *FACE_DETECT_RESULTX) *C.struct_FACE_DETECT_RESULT {
	var pFaceResult *C.struct_FACE_DETECT_RESULT = new(C.struct_FACE_DETECT_RESULT)
	pFaceResult.rcFace = pFaceResultX.GetRectFace()
	pFaceResult.ptLeftEye = pFaceResultX.GetLeftEye()
	pFaceResult.ptRightEye = pFaceResultX.GetRightEye()
	pFaceResult.ptMouth = pFaceResultX.GetMouth()
	pFaceResult.ptNose = pFaceResultX.GetNose()
	pFaceResult.FaceData = pFaceResultX.GetFaceData()
	pFaceResult.nAngleYaw, pFaceResult.nAnglePitch, pFaceResult.nAngleRoll, pFaceResult.nQuality = pFaceResultX.GetAttrs()
	return pFaceResult
}
func makeGO_FACE_DETECT_RESULT(pFaceResult *C.struct_FACE_DETECT_RESULT, pFaceResultx *FACE_DETECT_RESULTX) {
	pFaceResultx.SetRectFace(pFaceResult.rcFace)
	pFaceResultx.SetFaceData(pFaceResult.FaceData)
	var p C.struct_tagPOINT = pFaceResult.ptLeftEye
	pFaceResultx.SetLeftEye(POINTX{int(p.x), int(p.y)})
	p = pFaceResult.ptRightEye
	pFaceResultx.SetRightEye(POINTX{int(p.x), int(p.y)})
	p = pFaceResult.ptMouth
	pFaceResultx.SetMouth(POINTX{int(p.x), int(p.y)})
	p = pFaceResult.ptNose
	pFaceResultx.SetNose(POINTX{int(p.x), int(p.y)})
	pFaceResultx.SetAttrs(int(pFaceResult.nAngleYaw), int(pFaceResult.nAnglePitch), int(pFaceResult.nAngleRoll), int(pFaceResult.nQuality))
}
