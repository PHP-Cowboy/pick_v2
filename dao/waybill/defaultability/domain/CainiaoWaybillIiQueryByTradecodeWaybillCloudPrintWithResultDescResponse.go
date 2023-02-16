package domain

type CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse struct {
	/*
	   单查询错误code     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   单查询错误message     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   请求id     */
	ObjectId *string `json:"object_id,omitempty" `

	/*
	   是否查询     */
	Success *bool `json:"success,omitempty" `

	/*
	   面单信息     */
	WaybillCloudPrintResponse *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse `json:"waybill_cloud_print_response,omitempty" `
}

func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse) SetErrorCode(v string) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse) SetErrorMessage(v string) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse) SetObjectId(v string) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse {
	s.ObjectId = &v
	return s
}
func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse) SetSuccess(v bool) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse {
	s.Success = &v
	return s
}
func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse) SetWaybillCloudPrintResponse(v CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintWithResultDescResponse {
	s.WaybillCloudPrintResponse = &v
	return s
}
