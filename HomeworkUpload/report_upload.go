package HomeworkUpload

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/FileUpload"
	"main/Logger"
	"main/Msg"
	"main/QuickRes"
	"main/StatusCode"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func genTimeCode(Time time.Time) string {
	return fmt.Sprintf("%02d%02d%02d%09d", Time.Hour(), Time.Minute(), Time.Second(), Time.Nanosecond())
}

func studentNumNameCheck(num string, name string) bool {
	indexNum, exist := StudentNumMap[num]
	if !exist {
		return false
	}
	indexName, exist := StudentNameMap[name]
	if !exist {
		return false
	}
	if indexNum != indexName {
		return false
	}
	return true
}

func postValidation(c *gin.Context) (code int, msg string, ava bool, saveName string, studentNum string, studentName string, experimentNum int, file *multipart.FileHeader, hash string, group int) {
	// field names:
	// student_num, student_name, exp_num, report_file
	code = 0
	msg, saveName, studentNum, studentName, expNumStr, hash, sha256 := "", "", "", "", "", "", ""
	groupNum := -1
	groupExp := false
	experimentNum = -1
	file = nil
	var err error = nil
	regStudentNum := `^[A-Z][0-9]{8}$`
	regFileName := `\.(doc|docx)$`
	regFileNameGroup := `\.(zip)$`
	regNumMatch := regexp.MustCompile(regStudentNum)
	regFileMatch := regexp.MustCompile(regFileName)
	regFileGroupMatch := regexp.MustCompile(regFileNameGroup)
	expNumStr, exist := c.GetPostForm("exp_num")
	if !exist {
		goto errRet
	}
	expNumStr = strings.ReplaceAll(expNumStr, " ", "")
	experimentNum, err = strconv.Atoi(expNumStr)
	if err != nil {
		goto errRet
	}
	if experimentNum < 1 || experimentNum > HwControl.Control.MaxExp {
		goto errRet
	}
	if !ExpAva[experimentNum-1] || time.Now().Before(Starts[experimentNum-1]) {
		code = -1
		msg = "实验未开放"
		goto errRet
	}
	if time.Now().After(Ends[experimentNum-1]) {
		if time.Now().Before(ReStarts[experimentNum-1]) || time.Now().After(ReEnds[experimentNum-1]) {
			code = -1
			msg = "实验已结束"
			goto errRet
		}
	}
	groupExp = IsGroupExp(experimentNum)
	studentNum, exist = c.GetPostForm("student_num")
	if !exist {
		goto errRet
	}
	studentNum = strings.ToUpper(studentNum)
	studentName, exist = c.GetPostForm("student_name")
	if !exist {
		goto errRet
	}
	studentNum = strings.ReplaceAll(studentNum, " ", "")
	studentName = strings.ReplaceAll(studentName, " ", "")
	file, err = c.FormFile("report_file")
	if err != nil {
		goto errRet
	}
	// string validation
	// student num 'Y12345678'
	if !regNumMatch.MatchString(studentNum) {
		goto errRet
	}
	// file name *.doc / *.docx
	if !groupExp && !regFileMatch.MatchString(file.Filename) {
		goto errRet
	}
	if groupExp && !regFileGroupMatch.MatchString(file.Filename) {
		goto errRet
	}
	// database check
	if !studentNumNameCheck(studentNum, studentName) {
		goto errRet
	}
	// postFix = path.Ext(file.Filename)
	// build formatted name of files to be saved as
	if !groupExp {
		saveName = fmt.Sprintf("第%d次实验%s%s.doc",
			experimentNum, studentNum, studentName)
	} else {
		groupNum = Students[StudentNumMap[studentNum]].StudentGrp
		saveName = fmt.Sprintf("第%d组大作业.zip", groupNum)
	}
	// eg. "第1次实验Y02114xxx张三.doc"
	hash, err = FileUpload.CalculateSHA5256(file)
	// if sha256 has appended
	sha256, exist = c.GetPostForm("file_sha256")
	if exist && sha256 != hash {
		code = -1
		msg = "哈希校验失败，可能是传输错误，或许可以尝试刷新网页后重试"
		goto errRet
	}
	return code, msg, true, saveName, studentNum, studentName, experimentNum, file, hash, groupNum
errRet:
	return code, msg, false, saveName, studentNum, studentName, experimentNum, file, hash, groupNum
}

