package zdy_tools

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/lfxnxf/zdy_tools/inits"
	"github.com/lfxnxf/zdy_tools/inits/proxy"
	"github.com/lfxnxf/zdy_tools/zd_error"
	"github.com/lfxnxf/zdy_tools/zd_http"
)

type TagKey struct {
	ID           int64     `gorm:"column:id" json:"id"`                                 // 自增id
	Name         string    `gorm:"column:name" json:"name"`                             // 键
	DepartmentId int64     `gorm:"column:department_id" json:"department_id"`           // 部门id
	TenantId     int64     `gorm:"column:tenant_id" json:"tenant_id"`                   // 租户id
	Type         int64     `gorm:"column:type" json:"type"`                             // 类型，1：租户，2：公共
	Category     int64     `gorm:"column:category" json:"category"`                     // 分类 1、项目，2、技术，3、财务
	Creator      string    `gorm:"column:creator" json:"creator"`                       // 创建人
	Modifier     string    `gorm:"column:modifier" json:"modifier"`                     // 修改人
	IsDeleted    int64     `gorm:"column:is_deleted" json:"is_deleted"`                 // 是否被删除
	CreateTime   time.Time `gorm:"column:create_time; default:null" json:"create_time"` // 创建时间
}

func Test_Main(t *testing.T) {
	var cfg Config
	inits.Init(
		inits.ConfigPath("./config.yaml"),
		inits.Once(),
		inits.LoadLocalConfig(&cfg),
	)

	db := proxy.InitSQL("recourse-center.mysql")
	var list []TagKey
	err := db.Master().Table("tag_key").Find(&list).Error
	if err != nil {
		panic(err)
	}

	s := inits.NewHttpServer(cfg.Server)

	s.GET("/test", func(c *gin.Context) {
		zd_http.WriteJson(c, map[int]string{1: "word"}, zd_error.AddError("test error", "error"))
	})

	if err := s.StartHttp(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
