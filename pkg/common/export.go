/**
@author:Administrator
@date:2021/1/14
@note:
**/
package common

import (
	"errors"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strconv"
	"time"
)

// 下载路径配置
func setFilePath() string {
	// excel 下载路径
	exportPath := "./public/export/" + time.Now().Format("2006-01-02") + "/"
	if !IsDirExist(exportPath) {
		err := os.Mkdir(exportPath, 0666)
		if err != nil {
			return ""
		}
	}
	return Strval(exportPath)
}

/**
@参数说明
@data []map[string]interface{} 下载数据 如： [map[addTime:1609322361 avatar:http://test-img.jiaoshipai.com/avatar.php?size=big&uid=59661 content:发表课程评价 evaluationId:16090 nickname:user81820 star:1 userid:59661]]
@titles []string 表头  如：[]string{"ID", "用户", "昵称", "内容", "评星", "添加时间", "头像地址"}
@fields []string 字段名 导出顺序 和titiles保持一致 如：[]string{"evaluationId", "userid", "nickname", "content", "star", "addTime", "avatar"}
@fileName string 文件名 可传"" 默认 fileName + 时间戳
@SheetTitle string 工作表sheet页名 可传 "" 默认sheet1
@return filePath  ./public/export/2021-01-15/评价数据导出1610699333.xlsx
*/

//数据导出excel
func ExportToExcel(data []map[string]interface{}, titles []string, fields []string, fileName string, SheetTitle string) (string, error) {

	if len(data) < 1 || len(titles) < 1 || len(fields) < 1 {
		log.Println("data or titles or fileName is empty")
		return "", errors.New("data or titles or fileName is empty")
	}
	if Empty(SheetTitle) {
		SheetTitle = "Sheet1"
	}
	if Empty(fileName) {
		fileName = ""
	}
	// 存储路径
	exportPath := setFilePath()
	// 生成一个新的文件
	file := xlsx.NewFile()
	// 添加sheet页
	sheet, err := file.AddSheet(SheetTitle)
	if err != nil {
		return "", err
	}
	var cell *xlsx.Cell
	// 插入表头
	titleRow := sheet.AddRow()
	for _, v := range titles {
		cell = titleRow.AddCell()
		cell.Value = v
		//表头字体颜色
		cell.GetStyle().Font.Color = "000000"
		//居中显示
		cell.GetStyle().Alignment.Horizontal = "center"
		cell.GetStyle().Alignment.Vertical = "center"
	}
	// 插入内容
	//fieldName := []string{"evaluationId", "userid", "nickname", "content", "star", "addTime", "avatar"}
	for _, v := range data {
		titleRow = sheet.AddRow()
		for _, v2 := range fields {
			// 插入内容
			cell = titleRow.AddCell()
			cell.Value = Strval(v[v2])
		}
	}

	fullPath := exportPath + fileName + strconv.Itoa(int(time.Now().Unix())) + ".xlsx"
	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}
	return Strval(fullPath), nil
}

/**
@filePath  ./public/export/2021-01-15/评价数据导出1610699333.xlsx
@return [[ID 用户 昵称 内容 评星 添加时间 头像地址] [16114 508181 user05136  10 1610609740 http://test-img.jiaoshipai.com/avatar.php?size=big&uid=508181] [16103 507344 user76696 哈哈 10 1610365991 http://test-img.jiaoshipai.com/avatar.php?size=big&uid=507344]]
*/

// 数据从excel读取
func Import(filePath string) ([][]string, error) {
	var resultData [][]string
	if Empty(filePath) {
		log.Println("filePath is empty")
		return resultData, errors.New("filePath is empty")
	}
	// 打开文件
	xlsFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Println("filePath open false")
		return resultData, errors.New("filePath open false")
	}
	// 遍历sheet页读取
	for _, sheet := range xlsFile.Sheets {
		//遍历行读取
		for _, row := range sheet.Rows {
			var data []string
			// 遍历每行的列读取
			for _, cell := range row.Cells {
				data = append(data, cell.String())
			}
			resultData = append(resultData, data)
		}
	}
	return resultData, nil
}
