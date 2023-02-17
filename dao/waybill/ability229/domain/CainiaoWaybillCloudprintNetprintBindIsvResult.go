package domain

type CainiaoWaybillCloudprintNetprintBindIsvResult struct {
	/*
	   共享码     */
	Data *string `json:"data,omitempty" `

	/*
	   错误码     */
	ServerErrorCode *string `json:"server_error_code,omitempty" `

	/*
	   描述     */
	Describe *string `json:"describe,omitempty" `
}

func (s *CainiaoWaybillCloudprintNetprintBindIsvResult) SetData(v string) *CainiaoWaybillCloudprintNetprintBindIsvResult {
	s.Data = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintBindIsvResult) SetServerErrorCode(v string) *CainiaoWaybillCloudprintNetprintBindIsvResult {
	s.ServerErrorCode = &v
	return s
}
func (s *CainiaoWaybillCloudprintNetprintBindIsvResult) SetDescribe(v string) *CainiaoWaybillCloudprintNetprintBindIsvResult {
	s.Describe = &v
	return s
}
