package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoWaybillIiGetResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   系统自动生成
	*/
	Modules []domain.CainiaoWaybillIiGetWaybillCloudPrintResponse `json:"modules,omitempty" `
}
