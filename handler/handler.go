package handler

import (
	"encoding/json"
	"filerstore-server/meta"
	"filerstore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//static Fail to Creat file err:open /demo/1.jpg: no such file or directoryGoPath/src/loud/static/img  /home/unbuntu/GoPath/src/cloud/static/view/index.html
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//accept file and  local storage
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Fail to get data err%s\n", err.Error())
			return
		}
		defer file.Close()

		filemate := meta.FileMeta{
			FileName: head.Filename,
			Location: "/home/unbuntu/GolandProjects/tep/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-01 15:04:09"),
		}

		newFile, err := os.Create(filemate.Location)
		if err != nil {
			fmt.Printf("Fail to Creat file err:%s\n", err.Error())
			return
		}
		defer file.Close()

		filemate.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Fail to Copy data err%s\n", err.Error())
			return
		}
		newFile.Seek(0, 0)
		filemate.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(filemate)
		_ = meta.UpdateFileMetaDB(filemate)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)

	}

}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Finished")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	Fmate, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(Fmate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)

}

//query batch file meta info
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit")) // strconv.Atoi: conver string to int
	fileMetas := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandle(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)

}

// updatethe  metainfo
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

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//delete meta info
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)

}
