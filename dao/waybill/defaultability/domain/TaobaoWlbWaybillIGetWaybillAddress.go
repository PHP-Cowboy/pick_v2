package domain

type TaobaoWlbWaybillIGetWaybillAddress struct {
	/*
	   详细地址     */
	AddressDetail *string `json:"address_detail,omitempty" `

	/*
	   区名称（三级地址）     */
	Area *string `json:"area,omitempty" `

	/*
	   市名称（二级地址）     */
	City *string `json:"city,omitempty" `

	/*
	   省名称（一级地址）     */
	Province *string `json:"province,omitempty" `

	/*
	   街道\镇名称（四级地址）     */
	Town *string `json:"town,omitempty" `
}

func (s *TaobaoWlbWaybillIGetWaybillAddress) SetAddressDetail(v string) *TaobaoWlbWaybillIGetWaybillAddress {
	s.AddressDetail = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillAddress) SetArea(v string) *TaobaoWlbWaybillIGetWaybillAddress {
	s.Area = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillAddress) SetCity(v string) *TaobaoWlbWaybillIGetWaybillAddress {
	s.City = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillAddress) SetProvince(v string) *TaobaoWlbWaybillIGetWaybillAddress {
	s.Province = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillAddress) SetTown(v string) *TaobaoWlbWaybillIGetWaybillAddress {
	s.Town = &v
	return s
}
