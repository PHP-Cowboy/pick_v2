package domain

type CainiaoCloudprintMystdtemplatesGetKeyResult struct {
	/*
	   key的名称     */
	KeyName *string `json:"key_name,omitempty" `
}

func (s *CainiaoCloudprintMystdtemplatesGetKeyResult) SetKeyName(v string) *CainiaoCloudprintMystdtemplatesGetKeyResult {
	s.KeyName = &v
	return s
}
