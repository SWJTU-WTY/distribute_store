package router

import (
	"crypto/sha256"
	"distribute_store/meta"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func ShowUpLoad(ctx *gin.Context) {
	//ctx.JSON(200, gin.H{})
	ctx.HTML(http.StatusOK, "index.html", gin.H{})
	return
}

func FileUpLoad(ctx *gin.Context) {
	//接受文件流到本地目录
	// 单文件
	file, _ := ctx.FormFile("file")
	log.Println(file.Filename)

	//计算文件的哈希值
	hash := sha256.New()
	src, _ := file.Open()
	io.Copy(hash, src)

	//上传之前先保存文件元信息
	fileMeta := meta.FileMeta{
		FileHash:   string(hex.EncodeToString(hash.Sum(nil))),
		FileName:   file.Filename,
		FileSize:   file.Size,
		Location:   "memory/" + file.Filename,
		UploadTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	meta.UpdateFileMeta(fileMeta)
	// 上传文件至指定的完整文件路径
	ctx.SaveUploadedFile(file, fileMeta.Location)
	//fmt.Println(fileMeta.FileHash)
	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!\n %s", file.Filename, fileMeta.FileHash))

	return
}

func FileQuery(ctx *gin.Context) {
	fileHash := ctx.Query("filehash")
	//从缓存读取对应文件数据
	fileMeta, err := meta.GetFileMetaById(fileHash)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	//将其json格式返回给客户端
	res, _ := json.Marshal(fileMeta)
	ctx.String(http.StatusOK, string(res))
}

func FileDownload(ctx *gin.Context) {
	fileHash := ctx.Query("filehash")
	fileMeta, err := meta.GetFileMetaById(fileHash)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}

	file, err := os.Open(fileMeta.Location)
	defer file.Close()
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	data, _ := io.ReadAll(file)

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Type", "application/octect-stream")
	ctx.Header("Content-Disposition", "attachment; filename=\""+fileMeta.FileName+"\"")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(data)))
	ctx.Writer.Write([]byte(data))

	//将文件返回给客户端
	ctx.File(fileMeta.FileName)
	res, _ := json.Marshal(fileMeta)
	ctx.String(http.StatusOK, fmt.Sprint("文件已下载\n", string(res)))

}

// 删除文件以及元信息
func FileDelete(ctx *gin.Context) {
	fileHash := ctx.Query("filehash")
	//具体的操作
	fileMeta, err := meta.GetFileMetaById(fileHash)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	//先删除缓存中的
	meta.DeleteFileMeta(fileMeta)
	//删除本地的
	os.Remove(fileMeta.Location)
	res, _ := json.Marshal(fileMeta)
	ctx.String(http.StatusOK, fmt.Sprint("文件已删除\n", string(res)))
	return
}

func FileUpdate(ctx *gin.Context) {
	fileHash := ctx.Query("filehash")
	//具体的操作
	op := ctx.Query("op")
	newFileName := ctx.Query("filename")
	fileMeta, err := meta.GetFileMetaById(fileHash)
	if err != nil {
		ctx.String(http.StatusNotFound, err.Error())
		return
	}
	if op == "0" {
		//重命名
		fileMeta.FileName = newFileName

		// TODO
		// 本地的名称还没有更新

		meta.UpdateFileMeta(fileMeta)
		res, _ := json.Marshal(fileMeta)
		ctx.String(http.StatusOK, fmt.Sprint("文件已更新名称\n", string(res)))
		return
	}

}

func InitRouter() {
	engine := gin.Default()
	engine.LoadHTMLGlob("static/view/*")
	engine.Static("/static", "./static")
	// 为 multipart forms 设置较低的内存限制 (默认是 32 MiB)
	engine.MaxMultipartMemory = 8 << 20 // 8 MiB

	//文件上传的界面与上传接口
	engine.POST("/file/upload", FileUpLoad)
	engine.GET("/file/upload", ShowUpLoad)

	//文件查询
	engine.GET("file/query", FileQuery)

	//文件下载
	engine.GET("file/download", FileDownload)

	//文件删除
	engine.POST("file/delete", FileDelete)

	//文件更新
	engine.POST("file/update", FileUpdate)

	engine.Run(":8080")
}
