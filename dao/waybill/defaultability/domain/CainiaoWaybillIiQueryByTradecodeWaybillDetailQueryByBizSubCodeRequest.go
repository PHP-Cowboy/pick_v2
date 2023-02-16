package domain

type CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest struct {
	/*
	   订单号     */
	BizSubCode *string `json:"biz_sub_code,omitempty" `

	/*
	   请求id     */
	ObjectId *string `json:"object_id,omitempty" `
}

func (s *CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest) SetBizSubCode(v string) *CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest {
	s.BizSubCode = &v
	return s
}
func (s *CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest) SetObjectId(v string) *CainiaoWaybillIiQueryByTradecodeWaybillDetailQueryByBizSubCodeRequest {
	s.ObjectId = &v
	return s
}
