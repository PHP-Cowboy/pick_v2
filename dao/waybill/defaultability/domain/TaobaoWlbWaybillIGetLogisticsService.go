package domain

type TaobaoWlbWaybillIGetLogisticsService struct {
	/*
	   服务编码     */
	ServiceCode *string `json:"service_code,omitempty" `

	/*
	   服务类型值，json格式表示     */
	ServiceValue4Json *string `json:"service_value4_json,omitempty" `
}

func (s *TaobaoWlbWaybillIGetLogisticsService) SetServiceCode(v string) *TaobaoWlbWaybillIGetLogisticsService {
	s.ServiceCode = &v
	return s
}
func (s *TaobaoWlbWaybillIGetLogisticsService) SetServiceValue4Json(v string) *TaobaoWlbWaybillIGetLogisticsService {
	s.ServiceValue4Json = &v
	return s
}
