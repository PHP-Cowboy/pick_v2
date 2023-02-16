package domain

type TaobaoWlbWaybillIGetTradeOrderInfo struct {
	/*
	   收\发货地址     */
	ConsigneeAddress *TaobaoWlbWaybillIGetWaybillAddress `json:"consignee_address,omitempty" `

	/*
	   收货人     */
	ConsigneeName *string `json:"consignee_name,omitempty" `

	/*
	   收货人联系方式     */
	ConsigneePhone *string `json:"consignee_phone,omitempty" `

	/*
	   物流服务能力集合     */
	LogisticsServiceList *[]TaobaoWlbWaybillIGetLogisticsService `json:"logistics_service_list,omitempty" `

	/*
	   订单渠道     */
	OrderChannelsType *string `json:"order_channels_type,omitempty" `

	/*
	   包裹里面的商品名称     */
	PackageItems *[]TaobaoWlbWaybillIGetPackageItem `json:"package_items,omitempty" `

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
	   包裹体积（立方厘米） defalutValue:0    */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   包裹重量（克） defalutValue:0    */
	Weight *int64 `json:"weight,omitempty" `

	/*
	   包裹号(或者ERP订单号) defalutValue:-1    */
	PackageId *string `json:"package_id,omitempty" `
}

func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetConsigneeAddress(v TaobaoWlbWaybillIGetWaybillAddress) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.ConsigneeAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetConsigneeName(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.ConsigneeName = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetConsigneePhone(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.ConsigneePhone = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetLogisticsServiceList(v []TaobaoWlbWaybillIGetLogisticsService) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.LogisticsServiceList = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetOrderChannelsType(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.OrderChannelsType = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetPackageItems(v []TaobaoWlbWaybillIGetPackageItem) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.PackageItems = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetProductType(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.ProductType = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetSendName(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.SendName = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetSendPhone(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.SendPhone = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetTradeOrderList(v []string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.TradeOrderList = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetRealUserId(v int64) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.RealUserId = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetVolume(v int64) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.Volume = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetWeight(v int64) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.Weight = &v
	return s
}
func (s *TaobaoWlbWaybillIGetTradeOrderInfo) SetPackageId(v string) *TaobaoWlbWaybillIGetTradeOrderInfo {
	s.PackageId = &v
	return s
}
