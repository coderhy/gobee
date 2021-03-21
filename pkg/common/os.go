package common

import (
	"errors"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// Chmod 修改文件、目录的权限
func Chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// Chown 将指定文件的拥有者改为指定的用户或组
func Chown(filename string, uid, gid int) bool {
	return os.Chown(filename, uid, gid) == nil
}

// FuncName 获取正在运行的函数名
// Callers	0为本身函数(runtime.Callers ) 1 为函数名，2 为调用者名(上一级调用方的函数名称)，3 未识别的调用，设置成1
// 其实0 就是函数本身 FuncName，为一个函数调用栈
// 一法通，万法通，诸般深奥的学问到了极处，本是殊途同归
func FuncName() map[string]interface{} {

	// pc := make([]uintptr, 1)
	// runtime.Callers(2, pc)
	// return runtime.FuncForPC(pc[0]).Name()

	// funcName, file, line, ok := runtime.Caller(1)
	// fmt.Println("pc:", funcName, file, line, ok, runtime.FuncForPC(funcName).Name())
	pc, filePath, line, ok := runtime.Caller(1)
	//目录结构及函数名
	dirFuncName := runtime.FuncForPC(pc).Name()

	//真实文件名
	fileName := ""
	fileArgs := Explode("/", filePath)
	if len(fileArgs) > 0 {
		lastSerialNumber := len(fileArgs) - 1
		fileName = fileArgs[lastSerialNumber]

		if fileName != "" {
			fileNameArgs := Explode(".", fileName)
			if len(fileNameArgs) == 2 {
				fileName = fileNameArgs[0]
			}
		}
	}

	//函数名
	funcName := ""
	funcArgs := Explode(".", dirFuncName)
	if len(funcArgs) > 0 {
		lastSerialNumber := len(funcArgs) - 1
		funcName = funcArgs[lastSerialNumber]
	}
	result := map[string]interface{}{
		"filePath":    filePath,
		"fileName":    fileName,
		"line":        line,
		"dirFuncName": dirFuncName, //目录结构方法名
		"funcName":    funcName,
		"viewPath":    strings.ToLower(fileName) + "/" + Lcfirst(funcName), //模板路径
		"ok":          ok,
	}
	return result
}

//Display 视图模板路径(.html后缀)
func Display() string {

	// pc, filePath, line, ok := runtime.Caller(1)
	pc, filePath, _, _ := runtime.Caller(1)
	//目录结构及函数名
	dirFuncName := runtime.FuncForPC(pc).Name()

	//真实文件名
	fileName := ""
	fileArgs := Explode("/", filePath)
	if len(fileArgs) > 0 {
		lastSerialNumber := len(fileArgs) - 1
		fileName = fileArgs[lastSerialNumber]

		if fileName != "" {
			fileNameArgs := Explode(".", fileName)
			if len(fileNameArgs) == 2 {
				fileName = fileNameArgs[0]
			}
		}
	}

	//函数名
	funcName := ""
	funcArgs := Explode(".", dirFuncName)
	if len(funcArgs) > 0 {
		lastSerialNumber := len(funcArgs) - 1
		funcName = funcArgs[lastSerialNumber]
	}
	return strings.ToLower(fileName) + "/" + Lcfirst(funcName) + ".html" //模板路径
}

// GetFuncName 获取正在运行的函数名
func GetFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// FuncCall 反射调用
func FuncCall(obj interface{}, method string, args ...interface{}) ([]reflect.Value, error) {

	//是否存在验证
	fn := reflect.ValueOf(obj).MethodByName(method)
	if !fn.IsValid() { //fn.String() == "<invalid Value>"
		err := errors.New("method does not exist ")
		return nil, err
	}

	//参数处理
	if len(args) != fn.Type().NumIn() {
		//传入参数与需要的参数数量不一致
		err := errors.New("wrong number of parameters")
		return nil, err
	}

	params := make([]reflect.Value, len(args))
	for i := range args {
		// 判断应位置的参数是否与需要的类型相同
		if reflect.TypeOf(args[i]) != fn.Type().In(i) {
			err := errors.New("wrong parameter type")
			return nil, err
		}
		// 加入切片
		params[i] = reflect.ValueOf(args[i])
	}

	return fn.Call(params), nil
}
