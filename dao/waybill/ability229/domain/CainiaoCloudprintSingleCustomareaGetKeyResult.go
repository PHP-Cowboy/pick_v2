package domain

type CainiaoCloudprintSingleCustomareaGetKeyResult struct {
	/*
	   动态key的名称     */
	KeyName *string `json:"key_name,omitempty" `
}

func (s *CainiaoCloudprintSingleCustomareaGetKeyResult) SetKeyName(v string) *CainiaoCloudprintSingleCustomareaGetKeyResult {
	s.KeyName = &v
	return s
}
