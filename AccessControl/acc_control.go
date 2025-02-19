package AccessControl

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"main/QuickRes"
	"main/Sql"
	"main/StatusCode"
	"main/Strings"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	Admin = iota
	User
	Anonymous
)

func CheckSha256(str string) (string, bool) {
	lower := strings.ToLower(str)
	matchPattern := `^[a-f0-9]{64}$`
	re := regexp.MustCompile(matchPattern)
	return lower, re.MatchString(lower)
}

func GetUserGroup(c *gin.Context) (int, int) {
	cookie, err := c.Request.Cookie("identity")
	if err != nil {
		return 2, Anonymous
	}
	au, cookieAva := CheckSha256(cookie.Value)
	if !cookieAva {
		return 2, Anonymous
	}
	c.Set("auth", au)
	row, err := Sql.Db.Query(fmt.Sprintf("SELECT * FROM check_authorization('%s')", Strings.FmtQuery(au)))
	defer func(row *sql.Rows) {
		err := row.Close()
		if err != nil {
			return
		}
	}(row)
	if err != nil {
		return 2, Anonymous
	}
	row.Next()
	var uid int
	var groupName string
	err = row.Scan(&uid, &groupName)
	if err != nil {
		return 2, Anonymous
	}
	switch groupName {
	case "admin":
		return uid, Admin
	case "user":
		return uid, User
	case "anonymous":
		return 2, Anonymous
	default:
		return 2, Anonymous
	}
}

func GetIpPort(c *gin.Context) (string, int) {
	ip := c.ClientIP()
	RemoteAddr := c.Request.RemoteAddr
	_, port, err := net.SplitHostPort(RemoteAddr)
	if err != nil {
		return "0.0.0.0", 0
	}
	numPort, err := strconv.Atoi(port)
	if err != nil {
		numPort = 0
	}
	return ip, numPort
}

func GetUrl(c *gin.Context) string {
	url := c.Request.URL.Path
	return url
}

func GetControl(ip string, port int, url string, uid int) bool {
	rows, err := Sql.Db.Query(fmt.Sprintf("SELECT * FROM acc_control('%s', '%s', '%d')", Strings.FmtQuery(ip), Strings.FmtQuery(url), uid))
	if err != nil {
		return false
	}
	rows.Next()
	var acc bool
	err = rows.Scan(&acc)
	if err != nil {
		return false
	}
	err = rows.Close()
	if err != nil {
		return false
	}
	if !acc {
		return false
	}
	return true
}

func AccControl(c *gin.Context) (acc bool, ip string, port int, url string, uid int, group int) {
	QuickRes.SetOrigin(c)
	uid, group = GetUserGroup(c)
	ip, port = GetIpPort(c)
	url = GetUrl(c)
	acc = GetControl(ip, port, url, uid)
	c.Set("uid", uid)
	c.Set("ip", ip)
	c.Set("port", port)
	if !acc {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  StatusCode.NotPermitted,
			"message": ">_<",
		})
	}
	return acc, ip, port, url, uid, group
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type ResponseBody struct {
	Status int `json:"status"`
}

func ReturnLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		statusCode := c.Writer.Status()
		domain := c.Request.Host
		url := c.Request.RequestURI
		var uid, port, auth, ip string
		va, _ := c.Get("uid")
		if v, ok := va.(int); ok {
			uid = fmt.Sprintf("'%d'", v)
		} else {
			uid = "NULL"
		}
		va, _ = c.Get("port")
		if v, ok := va.(int); ok {
			port = fmt.Sprintf("'%d'", v)
		} else {
			port = "'0'"
		}
		va, _ = c.Get("auth")
		if v, ok := va.(string); ok {
			auth = fmt.Sprintf("'%s'", v)
		} else {
			auth = "NULL"
		}
		va, _ = c.Get("ip")
		if v, ok := va.(string); ok {
			ip = fmt.Sprintf("'%s'", v)
		} else {
			ip = "'::1'"
		}
		var responseBody ResponseBody
		responseBody.Status = -1
		stCode := "NULL"
		if err := json.NewDecoder(blw.body).Decode(&responseBody); err == nil && responseBody.Status != -1 {
			stCode = fmt.Sprintf("'%d'", responseBody.Status)
		}
		query := fmt.Sprintf("SELECT * FROM add_inet_action(%s, %s, %s, %s, '%s', '%s', '%d', %s)", uid, ip, port, auth, Strings.FmtQuery(Strings.FmtLength(domain, 50)), Strings.FmtQuery(Strings.FmtLength(url, 128)), statusCode, stCode)
		rows, err := Sql.Db.Query(query)
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				return
			}
		}(rows)
		if err != nil {
			return
		}
	}
}

func AccMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		QuickRes.SetOrigin(c)
		c.Next()
	}
}
