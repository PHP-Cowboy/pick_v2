package domain

type CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult struct {
	/*
	   data     */
	ResourceList *[]CainiaoCloudprintIsvResourcesGetIsvResourceDo `json:"resource_list,omitempty" `

	/*
	   错误码     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   错误消息     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   状态     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult) SetResourceList(v []CainiaoCloudprintIsvResourcesGetIsvResourceDo) *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult {
	s.ResourceList = &v
	return s
}
func (s *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintIsvResourcesGetCloudPrintBaseResult {
	s.Success = &v
	return s
}
