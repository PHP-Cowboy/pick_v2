package domain

type CainiaoCloudprintClientinfoPutCloudPrintBaseResult struct {
	/*
	   errorMessage     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   errorCode     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   data     */
	Data *string `json:"data,omitempty" `

	/*
	   success     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintClientinfoPutCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintClientinfoPutCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintClientinfoPutCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintClientinfoPutCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintClientinfoPutCloudPrintBaseResult) SetData(v string) *CainiaoCloudprintClientinfoPutCloudPrintBaseResult {
	s.Data = &v
	return s
}
func (s *CainiaoCloudprintClientinfoPutCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintClientinfoPutCloudPrintBaseResult {
	s.Success = &v
	return s
}
