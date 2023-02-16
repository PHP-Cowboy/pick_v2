package domain

type CainiaoWaybillIiProductWaybillServiceType struct {
	/*
	   code     */
	Code *string `json:"code,omitempty" `

	/*
	   name     */
	Name *string `json:"name,omitempty" `
}

func (s *CainiaoWaybillIiProductWaybillServiceType) SetCode(v string) *CainiaoWaybillIiProductWaybillServiceType {
	s.Code = &v
	return s
}
func (s *CainiaoWaybillIiProductWaybillServiceType) SetName(v string) *CainiaoWaybillIiProductWaybillServiceType {
	s.Name = &v
	return s
}
