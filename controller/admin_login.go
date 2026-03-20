package controller

import (
	"encoding/json"
	"time"

	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/int-xiaoli/go_gateway/dao"
	"github.com/int-xiaoli/go_gateway/dto"
	"github.com/int-xiaoli/go_gateway/middleware"
	"github.com/int-xiaoli/go_gateway/public"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLoginOut)

}

// adminLogin godoc
// @Summary 管理员登录
// @Description 管理员登录
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil { //这一步是请求内容绑定结构体params，执行到c.ShouldBind(params)才是具体的绑定
		middleware.ResponseError(c, 2000, err)
		return
	}
	//todo 具体登录的业务逻辑
	//1.params.UserName取得管理员信息 admininfo
	//2.admininfo.salt+params.Password 进行sha256加密=》saltPassword
	//3.saltPassword ==admininfo.password

	tx, err := lib.GetGormPool("default") //连接数据库
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params) //登录密码校验,
	if err != nil {                              // 同时返回的admin中也是数据表中符合输入用户名的数据库中管理员信息
		middleware.ResponseError(c, 2002, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName} //输出
	middleware.ResponseSuccess(c, out)
} //这里是/login的接口实现方法

// adminLogin godoc
// @Summary 管理员退出
// @Description 管理员退出
// @Tags 管理员接口
// @ID /admin_login/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/login_out [get]
func (adminLogin *AdminLoginController) AdminLoginOut(c *gin.Context) {

	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()

	middleware.ResponseSuccess(c, "")
} //退出接口
