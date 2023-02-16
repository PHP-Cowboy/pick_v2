package domain

type CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest struct {
	/*
	   物流公司CODE     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   物流服务内容<a href="http://open.taobao.com/doc2/detail.htm?spm=a219a.7629140.0.0.eK8aZm&treeId=17&articleId=26765&docType=2">链接</a>     */
	LogisticsServices *string `json:"logistics_services,omitempty" `

	/*
	   包裹信息     */
	PackageInfo *CainiaoWaybillIiUpdatePackageInfoDto `json:"package_info,omitempty" `

	/*
	   收件信息     */
	Recipient *CainiaoWaybillIiUpdateUserInfoDto `json:"recipient,omitempty" `

	/*
	   发件信息     */
	Sender *CainiaoWaybillIiUpdateUserInfoDto `json:"sender,omitempty" `

	/*
	   模板URL     */
	TemplateUrl *string `json:"template_url,omitempty" `

	/*
	   面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   请求表示id     */
	ObjectId *string `json:"object_id,omitempty" `
}

func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetCpCode(v string) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.CpCode = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetLogisticsServices(v string) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.LogisticsServices = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetPackageInfo(v CainiaoWaybillIiUpdatePackageInfoDto) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.PackageInfo = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetRecipient(v CainiaoWaybillIiUpdateUserInfoDto) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.Recipient = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetSender(v CainiaoWaybillIiUpdateUserInfoDto) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.Sender = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetTemplateUrl(v string) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.TemplateUrl = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetWaybillCode(v string) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.WaybillCode = &v
	return s
}
func (s *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest) SetObjectId(v string) *CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest {
	s.ObjectId = &v
	return s
}
