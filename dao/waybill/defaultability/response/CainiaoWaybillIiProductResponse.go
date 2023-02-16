package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoWaybillIiProductResponse struct {

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
	ProductTypes []domain.CainiaoWaybillIiProductWaybillProductType `json:"product_types,omitempty" `
}
