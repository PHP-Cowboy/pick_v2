package response

import (
	"pick_v2/dao/waybill/ability229/domain"
)

type CainiaoWaybillCloudprintNetprintBindResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   返回值
	*/
	Result domain.CainiaoWaybillCloudprintNetprintBindIsvResult `json:"result,omitempty" `
}
