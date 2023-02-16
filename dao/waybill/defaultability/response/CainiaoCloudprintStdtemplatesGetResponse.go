package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoCloudprintStdtemplatesGetResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   结果集
	*/
	Result domain.CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult `json:"result,omitempty" `
}
