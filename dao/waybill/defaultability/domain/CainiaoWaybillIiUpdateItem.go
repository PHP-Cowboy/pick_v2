package domain

type CainiaoWaybillIiUpdateItem struct {
	/*
	   数量 defalutValue:0    */
	Count *int64 `json:"count,omitempty" `

	/*
	   名称     */
	Name *string `json:"name,omitempty" `
}

func (s *CainiaoWaybillIiUpdateItem) SetCount(v int64) *CainiaoWaybillIiUpdateItem {
	s.Count = &v
	return s
}
func (s *CainiaoWaybillIiUpdateItem) SetName(v string) *CainiaoWaybillIiUpdateItem {
	s.Name = &v
	return s
}
