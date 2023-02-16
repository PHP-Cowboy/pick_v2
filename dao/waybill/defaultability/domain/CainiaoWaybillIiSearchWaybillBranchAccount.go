package domain

type CainiaoWaybillIiSearchWaybillBranchAccount struct {
	/*
	   已用面单数量     */
	AllocatedQuantity *int64 `json:"allocated_quantity,omitempty" `

	/*
	   网点Code     */
	BranchCode *string `json:"branch_code,omitempty" `

	/*
	   网点名称     */
	BranchName *string `json:"branch_name,omitempty" `

	/*
	   网点状态     */
	BranchStatus *int64 `json:"branch_status,omitempty" `

	/*
	   取消的面单总数     */
	CancelQuantity *int64 `json:"cancel_quantity,omitempty" `

	/*
	   已经打印的面单总数     */
	PrintQuantity *int64 `json:"print_quantity,omitempty" `

	/*
	   电子面单余额数量     */
	Quantity *int64 `json:"quantity,omitempty" `

	/*
	   当前网点下的发货地址     */
	ShippAddressCols *[]CainiaoWaybillIiSearchAddressDto `json:"shipp_address_cols,omitempty" `

	/*
	   可用的服务信息列表     */
	ServiceInfoCols *[]CainiaoWaybillIiSearchServiceInfoDto `json:"service_info_cols,omitempty" `

	/*
	   号段信息     */
	SegmentCode *string `json:"segment_code,omitempty" `

	/*
	   品牌code     */
	BrandCode *string `json:"brand_code,omitempty" `

	/*
	   月结卡号列表     */
	CustomerCodeList *[]string `json:"customer_code_list,omitempty" `

	/*
	   月结卡号map，key为shipp_address_cols.waybill_address_id,value为月结卡号。jsonString     */
	CustomerCodeMap *string `json:"customer_code_map,omitempty" `
}

func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetAllocatedQuantity(v int64) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.AllocatedQuantity = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetBranchCode(v string) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.BranchCode = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetBranchName(v string) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.BranchName = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetBranchStatus(v int64) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.BranchStatus = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetCancelQuantity(v int64) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.CancelQuantity = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetPrintQuantity(v int64) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.PrintQuantity = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetQuantity(v int64) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.Quantity = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetShippAddressCols(v []CainiaoWaybillIiSearchAddressDto) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.ShippAddressCols = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetServiceInfoCols(v []CainiaoWaybillIiSearchServiceInfoDto) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.ServiceInfoCols = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetSegmentCode(v string) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.SegmentCode = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetBrandCode(v string) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.BrandCode = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetCustomerCodeList(v []string) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.CustomerCodeList = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillBranchAccount) SetCustomerCodeMap(v string) *CainiaoWaybillIiSearchWaybillBranchAccount {
	s.CustomerCodeMap = &v
	return s
}
