package domain

type CainiaoWaybillIiUpdatePackageInfoDto struct {
	/*
	   商品     */
	Items *[]CainiaoWaybillIiUpdateItem `json:"items,omitempty" `

	/*
	   体积 defalutValue:0    */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   重量 defalutValue:0    */
	Weight *int64 `json:"weight,omitempty" `
}

func (s *CainiaoWaybillIiUpdatePackageInfoDto) SetItems(v []CainiaoWaybillIiUpdateItem) *CainiaoWaybillIiUpdatePackageInfoDto {
	s.Items = &v
	return s
}
func (s *CainiaoWaybillIiUpdatePackageInfoDto) SetVolume(v int64) *CainiaoWaybillIiUpdatePackageInfoDto {
	s.Volume = &v
	return s
}
func (s *CainiaoWaybillIiUpdatePackageInfoDto) SetWeight(v int64) *CainiaoWaybillIiUpdatePackageInfoDto {
	s.Weight = &v
	return s
}
