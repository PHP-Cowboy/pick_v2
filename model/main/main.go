package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"pick_v2/model"
	"time"
)

func main() {

	dsn := "root:a900614f907ceae1@tcp(192.168.1.40)/pick_v2?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:bnskdfglnbbgf@tcp(121.196.60.92)/pick_v2?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "pickv2user:whpoJJTEM7N0@tcp(rm-bp1v01uw93jftaj8p8o.mysql.rds.aliyuncs.com)/pick_v2?charset=utf8mb4&parseTime=True&loc=Local"

	logger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      true,        //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})

	if err != nil {
		panic(err)
	}

	_ = db.Set(model.TableOptions, model.GetOptions("订单表")).AutoMigrate(&model.Order{})

	_ = db.Set(model.TableOptions, model.GetOptions("订单商品表")).AutoMigrate(&model.OrderGoods{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货单")).AutoMigrate(&model.PickOrder{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货单商品")).AutoMigrate(&model.PickOrderGoods{})

	_ = db.Set(model.TableOptions, model.GetOptions("批次生成条件表")).AutoMigrate(&model.BatchCondition{})

	_ = db.Set(model.TableOptions, model.GetOptions("批次表")).AutoMigrate(&model.Batch{})

	_ = db.Set(model.TableOptions, model.GetOptions("预拣货列表")).AutoMigrate(&model.PrePick{})

	_ = db.Set(model.TableOptions, model.GetOptions("预拣货商品明细表")).AutoMigrate(&model.PrePickGoods{})

	_ = db.Set(model.TableOptions, model.GetOptions("预拣货备注明细表")).AutoMigrate(&model.PrePickRemark{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货列表")).AutoMigrate(&model.Pick{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货商品明细表")).AutoMigrate(&model.PickGoods{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货备注明细表")).AutoMigrate(&model.PickRemark{})

	_ = db.Set(model.TableOptions, model.GetOptions("店铺表")).AutoMigrate(&model.Shop{})

	_ = db.Set(model.TableOptions, model.GetOptions("分类表")).AutoMigrate(&model.Classification{})

	_ = db.Set(model.TableOptions, model.GetOptions("仓库表")).AutoMigrate(&model.Warehouse{})

	_ = db.Set(model.TableOptions, model.GetOptions("字典类型表")).AutoMigrate(&model.DictType{})

	_ = db.Set(model.TableOptions, model.GetOptions("字典表")).AutoMigrate(&model.Dict{})

	_ = db.Set(model.TableOptions, model.GetUserOptions("用户表")).AutoMigrate(&model.User{})

	_ = db.Set(model.TableOptions, model.GetOptions("角色表")).AutoMigrate(&model.Role{})

	_ = db.Set(model.TableOptions, model.GetOptions("菜单表")).AutoMigrate(&model.Menu{})

	_ = db.Set(model.TableOptions, model.GetOptions("角色菜单权限表")).AutoMigrate(&model.RoleMenu{})

	_ = db.Set(model.TableOptions, model.GetOptions("完成订单表")).AutoMigrate(&model.CompleteOrder{})

	_ = db.Set(model.TableOptions, model.GetOptions("完成订单明细表")).AutoMigrate(&model.CompleteOrderDetail{})

	_ = db.Set(model.TableOptions, model.GetOptions("限制发货")).AutoMigrate(&model.RestrictedShipment{})

	_ = db.Set(model.TableOptions, model.GetOptions("限制发货")).AutoMigrate(&model.LimitShipment{})

	_ = db.Set(model.TableOptions, model.GetOptions("推送日志")).AutoMigrate(&model.StockLog{})

	_ = db.Set(model.TableOptions, model.GetOptions("盘点任务表")).AutoMigrate(&model.InvTask{})

	_ = db.Set(model.TableOptions, model.GetOptions("自建盘点任务表")).AutoMigrate(&model.InvTaskSelfBuilt{})

	_ = db.Set(model.TableOptions, model.GetOptions("盘点任务商品记录表")).AutoMigrate(&model.InvTaskRecord{})

	_ = db.Set(model.TableOptions, model.GetOptions("用户盘点记录表")).AutoMigrate(&model.InventoryRecord{})

	_ = db.Set(model.TableOptions, model.GetOptions("出库单任务表")).AutoMigrate(&model.OutboundTask{})

	_ = db.Set(model.TableOptions, model.GetOptions("出库单订单表")).AutoMigrate(&model.OutboundOrder{})

	_ = db.Set(model.TableOptions, model.GetOptions("出库单商品表")).AutoMigrate(&model.OutboundGoods{})

}
