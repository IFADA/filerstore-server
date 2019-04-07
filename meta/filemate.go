package meta

import (
	mydb "filerstore-server/db"
	"sort"
)

type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(fMate FileMeta) {
	fileMetas[fMate.FileSha1] = fMate
}

func UpdateFileMetaDB(fMate FileMeta) bool {
	return mydb.OnFileUploadFinished(fMate.FileSha1, fMate.FileName, fMate.FileSize, fMate.Location)
}

//get the meta-info object from the sha1 value
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return fmeta, nil
}

//Gets a list of batch file meta infomation
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count] //返回个切片
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)

}
