package response

type CainiaoWaybillCloudprintNetprintPrintResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   错误码，0 为成功
	*/
	ServerErrorCode string `json:"server_error_code,omitempty" `
	/*
	   错误描述
	*/
	Describe string `json:"describe,omitempty" `
}
