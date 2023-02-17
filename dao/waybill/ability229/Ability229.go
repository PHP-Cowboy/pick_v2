package ability229

import (
	"errors"
	"log"
	"pick_v2/dao/waybill"
	"pick_v2/dao/waybill/ability229/request"
	"pick_v2/dao/waybill/ability229/response"
	"pick_v2/dao/waybill/util"
)

type Ability229 struct {
	Client *waybill.TopClient
}

func NewAbility229(client *waybill.TopClient) *Ability229 {
	return &Ability229{client}
}

/*
网络打印机打印接口
*/
func (ability *Ability229) CainiaoWaybillCloudprintNetprintPrint(req *request.CainiaoWaybillCloudprintNetprintPrintRequest) (*response.CainiaoWaybillCloudprintNetprintPrintResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Ability229 topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.waybill.cloudprint.netprint.print", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoWaybillCloudprintNetprintPrintResponse{}
	if err != nil {
		log.Println("cainiaoWaybillCloudprintNetprintPrint error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
网络打印机绑定
*/
func (ability *Ability229) CainiaoWaybillCloudprintNetprintBind(req *request.CainiaoWaybillCloudprintNetprintBindRequest) (*response.CainiaoWaybillCloudprintNetprintBindResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Ability229 topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.waybill.cloudprint.netprint.bind", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoWaybillCloudprintNetprintBindResponse{}
	if err != nil {
		log.Println("cainiaoWaybillCloudprintNetprintBind error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
打印验证码
*/
func (ability *Ability229) CainiaoWaybillCloudprintNetprintVerifycode(req *request.CainiaoWaybillCloudprintNetprintVerifycodeRequest) (*response.CainiaoWaybillCloudprintNetprintVerifycodeResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Ability229 topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.waybill.cloudprint.netprint.verifycode", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoWaybillCloudprintNetprintVerifycodeResponse{}
	if err != nil {
		log.Println("cainiaoWaybillCloudprintNetprintVerifycode error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
获取商家单一自定义区
*/
func (ability *Ability229) CainiaoCloudprintSingleCustomareaGet(req *request.CainiaoCloudprintSingleCustomareaGetRequest) (*response.CainiaoCloudprintSingleCustomareaGetResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Ability229 topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.cloudprint.single.customarea.get", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoCloudprintSingleCustomareaGetResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintSingleCustomareaGet error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
生成打印机渲染命令（通过打印机命令打印）
*/
func (ability *Ability229) CainiaoCloudprintCmdprintRender(req *request.CainiaoCloudprintCmdprintRenderRequest, session string) (*response.CainiaoCloudprintCmdprintRenderResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Ability229 topClient is nil")
	}
	var jsonStr, err = ability.Client.ExecuteWithSession("cainiao.cloudprint.cmdprint.render", req.ToMap(), req.ToFileMap(), session)
	var respStruct = response.CainiaoCloudprintCmdprintRenderResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintCmdprintRender error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}

/*
云打印客户端监控信息收集
*/
func (ability *Ability229) CainiaoCloudprintClientinfoPut(req *request.CainiaoCloudprintClientinfoPutRequest) (*response.CainiaoCloudprintClientinfoPutResponse, error) {
	if ability.Client == nil {
		return nil, errors.New("Ability229 topClient is nil")
	}
	var jsonStr, err = ability.Client.Execute("cainiao.cloudprint.clientinfo.put", req.ToMap(), req.ToFileMap())
	var respStruct = response.CainiaoCloudprintClientinfoPutResponse{}
	if err != nil {
		log.Println("cainiaoCloudprintClientinfoPut error", err)
		return &respStruct, err
	}
	err = util.HandleJsonResponse(jsonStr, &respStruct)
	if respStruct.Body == "" || len(respStruct.Body) == 0 {
		respStruct.Body = jsonStr
	}
	return &respStruct, err
}
