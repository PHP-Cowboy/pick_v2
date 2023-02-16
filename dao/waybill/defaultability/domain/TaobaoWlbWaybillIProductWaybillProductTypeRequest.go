package domain

type TaobaoWlbWaybillIProductWaybillProductTypeRequest struct {
	/*
	   物流商编码CODE     */
	CpCode *string `json:"cp_code,omitempty" `
}

func (s *TaobaoWlbWaybillIProductWaybillProductTypeRequest) SetCpCode(v string) *TaobaoWlbWaybillIProductWaybillProductTypeRequest {
	s.CpCode = &v
	return s
}
