package domain

type CainiaoWaybillIiProductWaybillProductType struct {
	/*
	   产品code     */
	Code *string `json:"code,omitempty" `

	/*
	   产品名称     */
	Name *string `json:"name,omitempty" `

	/*
	   物流服务     */
	ServiceTypes *[]CainiaoWaybillIiProductWaybillServiceType `json:"service_types,omitempty" `
}

func (s *CainiaoWaybillIiProductWaybillProductType) SetCode(v string) *CainiaoWaybillIiProductWaybillProductType {
	s.Code = &v
	return s
}
func (s *CainiaoWaybillIiProductWaybillProductType) SetName(v string) *CainiaoWaybillIiProductWaybillProductType {
	s.Name = &v
	return s
}
func (s *CainiaoWaybillIiProductWaybillProductType) SetServiceTypes(v []CainiaoWaybillIiProductWaybillServiceType) *CainiaoWaybillIiProductWaybillProductType {
	s.ServiceTypes = &v
	return s
}
