package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoCloudprintCustomaresGetResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   结果
	*/
	Result domain.CainiaoCloudprintCustomaresGetCloudPrintBaseResult `json:"result,omitempty" `
}
