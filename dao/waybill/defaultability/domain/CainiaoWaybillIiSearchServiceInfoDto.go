package domain

type CainiaoWaybillIiSearchServiceInfoDto struct {
	/*
	   服务名称     */
	ServiceName *string `json:"service_name,omitempty" `

	/*
	   服务编码     */
	ServiceCode *string `json:"service_code,omitempty" `

	/*
	   服务属性定义     */
	ServiceAttributes *[]CainiaoWaybillIiSearchServiceAttributeDto `json:"service_attributes,omitempty" `

	/*
	   服务的官方描述，可以用作前端展示     */
	ServiceDesc *string `json:"service_desc,omitempty" `

	/*
	   该服务是否为必选服务     */
	Required *bool `json:"required,omitempty" `
}

func (s *CainiaoWaybillIiSearchServiceInfoDto) SetServiceName(v string) *CainiaoWaybillIiSearchServiceInfoDto {
	s.ServiceName = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceInfoDto) SetServiceCode(v string) *CainiaoWaybillIiSearchServiceInfoDto {
	s.ServiceCode = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceInfoDto) SetServiceAttributes(v []CainiaoWaybillIiSearchServiceAttributeDto) *CainiaoWaybillIiSearchServiceInfoDto {
	s.ServiceAttributes = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceInfoDto) SetServiceDesc(v string) *CainiaoWaybillIiSearchServiceInfoDto {
	s.ServiceDesc = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceInfoDto) SetRequired(v bool) *CainiaoWaybillIiSearchServiceInfoDto {
	s.Required = &v
	return s
}
