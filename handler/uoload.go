package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"yun/meta"
	"yun/util"
)

func UploadHandler(w http.ResponseWriter,r *http.Request){
	if r.Method == "GET"{
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err!=nil{
			io.WriteString(w,"internel server error")
			return
		}
		io.WriteString(w,string(data))
	}else if r.Method =="POST"{
		file,head,err := r.FormFile("file")
		if err!=nil{
			fmt.Printf("Failed to get data,err:%s",err.Error())
			return
		}
		defer file.Close()
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "D:/test/+head.Filename",
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile,err := os.Create(fileMeta.Location)
		if err!=nil{
			fmt.Printf("Failed to get data,err:%s",err.Error())
			return
		}
		fileMeta.FileSize,err = io.Copy(newFile,file)
		if err!=nil{
			fmt.Printf("failed")
		}
		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Printf("%s",fileMeta.FileSha1)
		_ = meta.UpdateFileMetaDB(fileMeta)
		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)

	}
}

func UploadSucHandler(w http.ResponseWriter,r *http.Request){
	io.WriteString(w,"upload success！")
}

func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data,err := json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f,err := os.Open(fm.Location)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer f.Close()
	data,err := ioutil.ReadAll(f)
	if err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	// attachment表示文件将会提示下载到本地，而不是直接在浏览器中打开
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	// TODO: 更新文件表中的元信息记录

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler : 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")

	fMeta := meta.GetFileMeta(fileSha1)
	// 删除文件
	os.Remove(fMeta.Location)
	// 删除文件元信息
	meta.RemoveFileMeta(fileSha1)
	// TODO: 删除表文件信息
	w.WriteHeader(http.StatusOK)
}
