package domain

type CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult struct {
	/*
	   所有cp的标准模板     */
	Datas *[]CainiaoCloudprintStdtemplatesGetStandardTemplateResult `json:"datas,omitempty" `

	/*
	   错误码     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   错误信息     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   是否成功     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult) SetDatas(v []CainiaoCloudprintStdtemplatesGetStandardTemplateResult) *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult {
	s.Datas = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintStdtemplatesGetCloudPrintBaseResult {
	s.Success = &v
	return s
}
