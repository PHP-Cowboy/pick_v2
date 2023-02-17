package domain

type CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult struct {
	/*
	   data     */
	Data *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult `json:"data,omitempty" `

	/*
	   errorCode     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   错误消息     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   success     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult) SetData(v CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult) *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult {
	s.Data = &v
	return s
}
func (s *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintSingleCustomareaGetCloudPrintBaseResult {
	s.Success = &v
	return s
}
