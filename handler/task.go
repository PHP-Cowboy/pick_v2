package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"pick_v2/common/constant"
	"pick_v2/dao"
	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/model"
	"pick_v2/utils/cache"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"pick_v2/utils/xsq_net"
	"strings"
	"time"
)

// 置顶
func PickTopping(c *gin.Context) {
	var form req.PickToppingForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	val, err := cache.Incr(constant.PICK_TOPPING)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	result := global.DB.Model(model.Pick{}).Where("id = ?", form.Id).Update("sort", val)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 拣货池列表
func PickList(c *gin.Context) {
	var form req.PickListForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		ids        []int
		res        rsp.PickListRsp
		pick       []model.Pick
		pickGoods  []model.PickGoods
		pickRemark []model.PickRemark
		result     *gorm.DB
	)

	if form.Goods != "" || form.Number != "" || form.ShopId > 0 {
		result = db.Where(model.PickGoods{BatchId: form.BatchId, GoodsName: form.Goods, Number: form.Number, ShopId: form.ShopId}).Find(&pickGoods)
		if result.Error != nil {
			xsq_net.ErrorJSON(c, result.Error)
			return
		}

		for _, pg := range pickGoods {
			ids = append(ids, pg.Id)
		}
	}

	localDb := db.Table("t_pick")

	if len(ids) > 0 {
		localDb.Where("id in (?)", ids)
	}

	localDb.Where("batch_id = ? and status = ?", form.BatchId, form.Status)

	result = localDb.Find(&pick)

	res.Total = result.RowsAffected

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = localDb.Order("sort desc").Scopes(model.Paginate(form.Page, form.Size)).Find(&pick)

	//拣货ids
	pickIds := make([]int, 0, len(pick))
	for _, p := range pick {
		pickIds = append(pickIds, p.Id)
	}

	//拣货ids 的订单备注
	result = db.Where("pick_id in (?)", pickIds).Find(&pickRemark)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	query := "pick_id,count(distinct(shop_id)) as shop_num,count(distinct(number)) as order_num,sum(need_num) as need_num,sum(complete_num) as complete_num,sum(review_num) as review_num"

	err, numsMp := model.CountPickPoolNumsByPickIds(db, pickIds, query)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//构建pickId 对应的订单 是否有备注map
	remarkMp := make(map[int]struct{}, 0) //key 存在即为有
	for _, remark := range pickRemark {
		remarkMp[remark.PickId] = struct{}{}
	}

	isRemark := false

	list := make([]rsp.Pick, 0)
	for _, p := range pick {
		_, ok := remarkMp[p.Id]
		if ok { //拣货id在拣货备注中存在，即为有备注
			isRemark = true
		}

		nums, numsOk := numsMp[p.Id]

		if !numsOk {
			err = errors.New("拣货相关数量统计有误")
			return
		}

		list = append(list, rsp.Pick{
			Id:             p.Id,
			ShopCode:       p.ShopCode,
			ShopName:       p.ShopName,
			ShopNum:        nums.ShopNum,
			OrderNum:       nums.OrderNum,
			NeedNum:        nums.NeedNum,
			PickUser:       p.PickUser,
			TakeOrdersTime: p.TakeOrdersTime,
			IsRemark:       isRemark,
			Status:         p.Status,
			UpdateTime:     p.UpdateTime.Format(timeutil.TimeFormat),
			PickNum:        nums.CompleteNum,
			ReviewNum:      nums.ReviewNum,
			Num:            p.Num,
			PrintNum:       p.PrintNum,
			ReviewUser:     p.ReviewUser,
			ReviewTime:     p.ReviewTime,
		})
	}

	res.List = list

	xsq_net.SucJson(c, res)
}

// 快递拣货列表
func CentralizedAndSecondaryList(c *gin.Context) {

	userInfo := GetUserInfo(c)

	if userInfo == nil {
		xsq_net.ErrorJSON(c, ecode.GetContextUserInfoFailed)
		return
	}

	err, list := dao.CentralizedAndSecondaryList(global.DB, userInfo.Name)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.SucJson(c, gin.H{"list": list})
}

