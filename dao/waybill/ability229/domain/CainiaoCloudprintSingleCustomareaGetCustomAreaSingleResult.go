package domain

type CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult struct {
	/*
	   自定义区id     */
	CustomAreaId *int64 `json:"custom_area_id,omitempty" `

	/*
	   自定义区内容的URL     */
	CustomAreaUrl *string `json:"custom_area_url,omitempty" `

	/*
	   keys     */
	Keys *[]CainiaoCloudprintSingleCustomareaGetKeyResult `json:"keys,omitempty" `
}

func (s *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult) SetCustomAreaId(v int64) *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult {
	s.CustomAreaId = &v
	return s
}
func (s *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult) SetCustomAreaUrl(v string) *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult {
	s.CustomAreaUrl = &v
	return s
}
func (s *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult) SetKeys(v []CainiaoCloudprintSingleCustomareaGetKeyResult) *CainiaoCloudprintSingleCustomareaGetCustomAreaSingleResult {
	s.Keys = &v
	return s
}
