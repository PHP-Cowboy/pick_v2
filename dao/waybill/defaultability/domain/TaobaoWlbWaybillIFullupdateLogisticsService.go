package domain

type TaobaoWlbWaybillIFullupdateLogisticsService struct {
	/*
	   服务编码     */
	ServiceCode *string `json:"service_code,omitempty" `

	/*
	   服务类型值，json格式表示     */
	ServiceValue4Json *string `json:"service_value4_json,omitempty" `
}

func (s *TaobaoWlbWaybillIFullupdateLogisticsService) SetServiceCode(v string) *TaobaoWlbWaybillIFullupdateLogisticsService {
	s.ServiceCode = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateLogisticsService) SetServiceValue4Json(v string) *TaobaoWlbWaybillIFullupdateLogisticsService {
	s.ServiceValue4Json = &v
	return s
}
