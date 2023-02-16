package domain

type CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult struct {
	/*
	   data     */
	Datas *[]CainiaoCloudprintIsvtemplatesGetCustomTemplateResult `json:"datas,omitempty" `

	/*
	   错误码     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   错误消息     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   是否成功     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult) SetDatas(v []CainiaoCloudprintIsvtemplatesGetCustomTemplateResult) *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult {
	s.Datas = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintIsvtemplatesGetCloudPrintBaseResult {
	s.Success = &v
	return s
}
