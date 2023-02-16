package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoCloudprintIsvResourcesGetResponse struct {

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
	Result domain.CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult `json:"result,omitempty" `
}
