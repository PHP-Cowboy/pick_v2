package domain

type CainiaoWaybillIiGetItem struct {
	/*
	   数量     */
	Count *int64 `json:"count,omitempty" `

	/*
	   名称     */
	Name *string `json:"name,omitempty" `
}

func (s *CainiaoWaybillIiGetItem) SetCount(v int64) *CainiaoWaybillIiGetItem {
	s.Count = &v
	return s
}
func (s *CainiaoWaybillIiGetItem) SetName(v string) *CainiaoWaybillIiGetItem {
	s.Name = &v
	return s
}
