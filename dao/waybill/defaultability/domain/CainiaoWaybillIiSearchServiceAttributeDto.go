package domain

type CainiaoWaybillIiSearchServiceAttributeDto struct {
	/*
	   属性的值，用户实际传入的值     */
	AttributeCode *string `json:"attribute_code,omitempty" `

	/*
	   属性的名称，可以用于前端的展示     */
	AttributeName *string `json:"attribute_name,omitempty" `

	/*
	   属性的类型，可能值有 [number, string, enum]     */
	AttributeType *string `json:"attribute_type,omitempty" `

	/*
	   枚举类型的枚举值，key为用户选中的需要传值的数据，value为对应的描述，可以作为前端的展示     */
	TypeDesc *string `json:"type_desc,omitempty" `
}

func (s *CainiaoWaybillIiSearchServiceAttributeDto) SetAttributeCode(v string) *CainiaoWaybillIiSearchServiceAttributeDto {
	s.AttributeCode = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceAttributeDto) SetAttributeName(v string) *CainiaoWaybillIiSearchServiceAttributeDto {
	s.AttributeName = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceAttributeDto) SetAttributeType(v string) *CainiaoWaybillIiSearchServiceAttributeDto {
	s.AttributeType = &v
	return s
}
func (s *CainiaoWaybillIiSearchServiceAttributeDto) SetTypeDesc(v string) *CainiaoWaybillIiSearchServiceAttributeDto {
	s.TypeDesc = &v
	return s
}
