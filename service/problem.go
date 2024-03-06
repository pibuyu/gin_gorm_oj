package service

import (
	"gin_gorm_o/define"
	"gin_gorm_o/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetProoblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "输入页码，默认第一页"
// @Param size query int false "页面大小，默认为20"
// @Param keyword query string false "查询关键词，进行模糊查询"
// @Param category_identity query string false "分类的唯一标识"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetProbelmList Page strconv Error:", err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProbelmList Size strconv Error:", err)
		return
	}

	page = (page - 1) * size
	keyword := c.Query("keyword")
	categoryIdentity := c.Query("category_identity")
	var count int64

	problemList := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
	/*
		tx.Count()对查询结果计数，并复制给count
		tx.Omit()在查询结果里过滤掉比较长的content，返回简短的结果
	*/
	err = tx.Count(&count).Omit("content").Offset(page).Limit(size).Find(&problemList).Error
	if err != nil {
		log.Println("Get Probelm List Error:", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"data":             problemList,
			"TotalResultCount": count,
		},
	})
}

// GetProoblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	problemBasic := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&problemBasic).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Detail Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": problemBasic,
	})
}
