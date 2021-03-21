package common

import "os"

// Mkdir 创建文件
func Mkdir(filename string, mode os.FileMode) error {
	return os.Mkdir(filename, mode)
}

//IsDirExist 目录是否存在
func IsDirExist(filePath string) bool {
	result, erro := os.Stat(filePath)
	if erro != nil {
		return false
	}
	return result.IsDir()
}
