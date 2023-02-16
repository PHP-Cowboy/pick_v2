package domain

type TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest struct {
	/*
	   需要查询的订单号     */
	TradeOrderList *[]string `json:"trade_order_list,omitempty" `

	/*
	   CP快递公司编码     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   0:根据cp_code和waybil_code查询;1:根据订单号查询(默认根据cp_code和waybill_code查询) defalutValue:0    */
	QueryBy *int64 `json:"query_by,omitempty" `

	/*
	   电子面单单号     */
	WaybillCodes *[]string `json:"waybill_codes,omitempty" `
}

func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest) SetTradeOrderList(v []string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest {
	s.TradeOrderList = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest) SetCpCode(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest {
	s.CpCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest) SetQueryBy(v int64) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest {
	s.QueryBy = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest) SetWaybillCodes(v []string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryRequest {
	s.WaybillCodes = &v
	return s
}
