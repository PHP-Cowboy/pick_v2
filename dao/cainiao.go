package dao

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"pick_v2/dao/cainiao"
)

type WayBillCommon struct {
	MsgType            string `json:"msg_type"`
	LogisticProviderId string `json:"logistic_provider_id"`
	DataDigest         string `json:"data_digest"`
	ToCode             string `json:"to_code"`
	LogisticsInterface string `json:"logistics_interface"`
}

type TmsWayBillGet struct {
	//  物流公司Code 长度小于20
	CpCode string `json:"cp_code,omitempty" xml:"cp_code,omitempty"`
	// 发货人信息
	Sender *UserInfoDto `json:"sender,omitempty" xml:"sender,omitempty"`
	// 请求面单信息，数量限制为10
	TradeOrderInfoDtos []TradeOrderInfoDto `json:"tradeOrderInfoDtos,omitempty" xml:"tradeOrderInfoDtos>tradeOrderInfoDto,omitempty"`
	// 仓code， 仓库WMS系统对接落地配业务，其它场景请不要使用
	StoreCode string `json:"storeCode,omitempty" xml:"storeCode,omitempty"`
	// 配送资源code， 仓库WMS系统对接落地配业务，其它场景请不要使用
	ResourceCode string `json:"resourceCode,omitempty" xml:"resourceCode,omitempty"`
	// 月结卡号
	CustomerCode string `json:"customerCode,omitempty" xml:"customerCode,omitempty"`
}

// UserInfoDto 结构体
type UserInfoDto struct {
	// 手机号码（手机号和固定电话不能同时为空），长度小于20
	Mobile string `json:"mobile,omitempty" xml:"mobile,omitempty"`
	// 姓名，长度小于40
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	// 固定电话（手机号和固定电话不能同时为空），长度小于20
	Phone string `json:"phone,omitempty" xml:"phone,omitempty"`
	// 发货地址需要通过&lt;a href=&#34;http://open.taobao.com/doc2/detail.htm?spm=a219a.7629140.0.0.3OFCPk&amp;treeId=17&amp;articleId=104860&amp;docType=1&#34;&gt;search接口&lt;/a&gt;
	Address *AddressDto `json:"address,omitempty" xml:"address,omitempty"`
}

// AddressDto 结构体
type AddressDto struct {
	// 城市，长度小于20
	City string `json:"city,omitempty" xml:"city,omitempty"`
	// 详细地址，长度小于256
	Detail string `json:"detail,omitempty" xml:"detail,omitempty"`
	// 区，长度小于20
	District string `json:"district,omitempty" xml:"district,omitempty"`
	// 省，长度小于20
	Province string `json:"province,omitempty" xml:"province,omitempty"`
	// 街道，长度小于30
	Town string `json:"town,omitempty" xml:"town,omitempty"`
}

// TradeOrderInfoDto 结构体
type TradeOrderInfoDto struct {
	// 物流服务值（详见https://support-cnkuaidi.taobao.com/doc.htm#?docId=106156&amp;docType=1，如无特殊服务请置空）
	LogisticsServices string `json:"logisticsServices,omitempty" xml:"logisticsServices,omitempty"`
	// &lt;a href=&#34;http://open.taobao.com/docs/doc.htm?docType=1&amp;articleId=105086&amp;treeId=17&amp;platformId=17#6&#34;&gt;请求ID&lt;/a&gt;
	ObjectId string `json:"object_id,omitempty" xml:"object_id,omitempty"`
}

func Sign(content, key string) string {
	signature := md5.Sum([]byte(content + key))

	return base64.StdEncoding.EncodeToString(signature[:])
}

type CreateRequest struct {
	BaseURL     string
	AppKey      string
	FromCode    string
	PartnerCode string
	Data        struct {
		DeliveryOrder DeliveryOrder `json:"deliveryOrder"`
		OrderLines    []OrderLine   `json:"orderLines"`
	}
}

type OrderLine struct {
	OrderLineNo string `json:"orderLineNo"`
	OwnerCode   string `json:"ownerCode"`
	ItemCode    string `json:"itemCode"`
	ItemId      string `json:"itemId"`
	PlanQty     string `json:"planQty"`
}

type DeliveryOrder struct {
	DeliveryOrderCode    string  `json:"deliveryOrderCode"`
	PreDeliveryOrderCode string  `json:"preDeliveryOrderCode"`
	OrderType            string  `json:"orderType"` // 出库单类型，JYCK=一般交易出库单, HHCK=换货出库单, BFCK=补发出库单，QTCK=其他出库单
	WarehouseCode        string  `json:"warehouseCode"`
	SourcePlatformCode   string  `json:"sourcePlatformCode"`
	CreateTime           string  `json:"createTime"`     // 发货单创建时间 2015-06-12 20:26:32
	PlaceOrderTime       string  `json:"placeOrderTime"` // 前台订单 (店铺订单) 创建时间 (下单时间) 2015-06-12 20:26:32
	OperateTime          string  `json:"operateTime"`
	ShopNick             string  `json:"shopNick"`
	LogisticsCode        string  `json:"logisticsCode"` // 物流公司编码
	SenderInfo           Contact `json:"senderInfo"`
	ReceiverInfo         Contact `json:"receiverInfo"`
}

type Contact struct {
	Name          string `json:"name"`
	Mobile        string `json:"mobile"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Area          string `json:"area"`
	Town          string `json:"town"`
	DetailAddress string `json:"detailAddress"`
}

type CreateResponse struct {
	Flag            string `json:"flag"`
	Code            string `json:"code"`
	Message         string `json:"message"`
	DeliveryOrderId string `json:"deliveryOrderId"`
	WarehouseCode   string `json:"warehouseCode"`
	LogisticsCode   string `json:"logisticsCode"`
	DeliveryOrder   []struct {
		DeliveryOrderId string `json:"deliveryOrderId"`
		WarehouseCode   string `json:"warehouseCode"`
		LogisticsCode   string `json:"logisticsCode"`
		OrderLines      []struct {
			OrderLineNo string `json:"orderLineNo"`
			ItemCode    string `json:"itemCode"`
			ItemId      string `json:"itemId"`
			Quantity    string `json:"quantity"`
		} `json:"orderLines"`
	} `json:"deliveryOrders"`
}

var (
	baseURL   = "http://linkdaily.tbsandbox.com/gateway/link.do"
	appKey    = "940356"
	appSecret = "olmT3icH12h5n2Q0978x3BZnD89g722q"
)

func Create(req CreateRequest) (CreateResponse, error) {
	apiReq := cainiao.New(cainiao.Config{
		BaseURL:            baseURL,
		AppKey:             appKey,
		AppSecret:          appSecret,
		MsgType:            "TMS_WAYBILL_GET",
		LogisticProviderId: "",
		ToCode:             "",
		PartnerCode:        req.PartnerCode,
		FromCode:           req.FromCode,
	})

	vals, err := apiReq.Post(req.Data)
	if err != nil {
		return CreateResponse{}, err
	}

	result, err := vals.GetResult(CreateResponse{})
	if err != nil {
		return CreateResponse{}, fmt.Errorf("%w: %s", cainiao.ErrCallAPIFailed, err.Error())
	}

	resp, ok := result.(*CreateResponse)
	if !ok {
		return CreateResponse{}, fmt.Errorf("%w: result is not CreateResponse", cainiao.ErrCallAPIFailed)
	}

	if resp.Flag == "failure" {
		return *resp, fmt.Errorf("%w: (%s)%s", cainiao.ErrAPIBizError, resp.Code, resp.Message)
	}

	return *resp, nil
}