// 拣货明细
func GetPickDetail(c *gin.Context) {
	var (
		form req.GetPickDetailForm
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err, res := dao.GetPickDetail(global.DB, form)

	if err != nil {
		return
	}

	xsq_net.SucJson(c, res)
}

// 集中拣货明细
func CentralizedPickDetailPDA(c *gin.Context) {
	var form req.CentralizedPickDetailPDAForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	err, res := dao.CentralizedPickDetailPDA(global.DB, form)
	if err != nil {
		return
	}

	xsq_net.SucJson(c, res)
}

// 修改已出库件数
func ChangeNum(c *gin.Context) {
	var form req.ChangeNumReq

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var pick model.Pick

	db := global.DB

	result := db.First(&pick, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if pick.Status != 2 {
		xsq_net.ErrorJSON(c, errors.New("只允许修改已出库的"))
		return
	}

	result = db.Model(&model.Pick{}).Where("id = ?", form.Id).Update("num", form.Num)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 打印
func PushPrint(c *gin.Context) {
	var form req.PrintReq

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		pick      []model.Pick
		pickGoods []model.PickGoods
	)

	db := global.DB

	//复核完成
	result := db.Where("id in (?) and status = 2", form.Ids).Find(&pick)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	pickMp := make(map[int]string, 0)

	for i, p := range pick {
		pickMp[p.Id] = strings.Join(p.DeliveryOrderNo, ",")
		pick[i].PrintNum += 1
	}

	result = db.Where("pick_id in (?)", form.Ids).Find(&pickGoods)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	shopAndPickGoodsIdMp := make(map[string]struct{})

	printChSlice := make([]global.PrintCh, 0)
	for _, good := range pickGoods {
		deliveryOrderNo, pickOk := pickMp[good.PickId]

		if !pickOk {
			continue
		}

		mpKey := fmt.Sprintf("%d%d", good.PickId, good.ShopId)
		_, ok := shopAndPickGoodsIdMp[mpKey]
		if ok {
			continue
		}
		shopAndPickGoodsIdMp[mpKey] = struct{}{}
		printChSlice = append(printChSlice, global.PrintCh{
			DeliveryOrderNo: deliveryOrderNo,
			ShopId:          good.ShopId,
		})
	}

	for _, ch := range printChSlice {
		dao.AddPrintJobMap(constant.JH_HUOSE_CODE, &global.PrintCh{
			DeliveryOrderNo: ch.DeliveryOrderNo,
			ShopId:          ch.ShopId,
			Type:            form.Type, //打印类型
		})
	}

	result = db.Select("id", "update_time", "shop_code", "shop_name", "line", "print_num").Save(&pick)

	xsq_net.Success(c)
}

// 指派
func Assign(c *gin.Context) {
	var form req.AssignReq

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		user  model.User
		picks []model.Pick
	)

	db := global.DB

	result := db.First(&user, form.UserId)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if user.RoleId != 3 { //不是拣货员
		xsq_net.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	result = db.Where("id in (?)", form.PickIds).Find(&picks)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, p := range picks {
		if p.Status != 0 {
			xsq_net.ErrorJSON(c, errors.New("只能分配待拣货的任务"))
			return
		}
	}

	now := time.Now()

	//已有拣货员可以强转
	result = db.Model(&model.Pick{}).Where("id in (?)", form.PickIds).Updates(map[string]interface{}{"pick_user": user.Name, "take_orders_time": &now})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 修改复核数
func ChangeReviewNum(c *gin.Context) {
	var form req.ChangeReviewNumForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := dao.ChangeReviewNum(global.DB, form)

	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}

// 取消拣货
func CancelPick(c *gin.Context) {
	var form req.CancelPickForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	err := dao.CancelPick(global.DB, form)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	xsq_net.Success(c)
}
