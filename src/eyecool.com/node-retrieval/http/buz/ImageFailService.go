package buz

import (
	"eyecool.com/node-retrieval/logic"
)

type ImageFailService struct {
}

type ImageFailRequest struct {
	Repository_id string //库ID($编号@$集群号)
	Start         int    //从第几个结果开始返回
	Limit         int    //返回至多多少个结果
}

type ImageFailResponse struct {
	Total   int              `json:"total,omitempty"`
	Rtn     int              `json:"rtn"`
	Message string           `json:"message,omitempty"`
	Results ImageFailResults `json:"results,omitempty"`
}

type ImageFailResult struct {
	PictureUri    string `json:"picture_uri,omitempty"`    //导入图片的 uri
	FailedMessage string `json:"failed_message,omitempty"` //失败原因
	FailedRtn     int    `json:"failed_rtn,omitempty"`     //失败错误码
}
type ImageFailResults []*ImageFailResult

var imageFailLogic = new(logic.ImageFailLogic)

func (service *ImageFailService) GetFailImage(request *ImageFailRequest) *ImageFailResponse {
	response := new(ImageFailResponse)
	imageFails := imageFailLogic.FindByRepositoryId(request.Start, request.Limit, request.Repository_id)
	imageFailResults := ImageFailResults{}
	if len(imageFails) > 0 {
		for _, imageFail := range imageFails {
			imageFailResult := &ImageFailResult{
				PictureUri:    imageFail.ImageContextPath + imageFail.ImageUri,
				FailedMessage: imageFail.ImageDesc,
				FailedRtn:     imageFail.FailCode,
			}
			imageFailResults = append(imageFailResults, imageFailResult)
		}
	}
	response.Rtn = 0
	response.Message = "查询成功！"
	response.Total = len(imageFails)
	response.Results = imageFailResults
	return response
}
