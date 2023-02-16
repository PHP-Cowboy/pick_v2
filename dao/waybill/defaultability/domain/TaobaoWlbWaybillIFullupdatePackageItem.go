package domain

type TaobaoWlbWaybillIFullupdatePackageItem struct {
	/*
	   商品数量     */
	Count *int64 `json:"count,omitempty" `

	/*
	   商品名称     */
	ItemName *string `json:"item_name,omitempty" `
}

func (s *TaobaoWlbWaybillIFullupdatePackageItem) SetCount(v int64) *TaobaoWlbWaybillIFullupdatePackageItem {
	s.Count = &v
	return s
}
func (s *TaobaoWlbWaybillIFullupdatePackageItem) SetItemName(v string) *TaobaoWlbWaybillIFullupdatePackageItem {
	s.ItemName = &v
	return s
}
