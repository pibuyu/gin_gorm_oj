package service

import (
	"context"
	"encoding/json"
	"fmt"
	"gin_gorm_o/consts"
	"gin_gorm_o/define"
	"gin_gorm_o/helper"
	"gin_gorm_o/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

// GetProoblemDetail
// @Tags 公共方法
// @Summary 用户详情
// @Param identity query string false "user identity"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user-detail [get]
func GetUserDeatil(c *gin.Context) {
	identity := c.Query("identity")

	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "用户唯一标识不能为空",
		})
		return
	}

	data := new(models.UserBasic)
	err := models.DB.Where("identity = ?", identity).Find(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User Detail By Identity " + identity + " Error:" + err.Error(),
		})
		return
	}
	//return results
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": data,
	})
}

// GetProblemDetail
// @Tags 公共方法
// @Summary 发送验证码
// @Param toUserEmail query string false "to user email"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /send-code [get]
func SendCode(c *gin.Context) {
	toUserEmail := c.Query("toUserEmail")
	if toUserEmail == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "To user email can`t be null!",
		})
		return
	}

	//生成6位随机验证码
	generateCode := func(width int) string {
		numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		r := len(numeric)
		rand.Seed(time.Now().UnixNano())

		var sb strings.Builder
		for i := 0; i < width; i++ {
			fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
		}
		return sb.String()
	}
	code := generateCode(6)

	err := helper.SendVerifyCode(toUserEmail, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Send Code Error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Send Code Success,the code is " + code,
	})

	models.Redis.Set(ctx, "code", code, time.Minute*5)
	return

}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param name formData string true "userName"
// @Param password formData string true "password"
// @Param email formData string true "email"
// @Param code formData string true "verifyCode"
// @Param phone formData string false "userPhoneNumber"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	name := c.PostForm("name")
	password := c.PostForm("password")
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	phone := c.PostForm("phone")
	if name == "" || password == "" || email == "" || userCode == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "required information can`t be null!",
		})
		return
	}
	sysCode, err := models.Redis.Get(ctx, "code").Result()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "get code from redis error: " + err.Error(),
		})
		return
	}
	if userCode != sysCode {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "input code isn`t correct!",
		})
		return
	}

	//验证一下是否已经存在同名用户，如果没有，会报一个record not found，但不是错误
	oneUser := new(models.UserBasic)
	err = models.DB.Model(&models.UserBasic{}).Where("name = ?", name).First(&oneUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询是否有同名用户的过程中出错：" + err.Error(),
		})
		return
	}
	if oneUser.Password != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该用户名已被使用！",
		})
		return
	}

	uuid := helper.GetUUID()
	insertUser := &models.UserBasic{
		Identity: uuid,
		Name:     name,
		Password: helper.MD5(password),
		Phone:    phone,
		Mail:     email,
	}
	err = models.DB.Create(insertUser).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create user error:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "register success!",
		"data": map[string]interface{}{
			"userInfo": insertUser,
		},
	})
	return
}

// LogIn
// @Tags 公共方法
// @Summary 用户登录
// @Param name formData string true "userName"
// @Param password formData string true "password"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /login [post]
func LogIn(c *gin.Context) {
	name := c.PostForm("name")
	inputPassword := c.PostForm("password")
	if name == "" || inputPassword == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "required information can`t be null!",
		})
		return
	}
	inputPassword = helper.MD5(inputPassword)

	user := &models.UserBasic{}
	/*可以通过比较密码看用户是否存在*/
	//models.DB.Where("name = ?", name).Select("password").Find(&user)
	//if user.Password != inputPassword {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"msg":  "name or password is wrong!",
	//	})
	//	return
	//}

	/*也可以看能不能根据name和password查出记录来判断*/
	err := models.DB.Where("name =? and password =?", name, inputPassword).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound { //看是否返回ErrRecordNotFound错误来判断有没有查询到结果
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "User doesn`t exist!It may be caused by wrong name or password.",
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Query user from DB error:" + err.Error(),
			})
			return
		}
	}

	token, err := helper.GenerateToken(user.Name, user.Identity, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Generate token error:" + err.Error(),
		})
		return
	}

	//登录的时候把token写进redis
	redisKey := consts.LOG_IN_TOKEN_KEY + user.Name
	err = models.Redis.Set(ctx, redisKey, token, time.Hour*24*7).Err()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Set token to redis error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  200,
		"token": token,
		"msg":   "Login success!",
		"userInfo": map[string]interface{}{
			"name":     user.Name,
			"identity": user.Identity,
			"phone":    user.Phone,
			"email":    user.Mail,
		},
	})
	return
}

