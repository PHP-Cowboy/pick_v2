package domain

type TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo struct {
	/*
	   打印提示信息编码     */
	NoticeCode *string `json:"notice_code,omitempty" `

	/*
	   打印提示信息     */
	NoticeMessage *string `json:"notice_message,omitempty" `

	/*
	   打印次数     */
	PrintQuantity *int64 `json:"print_quantity,omitempty" `

	/*
	   电子面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `
}

func (s *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo) SetNoticeCode(v string) *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo {
	s.NoticeCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo) SetNoticeMessage(v string) *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo {
	s.NoticeMessage = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo) SetPrintQuantity(v int64) *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo {
	s.PrintQuantity = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo) SetWaybillCode(v string) *TaobaoWlbWaybillIPrintWaybillApplyPrintCheckInfo {
	s.WaybillCode = &v
	return s
}
