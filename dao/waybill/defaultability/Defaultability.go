package defaultability

import (
	"errors"
	"log"
	"pick_v2/dao/waybill"
	"pick_v2/dao/waybill/defaultability/request"
	"pick_v2/dao/waybill/defaultability/response"
	"pick_v2/dao/waybill/util"
)

type Defaultability struct {
	Client *waybill.TopClient
}

func NewDefaultability(client *waybill.TopClient) *Defaultability {
	return &Defaultability{client}
}

/*
获取所有的菜鸟标准电子面单模板
*/
func (ability *Defaultability) CainiaoCloudprintStdtemplatesGet(req *request.CainiaoCloudprintStdtemplatesGetRequest) (*response.CainiaoCloudprintStdtemplatesGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.cloudprint.stdtemplates.get", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoCloudprintStdtemplatesGetResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintStdtemplatesGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
获取用户使用的菜鸟电子面单模板信息
*/
func (ability *Defaultability) CainiaoCloudprintMystdtemplatesGet(req *request.CainiaoCloudprintMystdtemplatesGetRequest, session string) (*response.CainiaoCloudprintMystdtemplatesGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.cloudprint.mystdtemplates.get", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoCloudprintMystdtemplatesGetResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintMystdtemplatesGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
商家查询物流商产品类型接口
*/
func (ability *Defaultability) CainiaoWaybillIiProduct(req *request.CainiaoWaybillIiProductRequest, session string) (*response.CainiaoWaybillIiProductResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.product", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiProductResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiProduct error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
商家取消获取的电子面单号
*/
func (ability *Defaultability) CainiaoWaybillIiCancel(req *request.CainiaoWaybillIiCancelRequest, session string) (*response.CainiaoWaybillIiCancelResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.cancel", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiCancelResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiCancel error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
关键词过滤匹配
*/
func (ability *Defaultability) TaobaoKfcKeywordSearch(req *request.TaobaoKfcKeywordSearchRequest, session string) (*response.TaobaoKfcKeywordSearchResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.kfc.keyword.search", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoKfcKeywordSearchResponse{}
	if err != nil {
		log.Println("taobaoKfcKeywordSearch error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
isv资源查询
*/
func (ability *Defaultability) CainiaoCloudprintIsvResourcesGet(req *request.CainiaoCloudprintIsvResourcesGetRequest) (*response.CainiaoCloudprintIsvResourcesGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.cloudprint.isv.resources.get", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoCloudprintIsvResourcesGetResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintIsvResourcesGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
获取商家的自定义区模板信息
*/
func (ability *Defaultability) CainiaoCloudprintCustomaresGet(req *request.CainiaoCloudprintCustomaresGetRequest, session string) (*response.CainiaoCloudprintCustomaresGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.cloudprint.customares.get", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoCloudprintCustomaresGetResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintCustomaresGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
获取物流服务商电子面单号v1.0
*/
func (ability *Defaultability) TaobaoWlbWaybillIGet(req *request.TaobaoWlbWaybillIGetRequest, session string) (*response.TaobaoWlbWaybillIGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.wlb.waybill.i.get", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoWlbWaybillIGetResponse{}
	if err != nil {
		log.Println("taobaoWlbWaybillIGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
面单信息更新接口v1.0
*/
func (ability *Defaultability) TaobaoWlbWaybillIFullupdate(req *request.TaobaoWlbWaybillIFullupdateRequest, session string) (*response.TaobaoWlbWaybillIFullupdateResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.wlb.waybill.i.fullupdate", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoWlbWaybillIFullupdateResponse{}
	if err != nil {
		log.Println("taobaoWlbWaybillIFullupdate error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
打印确认接口v1.0
*/
func (ability *Defaultability) TaobaoWlbWaybillIPrint(req *request.TaobaoWlbWaybillIPrintRequest, session string) (*response.TaobaoWlbWaybillIPrintResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.wlb.waybill.i.print", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoWlbWaybillIPrintResponse{}
	if err != nil {
		log.Println("taobaoWlbWaybillIPrint error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
查面单号状态v1.0
*/
func (ability *Defaultability) TaobaoWlbWaybillIQuerydetail(req *request.TaobaoWlbWaybillIQuerydetailRequest, session string) (*response.TaobaoWlbWaybillIQuerydetailResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.wlb.waybill.i.querydetail", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoWlbWaybillIQuerydetailResponse{}
	if err != nil {
		log.Println("taobaoWlbWaybillIQuerydetail error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
商家取消获取的电子面单号v1.0
*/
func (ability *Defaultability) TaobaoWlbWaybillICancel(req *request.TaobaoWlbWaybillICancelRequest, session string) (*response.TaobaoWlbWaybillICancelResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.wlb.waybill.i.cancel", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoWlbWaybillICancelResponse{}
	if err != nil {
		log.Println("taobaoWlbWaybillICancel error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
商家查询物流商产品类型接口
*/
func (ability *Defaultability) TaobaoWlbWaybillIProduct(req *request.TaobaoWlbWaybillIProductRequest, session string) (*response.TaobaoWlbWaybillIProductResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("taobao.wlb.waybill.i.product", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.TaobaoWlbWaybillIProductResponse{}
	if err != nil {
		log.Println("taobaoWlbWaybillIProduct error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
通过订单号查询电子面单通接口
*/
func (ability *Defaultability) CainiaoWaybillIiQueryByTradecode(req *request.CainiaoWaybillIiQueryByTradecodeRequest, session string) (*response.CainiaoWaybillIiQueryByTradecodeResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.query.by.tradecode", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiQueryByTradecodeResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiQueryByTradecode error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
通过面单号查询面单打印报文
*/
func (ability *Defaultability) CainiaoWaybillIiQueryByWaybillcode(req *request.CainiaoWaybillIiQueryByWaybillcodeRequest, session string) (*response.CainiaoWaybillIiQueryByWaybillcodeResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.query.by.waybillcode", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiQueryByWaybillcodeResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiQueryByWaybillcode error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
云打印模板迁移接口
*/
func (ability *Defaultability) CainiaoCloudprintTemplatesMigrate(req *request.CainiaoCloudprintTemplatesMigrateRequest, session string) (*response.CainiaoCloudprintTemplatesMigrateResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.cloudprint.templates.migrate", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoCloudprintTemplatesMigrateResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintTemplatesMigrate error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
获取商家使用的标准模板
*/
func (ability *Defaultability) CainiaoCloudprintIsvtemplatesGet(req *request.CainiaoCloudprintIsvtemplatesGetRequest, session string) (*response.CainiaoCloudprintIsvtemplatesGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.cloudprint.isvtemplates.get", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoCloudprintIsvtemplatesGetResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintIsvtemplatesGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
查询面单服务订购及面单使用情况
*/
func (ability *Defaultability) CainiaoWaybillIiSearch(req *request.CainiaoWaybillIiSearchRequest, session string) (*response.CainiaoWaybillIiSearchResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.search", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiSearchResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiSearch error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
电子面单云打印接口
*/
func (ability *Defaultability) CainiaoWaybillIiGet(req *request.CainiaoWaybillIiGetRequest, session string) (*response.CainiaoWaybillIiGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.get", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiGetResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
电子面单云打印更新接口
*/
func (ability *Defaultability) CainiaoWaybillIiUpdate(req *request.CainiaoWaybillIiUpdateRequest, session string) (*response.CainiaoWaybillIiUpdateResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Defaultability topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.waybill.ii.update", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoWaybillIiUpdateResponse{}
	if err != nil {
		log.Println("cainiaoWaybillIiUpdate error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}
