package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoCloudprintIsvtemplatesGetResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   result
	*/
	Result domain.CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult `json:"result,omitempty" `
}
