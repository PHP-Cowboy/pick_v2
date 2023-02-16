package domain

type CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse struct {
	/*
	   面单信息     */
	PrintData *string `json:"print_data,omitempty" `

	/*
	   面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `

	/*
	   快运子母件中的母单号     */
	ParentWaybillCode *string `json:"parent_waybill_code,omitempty" `
}

func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse) SetPrintData(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse {
	s.PrintData = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse) SetWaybillCode(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse {
	s.WaybillCode = &v
	return s
}
func (s *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse) SetParentWaybillCode(v string) *CainiaoWaybillIiQueryByWaybillcodeWaybillCloudPrintResponse {
	s.ParentWaybillCode = &v
	return s
}
