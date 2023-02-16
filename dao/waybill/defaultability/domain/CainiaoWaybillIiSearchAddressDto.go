package domain

type CainiaoWaybillIiSearchAddressDto struct {
	/*
	   详细地址     */
	Detail *string `json:"detail,omitempty" `

	/*
	   区名称（三级地址）     */
	District *string `json:"district,omitempty" `

	/*
	   市名称（二级地址）     */
	City *string `json:"city,omitempty" `

	/*
	   省名称（一级地址）     */
	Province *string `json:"province,omitempty" `

	/*
	   街道\镇名称（四级地址）     */
	Town *string `json:"town,omitempty" `

	/*
	   订购关系id     */
	WaybillAddressId *string `json:"waybill_address_id,omitempty" `
}

func (s *CainiaoWaybillIiSearchAddressDto) SetDetail(v string) *CainiaoWaybillIiSearchAddressDto {
	s.Detail = &v
	return s
}
func (s *CainiaoWaybillIiSearchAddressDto) SetDistrict(v string) *CainiaoWaybillIiSearchAddressDto {
	s.District = &v
	return s
}
func (s *CainiaoWaybillIiSearchAddressDto) SetCity(v string) *CainiaoWaybillIiSearchAddressDto {
	s.City = &v
	return s
}
func (s *CainiaoWaybillIiSearchAddressDto) SetProvince(v string) *CainiaoWaybillIiSearchAddressDto {
	s.Province = &v
	return s
}
func (s *CainiaoWaybillIiSearchAddressDto) SetTown(v string) *CainiaoWaybillIiSearchAddressDto {
	s.Town = &v
	return s
}
func (s *CainiaoWaybillIiSearchAddressDto) SetWaybillAddressId(v string) *CainiaoWaybillIiSearchAddressDto {
	s.WaybillAddressId = &v
	return s
}
