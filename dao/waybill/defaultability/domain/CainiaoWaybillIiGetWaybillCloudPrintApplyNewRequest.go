package domain

type CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest struct {
	/*
	   <a href="http://open.taobao.com/doc2/detail.htm?spm=a219a.7629140.0.0.8cf9Nj&treeId=17&articleId=105085&docType=1#1">物流公司Code</a>，长度小于20     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   目前仅顺丰支持此字段，传入快递产品编码     */
	ProductCode *string `json:"product_code,omitempty" `

	/*
	   发货人信息     */
	Sender *CainiaoWaybillIiGetUserInfoDto `json:"sender,omitempty" `

	/*
	   请求面单信息，数量限制为10     */
	TradeOrderInfoDtos *[]CainiaoWaybillIiGetTradeOrderInfoDto `json:"trade_order_info_dtos,omitempty" `

	/*
	   仓code， 仓库WMS系统对接落地配业务，其它场景请不要使用     */
	StoreCode *string `json:"store_code,omitempty" `

	/*
	   配送资源code， 仓库WMS系统对接落地配业务，其它场景请不要使用     */
	ResourceCode *string `json:"resource_code,omitempty" `

	/*
	   是否使用智分宝预分拣， 仓库WMS系统对接落地配业务，其它场景请不要使用 defalutValue:false    */
	DmsSorting *bool `json:"dms_sorting,omitempty" `

	/*
	   订单上是否带3PLtiming属性, 该属性需要严格与订单上属性保持一致，如果不确定，请使用默认false。 defalutValue:false    */
	ThreePlTiming *bool `json:"three_pl_timing,omitempty" `

	/*
	   设定取号返回的云打印报文是否加密 defalutValue:false    */
	NeedEncrypt *bool `json:"need_encrypt,omitempty" `

	/*
	   快递公司支持一票多件，快运公司子母件请勿使用该参数 defalutValue:false    */
	MultiPackagesShipment *bool `json:"multi_packages_shipment,omitempty" `

	/*
	   品牌编码     */
	BrandCode *string `json:"brand_code,omitempty" `

	/*
	   扩展参数     */
	ExtraInfo *string `json:"extra_info,omitempty" `

	/*
	   月结卡号     */
	CustomerCode *string `json:"customer_code,omitempty" `

	/*
	   是否预约上门     */
	CallDoorPickUp *bool `json:"call_door_pick_up,omitempty" `

	/*
	   预约上门收件时间     */
	DoorPickUpTime *string `json:"door_pick_up_time,omitempty" `

	/*
	   预约上门截止时间     */
	DoorPickUpEndTime *string `json:"door_pick_up_end_time,omitempty" `

	/*
	   寄件网点编码     */
	ShippingBranchCode *string `json:"shipping_branch_code,omitempty" `
}

func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetCpCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.CpCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetProductCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.ProductCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetSender(v CainiaoWaybillIiGetUserInfoDto) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.Sender = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetTradeOrderInfoDtos(v []CainiaoWaybillIiGetTradeOrderInfoDto) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.TradeOrderInfoDtos = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetStoreCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.StoreCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetResourceCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.ResourceCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetDmsSorting(v bool) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.DmsSorting = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetThreePlTiming(v bool) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.ThreePlTiming = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetNeedEncrypt(v bool) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.NeedEncrypt = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetMultiPackagesShipment(v bool) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.MultiPackagesShipment = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetBrandCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.BrandCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetExtraInfo(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.ExtraInfo = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetCustomerCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.CustomerCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetCallDoorPickUp(v bool) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.CallDoorPickUp = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetDoorPickUpTime(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.DoorPickUpTime = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetDoorPickUpEndTime(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.DoorPickUpEndTime = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest) SetShippingBranchCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest {
	s.ShippingBranchCode = &v
	return s
}
