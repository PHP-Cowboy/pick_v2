package domain

type TaobaoWlbWaybillIGetWaybillApplyNewInfo struct {
	/*
	   包裹对应的派件（收件）物流服务商网点（分支机构）代码     */
	ConsigneeBranchCode *string `json:"consignee_branch_code,omitempty" `

	/*
	   包裹对应的派件（收件）物流服务商网点（分支机构）名称     */
	ConsigneeBranchName *string `json:"consignee_branch_name,omitempty" `

	/*
	   集包地代码     */
	PackageCenterCode *string `json:"package_center_code,omitempty" `

	/*
	   集包地名称     */
	PackageCenterName *string `json:"package_center_name,omitempty" `

	/*
	   打印配置项，传给ali-print组件     */
	PrintConfig *string `json:"print_config,omitempty" `

	/*
	   面单号对应的物流服务商网点（分支机构）代码     */
	ShippingBranchCode *string `json:"shipping_branch_code,omitempty" `

	/*
	   面单号对于的物流服务商网点（分支机构）名称     */
	ShippingBranchName *string `json:"shipping_branch_name,omitempty" `

	/*
	   根据收货地址返回大头笔信息     */
	ShortAddress *string `json:"short_address,omitempty" `

	/*
	   面单对应的订单列     */
	TradeOrderInfo *TaobaoWlbWaybillIGetTradeOrderInfo `json:"trade_order_info,omitempty" `

	/*
	   返回的面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `
}

func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetConsigneeBranchCode(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.ConsigneeBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetConsigneeBranchName(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.ConsigneeBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetPackageCenterCode(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.PackageCenterCode = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetPackageCenterName(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.PackageCenterName = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetPrintConfig(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.PrintConfig = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetShippingBranchCode(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.ShippingBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetShippingBranchName(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.ShippingBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetShortAddress(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.ShortAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetTradeOrderInfo(v TaobaoWlbWaybillIGetTradeOrderInfo) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.TradeOrderInfo = &v
	return s
}
func (s *TaobaoWlbWaybillIGetWaybillApplyNewInfo) SetWaybillCode(v string) *TaobaoWlbWaybillIGetWaybillApplyNewInfo {
	s.WaybillCode = &v
	return s
}
