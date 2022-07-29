package main

import (
	"Daisy/Data"
	"Daisy/Fight"
	"fmt"
)

func main() {
	//重新载入excel表数据
	Data.RootPath = "../../Server/"
	if err := Data.LoadDataTables(Data.RootPath + "res/DataTables"); err != nil {
		panic(fmt.Sprintf("读取DataTables配置出错，%s", err.Error()))
	}

	if !Data.VerifyExcel() {
		fmt.Println("\nExcel数据有错误, 罚策划钱")
		return
	}

	//战斗相关excel检测 报错会立马panic
	Fight.LoadConfig(Data.RootPath)

	fmt.Println("\nExcel检查通过")
}
