package domain

type TaobaoWlbWaybillIGetWaybillApplyNewRequest struct {
	/*
	   物流服务商编码     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   收\发货地址     */
	ShippingAddress *TaobaoWlbWaybillIGetWaybillAddress `json:"shipping_address,omitempty" `

	/*
	   订单数据     */
	TradeOrderInfoCols *[]TaobaoWlbWaybillIGetTradeOrderInfo `json:"trade_order_info_cols,omitempty" `
}

func (s *TaobaoWlbWaybillIGetWaybillApplyNewRequest) SetCpCode(v string) *TaobaoWlbWaybillIGetWaybillApplyNewRequest {
	s.CpCode = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewRequest) SetShippingAddress(v TaobaoWlbWaybillIGetWaybillAddress) *TaobaoWlbWaybillIGetWaybillApplyNewRequest {
	s.ShippingAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewRequest) SetTradeOrderInfoCols(v []TaobaoWlbWaybillIGetTradeOrderInfo) *TaobaoWlbWaybillIGetWaybillApplyNewRequest {
	s.TradeOrderInfoCols = &v
	return s
}
