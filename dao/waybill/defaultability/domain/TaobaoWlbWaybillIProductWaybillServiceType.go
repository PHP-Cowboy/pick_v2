package domain

type TaobaoWlbWaybillIProductWaybillServiceType struct {
	/*
	   物流服务能力编码     */
	Code *string `json:"code,omitempty" `

	/*
	   物流服务能力名称     */
	Name *string `json:"name,omitempty" `
}

func (s *TaobaoWlbWaybillIProductWaybillServiceType) SetCode(v string) *TaobaoWlbWaybillIProductWaybillServiceType {
	s.Code = &v
	return s
}
func (s *TaobaoWlbWaybillIProductWaybillServiceType) SetName(v string) *TaobaoWlbWaybillIProductWaybillServiceType {
	s.Name = &v
	return s
}
