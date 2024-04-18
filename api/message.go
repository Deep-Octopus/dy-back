package api

import (
	"dy/config"
	resp "dy/middleware"
	"dy/models"
	"dy/service"
	"dy/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func GetP2PMessages(c *gin.Context) {
	contact := models.Contact{}
	if err := c.ShouldBind(&contact); err != nil {
		c.JSON(http.StatusBadRequest, resp.ErrParam)
		return
	}
	msgs := service.GetMessagesByContact(contact)
	c.JSON(http.StatusOK, resp.OK.WithData(msgs))
}

func GetGroupMessages(c *gin.Context) {
	contact := models.Contact{}
	if err := c.ShouldBind(&contact); err != nil {
		c.JSON(http.StatusBadRequest, resp.ErrParam)
		return
	}
	msgs := service.GetGroupMessagesByContact(contact)
	c.JSON(http.StatusOK, resp.OK.WithData(msgs))
}

func SendMsg(ctx *gin.Context) {
	ws, err := config.UpGrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 关闭websocket连接
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)

	MsgHandler(ws, ctx)
}

func MsgHandler(ws *websocket.Conn, ctx *gin.Context) {
	for {
		msg, err := utils.Subscribe(ctx, "Octopus")
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Now().Format(config.DateTimeFormat)
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		if err := ws.WriteMessage(1, []byte(m)); err != nil {
			fmt.Println(err)
		}
	}
}

func SendUserMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}
