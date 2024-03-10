package router

import (
	_ "gin_gorm_o/docs"
	"gin_gorm_o/middleware"
	"gin_gorm_o/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()

	//swagger配置
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//配置路由规则
	r.GET("/", service.Hello)

	/*公用方法*/
	//问题
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)
	r.GET("/rank-list", service.GetRankList)

	//用户
	r.GET("/user-detail", service.GetUserDeatil)
	r.GET("/send-code", service.SendCode)
	r.POST("/register", service.Register) //要用post方法，用get方法报404，找半天错，淦
	r.POST("/login", service.LogIn)

	//提交记录
	r.GET("/submit-list", service.GetSubmitList)

	/*管理员私有方法*/
	authAdmin := r.Group("/admin").Use(middleware.AuthAdminCheck())

	//获取分类列表
	authAdmin.GET("/category-list", service.GetCategoryList)
	//分类创建
	authAdmin.POST("/category-create", service.CategoryCreate)
	//分类删除
	authAdmin.POST("/category-delete", service.CategoryDelete)
	//分类修改
	authAdmin.PUT("/category-update", service.CategoryUpdate)

	//创建问题
	authAdmin.POST("/problem-create", service.ProblemCreate)
	//删除问题
	authAdmin.POST("/problem-delete", service.ProblemDelete)
	//修改问题
	authAdmin.PUT("/problem-update", service.ProblemUpdate)

	/*用户私有方法*/
	authUser := r.Group("/user").Use(middleware.AuthUserCheck())
	//代码提交
	authUser.POST("/submit-code", service.SubmitCode)
	return r
}
