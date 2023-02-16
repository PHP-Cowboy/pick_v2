package domain

type CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult struct {
	/*
	   data     */
	Data *CainiaoCloudprintTemplatesMigrateCustomAreaResult `json:"data,omitempty" `

	/*
	   errorCode     */
	ErrorCode *string `json:"error_code,omitempty" `

	/*
	   errorMessage     */
	ErrorMessage *string `json:"error_message,omitempty" `

	/*
	   success     */
	Success *bool `json:"success,omitempty" `
}

func (s *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult) SetData(v CainiaoCloudprintTemplatesMigrateCustomAreaResult) *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult {
	s.Data = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult) SetErrorCode(v string) *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult {
	s.ErrorCode = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult) SetErrorMessage(v string) *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult {
	s.ErrorMessage = &v
	return s
}
func (s *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult) SetSuccess(v bool) *CainiaoCloudprintTemplatesMigrateCloudPrintBaseResult {
	s.Success = &v
	return s
}
