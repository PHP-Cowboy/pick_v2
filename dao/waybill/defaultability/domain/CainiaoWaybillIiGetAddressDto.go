package domain

type CainiaoWaybillIiGetAddressDto struct {
	/*
	   城市，长度小于20     */
	City *string `json:"city,omitempty" `

	/*
	   详细地址，长度小于256     */
	Detail *string `json:"detail,omitempty" `

	/*
	   区，长度小于20     */
	District *string `json:"district,omitempty" `

	/*
	   省，长度小于20     */
	Province *string `json:"province,omitempty" `

	/*
	   街道，长度小于30     */
	Town *string `json:"town,omitempty" `
}

func (s *CainiaoWaybillIiGetAddressDto) SetCity(v string) *CainiaoWaybillIiGetAddressDto {
	s.City = &v
	return s
}
func (s *CainiaoWaybillIiGetAddressDto) SetDetail(v string) *CainiaoWaybillIiGetAddressDto {
	s.Detail = &v
	return s
}
func (s *CainiaoWaybillIiGetAddressDto) SetDistrict(v string) *CainiaoWaybillIiGetAddressDto {
	s.District = &v
	return s
}
func (s *CainiaoWaybillIiGetAddressDto) SetProvince(v string) *CainiaoWaybillIiGetAddressDto {
	s.Province = &v
	return s
}
func (s *CainiaoWaybillIiGetAddressDto) SetTown(v string) *CainiaoWaybillIiGetAddressDto {
	s.Town = &v
	return s
}
