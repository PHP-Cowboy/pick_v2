package domain

type CainiaoCloudprintCustomaresGetCustomAreaResult struct {
	/*
	   自定义区id     */
	CustomAreaId *int64 `json:"custom_area_id,omitempty" `

	/*
	   自定义区url     */
	CustomAreaUrl *string `json:"custom_area_url,omitempty" `

	/*
	   keys     */
	Keys *[]CainiaoCloudprintCustomaresGetKeyResult `json:"keys,omitempty" `
}

func (s *CainiaoCloudprintCustomaresGetCustomAreaResult) SetCustomAreaId(v int64) *CainiaoCloudprintCustomaresGetCustomAreaResult {
	s.CustomAreaId = &v
	return s
}
func (s *CainiaoCloudprintCustomaresGetCustomAreaResult) SetCustomAreaUrl(v string) *CainiaoCloudprintCustomaresGetCustomAreaResult {
	s.CustomAreaUrl = &v
	return s
}
func (s *CainiaoCloudprintCustomaresGetCustomAreaResult) SetKeys(v []CainiaoCloudprintCustomaresGetKeyResult) *CainiaoCloudprintCustomaresGetCustomAreaResult {
	s.Keys = &v
	return s
}
