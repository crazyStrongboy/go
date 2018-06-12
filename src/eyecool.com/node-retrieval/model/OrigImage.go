package model



type OrigImage struct {
	Id         int  `json:"id" xorm:"pk autoincr"`
	FaceNum    int
	FeatList   string
	Feat   string `xorm:"-"`
}



func (self *OrigImage) TableName() string {
	return "buz_orig_image"
}



type OrigImageFull struct {
	Id         int  `json:"id" xorm:"pk autoincr"`
	FaceNum    int
	FeatList   string
	Uuid string
	CameraId  string
	ImageName  string
	ImageRealPath string
	FaceRect  string
	FaceProp    string
	Timestamp int64
	ImageContextPath string
	FaceImageUri string
	PictureUri string
}



func (self *OrigImageFull) TableName() string {
	return "buz_orig_image"
}



