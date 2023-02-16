package domain

type TaobaoWlbWaybillIQuerydetailPackageItem struct {
	/*
	   商品数量     */
	Count *int64 `json:"count,omitempty" `

	/*
	   商品名称     */
	ItemName *string `json:"item_name,omitempty" `
}

func (s *TaobaoWlbWaybillIQuerydetailPackageItem) SetCount(v int64) *TaobaoWlbWaybillIQuerydetailPackageItem {
	s.Count = &v
	return s
}
func (s *TaobaoWlbWaybillIQuerydetailPackageItem) SetItemName(v string) *TaobaoWlbWaybillIQuerydetailPackageItem {
	s.ItemName = &v
	return s
}
