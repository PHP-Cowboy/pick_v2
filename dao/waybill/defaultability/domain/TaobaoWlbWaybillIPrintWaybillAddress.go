package domain

type TaobaoWlbWaybillIPrintWaybillAddress struct {
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
	   末级地址     */
	DivisionId *int64 `json:"division_id,omitempty" `

	/*
	   省名称（一级地址）     */
	Province *string `json:"province,omitempty" `

	/*
	   街道\镇名称（四级地址）     */
	Town *string `json:"town,omitempty" `
}

func (s *TaobaoWlbWaybillIPrintWaybillAddress) SetAddressDetail(v string) *TaobaoWlbWaybillIPrintWaybillAddress {
	s.AddressDetail = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillAddress) SetArea(v string) *TaobaoWlbWaybillIPrintWaybillAddress {
	s.Area = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillAddress) SetCity(v string) *TaobaoWlbWaybillIPrintWaybillAddress {
	s.City = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillAddress) SetDivisionId(v int64) *TaobaoWlbWaybillIPrintWaybillAddress {
	s.DivisionId = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillAddress) SetProvince(v string) *TaobaoWlbWaybillIPrintWaybillAddress {
	s.Province = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillAddress) SetTown(v string) *TaobaoWlbWaybillIPrintWaybillAddress {
	s.Town = &v
	return s
}
