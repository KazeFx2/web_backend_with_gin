package Login

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
)

func PasswdLoginCallback(r *gin.Engine) {
	r.POST(Router.FatherPath+Router.Login, func(c *gin.Context) {
		acc, ip, port, _, _, _ := AccessControl.AccControl(c)
		if !acc {
			return
		}
		user, exist := c.GetPostForm("user")
		if !exist {
			QuickRes.BadRequest(c)
			return
		}
		passwd, exist := c.GetPostForm("password")
		if !exist {
			QuickRes.BadRequest(c)
			return
		}
		rows, err := Sql.Db.Query(fmt.Sprintf("SELECT * FROM login_account('%s', '%s', '%s', '%d')", Strings.FmtQuery(user), Strings.FmtQuery(passwd), Strings.FmtQuery(ip), port))
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				return
			}
		}(rows)
		if err != nil {
			Logger.LogE(err.Error())
			QuickRes.InternalError(c)
			return
		}
		rows.Next()
		var st int
		var au string
		var msg string
		err = rows.Scan(&st, &msg, &au)
		if err != nil {
			QuickRes.InternalError(c)
			return
		}
		if st != 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  StatusCode.Error,
				"message": msg,
				"data":    au,
			})
		} else {
			c.SetCookie("identity", au, 3600*24*7, "/", c.Request.Host, false, false)
			c.JSON(http.StatusOK, gin.H{
				"status": StatusCode.Success,
				"data":   au,
			})
		}
	})
}
