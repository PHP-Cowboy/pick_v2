package domain

type CainiaoWaybillIiSearchWaybillApplySubscriptionInfo struct {
	/*
	   CP网点信息及对应的商家的发货信息     */
	BranchAccountCols *[]CainiaoWaybillIiSearchWaybillBranchAccount `json:"branch_account_cols,omitempty" `

	/*
	   物流服务商ID     */
	CpCode *string `json:"cp_code,omitempty" `

	/*
	   物流服务商业务类型 1：直营 2：加盟 3：落地配 4：直营带网点     */
	CpType *int64 `json:"cp_type,omitempty" `
}

func (s *CainiaoWaybillIiSearchWaybillApplySubscriptionInfo) SetBranchAccountCols(v []CainiaoWaybillIiSearchWaybillBranchAccount) *CainiaoWaybillIiSearchWaybillApplySubscriptionInfo {
	s.BranchAccountCols = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillApplySubscriptionInfo) SetCpCode(v string) *CainiaoWaybillIiSearchWaybillApplySubscriptionInfo {
	s.CpCode = &v
	return s
}
func (s *CainiaoWaybillIiSearchWaybillApplySubscriptionInfo) SetCpType(v int64) *CainiaoWaybillIiSearchWaybillApplySubscriptionInfo {
	s.CpType = &v
	return s
}
