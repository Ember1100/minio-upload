package routers

import (
	"fmt"
	"net/http"
	"strconv"

	"minio-upload-go-api/core"
	"minio-upload-go-api/models"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	arr, err := core.FindData()
	if err != nil {
		fmt.Print(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询学生信息成功",
		"data":    arr,
	})
}

func TaskByIdentifierHandler(c *gin.Context) {
	identifier := c.Param("identifier")
	data, err := core.GetByIdentifier(identifier)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "成功",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "查询学生信息成功",
		"data":    data,
	})
}

func TestUploadToMinIO(c *gin.Context) {
	objectKey := "富邦财险_新核心系统实施项目_反洗钱系统设计说明书_v0.2_20230228.doc"
	filePath := "C:/Home/File/huangqi/work/file/富邦财险_新核心系统实施项目_反洗钱系统设计说明书_v0.2_20230228.doc"
	core.UploadToMinIO(objectKey, filePath)
	c.JSON(http.StatusOK, gin.H{
		"message": "上传完成文件成功",
		"code":    "000",
	})
}

func InitTask(c *gin.Context) {
	var param models.InitTaskParam

	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%+v\n", param)
	data, err := core.InitTask(param)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"code":    200000,
		"data":    data,
	})
}

func GenPreSignUploadURL(c *gin.Context) {
	identifier := c.Param("identifier")
	task, err := core.GetByIdentifier(identifier)
	acc := models.SysUploadTask{}
	if task == acc || err != nil {
		fmt.Print(task)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "分片任务不存在",
			"error":   200001,
		})
		return
	}
	partNumber := c.Param("partNumber")
	number, err := strconv.ParseInt(partNumber, 10, 64)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "失败",
			"error":   err.Error(),
		})
		return
	}

	data, err := core.GenPreSignUploadURL(task.ObjectKey, task.UploadID, number)
	// 使用ShouldBindJSON方法将请求体中的JSON数据绑定到param变量
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "失败",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"code":    200000,
		"data":    data,
	})
}

func GetTaskInfo(c *gin.Context) {
	identifier := c.Param("identifier")
	data := core.GetTaskInfo(identifier)
	// 使用ShouldBindJSON方法将请求体中的JSON数据绑定到param变量

	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"code":    200000,
		"data":    data,
	})
}
func Merge(c *gin.Context) {
	identifier := c.Param("identifier")
	core.Merge(identifier)
	// 使用ShouldBindJSON方法将请求体中的JSON数据绑定到param变量

	c.JSON(http.StatusOK, gin.H{
		"message": "成功",
		"code":    200000,
		"data":    nil,
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(Cors())
	// router.LoadHTMLGlob("./templates/*")
	// router.GET("/demo", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{})
	// })
	router.GET("/GetTasks", PingHandler)
	router.GET("test_upload", TestUploadToMinIO)
	router.POST("/v1/minio/tasks", InitTask)
	router.GET("/v1/minio/tasks/:identifier/:partNumber", GenPreSignUploadURL)
	router.GET("/v1/minio/tasks/:identifier", GetTaskInfo)
	router.POST("/v1/minio/tasks/merge/:identifier", Merge)

	return router
}
