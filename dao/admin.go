package dao

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/int-xiaoli/go_gateway/dto"
	"github.com/int-xiaoli/go_gateway/public"
	"gorm.io/gorm"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}
func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB,
	params *dto.AdminLoginInput) (*Admin, error) {

	search := &Admin{
		UserName: params.UserName,
		IsDelete: 0,
	}
	adminInfo, err := t.Find(c, tx, search) //这里其实是为了查询后台数据中用户信息的密码hash值
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	saltPassword := public.GetSaltPassword(adminInfo.Salt, params.Password)
	if saltPassword != adminInfo.Password {
		return nil, errors.New("密码错误,请重新输入")
	}
	return adminInfo, nil
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	result := tx.WithContext(c.Request.Context()).Where(search).Find(out)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return out, nil
} //查询接口 支持结构体类型查询

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.WithContext(c).Save(t).Error
}
