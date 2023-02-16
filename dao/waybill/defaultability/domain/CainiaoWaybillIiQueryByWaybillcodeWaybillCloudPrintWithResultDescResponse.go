package domain

type CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse struct {
	/*
	   错误码     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   错误信息     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   请求id     */
	ObjectId *string `json:"object_id,omitempty" `

	/*
	   请求成功     */
	Success *bool `json:"success,omitempty" `

	/*
	   面单查询结构体     */
	WaybillCloudPrintResponse *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse `json:"waybill_cloud_print_response,omitempty" `
}

func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse) SetErrorCode(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse) SetErrorMessage(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse) SetObjectId(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse {
	s.ObjectId = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse) SetSuccess(v bool) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse {
	s.Success = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse) SetWaybillCloudPrintResponse(v CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintWithResultDescResponse {
	s.WaybillCloudPrintResponse = &v
	return s
}