// GetRankList
// @Tags 公共方法
// @Summary 提交排行榜
// @Param page query string false "page"
// @Param size query string false "size"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetRankList Page strconv Error:", err)
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetRankList Size strconv Error:", err)
		return
	}

	page = (page - 1) * size
	var count int64

	userList := make([]*models.UserBasic, 0)
	err = models.DB.Model(new(models.UserBasic)).Count(&count).Order("finish_problem_num DESC,submit_num ASC").
		Offset(page).Limit(size).Find(&userList).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Rank List Error:" + err.Error(),
		})
		return
	}

	//封装成更简单的格式返回，只保留name、finish_problem_num和submit_num字段
	rankList := make([]*models.UserRank, 0)
	for _, user := range userList {
		user_rank := &models.UserRank{
			Name:             user.Name,
			FinishProblemNum: user.FinishProblemNum,
			SubmitNum:        user.SubmitNum,
		}
		rankList = append(rankList, user_rank)
	}
	//封装完毕，返回简单结果
	c.JSON(http.StatusOK, gin.H{
		"code":                   -1,
		"currentPageResultCount": count,
		"data": map[string]interface{}{
			"list": rankList,
		},
	})
	return

}

// ProblemCreate
// @Tags 管理员私有方法
// @Summary 创建问题
// @Param authorization header string true "token"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @param max_mem formData int false "max_mem"
// @param max_runtime formData int false "max_runtime"
// @Param category_ids formData array false "category_ids"
// @Param test_cases formData array true "test_cases"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /admin/problem-create [post]
func ProblemCreate(c *gin.Context) {
	//log.Println("problem create开始执行")
	title := c.PostForm("title")
	content := c.PostForm("content")
	max_mem, _ := strconv.Atoi(c.PostForm("max_mem"))
	max_runtime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	category_ids := c.PostFormArray("category_ids")
	test_cases := c.PostFormArray("test_cases")
	if title == "" || content == "" || test_cases == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  " create problem : required information can`t be null!",
		})
		return
	}

	problemIdentity := helper.GetUUID()
	data := &models.ProblemBasic{
		Identity:   problemIdentity,
		Title:      title,
		Content:    content,
		MaxMem:     max_mem,
		MaxRuntime: max_runtime,
	}

	//处理问题所属分类
	categoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range category_ids {
		//一个个id就是用户输入的问题所属分类的id
		id, _ := strconv.Atoi(id)
		categoryBasics = append(categoryBasics, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(id),
		})
	}
	data.ProblemCategories = categoryBasics

	//处理问题的测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, testCase := range test_cases {
		// {"input":"1 2\n","output":"3"}
		caseMap := make(map[string]string)
		err := json.Unmarshal([]byte(testCase), &caseMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "testCase to map error:" + err.Error(),
			})
			return
		}

		if _, ok := caseMap["input"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例的格式错误",
			})
			return
		}
		if _, ok := caseMap["output"]; !ok {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例的格式错误",
			})
			return
		}

		oneCaseBasic := &models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: problemIdentity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		}
		testCaseBasics = append(testCaseBasics, oneCaseBasic)
	}
	data.TestCases = testCaseBasics

	//创建问题
	err := models.DB.Create(&data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "create problem error:" + err.Error(),
		})
		return
	}
	//创建问题成功，返回新建问题的identity
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"Create Problem Success!Problem Identity is": data.Identity,
		},
	})
	return
}
