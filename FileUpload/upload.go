package FileUpload

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"main/AccessControl"
	"main/Logger"
	"main/ParamTools"
	"main/QuickRes"
	"main/Router"
	"main/Sql"
	"main/StatusCode"
	"main/Strings"
	"main/Vars"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func CalculateSHA5256(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	Hash := sha256.New()

	if _, err := io.Copy(Hash, src); err != nil {
		return "", err
	}

	hashInBytes := Hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashInBytes)

	return hashString, nil
}

func FileProcessSave(c *gin.Context, f *multipart.FileHeader, fun func(f *multipart.FileHeader) (string, error), path string) (string, error) {
	filename := ""
	if fun != nil {
		var err error
		filename, err = fun(f)
		if err != nil {
			return "", err
		}
		err = c.SaveUploadedFile(f, filepath.Join(path, filename))
		if err != nil {
			return "", err
		}
	}
	return filename, nil
}

func UploadCallback(r *gin.Engine) {
	r.POST(Router.FatherPath+Router.FileUpload, func(c *gin.Context) {
		acc, _, _, _, uid, group := AccessControl.AccControl(c)
		if !acc {
			return
		}
		f, err := c.FormFile("file")
		if err != nil {
			QuickRes.InternalError(c)
			return
		}
		_type, _ := ParamTools.GetParam(c, "type")
		if group == AccessControl.Anonymous && _type != "image" && f.Size > 5*1024*1024 {
			QuickRes.NotPermitted(c)
			return
		}
		if _type == "" {
			_type = "NULL"
		}
		hashName, err := FileProcessSave(c, f, func(f *multipart.FileHeader) (string, error) {
			name, err := CalculateSHA5256(f)
			if err != nil {
				return "", err
			}
			rows, err := Sql.Db.Query(fmt.Sprintf("SELECT * FROM file_upload('%d', '%s', '%s', '%s')", uid, name, Strings.FmtQuery(f.Filename), Strings.FmtQuery(_type)))
			defer func(rows *sql.Rows) {
				err := rows.Close()
				if err != nil {
					return
				}
			}(rows)
			if err != nil {
				return "", err
			}
			rows.Next()
			var st int
			var msg string
			err = rows.Scan(&st, &msg)
			if err != nil {
				return "", err
			}
			if st == -1 {
				return "", err
			}
			return name, nil
		}, Vars.UploadSavePath)
		if err != nil {
			Logger.LogE(err.Error())
			QuickRes.InternalError(c)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": StatusCode.Success,
			"data":   hashName,
		})
	})
}
