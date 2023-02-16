package domain

type CainiaoCloudprintIsvtemplatesGetKeyResult struct {
	/*
	   key名称     */
	KeyName *string `json:"key_name,omitempty" `
}

func (s *CainiaoCloudprintIsvtemplatesGetKeyResult) SetKeyName(v string) *CainiaoCloudprintIsvtemplatesGetKeyResult {
	s.KeyName = &v
	return s
}
