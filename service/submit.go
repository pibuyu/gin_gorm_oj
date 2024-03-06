package service

import (
	"gin_gorm_o/define"
	"gin_gorm_o/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetProoblemDetail
// @Tags 公共方法
// @Summary 提交记录列表
// @Param page query int false "输入页码，默认第一页"
// @Param size query int false "页面大小，默认为20"
// @Param problem_identity query string false "问题的唯一标识"
// @Param user_identity query string false "用户的唯一标识"
// @Param status query int false "提交状态"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetSubmitList Page strconv Error:", err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetSubmitList Size strconv Error:", err)
		return
	}

	page = (page - 1) * size
	var count int64
	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)

	data := make([]models.SubmitBasic, 0)
	err = tx.Count(&count).Offset(page).Limit(size).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Submit List Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"TotalResultCount": count,
			"data":             data,
		},
	})
}
