package HomeworkUpload

import (
	"archive/zip"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/QuickRes"
	"main/StatusCode"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

const (
	tmpFile = "tmp.zip"
)

func authCheck(c *gin.Context) bool {
	auth := c.Query("auth")
	if auth != HwControl.Control.Auth {
		return false
	}
	return true
}

func expCheck(c *gin.Context) (int, bool) {
	expStr := c.Query("exp")
	if expStr == "" {
		return -1, false
	}
	exp, err := strconv.Atoi(expStr)
	if err != nil {
		return -1, false
	}
	if exp < 1 || exp > HwControl.Control.MaxExp {
		return -1, false
	}
	return exp, true
}

func fileListCheck(c *gin.Context) (bool, string) {
	if !authCheck(c) {
		return false, ""
	}
	exp := c.Query("exp")
	if exp == "" {
		return false, ""
	}
	_, err := strconv.Atoi(exp)
	if err != nil {
		return false, ""
	}
	return true, exp
}

func downloadCheck(c *gin.Context) (bool, string) {
	if !authCheck(c) {
		return false, ""
	}
	exp := c.Query("exp")
	if exp == "" {
		return false, ""
	}
	fileName := c.Query("name")
	if fileName == "" {
		return false, ""
	}
	return true, path.Join("./", exp, fileName)
}

func ZipAll(srcDir string, zipFileName string, ext bool, late bool, reStart time.Time, reEnd time.Time) error {
	err := os.RemoveAll(zipFileName)
	if err != nil {
		return err
	}

	zipFile, _ := os.Create(zipFileName)
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			return
		}
	}(zipFile)

	archive := zip.NewWriter(zipFile)
	defer func(archive *zip.Writer) {
		err := archive.Close()
		if err != nil {
			return
		}
	}(archive)

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {

		if path == srcDir {
			return nil
		}

		header, _ := zip.FileInfoHeader(info)
		_, name := filepath.Split(path)
		header.Name = name

		if info.IsDir() {
			header.Name += `/`
		} else {
			header.Method = zip.Deflate
			if ext {
				inLateTime := info.ModTime().After(reStart) && info.ModTime().Before(reEnd)
				if late && !inLateTime {
					return nil
				}
				if !late && inLateTime {
					return nil
				}
			}
		}

		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					return
				}
			}(file)
			_, err := io.Copy(writer, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func FileDownloadCallback(r *gin.Engine) {
	// file_download
	r.GET("/api/report_download", func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		ava, filePath := downloadCheck(c)
		if !ava {
			QuickRes.BadRequest(c)
			return
		}
		filePath = path.Join(HwControl.Control.SaveDir, filePath)
		stat, err := os.Stat(filePath)
		if err != nil || stat.IsDir() {
			QuickRes.BadRequest(c)
			return
		}
		c.File(filePath)
	})
	// Download all
	r.GET("/api/report_download_all", func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		ava := authCheck(c)
		if !ava {
			QuickRes.NotPermitted(c)
			return
		}
		dir := c.Query("exp")
		ext := c.Query("ext") == "true"
		late := c.Query("late") == "true"
		val, err := strconv.Atoi(dir)
		if err != nil || val < 1 || val > 8 {
			QuickRes.BadRequest(c)
			return
		}
		err = ZipAll(path.Join(HwControl.Control.SaveDir, dir), tmpFile, ext, late, ReStarts[val-1], ReEnds[val-1])
		if err != nil {
			QuickRes.InternalError(c)
			return
		}
		c.File(tmpFile)
	})
	// file list request
	r.GET("/api/report_list", func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		ava, dir := fileListCheck(c)
		if !ava {
			QuickRes.BadRequest(c)
			return
		}
		targetDir := path.Join(HwControl.Control.SaveDir, dir)
		dirEt, err := os.ReadDir(targetDir)
		if err != nil {
			QuickRes.BadRequest(c)
			return
		}
		fileNames := make([]string, 0)
		// eg. "第1次实验Y02114xxx张三.doc"
		regNameExp := `^第[0-9]*次实验[A-Z][0-9]{8}.*\.(doc|docx)$`
		regMatch := regexp.MustCompile(regNameExp)
		for i := 0; i < len(dirEt); i++ {
			f := dirEt[i]
			if regMatch.MatchString(f.Name()) && !f.IsDir() {
				fileNames = append(fileNames, f.Name())
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": StatusCode.Success,
			"data":   fileNames,
		})
	})
	// au validate
	r.GET("/api/report_validate", func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		ava := authCheck(c)
		if ava {
			QuickRes.ProcessOK(c)
		} else {
			QuickRes.NotPermitted(c)
		}
	})
	// get not upload list
	r.GET("/api/report_not_upload", func(c *gin.Context) {
		if !authCheck(c) {
			QuickRes.NotPermitted(c)
			return
		}
		exp := -1
		ava := false
		if exp, ava = expCheck(c); !ava {
			QuickRes.BadRequest(c)
			return
		}
		name := false
		withName := c.Query("name")
		if withName == "true" {
			name = true
		}
		ret := make([]string, 0)
		for inx, rc := range DynamicExpInfo[exp-1].Records {
			if !rc.Uploaded {
				str := ""
				if IsGroupExp(exp) {
					str = fmt.Sprintf("第%d组", inx+1)
					if name {
						for _, tmp := range GroupStu[inx+1] {
							str += " " + Students[tmp].StudentName
						}
					}
				} else {
					str = rc.Stu.StudentNum
					if name {
						str += " " + rc.Stu.StudentName
					}
				}
				ret = append(ret, str)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": StatusCode.Success,
			"data": gin.H{
				"students": ret,
				"total":    len(ret),
			},
		})
	})
}
