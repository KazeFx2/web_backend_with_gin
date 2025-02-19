package HomeworkUpload

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"main/Logger"
	"main/Msg"
	"main/QuickRes"
)

var WsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HwWsCallback(r *gin.Engine) {
	r.GET("/ws", func(c *gin.Context) {
		if !authCheck(c) {
			QuickRes.NotPermitted(c)
			return
		}
		conn, err := WsUpgrade.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			Logger.LogE("websocket WsUpgrade error: %v", err)
			return
		}
		defer func(conn *websocket.Conn) {
			err := conn.Close()
			if err != nil {
				return
			}
		}(conn)

		for {
			flag := false
			for job := range Msg.MessageChan {
				str := job.Stu.StudentNum + "&&" + job.Stu.StudentName + "&&" + job.Msg
				if err := conn.WriteMessage(websocket.TextMessage, []byte(str)); err != nil {
					Logger.LogE("write error: %v", err)
					Msg.MessageChan <- job
					flag = true
					break
				}
			}
			if flag {
				break
			}
			type_, _, err := conn.ReadMessage()
			if err != nil {
				Logger.LogE("read error: %v", err)
				break
			} else if type_ == websocket.PingMessage {
				Logger.LogI("received ping")
				err = conn.WriteMessage(websocket.PongMessage, []byte{})
				if err != nil {
					Logger.LogE("write 'pong' error: %v", err)
					break
				}
			}
		}
	})
}
