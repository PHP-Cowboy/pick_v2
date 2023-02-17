package response

type CainiaoCloudprintCmdprintRenderResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   0成功，非0失败
	*/
	RetCode string `json:"ret_code,omitempty" `
	/*
	   错误信息
	*/
	RetMsg string `json:"ret_msg,omitempty" `
	/*
	   指令集内容串
	*/
	CmdContent string `json:"cmd_content,omitempty" `
	/*
	   指令集编码方式：origin-原串 gzip-采用gzip压缩并base64编码
	*/
	CmdEncoding string `json:"cmd_encoding,omitempty" `
}
