package Account

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/AccessControl"
	"main/Logger"
	"main/QuickRes"
	"main/Router"
	"main/Sql"
	"main/StatusCode"
	"main/Strings"
	"net/http"
	"regexp"
)

func Register(username string, avatar string, email string, passwd string) (int, int, string) {
	if len(username) > 50 {
		return StatusCode.Error, -1, "用户名长于50字符"
	}
	forbiddenCharsPattern := `[!@#$%^&*(),.?":{}|<>]`
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	startWithNumPattern := `^[0-9]`
	forbiddenRe := regexp.MustCompile(forbiddenCharsPattern)
	emailRe := regexp.MustCompile(emailPattern)
	startWithNumRe := regexp.MustCompile(startWithNumPattern)
	avatar, avatarAva := AccessControl.CheckSha256(avatar)
	passwd, passwdAva := AccessControl.CheckSha256(passwd)
	if forbiddenRe.MatchString(username) {
		return StatusCode.Error, -1, fmt.Sprintf("用户名不可包含 '%s'", forbiddenCharsPattern)
	}
	if startWithNumRe.MatchString(username) {
		return StatusCode.Error, -1, "用户名不可以数字开头"
	}
	if !emailRe.MatchString(email) {
		return StatusCode.Error, -1, "非法的邮箱地址"
	}
	if !avatarAva {
		return StatusCode.Error, -1, "无效的头像地址"
	}
	if !passwdAva {
		return StatusCode.Error, -1, "无效的密码"
	}
	query := fmt.Sprintf("SELECT * FROM register_user('%s', '%s', '%s', '%s')", Strings.FmtQuery(username), Strings.FmtQuery(avatar), Strings.FmtQuery(email), Strings.FmtQuery(passwd))
	row, err := Sql.Db.Query(query)
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			return
		}
	}(row)
	if err != nil {
		Logger.Log(Logger.LogError, err.Error())
		return StatusCode.Error, -1, "内部错误"
	}
	var _uid int
	var _msg string
	row.Next()
	err = row.Scan(&_uid, &_msg)
	if err != nil {
		Logger.Log(Logger.LogError, err.Error())
		return StatusCode.Error, -1, "内部错误"
	}
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			return
		}
	}(row)
	if _uid == -1 {
		return StatusCode.Error, -1, _msg
	}
	return StatusCode.Success, _uid, _msg
}

func RegisterCallback(r *gin.Engine) {
	r.POST(Router.FatherPath+Router.Register, func(c *gin.Context) {
		acc, _, _, _, uid, _ := AccessControl.AccControl(c)
		if !acc {
			return
		}
		var exist bool
		var username string
		var avatar string
		var email string
		var passwd string
		if username, exist = c.GetPostForm("username"); !exist {
			QuickRes.BadRequest(c)
			return
		}
		if avatar, exist = c.GetPostForm("avatar"); !exist {
			QuickRes.BadRequest(c)
			return
		}
		if email, exist = c.GetPostForm("email"); !exist {
			QuickRes.BadRequest(c)
			return
		}
		if passwd, exist = c.GetPostForm("password"); !exist {
			QuickRes.BadRequest(c)
			return
		}
		st, uid, msg := Register(username, avatar, email, passwd)
		if st == StatusCode.Success {
			c.JSON(http.StatusOK, gin.H{
				"status": StatusCode.Success,
				"data":   uid,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  StatusCode.Error,
				"message": msg,
			})
		}
	})
}
