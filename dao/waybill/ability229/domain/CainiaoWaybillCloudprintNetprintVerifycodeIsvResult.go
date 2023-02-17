package domain

type CainiaoWaybillCloudprintNetprintVerifycodeIsvResult struct {
	/*
	   错误码，0 表示正常     */
	ServerErrorCode *string `json:"server_error_code,omitempty" `

	/*
	   描述     */
	Describe *string `json:"describe,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintVerifycodeIsvResult) SetServerErrorCode(v string) *CainiaoWaybillCloudprintNetprintVerifycodeIsvResult {
	s.ServerErrorCode = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintVerifycodeIsvResult) SetDescribe(v string) *CainiaoWaybillCloudprintNetprintVerifycodeIsvResult {
	s.Describe = &v
	return s
}
