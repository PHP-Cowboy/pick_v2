package domain

type CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse struct {
	/*
	   面单信息     */
	PrintData *string `json:"print_data,omitempty" `

	/*
	   面单号     */
	WaybillCode *string `json:"waybill_code,omitempty" `
}

func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse) SetPrintData(v string) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse {
	s.PrintData = &v
	return s
}
func (s *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse) SetWaybillCode(v string) *CainiaoWaybillIiQueryByTradecodeWaybillCloudPrintResponse {
	s.WaybillCode = &v
	return s
}
