package dao

import (
	"gorm.io/gorm"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/model"
	"pick_v2/utils/ecode"
	"pick_v2/utils/slice"
)

func DeliveryMethodInfo(db *gorm.DB, form req.DeliveryMethodInfoForm) (err error, res rsp.DeliveryMethodInfoRsp) {
	var (
		outboundOrder model.OutboundOrder
	)

	result := db.Model(&model.OutboundOrder{}).
		Where("task_id = ? and number = ?", form.TaskId, form.Number).
		First(&outboundOrder)

	if result.Error != nil {
		return result.Error, res
	}

	res.UserName = outboundOrder.ConsigneeName
	res.Tel = outboundOrder.ConsigneeTel
	res.Province = outboundOrder.Province
	res.City = outboundOrder.City
	res.District = outboundOrder.District
	res.Address = outboundOrder.Address

	return nil, res
}

func ChangeDeliveryMethod(db *gorm.DB, form req.ChangeDeliveryMethodForm) (err error) {
	var dict []int

	result := db.Model(&model.Dict{}).
		Select("value").
		Where("type_code = ?", "delivery_method").
		Find(&dict)

	if result.Error != nil {
		return result.Error
	}

	//校验是否在字典中
	if ok, _ := slice.InArray(form.DistributionType, dict); !ok {
		return ecode.DataCannotBeModified
	}

	return model.OutboundOrderBatchUpdate(db, &model.OutboundOrder{TaskId: form.TaskId, Number: form.Number}, map[string]interface{}{"distribution_type": form.DistributionType})
}
