package domain

type TaobaoWlbWaybillIQuerydetailWaybillAddress struct {
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

func (s *TaobaoWlbWaybillIQuerydetailWaybillAddress) SetAddressDetail(v string) *TaobaoWlbWaybillIQuerydetailWaybillAddress {
	s.AddressDetail = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillAddress) SetArea(v string) *TaobaoWlbWaybillIQuerydetailWaybillAddress {
	s.Area = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillAddress) SetCity(v string) *TaobaoWlbWaybillIQuerydetailWaybillAddress {
	s.City = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillAddress) SetProvince(v string) *TaobaoWlbWaybillIQuerydetailWaybillAddress {
	s.Province = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillAddress) SetTown(v string) *TaobaoWlbWaybillIQuerydetailWaybillAddress {
	s.Town = &v
	return s
}
