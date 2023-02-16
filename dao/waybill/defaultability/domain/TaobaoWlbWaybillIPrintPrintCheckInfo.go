package domain

type TaobaoWlbWaybillIPrintPrintCheckInfo struct {
	/*
	   收\发货地址     */
	ConsigneeAddress *TaobaoWlbWaybillIPrintWaybillAddress `json:"consignee_address,omitempty" `

	/*
	   收货网点编码     */
	ConsigneeBranchCode *string `json:"consignee_branch_code,omitempty" `

	/*
	   收货网点信息     */
	ConsigneeBranchName *string `json:"consignee_branch_name,omitempty" `

	/*
	   收件人姓名     */
	ConsigneeName *string `json:"consignee_name,omitempty" `

	/*
	   consigneePhone     */
	ConsigneePhone *string `json:"consignee_phone,omitempty" `

	/*
	   物流服务能力集合     */
	LogisticsServiceList *[]TaobaoWlbWaybillIPrintLogisticsService `json:"logistics_service_list,omitempty" `

	/*
	   集包地、目的地中心代码。打 印时根据该 code 生成目的地 中心的条码，条码生成的算法 与对应的电子面单条码一致     */
	PackageCenterCode *string `json:"package_center_code,omitempty" `

	/*
	   集包地、目的地中心名称     */
	PackageCenterName *string `json:"package_center_name,omitempty" `

	/*
	   打标设置字段，直接传给ali-lodop。不用管具体含义。     */
	PrintConfig *string `json:"print_config,omitempty" `

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
	   收\发货地址     */
	ShippingAddress *TaobaoWlbWaybillIPrintWaybillAddress `json:"shipping_address,omitempty" `

	/*
	   发货网点编码     */
	ShippingBranchCode *string `json:"shipping_branch_code,omitempty" `

	/*
	   发货网点信息     */
	ShippingBranchName *string `json:"shipping_branch_name,omitempty" `

	/*
	   拣货规则（大头笔信息）     */
	ShortAddress *string `json:"short_address,omitempty" `

	/*
	   使用者ID     */
	RealUserId *int64 `json:"real_user_id,omitempty" `

	/*
	   包裹体积 单位为ML(毫升)或立方厘米 defalutValue:0    */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   电子面单单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   包裹重量 单位为G(克) defalutValue:0    */
	Weight *int64 `json:"weight,omitempty" `
}

func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetConsigneeAddress(v TaobaoWlbWaybillIPrintWaybillAddress) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ConsigneeAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetConsigneeBranchCode(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ConsigneeBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetConsigneeBranchName(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ConsigneeBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetConsigneeName(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ConsigneeName = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetConsigneePhone(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ConsigneePhone = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetLogisticsServiceList(v []TaobaoWlbWaybillIPrintLogisticsService) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.LogisticsServiceList = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetPackageCenterCode(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.PackageCenterCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetPackageCenterName(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.PackageCenterName = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetPrintConfig(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.PrintConfig = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetProductType(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ProductType = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetSendName(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.SendName = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetSendPhone(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.SendPhone = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetShippingAddress(v TaobaoWlbWaybillIPrintWaybillAddress) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ShippingAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetShippingBranchCode(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ShippingBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetShippingBranchName(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ShippingBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetShortAddress(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.ShortAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetRealUserId(v int64) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.RealUserId = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetVolume(v int64) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.Volume = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetWaybillCode(v string) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.WaybillCode = &v
	return s
}
func (s *TaobaoWlbWaybillIPrintPrintCheckInfo) SetWeight(v int64) *TaobaoWlbWaybillIPrintPrintCheckInfo {
	s.Weight = &v
	return s
}
