package domain

type CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest struct {
	/*
	   快递公司code     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   请求id     */
	ObjectId *string `json:"object_id,omitempty" `

	/*
	   电子面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `
}

func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest) SetCpCode(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest {
	s.CpCode = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest) SetObjectId(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest {
	s.ObjectId = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest) SetWaybillCode(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillDetailQueryByWaybillCodeRequest {
	s.WaybillCode = &v
	return s
}
