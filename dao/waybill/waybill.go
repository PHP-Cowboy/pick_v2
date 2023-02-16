package waybill

import (
	"fmt"
	"pick_v2/dao/waybill/defaultability"
	"pick_v2/dao/waybill/defaultability/domain"
	"pick_v2/dao/waybill/defaultability/request"
)

func WaybillIiGet() {
	client := NewDefaultTopClient(AppKey, AppSecret, GatewayUrl, 20000, 20000)
	ability := defaultability.NewDefaultability(&client)

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
