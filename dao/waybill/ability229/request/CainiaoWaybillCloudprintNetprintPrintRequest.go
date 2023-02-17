package request

import (
	"pick_v2/dao/waybill/ability229/domain"
	"pick_v2/dao/waybill/util"
)

type CainiaoWaybillCloudprintNetprintPrintRequest struct {
	/*
	   请求     */
	PrinterPrintData *domain.CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest `json:"printer_print_data" required:"true" `
}

func (s *CainiaoWaybillCloudprintNetprintPrintRequest) SetPrinterPrintData(v domain.CainiaoWaybillCloudprintNetprintPrintCloudPrinterPrintRequest) *CainiaoWaybillCloudprintNetprintPrintRequest {
	s.PrinterPrintData = &v
	return s
}

func (req *CainiaoWaybillCloudprintNetprintPrintRequest) ToMap() map[string]interface{} {
	paramMap := make(map[string]interface{})
	if req.PrinterPrintData != nil {
		paramMap["printer_print_data"] = util.ConvertStruct(*req.PrinterPrintData)
	}
	return paramMap
}

func (req *CainiaoWaybillCloudprintNetprintPrintRequest) ToFileMap() map[string]interface{} {
	fileMap := make(map[string]interface{})
	return fileMap
}
