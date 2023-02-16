package domain

type CainiaoCloudprintTemplatesMigrateKeyResult struct {
	/*
	   变量名称     */
	KeyName *string `json:"key_name,omitempty" `
}

func (s *CainiaoCloudprintTemplatesMigrateKeyResult) SetKeyName(v string) *CainiaoCloudprintTemplatesMigrateKeyResult {
	s.KeyName = &v
	return s
}
