package main

import (
	"github.com/gin-gonic/gin"
	"main/AccessControl"
	"main/HomeworkUpload"
	"main/Logger"
	"net/http"
)

func helloWorldCallback(r *gin.Engine) {
	r.GET("/hello", func(c *gin.Context) {
		acc, _, _, _, _, _ := AccessControl.AccControl(c)
		if !acc {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})
}

const debug = false

func main() {
	var err error
	//err := Sql.ConnectPSQL()
	//if err != nil {
	//	Logger.Log(Logger.LogMerge, err.Error())
	//	return
	//}
	Logger.SetLogLevel(Logger.LogDebug)
	var r *gin.Engine
	if debug {
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	//r := gin.Default()
	r.Use(AccessControl.AccMiddleware())
	//OAuth.AuGithubCallback(r)
	//FileUpload.UploadCallback(r)
	//Account.RegisterCallback(r)
	//Login.PasswdLoginCallback(r)
	err = HomeworkUpload.InitLoad("conf.yaml")
	if err != nil {
		return
	}
	//err = HomeworkUpload.LoadStudents("tmp.xls")
	//if err != nil {
	//	return
	//}
	//err = Config.WriteYaml("conf.yaml", &HomeworkUpload.HwControl)
	//if err != nil {
	//	return
	//}
	for index, st := range HomeworkUpload.HwControl.Students {
		HomeworkUpload.Students = append(HomeworkUpload.Students, st)
		HomeworkUpload.StudentNumMap[st.StudentNum] = index
		HomeworkUpload.StudentNameMap[st.StudentName] = index
	}
	Logger.LogI("Total %d students", len(HomeworkUpload.Students))
	err = HomeworkUpload.InitDynamicInfo()
	if err != nil {
		return
	}
	HomeworkUpload.GenTimer()
	HomeworkUpload.ReportUploadCallback(r)
	HomeworkUpload.FileDownloadCallback(r)
	HomeworkUpload.HwWsCallback(r)
	//helloWorldCallback(r)
	// runTLS
	Logger.LogI("Started!")
	if debug {
		err = r.RunTLS("0.0.0.0:3001", "./ssl/loc.kazefx.top.pem", "./ssl/loc.kazefx.top.key")
	} else {
		err = r.RunTLS("0.0.0.0:3001", "./ssl/kazefx.top.pem", "./ssl/kazefx.top.key")
	}
	Logger.LogE("Error Closed!\n%v", err)
	// err = r.Run("0.0.0.0:3001")
	if err != nil {
		return
	}
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
