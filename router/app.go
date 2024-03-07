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
	r.POST("/problem-create", middleware.AuthAdminCheck(), service.ProblemCreate)
	r.GET("/category-list", middleware.AuthAdminCheck(), service.GetCategoryList)

	//分类创建
	r.POST("/category-create", middleware.AuthAdminCheck(), service.CategoryCreate)
	//分类删除
	r.POST("/category-delete", middleware.AuthAdminCheck(), service.CategoryDelete)
	//分类修改
	r.PUT("/category-update", middleware.AuthAdminCheck(), service.CategoryUpdate)

	//删除问题
	r.POST("/problem-delete", middleware.AuthAdminCheck(), service.ProblemDelete)
	//修改问题
	r.PUT("/problem-update", middleware.AuthAdminCheck(), service.ProblemUpdate)

	return r
}
