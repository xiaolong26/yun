package main

import (
	"fmt"
	"net/http"
	"yun/handler"
)

func main(){
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	err := http.ListenAndServe(":8080",nil)
	if err!=nil{
		fmt.Printf("连接失败,%s",err.Error())
	}
}
