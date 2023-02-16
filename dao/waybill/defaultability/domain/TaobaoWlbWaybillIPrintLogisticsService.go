package domain

type TaobaoWlbWaybillIPrintLogisticsService struct {
	/*
	   服务编码     */
	ServiceCode *string `json:"service_code,omitempty" `

	/*
	   服务类型值，json格式表示     */
	ServiceValue4Json *string `json:"service_value4_json,omitempty" `
}

func (s *TaobaoWlbWaybillIPrintLogisticsService) SetServiceCode(v string) *TaobaoWlbWaybillIPrintLogisticsService {
	s.ServiceCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintLogisticsService) SetServiceValue4Json(v string) *TaobaoWlbWaybillIPrintLogisticsService {
	s.ServiceValue4Json = &v
	return s
}
