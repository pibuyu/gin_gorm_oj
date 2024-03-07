package service

import (
	"gin_gorm_o/define"
	"gin_gorm_o/helper"
	"gin_gorm_o/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 分类列表
// @Param authorization header string true "authorization"
// @Param page query int false "输入页码，默认第一页"
// @Param size query int false "页面大小，默认为20"
// @Param keyword query string false "查询关键词，进行模糊查询"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-list [get]
func GetCategoryList(c *gin.Context) {
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
	var count int64

	categoryList := make([]*models.CategoryBasic, 0)
	/*不用判断keyword是否为空，然后tx.where()，直接连着判断就行，keyword为空不影响查询结果*/
	err = models.DB.Model(new(models.CategoryBasic)).Count(&count).Where("name like ?", "%"+keyword+"%").Offset(page).Limit(size).Find(&categoryList).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Category List Error:" + err.Error(),
		})
		return
	}
	log.Println(categoryList)
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"list":             categoryList,
			"TotalResultCount": count,
		},
	})
	return
}

// CategoryCreate
// @Tags 管理员私有方法
// @Summary 分类创建
// @Param authorization header string true "authorization"
// @Param name formData string true "分类名"
// @Param parent_id formData int false "父类名"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-create [post]
func CategoryCreate(c *gin.Context) {
	categoryName := c.PostForm("name")
	categoryParentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	categoryIdentity := helper.GetUUID()
	oneCategory := &models.CategoryBasic{
		Name:     categoryName,
		ParentId: categoryParentId,
		Identity: categoryIdentity,
	}
	err := models.DB.Debug().Model(new(models.CategoryBasic)).Create(oneCategory).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Create Category Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"msg":          "category create success",
			"categoryInfo": oneCategory,
		},
	})
	return
}

// CategoryDelete
// @Tags 管理员私有方法
// @Summary 分类删除
// @Param authorization header string true "authorization"
// @Param name formData string true "分类名"
// @Param parent_id formData int false "父类名"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-delete [post]
func CategoryDelete(c *gin.Context) {
	categoryName := c.PostForm("name")
	categoryParentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	err := models.DB.Model(new(models.CategoryBasic)).Where("name = ? and parent_id = ?", categoryName, categoryParentId).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Delete Category Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "category delete success,deleted category name: " + categoryName + ", parent_id:" + strconv.Itoa(categoryParentId),
	})
	return
}

// CategoryUpdate
// @Tags 管理员私有方法
// @Summary 分类更新
// @Param authorization header string true "authorization"
// @Param name formData string true "分类名"
// @Param identity formData string false "分类唯一标识"
// @Param parent_id formData int false "父类名"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /category-update [post]
func CategoryUpdate(c *gin.Context) {
	categoryName := c.PostForm("name")
	categoryIdentity := c.PostForm("identity")
	categoryParentId, _ := strconv.Atoi(c.PostForm("parent_id"))
	newCategory := &models.CategoryBasic{
		Name:     categoryName,
		ParentId: categoryParentId,
		Identity: categoryIdentity,
	}
	log.Println(newCategory)
	err := models.DB.Debug().Model(new(models.CategoryBasic)).Where("name = ?", categoryName).
		Updates(newCategory).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Update Category Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"msg":          "category update success",
			"categoryInfo": newCategory,
		},
	})
	return
}
