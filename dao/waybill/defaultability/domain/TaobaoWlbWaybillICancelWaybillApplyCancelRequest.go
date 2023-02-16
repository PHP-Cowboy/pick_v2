package domain

type TaobaoWlbWaybillICancelWaybillApplyCancelRequest struct {
	/*
	   CP快递公司编码     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   交易订单列表     */
	TradeOrderList *[]string `json:"trade_order_list,omitempty" `

	/*
	   面单使用者编号     */
	RealUserId *int64 `json:"real_user_id,omitempty" `

	/*
	   电子面单号码     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   ERP订单号或包裹号 defalutValue:-1    */
	PackageId *string `json:"package_id,omitempty" `
}

func (s *TaobaoWlbWaybillICancelWaybillApplyCancelRequest) SetCpCode(v string) *TaobaoWlbWaybillICancelWaybillApplyCancelRequest {
	s.CpCode = &v
	return s
}
func (s *TaobaoWlbWaybillICancelWaybillApplyCancelRequest) SetTradeOrderList(v []string) *TaobaoWlbWaybillICancelWaybillApplyCancelRequest {
	s.TradeOrderList = &v
	return s
}
func (s *TaobaoWlbWaybillICancelWaybillApplyCancelRequest) SetRealUserId(v int64) *TaobaoWlbWaybillICancelWaybillApplyCancelRequest {
	s.RealUserId = &v
	return s
}
func (s *TaobaoWlbWaybillICancelWaybillApplyCancelRequest) SetWaybillCode(v string) *TaobaoWlbWaybillICancelWaybillApplyCancelRequest {
	s.WaybillCode = &v
	return s
}
func (s *TaobaoWlbWaybillICancelWaybillApplyCancelRequest) SetPackageId(v string) *TaobaoWlbWaybillICancelWaybillApplyCancelRequest {
	s.PackageId = &v
	return s
}
