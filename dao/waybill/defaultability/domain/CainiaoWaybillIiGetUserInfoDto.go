package domain

type CainiaoWaybillIiGetUserInfoDto struct {
	/*
	   发货地址需要通过<a href="http://open.taobao.com/doc2/detail.htm?spm=a219a.7629140.0.0.3OFCPk&treeId=17&articleId=104860&docType=1">search接口</a>     */
	Address *CainiaoWaybillIiGetAddressDto `json:"address,omitempty" `

	/*
	   手机号码（手机号和固定电话不能同时为空），长度小于20     */
	Mobile *string `json:"mobile,omitempty" `

	/*
	   姓名，长度小于40     */
	Name *string `json:"name,omitempty" `

	/*
	   固定电话（手机号和固定电话不能同时为空），长度小于20     */
	Phone *string `json:"phone,omitempty" `
}

func (s *CainiaoWaybillIiGetUserInfoDto) SetAddress(v CainiaoWaybillIiGetAddressDto) *CainiaoWaybillIiGetUserInfoDto {
	s.Address = &v
	return s
}
func (s *CainiaoWaybillIiGetUserInfoDto) SetMobile(v string) *CainiaoWaybillIiGetUserInfoDto {
	s.Mobile = &v
	return s
}
func (s *CainiaoWaybillIiGetUserInfoDto) SetName(v string) *CainiaoWaybillIiGetUserInfoDto {
	s.Name = &v
	return s
}
func (s *CainiaoWaybillIiGetUserInfoDto) SetPhone(v string) *CainiaoWaybillIiGetUserInfoDto {
	s.Phone = &v
	return s
}
