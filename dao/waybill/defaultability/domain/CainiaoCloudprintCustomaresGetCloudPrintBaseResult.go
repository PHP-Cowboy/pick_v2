package domain

type CainiaoCloudprintCustomaresGetCloudPrintBaseResult struct {
	/*
	   数据     */
	Datas *[]CainiaoCloudprintCustomaresGetCustomAreaResult `json:"datas,omitempty" `

	/*
	   系统自动生成     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   系统自动生成     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   系统自动生成     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintCustomaresGetCloudPrintBaseResult) SetDatas(v []CainiaoCloudprintCustomaresGetCustomAreaResult) *CainiaoCloudprintCustomaresGetCloudPrintBaseResult {
	s.Datas = &v
	return s
}
func (s *CainiaoCloudprintCustomaresGetCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintCustomaresGetCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintCustomaresGetCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintCustomaresGetCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintCustomaresGetCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintCustomaresGetCloudPrintBaseResult {
	s.Success = &v
	return s
}
