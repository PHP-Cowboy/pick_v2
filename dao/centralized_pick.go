package dao

import (
	"gorm.io/gorm"
	"pick_v2/forms/rsp"
	"pick_v2/model"
)

// 生成集中拣货
func CreateCentralizedPick(db *gorm.DB, outboundGoodsJoinOrder []model.OutboundGoodsJoinOrder, batchId int) error {

	mpCentralized := make(map[string]model.CentralizedPick, 0)

	//按sku归集数据
	for _, order := range outboundGoodsJoinOrder {
		cp, mpCentralizedOk := mpCentralized[order.Sku]

		hasRemark := 0

		if order.GoodsRemark != "" {
			hasRemark = 1
		}

		if !mpCentralizedOk {
			cp = model.CentralizedPick{
				BatchId:        batchId,
				Sku:            order.Sku,
				GoodsName:      order.GoodsName,
				GoodsType:      order.GoodsType,
				GoodsSpe:       order.GoodsSpe,
				PickUser:       "",
				TakeOrdersTime: nil,
				HasRemark:      hasRemark,
			}
		}

		cp.NeedNum += order.LackCount

		mpCentralized[order.Sku] = cp
	}

	//集中拣货数据构造
	centralizedPicks := make([]model.CentralizedPick, 0, len(mpCentralized))

	for _, pick := range mpCentralized {
		centralizedPicks = append(centralizedPicks, pick)
	}

	//集中拣货数据保存
	err := model.CentralizedPickSave(db, &centralizedPicks)

	if err != nil {
		return err
	}

	return nil
}

// 集中拣货列表
func CentralizedPickList(db *gorm.DB) (err error, res rsp.CentralizedPickListRsp) {
	err, total, centralizedPickList := model.GetCentralizedPickList(db)
	if err != nil {
		return err, res
	}

	res.Total = total

	list := make([]*rsp.CentralizedPickList, 0, len(centralizedPickList))

	for _, pick := range centralizedPickList {
		list = append(list, &rsp.CentralizedPickList{
			TaskName:  pick.GoodsName,
			GoodsName: pick.GoodsName,
			GoodsSpe:  pick.GoodsSpe,
			NeedNum:   pick.NeedNum,
			PickNum:   pick.PickNum,
			PickUser:  pick.PickUser,
			HasRemark: pick.HasRemark,
		})
	}

	res.List = list
	return
}
