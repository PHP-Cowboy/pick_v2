package response

import ()

type CainiaoWaybillIiUpdateResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   模板内容
	*/
	PrintData string `json:"print_data,omitempty" `
	/*
	   面单号
	*/
	WaybillCode string `json:"waybill_code,omitempty" `
}
