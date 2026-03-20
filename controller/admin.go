package controller

import (
	"encoding/json"
	"fmt"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/int-xiaoli/go_gateway/dao"
	"github.com/int-xiaoli/go_gateway/dto"
	"github.com/int-xiaoli/go_gateway/middleware"
	"github.com/int-xiaoli/go_gateway/public"
)

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) { //子路由注册
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

// adminLogin godoc
// @Summary 获取管理员信息
// @Description 获取当前管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminLogin *AdminController) AdminInfo(c *gin.Context) {
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), //这一步把sessinfo信息转换成结构体
		adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//TODO
	//1.读取sessionKey对应json转换为结构体
	//2.取出数据然后封装输出结构体

	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://image2url.com/r2/default/gifs/1773841228466-39b87e59-7101-4234-976e-e9bab6b99e4a.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	} //输出
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminLogin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil { //这一步是请求内容绑定结构体params，执行到c.ShouldBind(params)才是具体的绑定
		middleware.ResponseError(c, 2000, err)
		return
	} //这是获取请求内容中参数
	//TODO
	//1.session获取用户信息到结构体 sessInfo
	//2.sessInfo.ID  读取数据库信息 adminInfo
	//3.params.password+adminInfo.salt 进行sha256加密=》saltPassword
	//4.saltPassword==>adminInfo.password 执行数据保存

	//session获取用户信息
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), //这一步把sessinfo信息转换成结构体
		adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//获取数据连接池
	tx, err := lib.GetGormPool("default") //连接数据库
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//读取数据库信息
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//生成新密码 saltPassword
	saltPassword := public.GetSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword
	//保存数据
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	middleware.ResponseSuccess(c, "")

}
