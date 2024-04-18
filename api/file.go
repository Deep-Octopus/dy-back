package api

import (
	resp "dy/middleware"
	"dy/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func Upload(c *gin.Context) {
	srcFile, head, err := c.Request.FormFile("file") // 表单中文件字段的名字
	if err != nil {
		c.JSON(http.StatusBadRequest, resp.ErrParam)
		return
	}

	suffix := ".png"
	temp := strings.Split(head.Filename, ".")
	if len(temp) > 1 {
		suffix = "." + temp[len(temp)-1]
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	url := utils.CONF.App.Static.BaseSrc + "images/" + fileName
	dstFile, err := os.Create(url)
	// 将文件保存到指定路径
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Err.WithMsg("文件上传失败"))
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, resp.Err.WithMsg("文件上传失败"))
	}
	c.JSON(http.StatusOK, resp.OK.WithData(url))

}

func UploadFile(c *gin.Context) {
	srcFile, head, err := c.Request.FormFile("file") // 表单中文件字段的名字
	if err != nil {
		c.JSON(http.StatusOK, resp.ErrParam.WithMsg("参数错误"))
		return
	}

	temp := strings.Split(head.Filename, ".")
	if len(temp) <= 1 {
		c.JSON(http.StatusOK, resp.ErrParam.WithMsg("文件信息有误"))
		return
	}
	suffix := "." + temp[len(temp)-1]
	path := temp[len(temp)-1]
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	url := utils.CONF.App.Static.BaseSrc + path + "/" + fileName

	// 确保文件夹存在，如果不存在则创建
	if err := os.MkdirAll(utils.CONF.App.Static.BaseSrc+path, os.ModePerm); err != nil {
		// 处理创建文件夹失败的错误
		c.JSON(http.StatusOK, resp.Err.WithMsg("无法创建文件夹:"+err.Error()))
		return
	}
	dstFile, err := os.Create(url)
	// 将文件保存到指定路径
	if err != nil {
		c.JSON(http.StatusOK, resp.Err.WithMsg("文件上传失败"))
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		c.JSON(http.StatusOK, resp.Err.WithMsg("文件上传失败"))
		return
	}
	c.JSON(http.StatusOK, resp.OK.WithData(url))
}

func DeleteFile(c *gin.Context) {
	filePath := c.Query("filePath")
	if filePath == "" {
		c.JSON(http.StatusOK, resp.ErrParam.WithMsg("错误的路径"))
		return
	}

	// Assuming your files are stored in the same directory structure as in the UploadFile function
	path := filePath

	// Check if the file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusOK, resp.ErrNotFound.WithMsg("文件不存在"))
		return
	}

	// Attempt to remove the file
	err := os.Remove(path)
	if err != nil {
		c.JSON(http.StatusOK, resp.Err.WithMsg("文件删除失败"))
		return
	}

	c.JSON(http.StatusOK, resp.OK)
}