func uploadQueryCheck(c *gin.Context) (ava bool, exp int, studentNum string) {
	expStr := c.Query("exp")
	if expStr == "" {
		return false, -1, ""
	}
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return false, -1, ""
	}
	if exp < 1 || exp > HwControl.Control.MaxExp {
		return false, -1, ""
	}
	studentNum = c.Query("student")
	if studentNum == "" {
		return false, -1, ""
	}
	return true, exp, studentNum
}

func ReportUploadCallback(r *gin.Engine) {
	// upload
	r.POST("/api/report_upload", func(c *gin.Context) {
		// ===== WE NEED =====
		// 1. Student No.
		// 2. Student Name
		// 3. Experiment No.
		// 4. File (.doc/.docx)
		QuickRes.SetOrigin(c)
		code, msg, ava, saveName, stuNum, stuName, experimentNum, file, hash, group :=
			postValidation(c)
		if !ava {
			if code == 0 {
				QuickRes.BadRequest(c)
				return
			} else {
				c.JSON(http.StatusForbidden, gin.H{
					"status": StatusCode.NotPermitted,
					"msg":    msg,
				})
				return
			}
		}
		expSubDir := fmt.Sprintf("%d", experimentNum)
		targetPath := path.Join(HwControl.Control.SaveDir, expSubDir)
		// dir create
		stat, err := os.Stat(targetPath)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(targetPath, os.ModePerm)
				if err != nil {
					QuickRes.InternalError(c)
					return
				}
			}
		} else if !stat.IsDir() {
			err = os.Remove(targetPath)
			if err != nil {
				QuickRes.InternalError(c)
				return
			}
			err = os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				QuickRes.InternalError(c)
				return
			}
		}
		err = c.SaveUploadedFile(file, path.Join(targetPath, saveName))
		if err != nil {
			QuickRes.InternalError(c)
			return
		}
		st, err := os.Stat(path.Join(targetPath, saveName))
		if err != nil {
			QuickRes.InternalError(c)
			return
		}
		Time := st.ModTime()
		downCode := genTimeCode(Time)
		url := ""
		if group == -1 {
			if !updateStudent(stuNum, stuName, experimentNum, Time, hash) {
				QuickRes.BadRequest(c)
				return
			}
			url = fmt.Sprintf("https://kazefx.top:3001/api/report_download_stu?exp=%d&stu=%s&code=%s", experimentNum, stuNum, downCode)
			Logger.LogI("studentNum: %s, studentName: %s, experimentNum: %d, hash: %s", stuNum, stuName, experimentNum, hash)
			Msg.MessageChan <- Msg.Message{
				Stu: Msg.Student{
					StudentName: stuName,
					StudentNum:  stuNum,
				},
				Msg: fmt.Sprintf("学号: %s, 姓名: %s, 于 %s 上传了第 %d 次实验报告，如非本人操作请立即联系课代表并重新上传\n上传地址: https://kazefx.top/#/report_upload\n备用: http://cdn.kazefx.top/#/report_upload\n文件Sha256：%s\n文件下载：%s", stuNum, stuName, Time.Format("2006-01-02 15:04:05"), experimentNum, hash, url),
			}
		} else {
			if !updateGroup(group, experimentNum, Time, hash) {
				QuickRes.BadRequest(c)
				return
			}
			url = fmt.Sprintf("https://kazefx.top:3001/api/report_download_stu?exp=%d&grp=%d&code=%s", experimentNum, group, downCode)
			Logger.LogI("studentNum: %s, studentName: %s, experimentNum: %d, groupNum: %d, hash: %s", stuNum, stuName, experimentNum, group, hash)
		}
		msg = fmt.Sprintf("上传成功!<br>可以点击<a href=\"%s\">下载文件</a>来下载并检查已上传文件是否正确<br>如有需要请自行保存链接:%s<br><font color=\"red\">注: 该链接将会在再次上传该次实验后失效</font><br>", url, url)
		c.JSON(http.StatusOK, gin.H{
			"status": StatusCode.Success,
			"msg":    msg,
		})
	})
	// get upload_info
	r.GET("/api/report_upload_query", func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		ava, exp, stuNum := uploadQueryCheck(c)
		if !ava {
			QuickRes.BadRequest(c)
			return
		}
		var data FileRecord
		group := -1
		stu := make([]string, 0)
		if IsGroupExp(exp) {
			group = Students[StudentNumMap[stuNum]].StudentGrp
			data = DynamicExpInfo[exp-1].Records[group-1]
			for _, i := range GroupStu[group] {
				stu = append(stu, Students[i].StudentName)
			}
		} else {
			data = DynamicExpInfo[exp-1].Records[StudentNumMap[stuNum]]
		}
		if data.Uploaded {
			c.JSON(http.StatusOK, gin.H{
				"status": StatusCode.Success,
				"msg":    fmt.Sprintf("已上传！<br>最后一次上传时间为：%s<br>文件Sha256: %s", data.LastUpload.Format("2006-01-02 15:04"), data.Hash),
				"group":  group,
				"stu":    stu,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": StatusCode.Success,
				"msg":    "您尚未上传！",
				"group":  group,
				"stu":    stu,
			})
		}
	})
	// get exp time
	r.GET("/api/report_exp_time_query", func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		expStr := c.Query("exp")
		if expStr == "" {
			QuickRes.BadRequest(c)
			return
		}
		exp, err := strconv.Atoi(expStr)
		if err != nil {
			QuickRes.BadRequest(c)
			return
		}
		if exp < 1 || exp > HwControl.Control.MaxExp {
			QuickRes.BadRequest(c)
			return
		}
		left := "尚未开放或已失效"
		if ExpAva[exp-1] && time.Now().After(Starts[exp-1]) && time.Now().Before(Ends[exp-1]) {
			dur := Ends[exp-1].Sub(time.Now())
			left = fmt.Sprintf("%d小时%d分钟", int(dur.Hours()), int(dur.Minutes())%60)
		} else if ExpAva[exp-1] && time.Now().After(ReStarts[exp-1]) && time.Now().Before(ReEnds[exp-1]) {
			dur := ReEnds[exp-1].Sub(time.Now())
			left = fmt.Sprintf("%d小时%d分钟<font style='color: red'>（补交剩余）</font>", int(dur.Hours()), int(dur.Minutes())%60)
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   StatusCode.Success,
			"start":    fmt.Sprintf("%v", Starts[exp-1]),
			"end":      fmt.Sprintf("%v", Ends[exp-1]),
			"re_start": fmt.Sprintf("%v", ReStarts[exp-1]),
			"re_end":   fmt.Sprintf("%v", ReEnds[exp-1]),
			"left":     left,
		})
	})
	// download via code
	r.GET("/api/report_download_stu", func(c *gin.Context) {
		exp, ava := expCheck(c)
		if !ava {
			QuickRes.BadRequest(c)
			return
		}
		stuNum := c.Query("stu")
		grpNum := c.Query("grp")
		if stuNum == "" && grpNum == "" {
			QuickRes.BadRequest(c)
			return
		}
		var rc FileRecord
		if grpNum == "" {
			inx, ok := StudentNumMap[stuNum]
			if !ok {
				QuickRes.BadRequest(c)
				return
			}
			rc = DynamicExpInfo[exp-1].Records[inx]
		} else {
			groupNum, err := strconv.Atoi(grpNum)
			if err != nil || groupNum > len(GroupStu) || groupNum < 1 {
				QuickRes.BadRequest(c)
				return
			}
			rc = DynamicExpInfo[exp-1].Records[groupNum-1]
		}
		TimeCode := genTimeCode(rc.LastUpload)
		code := c.Query("code")
		if code != TimeCode {
			QuickRes.BadRequest(c)
			return
		}
		fileName := ""
		if grpNum == "" {
			fileName = fmt.Sprintf("第%d次实验%s%s.doc", exp, rc.Stu.StudentNum, rc.Stu.StudentName)
		} else {
			fileName = fmt.Sprintf("第%s组大作业.zip", grpNum)
		}
		_path := path.Join(HwControl.Control.SaveDir, strconv.Itoa(exp), fileName)
		c.File(_path)
	})
}
