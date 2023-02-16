package domain

type CainiaoWaybillIiUpdateAddressDto struct {
	/*
	   城市     */
	City *string `json:"city,omitempty" `

	/*
	   详细地址     */
	Detail *string `json:"detail,omitempty" `

	/*
	   区地址     */
	District *string `json:"district,omitempty" `

	/*
	   省     */
	Province *string `json:"province,omitempty" `

	/*
	   街道     */
	Town *string `json:"town,omitempty" `
}

func (s *CainiaoWaybillIiUpdateAddressDto) SetCity(v string) *CainiaoWaybillIiUpdateAddressDto {
	s.City = &v
	return s
}
func (s *CainiaoWaybillIiUpdateAddressDto) SetDetail(v string) *CainiaoWaybillIiUpdateAddressDto {
	s.Detail = &v
	return s
}
func (s *CainiaoWaybillIiUpdateAddressDto) SetDistrict(v string) *CainiaoWaybillIiUpdateAddressDto {
	s.District = &v
	return s
}
func (s *CainiaoWaybillIiUpdateAddressDto) SetProvince(v string) *CainiaoWaybillIiUpdateAddressDto {
	s.Province = &v
	return s
}
func (s *CainiaoWaybillIiUpdateAddressDto) SetTown(v string) *CainiaoWaybillIiUpdateAddressDto {
	s.Town = &v
	return s
}
