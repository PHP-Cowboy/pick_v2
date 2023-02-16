package response

import ()

type TaobaoWlbWaybillICancelResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   调用取消是否成功
	*/
	CancelResult bool `json:"cancel_result,omitempty" `
}
