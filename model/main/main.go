package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"pick_v2/model"
	"pick_v2/model/batch"
	"pick_v2/model/order"
	"pick_v2/model/other"
	"time"
)

func main() {

	dsn := "root:bnskdfglnbbgf@tcp(121.196.60.92)/pick_v2?charset=utf8mb4&parseTime=True&loc=Local"

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

	_ = db.Set(model.TableOptions, model.GetOptions("拣货订单数据表")).AutoMigrate(&order.OrderInfo{})

	_ = db.Set(model.TableOptions, model.GetOptions("批次生成条件表")).AutoMigrate(&batch.BatchCondition{})

	_ = db.Set(model.TableOptions, model.GetOptions("批次表")).AutoMigrate(&batch.Batch{})

	_ = db.Set(model.TableOptions, model.GetOptions("预拣货列表")).AutoMigrate(&batch.PrePick{})

	_ = db.Set(model.TableOptions, model.GetOptions("预拣货商品明细表")).AutoMigrate(&batch.PrePickGoods{})

	_ = db.Set(model.TableOptions, model.GetOptions("预拣货备注明细表")).AutoMigrate(&batch.PrePickRemark{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货列表")).AutoMigrate(&batch.Pick{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货商品明细表")).AutoMigrate(&batch.PickGoods{})

	_ = db.Set(model.TableOptions, model.GetOptions("拣货备注明细表")).AutoMigrate(&batch.PickRemark{})

	_ = db.Set(model.TableOptions, model.GetOptions("店铺表")).AutoMigrate(&other.Shop{})

	_ = db.Set(model.TableOptions, model.GetOptions("分类表")).AutoMigrate(&other.Classification{})

	_ = db.Set(model.TableOptions, model.GetOptions("仓库表")).AutoMigrate(&other.Warehouse{})

	_ = db.Set(model.TableOptions, model.GetOptions("字典类型表")).AutoMigrate(&other.DictType{})

	_ = db.Set(model.TableOptions, model.GetOptions("字典表")).AutoMigrate(&other.Dict{})

	_ = db.Set(model.TableOptions, model.GetOptions("用户表")).AutoMigrate(&model.User{})

	_ = db.Set(model.TableOptions, model.GetOptions("角色表")).AutoMigrate(&model.Role{})

	_ = db.Set(model.TableOptions, model.GetOptions("菜单表")).AutoMigrate(&model.Menu{})

	_ = db.Set(model.TableOptions, model.GetOptions("角色菜单权限表")).AutoMigrate(&model.RoleMenu{})

	_ = db.Set(model.TableOptions, model.GetOptions("完成订单表")).AutoMigrate(&order.CompleteOrder{})

	_ = db.Set(model.TableOptions, model.GetOptions("完成订单明细表")).AutoMigrate(&order.CompleteOrderDetail{})

}
