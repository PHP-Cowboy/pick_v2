package domain

type CainiaoWaybillIiGetTradeOrderInfoDto struct {
	/*
	   物流服务值（详见https://support-cnkuaidi.taobao.com/doc.htm#?docId=106156&docType=1，如无特殊服务请置空）     */
	LogisticsServices *string `json:"logistics_services,omitempty" `

	/*
	   <a href="http://open.taobao.com/docs/doc.htm?docType=1&articleId=105086&treeId=17&platformId=17#6">请求ID</a>     */
	ObjectId *string `json:"object_id,omitempty" `

	/*
	   订单信息     */
	OrderInfo *CainiaoWaybillIiGetOrderInfoDto `json:"order_info,omitempty" `

	/*
	   包裹信息     */
	PackageInfo *CainiaoWaybillIiGetPackageInfoDto `json:"package_info,omitempty" `

	/*
	   收件人信息     */
	Recipient *CainiaoWaybillIiGetRecipientInfoDto `json:"recipient,omitempty" `

	/*
	   云打印标准模板URL（组装云打印结果使用，值格式http://cloudprint.cainiao.com/template/standard/${模板ID}）     */
	TemplateUrl *string `json:"template_url,omitempty" `

	/*
	   使用者ID（使用电子面单账号的实际商家ID，如存在一个电子面单账号多个店铺使用时，请传入店铺的商家ID）     */
	UserId *int64 `json:"user_id,omitempty" `

	/*
	   带面单号模式取号，目前仅顺丰支持     */
	WaybillCode *string `json:"waybill_code,omitempty" `
}

func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetLogisticsServices(v string) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.LogisticsServices = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetObjectId(v string) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.ObjectId = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetOrderInfo(v CainiaoWaybillIiGetOrderInfoDto) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.OrderInfo = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetPackageInfo(v CainiaoWaybillIiGetPackageInfoDto) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.PackageInfo = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetRecipient(v CainiaoWaybillIiGetRecipientInfoDto) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.Recipient = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetTemplateUrl(v string) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.TemplateUrl = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetUserId(v int64) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.UserId = &v
	return s
}
func (s *CainiaoWaybillIiGetTradeOrderInfoDto) SetWaybillCode(v string) *CainiaoWaybillIiGetTradeOrderInfoDto {
	s.WaybillCode = &v
	return s
}
