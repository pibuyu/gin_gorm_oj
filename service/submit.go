package service

import (
	"bytes"
	"errors"
	"gin_gorm_o/define"
	"gin_gorm_o/helper"
	"gin_gorm_o/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
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

// SubmitCode
// @Tags 用户私有方法
// @Summary 提交代码
// @Param authorization header string true "token"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200","data":""}"
// @Router /user/submit-code [post]
func SubmitCode(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	code, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Read Code Error:" + err.Error(),
		})
		return
	}

	//代码保存
	path, err := helper.CodeSave(code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Save Code Error:" + err.Error(),
		})
		return
	}

	u, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get User from from request body Error",
		})
		return
	}
	userClaim := u.(*helper.UserClaims) //类型强转为UserClaims

	//生成一个submit_basic并插入
	submitBasic := &models.SubmitBasic{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaim.Identity,
		Path:            path,
	}

	//判断代码的正确性，返回WA,OOM等信息
	//首先查询这个problem对应的basic
	problemBasic := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(problemBasic).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get Problem Basic Error:" + err.Error(),
		})
		return
	}
	//遍历测试用例并执行
	WA := make(chan int)
	OOM := make(chan int)
	CE := make(chan int)
	passCount := 0
	//提示执行结果
	var msg string
	var lock sync.Mutex

	for _, testCase := range problemBasic.TestCases {
		testCase := testCase
		go func() {
			//执行测试
			cmd := exec.Command("go", "run", path)
			var out, stderr bytes.Buffer

			cmd.Stderr = &stderr
			cmd.Stdout = &out
			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}

			io.WriteString(stdinPipe, testCase.Input)
			//运行给定的code，比对拿到的输出结果和正确的输出结果
			var beginMem runtime.MemStats
			runtime.ReadMemStats(&beginMem) //记录运行前的内存状态
			err = cmd.Run()
			if err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					CE <- 1
					msg = stderr.String()
					return
				}
			}

			var endMem runtime.MemStats
			runtime.ReadMemStats(&endMem) //记录运行后的内存状态

			//答案错误
			log.Println("程序输出为: ", out.String())
			if testCase.Output != out.String() {

				WA <- 1
				msg = "答案错误"
				return
			}
			//运行超内存
			log.Println("运行所需内存为: ", endMem.Alloc/1024-(beginMem.Alloc/1024))
			if endMem.Alloc/1024-(beginMem.Alloc/1024) > uint64(problemBasic.MaxMem) {
				OOM <- 1
				msg = "运行超内存"
				return
			}

			lock.Lock()
			passCount++
			log.Println("通过的测试案例数：", passCount)
			lock.Unlock()
		}()
	}
	//-1-待判断；1-答案正确；2-答案错误；3-运行超时；4-运行超内存;5-编译错误
	select {
	case <-WA:
		submitBasic.Status = 2
	case <-OOM:
		submitBasic.Status = 4
	case <-CE:
		submitBasic.Status = 5
	case <-time.After(time.Millisecond * time.Duration(problemBasic.MaxRuntime)):
		if passCount == len(problemBasic.TestCases) {
			submitBasic.Status = 1
			msg = "答案正确"
		} else {
			submitBasic.Status = 3
			msg = "运行超时"
		}
	}

	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		//保存代码
		err = tx.Create(submitBasic).Error
		if err != nil {
			return errors.New("Save Submit Basic Error:" + err.Error())
		}
		//需要插入的数据
		m := make(map[string]interface{})
		m["submit_num"] = gorm.Expr("submit_num+?", 1) //无论对错，提交次数+1
		if submitBasic.Status == 1 {
			m["pass_num"] = gorm.Expr("pass_num+?", 1) //答案全部正确时，通过数+1
		}
		//更新user_basic中的pass_num和submit_num
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaim.Identity).Updates(m).Error
		if err != nil {
			return errors.New("Update User Basic submit_num and pass_num Error:" + err.Error())
		}
		//更新problem_basic中的pass_num和submit_num
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("Update Problem Basic submit_num and pass_num Error:" + err.Error())
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Transaction Error:" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": map[string]interface{}{
			"status": submitBasic.Status,
			"msg":    msg,
		},
	})

}
