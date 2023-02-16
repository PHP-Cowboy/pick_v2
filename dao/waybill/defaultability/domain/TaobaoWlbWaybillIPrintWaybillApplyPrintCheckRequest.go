package domain

type TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest struct {
	/*
	   物流服务商Code     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   面单详情信息     */
	PrintCheckInfoCols *[]TaobaoWlbWaybillIPrintPrintCheckInfo `json:"print_check_info_cols,omitempty" `
}

func (s *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest) SetCpCode(v string) *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest {
	s.CpCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest) SetPrintCheckInfoCols(v []TaobaoWlbWaybillIPrintPrintCheckInfo) *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckRequest {
	s.PrintCheckInfoCols = &v
	return s
}
