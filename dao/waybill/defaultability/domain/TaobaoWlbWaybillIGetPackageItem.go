package domain

type TaobaoWlbWaybillIGetPackageItem struct {
	/*
	   商品数量     */
	Count *int64 `json:"count,omitempty" `

	/*
	   商品名称     */
	ItemName *string `json:"item_name,omitempty" `
}

func (s *TaobaoWlbWaybillIGetPackageItem) SetCount(v int64) *TaobaoWlbWaybillIGetPackageItem {
	s.Count = &v
	return s
}
func (s *TaobaoWlbWaybillIGetPackageItem) SetItemName(v string) *TaobaoWlbWaybillIGetPackageItem {
	s.ItemName = &v
	return s
}
