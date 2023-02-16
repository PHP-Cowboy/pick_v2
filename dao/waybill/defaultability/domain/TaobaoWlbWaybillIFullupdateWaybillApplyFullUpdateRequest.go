package domain

type TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest struct {
	/*
	   收件人姓名     */
	ConsigneeName *string `json:"consignee_name,omitempty" `

	/*
	   收件人电话     */
	ConsigneePhone *string `json:"consignee_phone,omitempty" `

	/*
	   快递服务商CODE     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   物流服务能力集合     */
	LogisticsServiceList *[]TaobaoWlbWaybillIFullupdateLogisticsService `json:"logistics_service_list,omitempty" `

	/*
	   快递服务产品类型编码     */
	ProductType *string `json:"product_type,omitempty" `

	/*
	   发件人姓名     */
	SendName *string `json:"send_name,omitempty" `

	/*
	   发件人联系方式     */
	SendPhone *string `json:"send_phone,omitempty" `

	/*
	   交易订单号（组合表示合并订单）     */
	TradeOrderList *[]string `json:"trade_order_list,omitempty" `

	/*
	   使用者ID     */
	RealUserId *int64 `json:"real_user_id,omitempty" `

	/*
	   包裹重量 单位为G(克) defalutValue:0    */
	Weight *int64 `json:"weight,omitempty" `

	/*
	   收\发货地址     */
	ConsigneeAddress *TaobaoWlbWaybillIFullupdateWaybillAddress `json:"consignee_address,omitempty" `

	/*
	   订单渠道类型     */
	OrderChannelsType *string `json:"order_channels_type,omitempty" `

	/*
	   包裹里面的商品名称     */
	PackageItems *[]TaobaoWlbWaybillIFullupdatePackageItem `json:"package_items,omitempty" `

	/*
	   包裹体积 单位为ML(毫升)或立方厘米 defalutValue:0    */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   电子面单单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   ERP 订单号或包裹号 defalutValue:-1    */
	PackageId *string `json:"package_id,omitempty" `
}

func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetConsigneeName(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.ConsigneeName = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetConsigneePhone(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.ConsigneePhone = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetCpCode(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.CpCode = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetLogisticsServiceList(v []TaobaoWlbWaybillIFullupdateLogisticsService) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.LogisticsServiceList = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetProductType(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.ProductType = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetSendName(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.SendName = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetSendPhone(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.SendPhone = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetTradeOrderList(v []string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.TradeOrderList = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetRealUserId(v int64) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.RealUserId = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetWeight(v int64) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.Weight = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetConsigneeAddress(v TaobaoWlbWaybillIFullupdateWaybillAddress) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.ConsigneeAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetOrderChannelsType(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.OrderChannelsType = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetPackageItems(v []TaobaoWlbWaybillIFullupdatePackageItem) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.PackageItems = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetVolume(v int64) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.Volume = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetWaybillCode(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.WaybillCode = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest) SetPackageId(v string) *TaobaoWlbWaybillIFullupdateWaybillApplyFullUpdateRequest {
	s.PackageId = &v
	return s
}
