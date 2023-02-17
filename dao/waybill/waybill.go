package waybill

import (
	"fmt"
	"pick_v2/dao/waybill/ability229"
	ability229Domain "pick_v2/dao/waybill/ability229/domain"
	ability229req "pick_v2/dao/waybill/ability229/request"
	"pick_v2/dao/waybill/defaultability"
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/defaultability/request"
)

func NewAbility() *defaultability.Defaultability {
	client := NewDefaultTopClient(AppKey, AppSecret, GatewayUrl, 20000, 20000)
	return defaultability.NewDefaultability(&client)
}

func NewAbility229() *ability229.Ability229 {

	client := NewDefaultTopClient(AppKey, AppSecret, GatewayUrl, 20000, 20000)
	return ability229.NewAbility229(&client)
}

// 获取电子面单
func WaybillIiGet() {
	ability := NewAbility()

	cainiaoWaybillIiGetAddressDto := domain.CainiaoWaybillIiGetAddressDto{}
	cainiaoWaybillIiGetAddressDto.SetCity("北京市")
	cainiaoWaybillIiGetAddressDto.SetDetail("花家地社区卫生服务站")
	cainiaoWaybillIiGetAddressDto.SetDistrict("朝阳区")
	cainiaoWaybillIiGetAddressDto.SetProvince("北京")
	cainiaoWaybillIiGetAddressDto.SetTown("望京街道")
	cainiaoWaybillIiGetUserInfoDto := domain.CainiaoWaybillIiGetUserInfoDto{}
	cainiaoWaybillIiGetUserInfoDto.SetAddress(cainiaoWaybillIiGetAddressDto)
	cainiaoWaybillIiGetUserInfoDto.SetMobile("1326443654")
	cainiaoWaybillIiGetUserInfoDto.SetName("Bar")
	cainiaoWaybillIiGetUserInfoDto.SetPhone("057123222")
	cainiaoWaybillIiGetOrderInfoDto := domain.CainiaoWaybillIiGetOrderInfoDto{}
	cainiaoWaybillIiGetOrderInfoDto.SetOrderChannelsType("TB")
	// 1222221
	cainiaoWaybillIiGetOrderInfoDto.SetTradeOrderList([]string{})
	// 123456,456789
	cainiaoWaybillIiGetOrderInfoDto.SetOutTradeOrderList([]string{})
	// 12,34,56,78
	cainiaoWaybillIiGetOrderInfoDto.SetOutTradeSubOrderList([]string{})
	cainiaoWaybillIiGetItem := domain.CainiaoWaybillIiGetItem{}
	cainiaoWaybillIiGetItem.SetCount(1)
	cainiaoWaybillIiGetItem.SetName("衣服")
	cainiaoWaybillIiGetPackageInfoDto := domain.CainiaoWaybillIiGetPackageInfoDto{}
	cainiaoWaybillIiGetPackageInfoDto.SetId("1")
	//
	cainiaoWaybillIiGetPackageInfoDto.SetItems([]domain.CainiaoWaybillIiGetItem{})
	cainiaoWaybillIiGetPackageInfoDto.SetVolume(1)
	cainiaoWaybillIiGetPackageInfoDto.SetWeight(1)
	cainiaoWaybillIiGetPackageInfoDto.SetTotalPackagesCount(10)
	cainiaoWaybillIiGetPackageInfoDto.SetPackagingDescription("5纸3木2拖")
	cainiaoWaybillIiGetPackageInfoDto.SetGoodsDescription("服装")
	cainiaoWaybillIiGetPackageInfoDto.SetLength(30)
	cainiaoWaybillIiGetPackageInfoDto.SetWidth(30)
	cainiaoWaybillIiGetPackageInfoDto.SetHeight(50)
	cainiaoWaybillIiGetPackageInfoDto.SetGoodValue("34.3")
	cainiaoWaybillIiGetRecipientAddressDto := domain.CainiaoWaybillIiGetAddressDto{}
	cainiaoWaybillIiGetRecipientAddressDto.SetCity("北京市")
	cainiaoWaybillIiGetRecipientAddressDto.SetDetail("花家地社区卫生服务站")
	cainiaoWaybillIiGetRecipientAddressDto.SetDistrict("朝阳区")
	cainiaoWaybillIiGetRecipientAddressDto.SetProvince("北京")
	cainiaoWaybillIiGetRecipientAddressDto.SetTown("望京街道")

	cainiaoWaybillIiGetRecipientInfoDto := domain.CainiaoWaybillIiGetRecipientInfoDto{}
	cainiaoWaybillIiGetRecipientInfoDto.SetAddress(cainiaoWaybillIiGetRecipientAddressDto)
	cainiaoWaybillIiGetRecipientInfoDto.SetMobile("1326443654")
	cainiaoWaybillIiGetRecipientInfoDto.SetName("Bar")
	cainiaoWaybillIiGetRecipientInfoDto.SetPhone("057123222")
	cainiaoWaybillIiGetRecipientInfoDto.SetOaid("abcdefghijk")
	cainiaoWaybillIiGetRecipientInfoDto.SetTid("1527014522198024829")
	cainiaoWaybillIiGetRecipientInfoDto.SetCaid("As268woscee")
	cainiaoWaybillIiGetTradeOrderInfoDto := domain.CainiaoWaybillIiGetTradeOrderInfoDto{}
	cainiaoWaybillIiGetTradeOrderInfoDto.SetLogisticsServices("如不需要特殊服务，该值为空")
	cainiaoWaybillIiGetTradeOrderInfoDto.SetObjectId("1")
	cainiaoWaybillIiGetTradeOrderInfoDto.SetOrderInfo(cainiaoWaybillIiGetOrderInfoDto)
	cainiaoWaybillIiGetTradeOrderInfoDto.SetPackageInfo(cainiaoWaybillIiGetPackageInfoDto)
	cainiaoWaybillIiGetTradeOrderInfoDto.SetRecipient(cainiaoWaybillIiGetRecipientInfoDto)
	cainiaoWaybillIiGetTradeOrderInfoDto.SetTemplateUrl("http://cloudprint.cainiao.com/template/standard/101")
	cainiaoWaybillIiGetTradeOrderInfoDto.SetUserId(12)
	cainiaoWaybillIiGetTradeOrderInfoDto.SetWaybillCode("SF982933200")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest := domain.CainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest{}
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetCpCode("POSTB")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetProductCode("目前仅顺丰场景支持此字段，传入快递产品编码")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetSender(cainiaoWaybillIiGetUserInfoDto)
	//
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetTradeOrderInfoDtos([]domain.CainiaoWaybillIiGetTradeOrderInfoDto{})
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetStoreCode("553323")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetResourceCode("DISTRIBUTOR_978324")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetDmsSorting(false)
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetThreePlTiming(false)
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetNeedEncrypt(false)
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetMultiPackagesShipment(false)
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetBrandCode("FOP")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetExtraInfo(`{"isvClientCode":"ab12344"}`)
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetCustomerCode("adb123345")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetCallDoorPickUp(false)
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetDoorPickUpTime("2021-08-07 12:34:30")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetDoorPickUpEndTime("2021-08-07 12:34:30")
	cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest.SetShippingBranchCode("10001")

	req := request.CainiaoWaybillIiGetRequest{}
	req.SetParamWaybillCloudPrintApplyNewRequest(cainiaoWaybillIiGetWaybillCloudPrintApplyNewRequest)

	resp, err := ability.CainiaoWaybillIiGet(&req, UserSession)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Body)
	}
}

