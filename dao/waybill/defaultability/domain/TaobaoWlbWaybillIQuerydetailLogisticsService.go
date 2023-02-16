package domain

type TaobaoWlbWaybillIQuerydetailLogisticsService struct {
	/*
	   服务编码     */
	ServiceCode *string `json:"service_code,omitempty" `

	/*
	   服务类型值，json格式表示     */
	ServiceValue4Json *string `json:"service_value4_json,omitempty" `
}

func (s *TaobaoWlbWaybillIQuerydetailLogisticsService) SetServiceCode(v string) *TaobaoWlbWaybillIQuerydetailLogisticsService {
	s.ServiceCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailLogisticsService) SetServiceValue4Json(v string) *TaobaoWlbWaybillIQuerydetailLogisticsService {
	s.ServiceValue4Json = &v
	return s
}
