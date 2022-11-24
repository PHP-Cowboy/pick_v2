package dao

import (
	"errors"
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/model"
	"pick_v2/utils/slice"
)

// 取消拣货
func CancelPick(db *gorm.DB, form req.CancelPickForm) error {

	err, pickList := model.GetPickListByIds(db, form.Ids)

	if err != nil {
		return err
	}

	//拣货池状态
	for _, p := range pickList {
		if p.PickUser != "" {
			return errors.New("已有拣货任务被接单，无法取消")
		}
	}

	var numbers []string

	err, pickGoodsList := model.GetPickGoodsByPickIds(db, form.Ids)
	if err != nil {
		return err
	}

	for _, pg := range pickGoodsList {
		numbers = append(numbers, pg.Number)
	}

	numbers = slice.UniqueSlice(numbers)

	return nil
}