// 取消获取的电子面单号
func WaybillIiCancel() {

	ability := NewAbility()

	req := request.CainiaoWaybillIiCancelRequest{}
	req.SetCpCode("POSTB")
	req.SetWaybillCode("1111")

	resp, err := ability.CainiaoWaybillIiCancel(&req, UserSession)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Body)
	}
}

// 电子面单更新
func WaybillIiUpdate() {

	ability := NewAbility()

	cainiaoWaybillIiUpdateItem := domain.CainiaoWaybillIiUpdateItem{}
	cainiaoWaybillIiUpdateItem.SetCount(1)
	cainiaoWaybillIiUpdateItem.SetName("鞋子")
	cainiaoWaybillIiUpdatePackageInfoDto := domain.CainiaoWaybillIiUpdatePackageInfoDto{}
	//
	cainiaoWaybillIiUpdatePackageInfoDto.SetItems([]domain.CainiaoWaybillIiUpdateItem{})
	cainiaoWaybillIiUpdatePackageInfoDto.SetVolume(1)
	cainiaoWaybillIiUpdatePackageInfoDto.SetWeight(1)
	cainiaoWaybillIiUpdateAddressDto := domain.CainiaoWaybillIiUpdateAddressDto{}
	cainiaoWaybillIiUpdateAddressDto.SetCity("杭州市")
	cainiaoWaybillIiUpdateAddressDto.SetDetail("西溪园区")
	cainiaoWaybillIiUpdateAddressDto.SetDistrict("余杭区")
	cainiaoWaybillIiUpdateAddressDto.SetProvince("浙江省")
	cainiaoWaybillIiUpdateAddressDto.SetTown("文一西路")
	cainiaoWaybillIiUpdateUserInfoDto := domain.CainiaoWaybillIiUpdateUserInfoDto{}
	cainiaoWaybillIiUpdateUserInfoDto.SetAddress(cainiaoWaybillIiUpdateAddressDto)
	cainiaoWaybillIiUpdateUserInfoDto.SetMobile("132432323")
	cainiaoWaybillIiUpdateUserInfoDto.SetName("Foo")
	cainiaoWaybillIiUpdateUserInfoDto.SetPhone("05712323241")
	cainiaoWaybillIiUpdateUserInfoDto.SetOaid("abcdefghijklmn")
	cainiaoWaybillIiUpdateUserInfoDto.SetCaid("abcdefghijklmn")

	cainiaoWaybillIiUpdateUserInfoDto = domain.CainiaoWaybillIiUpdateUserInfoDto{}
	cainiaoWaybillIiUpdateUserInfoDto.SetMobile("1352353325")
	cainiaoWaybillIiUpdateUserInfoDto.SetName("Foo")
	cainiaoWaybillIiUpdateUserInfoDto.SetPhone("05714232523")
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest := domain.CainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest{}
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetCpCode("POSTB")
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetLogisticsServices(`{     "SVC-COD": {         "value": "200"     } }`)
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetPackageInfo(cainiaoWaybillIiUpdatePackageInfoDto)
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetRecipient(cainiaoWaybillIiUpdateUserInfoDto)
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetSender(cainiaoWaybillIiUpdateUserInfoDto)
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetTemplateUrl("http://cloudprint.cainiao.com/cloudprint/template/getStandardTemplate.json?template_id=1001")
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetWaybillCode("9890000160004")
	cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest.SetObjectId("x")

	req := request.CainiaoWaybillIiUpdateRequest{}
	req.SetParamWaybillCloudPrintUpdateRequest(cainiaoWaybillIiUpdateWaybillCloudPrintUpdateRequest)

	resp, err := ability.CainiaoWaybillIiUpdate(&req, UserSession)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Body)
	}
}

