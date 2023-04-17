package meta

import "distribute_store/util"

var hash map[string]FileMeta

// 文件元信息
type FileMeta struct {
	//计算文件的哈希值，这也是很多网盘实现秒传功能的原理，
	//将计算得到的哈希值与数据库进行比较，当两个文件的哈希值相同的时候，直接上传成功。
	FileHash string
	FileName string
	FileSize int64
	Location string
	//时间戳格式化之后的字符串
	UploadTime string
}

// 提供一个接口，通过Fileid获取对应的文件元信息
func GetFileMetaById(FileHash string) (FileMeta, error) {
	if v, ok := hash[FileHash]; ok {
		return v, nil
	}
	return hash[FileHash], util.MyError{"文件不存在"}
}

func init() {
	hash = make(map[string]FileMeta)
}

// 更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	hash[fmeta.FileHash] = fmeta
}

// 删除文件元信息
func DeleteFileMeta(fmeta FileMeta) {
	delete(hash, fmeta.FileHash)
}
