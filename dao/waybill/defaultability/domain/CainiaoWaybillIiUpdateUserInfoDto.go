package domain

type CainiaoWaybillIiUpdateUserInfoDto struct {
	/*
	   地址     */
	Address *CainiaoWaybillIiUpdateAddressDto `json:"address,omitempty" `

	/*
	   手机号码     */
	Mobile *string `json:"mobile,omitempty" `

	/*
	   姓名     */
	Name *string `json:"name,omitempty" `

	/*
	   固定电话     */
	Phone *string `json:"phone,omitempty" `

	/*
	   开放地址ID     */
	Oaid *string `json:"oaid,omitempty" `

	/*
	   菜鸟地址ID，针对电商平台加密订单场景使用，淘系订单使用oaid，非淘使用caid。     */
	Caid *string `json:"caid,omitempty" `
}

func (s *CainiaoWaybillIiUpdateUserInfoDto) SetAddress(v CainiaoWaybillIiUpdateAddressDto) *CainiaoWaybillIiUpdateUserInfoDto {
	s.Address = &v
	return s
}
func (s *CainiaoWaybillIiUpdateUserInfoDto) SetMobile(v string) *CainiaoWaybillIiUpdateUserInfoDto {
	s.Mobile = &v
	return s
}
func (s *CainiaoWaybillIiUpdateUserInfoDto) SetName(v string) *CainiaoWaybillIiUpdateUserInfoDto {
	s.Name = &v
	return s
}
func (s *CainiaoWaybillIiUpdateUserInfoDto) SetPhone(v string) *CainiaoWaybillIiUpdateUserInfoDto {
	s.Phone = &v
	return s
}
func (s *CainiaoWaybillIiUpdateUserInfoDto) SetOaid(v string) *CainiaoWaybillIiUpdateUserInfoDto {
	s.Oaid = &v
	return s
}
func (s *CainiaoWaybillIiUpdateUserInfoDto) SetCaid(v string) *CainiaoWaybillIiUpdateUserInfoDto {
	s.Caid = &v
	return s
}
