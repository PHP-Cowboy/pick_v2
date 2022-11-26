package dao

import (
	"errors"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/slice"
	"strconv"
	"strings"
)

// 验证是否有任务被接单
func CheckPickIsReceived(db *gorm.DB, ids []int) (err error, pickList []model.Pick) {
	err, pickList = model.GetPickListByIds(db, ids)

	if err != nil {
		return
	}

	//拣货池状态
	for _, p := range pickList {
		if p.PickUser != "" {
			return errors.New("已有拣货任务被接单，无法取消"), nil
		}
	}

	return
}

// 校验任务是否有订单已被接单拣货
func CheckPickNumberIsReceived(db *gorm.DB, ids []int) (err error) {

	var (
		numbers       []string
		exist         bool
		pickGoodsList []model.PickGoods
	)

	//拣货池根据所选任务查询全部商品
	err, pickGoodsList = model.GetPickGoodsByPickIds(db, ids)

	if err != nil {
		return
	}

	//全部商品的订单编号
	for _, pg := range pickGoodsList {
		numbers = append(numbers, pg.Number)
	}

	numbers = slice.UniqueSlice(numbers)

	//查询订单是否有已拣的
	err, exist = model.GetFirstPickGoodsByNumbers(db, numbers)

	if err != nil {
		return
	}

	if exist {
		return errors.New("该任务订单部分商品已拣，不可关闭")
	}

	return
}

// 拣货池任务返回预拣池
func PickReturnPrePick(db *gorm.DB, pickList []model.Pick) error {

	var (
		ids []int
	)

	for _, pl := range pickList {
		//获取预拣池id
		prePickIds := strings.Split(pl.PrePickIds, ",")

		for _, id := range prePickIds {
			prePickId, err := strconv.Atoi(id)
			if err != nil {
				return err
			}
			ids = append(ids, prePickId)
		}
	}

	tx := db.Begin()

	var (
		prePickIds []int
		pickIds    []int
	)

	//更新预拣池状态
	err := ChangePrePickStatus(tx, prePickIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	//更新拣货池状态
	err = ChangePickStatus(tx, pickIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// 更新预拣池状态
func ChangePrePickStatus(db *gorm.DB, prePickIds []int) error {
	//更新预拣池商品状态
	err := model.UpdatePrePickGoodsByPrePickIds(db, prePickIds, map[string]interface{}{"status": model.PrePickGoodsStatusProcessing})
	if err != nil {
		return err
	}

	//更新预拣池任务状态
	err = model.UpdatePrePickByIds(db, prePickIds, map[string]interface{}{"status": model.PrePickStatusUnhandled})
	if err != nil {
		return err
	}

	//更新预拣池备注状态
	err = model.UpdatePrePickRemarkByPrePickIds(db, prePickIds, map[string]interface{}{"status": model.PrePickRemarkStatusUnhandled})
	if err != nil {
		return err
	}

	return nil
}

// 更新拣货池状态
func ChangePickStatus(db *gorm.DB, pickIds []int) error {
	//拣货池商品和备注表没有状态
	err := model.UpdatePickByIds(db, pickIds, map[string]interface{}{"status": model.ReturnPrePickStatus})
	if err != nil {
		return err
	}

	return nil
}

// 取消拣货
func CancelPick(db *gorm.DB, form req.CancelPickForm) error {

	var (
		err      error
		pickList []model.Pick
	)

	//验证是否有任务被接单
	err, pickList = CheckPickIsReceived(db, form.Ids)

	if err != nil {
		return err
	}

	//校验任务是否有订单已被接单拣货
	err = CheckPickNumberIsReceived(db, form.Ids)

	if err != nil {
		return err
	}

	//拣货池任务返回预拣池
	err = PickReturnPrePick(db, pickList)
	if err != nil {
		return err
	}

	return nil
}
