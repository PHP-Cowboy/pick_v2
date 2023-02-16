package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoWaybillIiQueryByWaybillcodeResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   查询返回值
	*/
	Modules []domain.CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse `json:"modules,omitempty" `
}
