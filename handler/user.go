package handler

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"pick_v2/common/constant"
	"pick_v2/utils/ecode"
	"pick_v2/utils/timeutil"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"pick_v2/forms/req"
	"pick_v2/forms/rsp"
	"pick_v2/global"
	"pick_v2/middlewares"
	"pick_v2/model"
	"pick_v2/utils/xsq_net"
)

// 新增用户
func CreateUser(c *gin.Context) {
	var form req.AddUserForm
	err := c.ShouldBind(&form)
	if err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	db := global.DB

	var (
		user      model.User
		warehouse model.Warehouse
		role      model.Role
	)

	result := db.First(&warehouse, form.WarehouseId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.WarehouseNotFound)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result = db.First(&role, form.RoleId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			xsq_net.ErrorJSON(c, ecode.RoleNotFound)
			return
		}
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	user.Password = GenderPwd(form.Password)
	user.Name = form.Name
	user.RoleId = form.RoleId
	user.Role = role.Name
	user.WarehouseId = form.WarehouseId

	result = db.Save(&user)
	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	userRsp := rsp.AddUserRsp{
		Id:          user.Id,
		Account:     strconv.Itoa(user.Id),
		Name:        user.Name,
		Role:        user.Role,
		WarehouseId: user.WarehouseId,
	}

	xsq_net.SucJson(c, userRsp)
}

// 获取用户列表
func GetUserList(c *gin.Context) {
	var form req.GetUserListForm

	err := c.ShouldBind(&form)
	if err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		users []model.User
		res   rsp.UserListRsp
	)

	db := global.DB

	result := db.Where("delete_time is null").Find(&users)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	res.Total = result.RowsAffected

	db.Where("delete_time is null").Scopes(model.Paginate(form.Page, form.Size)).Find(&users)

	status := false

	for _, user := range users {
		if user.Status == 1 {
			status = true
		}
		res.List = append(res.List, &rsp.User{
			Id:          user.Id,
			Account:     strconv.Itoa(user.Id),
			Name:        user.Name,
			Role:        user.Role,
			WarehouseId: user.WarehouseId,
			Status:      status,
			CreateTime:  user.CreateTime.Format(timeutil.TimeFormat),
		})
	}

	xsq_net.SucJson(c, res)
}

// 登录
func Login(c *gin.Context) {
	var form req.LoginForm

	bindingBody := binding.Default(c.Request.Method, c.ContentType()).(binding.BindingBody)

	if err := c.ShouldBindBodyWith(&form, bindingBody); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	var (
		user model.User
	)

	db := global.DB
	result := db.Where("id = ? and status = 1 and delete_time is null", form.Id).First(&user)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		xsq_net.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	if user.WarehouseId != form.WarehouseId {
		xsq_net.ErrorJSON(c, ecode.WarehouseSelectError)
		return
	}

	options := &password.Options{16, 100, 32, sha512.New}

	pwdSlice := strings.Split(user.Password, "$")

	if !password.Verify(form.Password, pwdSlice[1], pwdSlice[2], options) {
		xsq_net.ErrorJSON(c, ecode.PasswordCheckFailed)
		return
	}

	account := strconv.Itoa(user.Id)

	hour := time.Duration(24)

	claims := middlewares.CustomClaims{
		ID:             user.Id,
		Account:        account,
		Name:           user.Name,
		WarehouseId:    user.WarehouseId,
		AuthorityId:    user.RoleId,
		StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add(hour * time.Hour).Unix()},
	}

	j := middlewares.NewJwt()
	token, err := j.CreateToken(claims)
	if err != nil {
		xsq_net.ErrorJSON(c, err)
		return
	}

	//token存入redis
	redisKey := constant.LOGIN_PREFIX + account
	err = global.Redis.Set(context.Background(), redisKey, token, hour*60*60*time.Second).Err()
	if err != nil {
		xsq_net.ErrorJSON(c, ecode.RedisFailedToSetData)
		return
	}

	xsq_net.SucJson(c, gin.H{
		"token":       token,
		"roleId":      user.RoleId,
		"userId":      user.Id,
		"account":     account,
		"name":        user.Name,
		"warehouseId": user.WarehouseId,
	})
}

// 修改 名称 密码 状态 组织
func ChangeUser(c *gin.Context) {
	var form req.CheckPwdForm

	err := c.ShouldBind(&form)

	if err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}
	var (
		user   model.User
		update model.User
	)

	db := global.DB

	result := db.First(&user, form.Id)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		xsq_net.ErrorJSON(c, ecode.UserNotFound)
		return
	}

	//更新用户密码
	if form.NewPassword != "" {
		update.Password = GenderPwd(form.NewPassword)
	}

	//更新用户名称
	if form.Name != "" {
		update.Name = form.Name
	}

	//更新用户状态
	if form.Status > 0 {
		update.Status = form.Status
	}

	//更新用户角色
	if form.RoleId > 0 {
		update.RoleId = form.RoleId
	}

	result = db.Model(&user).Updates(update)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 批量删除角色
func BatchDeleteUser(c *gin.Context) {
	var form req.BatchDeleteUserForm

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(model.User{}).Where("id in (?)", form.Ids).Updates(map[string]interface{}{"delete_time": time.Now().Format(timeutil.TimeFormat)})

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	xsq_net.Success(c)
}

// 获取仓库用户数列表
func GetWarehouseUserCountList(c *gin.Context) {
	var (
		form req.WarehouseUserCountForm
		res  []*rsp.GetWarehouseUserCountListRsp
	)

	err := c.ShouldBind(&form)

	if err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Raw("SELECT COUNT(u.id) as count,w.id as warehouse_id, w.warehouse_name FROM t_warehouse as w left join `t_user` as u on u.warehouse_id = w.id GROUP BY w.id")

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	result.Scan(&res)

	xsq_net.SucJson(c, res)
}

func GenderPwd(pwd string) string {
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(pwd, options)
	return fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, encodedPwd)
}

// 获取拣货员列表
func GetPickerList(c *gin.Context) {

	var (
		form  req.GetPickerListReq
		users []model.User
		res   []rsp.GetPickerListRsp
	)

	if err := c.ShouldBind(&form); err != nil {
		xsq_net.ErrorJSON(c, ecode.ParamInvalid)
		return
	}

	result := global.DB.Model(&model.User{}).Where("role_id = 3").Find(&users)

	if result.Error != nil {
		xsq_net.ErrorJSON(c, result.Error)
		return
	}

	for _, u := range users {
		res = append(res, rsp.GetPickerListRsp{
			Id:   u.Id,
			Name: u.Name,
		})
	}

	xsq_net.SucJson(c, res)
}

// 校验登录是否过期
func Check(c *gin.Context) {
	//只需要中间件校验登录是否过期即可
	xsq_net.Success(c)
}