// 获取所有的菜鸟标准电子面单模板
func StdTemplatesGet() {

}

// 通过打印机命令打印
func Print() {

	ability := NewAbility229()

	cainiaoCloudprintCmdprintRenderRenderContent := ability229Domain.CainiaoCloudprintCmdprintRenderRenderContent{}
	cainiaoCloudprintCmdprintRenderRenderContent.SetPrintData("{}")
	cainiaoCloudprintCmdprintRenderRenderContent.SetTemplateUrl("http://cloudprint.cainiao.com/template/standard/401")
	cainiaoCloudprintCmdprintRenderRenderContent.SetEncrypted(true)
	cainiaoCloudprintCmdprintRenderRenderContent.SetVer("waybill_print_secret_version_1")
	cainiaoCloudprintCmdprintRenderRenderContent.SetSignature("MD:8jyWc0A8m/4CkO9bw9oqHA==")
	cainiaoCloudprintCmdprintRenderRenderContent.SetAddData("{ sender:{ address:{ detail:蒋村街道西溪诚园小区2-1-101 } } }")
	cainiaoCloudprintCmdprintRenderRenderDocument := ability229Domain.CainiaoCloudprintCmdprintRenderRenderDocument{}
	//
	cainiaoCloudprintCmdprintRenderRenderDocument.SetContents([]ability229Domain.CainiaoCloudprintCmdprintRenderRenderContent{})
	cainiaoCloudprintCmdprintRenderRenderConfig := ability229Domain.CainiaoCloudprintCmdprintRenderRenderConfig{}
	cainiaoCloudprintCmdprintRenderRenderConfig.SetOrientation("normal")
	cainiaoCloudprintCmdprintRenderRenderConfig.SetNeedBottomLogo(true)
	cainiaoCloudprintCmdprintRenderRenderConfig.SetNeedMiddleLogo(true)
	cainiaoCloudprintCmdprintRenderRenderConfig.SetNeedTopLogo(true)
	cainiaoCloudprintCmdprintRenderCmdRenderParams := ability229Domain.CainiaoCloudprintCmdprintRenderCmdRenderParams{}
	cainiaoCloudprintCmdprintRenderCmdRenderParams.SetDocument(cainiaoCloudprintCmdprintRenderRenderDocument)
	cainiaoCloudprintCmdprintRenderCmdRenderParams.SetPrinterName("KM-300S-EB13")
	cainiaoCloudprintCmdprintRenderCmdRenderParams.SetClientId("abc123")
	cainiaoCloudprintCmdprintRenderCmdRenderParams.SetClientType("alipay")
	cainiaoCloudprintCmdprintRenderCmdRenderParams.SetConfig(cainiaoCloudprintCmdprintRenderRenderConfig)

	req := ability229req.CainiaoCloudprintCmdprintRenderRequest{}
	req.SetParams(cainiaoCloudprintCmdprintRenderCmdRenderParams)

	resp, err := ability.CainiaoCloudprintCmdprintRender(&req, UserSession)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(resp.Body)
	}
}
