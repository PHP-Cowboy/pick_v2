package domain

type CainiaoCloudprintCustomaresGetKeyResult struct {
	/*
	   key名称     */
	KeyName *string `json:"key_name,omitempty" `
}

func (s *CainiaoCloudprintCustomaresGetKeyResult) SetKeyName(v string) *CainiaoCloudprintCustomaresGetKeyResult {
	s.KeyName = &v
	return s
}
