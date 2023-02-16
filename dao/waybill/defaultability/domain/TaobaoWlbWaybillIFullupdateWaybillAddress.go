package domain

type TaobaoWlbWaybillIFullupdateWaybillAddress struct {
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
	   一级地址（省、直辖市     */
	Province *string `json:"province,omitempty" `

	/*
	   街道\镇名称（四级地址）     */
	Town *string `json:"town,omitempty" `
}

func (s *TaobaoWlbWaybillIFullupdateWaybillAddress) SetAddressDetail(v string) *TaobaoWlbWaybillIFullupdateWaybillAddress {
	s.AddressDetail = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillAddress) SetArea(v string) *TaobaoWlbWaybillIFullupdateWaybillAddress {
	s.Area = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillAddress) SetCity(v string) *TaobaoWlbWaybillIFullupdateWaybillAddress {
	s.City = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillAddress) SetProvince(v string) *TaobaoWlbWaybillIFullupdateWaybillAddress {
	s.Province = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillAddress) SetTown(v string) *TaobaoWlbWaybillIFullupdateWaybillAddress {
	s.Town = &v
	return s
}
