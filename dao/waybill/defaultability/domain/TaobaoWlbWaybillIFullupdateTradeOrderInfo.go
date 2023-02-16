package domain

type TaobaoWlbWaybillIFullupdateTradeOrderInfo struct {
	/*
	   收货人地址     */
	ConsigneeAddress *TaobaoWlbWaybillIFullupdateWaybillAddress `json:"consignee_address,omitempty" `

	/*
	   收货人     */
	ConsigneeName *string `json:"consignee_name,omitempty" `

	/*
	   收货人联系方式     */
	ConsigneePhone *string `json:"consignee_phone,omitempty" `

	/*
	   物流服务能力集合     */
	LogisticsServiceList *[]TaobaoWlbWaybillIFullupdateLogisticsService `json:"logistics_service_list,omitempty" `

	/*
	   订单渠道     */
	OrderChannelsType *string `json:"order_channels_type,omitempty" `

	/*
	   包裹中的商品类型     */
	PackageItems *[]TaobaoWlbWaybillIFullupdatePackageItem `json:"package_items,omitempty" `

	/*
	   快递服务产品类型编码     */
	ProductType *string `json:"product_type,omitempty" `

	/*
	   发货人姓名     */
	SendName *string `json:"send_name,omitempty" `

	/*
	   发货人联系方式     */
	SendPhone *string `json:"send_phone,omitempty" `

	/*
	   交易订单列表     */
	TradeOrderList *[]string `json:"trade_order_list,omitempty" `

	/*
	   使用者ID     */
	RealUserId *int64 `json:"real_user_id,omitempty" `

	/*
	   包裹体积（立方厘米）     */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   包裹重量（克）     */
	Weight *int64 `json:"weight,omitempty" `

	/*
	   包裹号(或者ERP订单号)     */
	PackageId *string `json:"package_id,omitempty" `
}

func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetConsigneeAddress(v TaobaoWlbWaybillIFullupdateWaybillAddress) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.ConsigneeAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetConsigneeName(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.ConsigneeName = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetConsigneePhone(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.ConsigneePhone = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetLogisticsServiceList(v []TaobaoWlbWaybillIFullupdateLogisticsService) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.LogisticsServiceList = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetOrderChannelsType(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.OrderChannelsType = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetPackageItems(v []TaobaoWlbWaybillIFullupdatePackageItem) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.PackageItems = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetProductType(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.ProductType = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetSendName(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.SendName = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetSendPhone(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.SendPhone = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetTradeOrderList(v []string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.TradeOrderList = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetRealUserId(v int64) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.RealUserId = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetVolume(v int64) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.Volume = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetWeight(v int64) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.Weight = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdateTradeOrderInfo) SetPackageId(v string) *TaobaoWlbWaybillIFullupdateTradeOrderInfo {
	s.PackageId = &v
	return s
}
