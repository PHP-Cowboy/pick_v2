package domain

type TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo struct {
	/*
	   挑拣规则（大头笔信息）     */
	ShortAddress *string `json:"short_address,omitempty" `

	/*
	   收货网点编码     */
	ConsigneeBranchName *string `json:"consignee_branch_name,omitempty" `

	/*
	   收货网点信息     */
	ConsigneeBranchCode *string `json:"consignee_branch_code,omitempty" `

	/*
	   --     */
	TradeOrderInfo *TaobaoWlbWaybillIFullupdateTradeOrderInfo `json:"trade_order_info,omitempty" `

	/*
	   --     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   集包地、目的地中心名称     */
	PackageCenterName *string `json:"package_center_name,omitempty" `

	/*
	   集包地、目的地中心代码。打印时根据该 code 生成目的地中心的条码，条码生成的算法与对应的电子面单条码一致     */
	PackageCenterCode *string `json:"package_center_code,omitempty" `
}

func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetShortAddress(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.ShortAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetConsigneeBranchName(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.ConsigneeBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetConsigneeBranchCode(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.ConsigneeBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetTradeOrderInfo(v TaobaoWlbWaybillIFullupdateTradeOrderInfo) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.TradeOrderInfo = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetWaybillCode(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.WaybillCode = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetPackageCenterName(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.PackageCenterName = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo) SetPackageCenterCode(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyUpdateInfo {
	s.PackageCenterCode = &v
	return s
}
