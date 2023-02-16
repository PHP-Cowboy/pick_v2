package response

import (
	"topsdk/defaultability/domain"
)

type CainiaoWaybillIiSearchResponse struct {

	/*
	   System request id
	*/
	RequestId string `json:"request_id,omitempty" `

	/*
	   System body
	*/
	Body string

	/*
	   CP网点信息及对应的商家的发货信息
	*/
	WaybillApplySubscriptionCols []domain.CainiaoWaybillIiSearchWaybillApplySubscriptionInfo `json:"waybill_apply_subscription_cols,omitempty" `
}
