package model

type FeatureEntry struct {
	PeopleId    int64
	FaceImageId string
	Pos         int
	Status      int
}

type ImageSource struct {
	TaskTuple
	Type                int // 0:image ,1:feature
	CameraId            string
	CameraIp            string
	CreateTime          int64
	CaptureTime         int64
	OrigPath            string
	ImageOriginalName   string
	ImageContextPath    string
	OrigImageUri        string
	OrigImageUuid       string
	Topk                int
	FaceNum             int
	FACE_DETECT_RESULTS []FACE_DETECT_RESULT
	FaceFeatureBufs     []byte
	FaceRects           string
	FaceFeatureBufsB64  string
}

type POINT struct {
	X, Y int
}

type RECT struct {
	Left, Top, Right, Bottom int
}
type FACE_DETECT_RESULT struct {
	FaceRect                                 RECT
	LeftEye, RightEye, Mouth, Nose           POINT
	AngleYaw, AnglePitch, AngleRoll, Quality int
}
