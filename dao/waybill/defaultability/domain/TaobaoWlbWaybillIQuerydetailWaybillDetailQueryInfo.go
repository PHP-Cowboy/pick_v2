package domain

import (
	"topsdk/util"
)

type TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo struct {
	/*
	   交易订单列表     */
	TradeOrderList *[]string `json:"trade_order_list,omitempty" `

	/*
	   收货人地址     */
	ConsigneeAddress *TaobaoWlbWaybillIQuerydetailWaybillAddress `json:"consignee_address,omitempty" `

	/*
	   包裹对应的派件（收件）物流服务商网点（分支机构）代码     */
	ConsigneeBranchCode *string `json:"consignee_branch_code,omitempty" `

	/*
	   包裹对应的派件（收件）物流服务商网点（分支机构）名称     */
	ConsigneeBranchName *string `json:"consignee_branch_name,omitempty" `

	/*
	   收件人姓名     */
	ConsigneeName *string `json:"consignee_name,omitempty" `

	/*
	   收件人联系方式     */
	ConsigneePhone *string `json:"consignee_phone,omitempty" `

	/*
	   物流商编码CODE     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   创建时间     */
	CreateTime *util.LocalTime `json:"create_time,omitempty" `

	/*
	   最后一次打印时间     */
	LastPrintTime *util.LocalTime `json:"last_print_time,omitempty" `

	/*
	   物流服务能力集合     */
	LogisticsServiceList *[]TaobaoWlbWaybillIQuerydetailLogisticsService `json:"logistics_service_list,omitempty" `

	/*
	   集包地、目的地中心代码。打印时根据该 code 生成目的地中心的条码，条码生成的算法与对应的电子面单条码一致     */
	PackageCenterCode *string `json:"package_center_code,omitempty" `

	/*
	   集包地、目的地中心名称     */
	PackageCenterName *string `json:"package_center_name,omitempty" `

	/*
	   ERP订单号或包裹号     */
	PackageId *string `json:"package_id,omitempty" `

	/*
	   包裹里面的商品类型     */
	PackageItems *[]TaobaoWlbWaybillIQuerydetailPackageItem `json:"package_items,omitempty" `

	/*
	   揽收时间     */
	PickupTime *util.LocalTime `json:"pickup_time,omitempty" `

	/*
	   打印配置项     */
	PrintConfig *string `json:"print_config,omitempty" `

	/*
	   打印次数     */
	PrintCount *int64 `json:"print_count,omitempty" `

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
	   发货地址     */
	ShippingAddress *TaobaoWlbWaybillIQuerydetailWaybillAddress `json:"shipping_address,omitempty" `

	/*
	   发货网点编码     */
	ShippingBranchCode *string `json:"shipping_branch_code,omitempty" `

	/*
	   发货网点信息     */
	ShippingBranchName *string `json:"shipping_branch_name,omitempty" `

	/*
	   大头笔信息     */
	ShortAddress *string `json:"short_address,omitempty" `

	/*
	   签收时间     */
	SignTime *util.LocalTime `json:"sign_time,omitempty" `

	/*
	   面单状态     */
	Status *int64 `json:"status,omitempty" `

	/*
	   使用者ID     */
	RealUserId *int64 `json:"real_user_id,omitempty" `

	/*
	   包裹重量 单位为G(克)     */
	Volume *int64 `json:"volume,omitempty" `

	/*
	   电子面单信息     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   包裹体积 单位为ML(毫升)或立方厘米     */
	Weight *int64 `json:"weight,omitempty" `
}

func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetTradeOrderList(v []string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.TradeOrderList = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetConsigneeAddress(v TaobaoWlbWaybillIQuerydetailWaybillAddress) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ConsigneeAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetConsigneeBranchCode(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ConsigneeBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetConsigneeBranchName(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ConsigneeBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetConsigneeName(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ConsigneeName = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetConsigneePhone(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ConsigneePhone = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetCpCode(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.CpCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetCreateTime(v util.LocalTime) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.CreateTime = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetLastPrintTime(v util.LocalTime) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.LastPrintTime = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetLogisticsServiceList(v []TaobaoWlbWaybillIQuerydetailLogisticsService) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.LogisticsServiceList = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPackageCenterCode(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PackageCenterCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPackageCenterName(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PackageCenterName = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPackageId(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PackageId = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPackageItems(v []TaobaoWlbWaybillIQuerydetailPackageItem) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PackageItems = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPickupTime(v util.LocalTime) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PickupTime = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPrintConfig(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PrintConfig = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetPrintCount(v int64) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.PrintCount = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetProductType(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ProductType = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetSendName(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.SendName = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetSendPhone(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.SendPhone = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetShippingAddress(v TaobaoWlbWaybillIQuerydetailWaybillAddress) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ShippingAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetShippingBranchCode(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ShippingBranchCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetShippingBranchName(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ShippingBranchName = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetShortAddress(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.ShortAddress = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetSignTime(v util.LocalTime) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.SignTime = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetStatus(v int64) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.Status = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetRealUserId(v int64) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.RealUserId = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetVolume(v int64) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.Volume = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetWaybillCode(v string) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.WaybillCode = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo) SetWeight(v int64) *TaobaoWlbWaybillIQuerydetailWaybillDetailQueryInfo {
	s.Weight = &v
	return s
}
