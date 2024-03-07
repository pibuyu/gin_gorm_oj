package service

import (
	"encoding/json"
	"gin_gorm_o/define"
	"gin_gorm_o/helper"
	"gin_gorm_o/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// GetProblemList
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

// ProblemDelete
// @Tags 管理员私有方法
// @Summary 删除问题
// @Param authorization header string true "authorization"
// @Param identity formData string true "问题的唯一标识"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-delete [post]
func ProblemDelete(c *gin.Context) {
	problemIdentity := c.PostForm("identity")
	if problemIdentity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem identity can not be empty",
		})
		return
	}

	/*
		删除问题之前先尝试删除problem_category中问题和分类的关联数据项.
		本来想对err额外做一次ErrNotFound判断，发现不需要,
		因为delete操作实际上是update语句，当没有问题和分类的关联项时，update影响的行数为0，不会报错。
	*/
	err := models.DB.Debug().Model(&models.ProblemCategory{}).
		Where("problem_id =(select id from problem_basic where identity = ?)", problemIdentity).
		Delete(&models.ProblemCategory{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Delete Problem-related Category Error:" + err.Error(),
		})
		return
	}

	//删除分类之后，再删除问题
	err = models.DB.Debug().Model(&models.ProblemBasic{}).
		Where("identity = ?", problemIdentity).Delete(&models.ProblemBasic{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Delete Problem Error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"msg":                     "delete problem success!",
			"deleted prblem identity": problemIdentity,
		},
	})
}

// ProblemUpdate
// @Tags 管理员私有方法
// @Summary 问题更新
// @Param authorization header string true "authorization"
// @Param identity formData string true "问题的唯一标识"
// @Param title formData string false "问题标题"
// @Param content formData string false "问题内容"
// @Param max_mem formData int false "最大运行内存"
// @Param max_runtime formData int false "最大运行时间"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array false "test_cases"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /problem-update [put]
func ProblemUpdate(c *gin.Context) {
	identity := c.PostForm("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))

	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem identity can not be empty",
		})
		return
	}

	oneProblem := &models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxMem:     maxMem,
		MaxRuntime: maxRuntime,
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		/*问题基础信息保存
		不需要管用户没传的字段，因为update方法只修改了update_at和用户传进来的字段。
		没传的字段保持不动的
		*/

		err := tx.Debug().Model(new(models.ProblemBasic)).
			Where("identity = ?", identity).Updates(oneProblem).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "in update problem transaction,save problem-basic error: " + err.Error(),
			})
			return err
		}

		/*问题所属分类保存*/
		//先获取problem的id
		err = tx.Debug().
			Select("id").Where("identity = ?", identity).Find(&oneProblem).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "get problem-id error: " + err.Error(),
			})
			return err
		}

		//然后去problem_category表插入problem_id和category_id关联项
		if categoryIds == nil {
			return nil
		}
		//1、删除已经存在的关联项目
		err = tx.Debug().Where("problem_id = ?", oneProblem.ID).Delete(new(models.ProblemCategory)).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "delete existed problem-category item error: " + err.Error(),
			})
			return err
		}
		//2、增加新的关联项
		problemCategories := make([]*models.ProblemCategory, 0)

		//for _, id := range categoryIds {
		//	intId, _ := strconv.Atoi(id)
		//	problemCategories = append(problemCategories, &models.ProblemCategory{
		//		ProblemId:  oneProblem.ID,
		//		CategoryId: uint(intId),
		//	})
		//}

		//categoryIds是一个[]string类型,里面只包含一个元素：[category_1,category_2...],我要用categoryIds[0]先取出除方括号的字符串，然后以逗号切割成slice，然后挨个取category_id
		idsSlice := strings.Split(categoryIds[0], ",")
		for _, id := range idsSlice {
			intId, _ := strconv.Atoi(id)
			problemCategories = append(problemCategories, &models.ProblemCategory{
				ProblemId:  oneProblem.ID,
				CategoryId: uint(intId),
			})
		}
		err = tx.Debug().Create(&problemCategories).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "create problem-category item error: " + err.Error(),
			})
			return err
		}

		/*问题测试用例保存*/
		if testCases == nil {
			return nil
		}
		//1、删除已有的问题和测试用例的关联项
		err = tx.Debug().Where("problem_identity = ?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "delete existed problem-testcase item error: " + err.Error(),
			})
			return err
		}
		//2、创建新的问题和测试用例的关联项
		testCaseBasics := make([]*models.TestCase, 0)
		for _, testCase := range testCases {
			caseMap := make(map[string]string)
			err := json.Unmarshal([]byte(testCase), &caseMap)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "testCase to map error:" + err.Error(),
				})
				return err
			}

			if _, ok := caseMap["input"]; !ok {
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "测试用例的格式错误",
				})
				return err
			}
			if _, ok := caseMap["output"]; !ok {
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  "测试用例的格式错误",
				})
				return err
			}

			oneCaseBasic := &models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			}
			testCaseBasics = append(testCaseBasics, oneCaseBasic)
		}
		err = tx.Debug().Create(&testCaseBasics).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "create problem_testcase error : " + err.Error(),
			})
		}

		//上面三项事务都没问题，返回nil
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "update problem transaction error: " + err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":                     200,
		"msg":                      "update problem success!",
		"updated problem identity": identity,
		"data": map[string]interface{}{
			"problemInfo": oneProblem,
		},
	})
}
