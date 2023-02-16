package domain

type CainiaoWaybillIiGetWaybillCloudPrintResponse struct {
	/*
	   请求id     */
	ObjectId *string `json:"object_id,omitempty" `

	/*
	   云打印内容（encryptedData表示加密结果，data表示非加密结果）;模板内容,具体解释见<a href="http://open.taobao.com/doc2/detail.htm?spm=a219a.7629140.0.0.8cf9Nj&treeId=17&articleId=105085&docType=1#12">链接</a>     */
	PrintData *string `json:"print_data,omitempty" `

	/*
	   面单号, 子母件模式下为子面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   子母件中的母单号，当为子母件模式时，需要此单号为实际挂载物流详情的单号，需要使用此单号进行发货，查询物流详情，非子母件，此字段为空     */
	ParentWaybillCode *string `json:"parent_waybill_code,omitempty" `

	/*
	   拓展信息，特殊场景下使用     */
	ExtraInfo *string `json:"extra_info,omitempty" `

	/*
	   cp_code，跟入参传入的cp_code保持一致     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   真实取号的cp_code，在淘特官方物流等虚拟cp的场景中real_cp_code和入参中的cp_code不一样，其他场景二者一样     */
	RealCpCode *string `json:"real_cp_code,omitempty" `
}

func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetObjectId(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.ObjectId = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetPrintData(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.PrintData = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetWaybillCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.WaybillCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetParentWaybillCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.ParentWaybillCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetExtraInfo(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.ExtraInfo = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetCpCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.CpCode = &v
	return s
}
func (s *CainiaoWaybillIiGetWaybillCloudPrintResponse) SetRealCpCode(v string) *CainiaoWaybillIiGetWaybillCloudPrintResponse {
	s.RealCpCode = &v
	return s
}
