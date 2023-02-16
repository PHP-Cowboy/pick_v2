package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoCloudprintTemplatesMigrateResponse struct {

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
	Result domain.CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult `json:"result,omitempty" `
}
