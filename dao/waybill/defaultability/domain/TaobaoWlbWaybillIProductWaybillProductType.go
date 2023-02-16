package domain

type TaobaoWlbWaybillIProductWaybillProductType struct {
	/*
	   物流产品类型编码     */
	Code *string `json:"code,omitempty" `

	/*
	   物流产品类型名称     */
	Name *string `json:"name,omitempty" `

	/*
	   物流服务类型列表     */
	ServiceTypes *[]TaobaoWlbWaybillIProductWaybillServiceType `json:"service_types,omitempty" `
}

func (s *TaobaoWlbWaybillIProductWaybillProductType) SetCode(v string) *TaobaoWlbWaybillIProductWaybillProductType {
	s.Code = &v
	return s
}
func (s *TaobaoWlbWaybillIProductWaybillProductType) SetName(v string) *TaobaoWlbWaybillIProductWaybillProductType {
	s.Name = &v
	return s
}
func (s *TaobaoWlbWaybillIProductWaybillProductType) SetServiceTypes(v []TaobaoWlbWaybillIProductWaybillServiceType) *TaobaoWlbWaybillIProductWaybillProductType {
	s.ServiceTypes = &v
	return s
}
