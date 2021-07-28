package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

const LINE1 string = "**************"
const LINE2 string = "--------------"
const LINE3 string = "++++++++++++++"
const JSONFILE string = "db.json"

// 当前时间
var curTime string = time.Now().Format("2006/01/02")

type ItemsList struct {
	Content string `json: "content"`
	Time    string `json: "time"`
}

type Todo struct {
	Items []ItemsList `json: "items"`
}

func main() {
	var no int8
	// 显示当天代办事项
	fmt.Println(LINE3)
	fmt.Println("当天代办事项")
	queryDayTodo()
	fmt.Println(LINE3)
	for {
		// 菜单选项
		menu()
		fmt.Printf("请选择功能编号: ")
		fmt.Scanf("%d", &no)
		switch no {
		case 1:
			addTodo()
		case 2:
			queryTodo()
		case 3:
			deleteTodo()
		case 4:
			updateTodo()
		case 5:
			os.Exit(0)
		}
	}
}

// 检查错误
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 菜单
func menu() {
	fmt.Println(LINE1, "代办事项", LINE1)
	// 输出当前时间
	fmt.Println("当前时间:", time.Now().Format("2006/01/02"))
	fmt.Printf("%11s\n", "1.新增代办事项")
	fmt.Printf("%11s\n", "2.查看代办事项")
	fmt.Printf("%11s\n", "3.删除代办事项")
	fmt.Printf("%11s\n", "4.更改代办事项")
	fmt.Printf("%7s\n", "5.退出")
}

// 读取json文件
func readJson(mapTodo *map[string]interface{}) {
	data, err := ioutil.ReadFile(JSONFILE)
	checkErr(err)
	err = json.Unmarshal(data, &mapTodo)
	checkErr(err)
}

// 写入json文件
func writeJson(item *[]map[string]string) {
	// 写入json
	db := make(map[string]interface{})
	db["Items"] = *item
	data, err := json.MarshalIndent(db, "", "	")
	checkErr(err)
	err = ioutil.WriteFile(JSONFILE, data, 0755)
	checkErr(err)
}

// 格式化输出
func formatInput(str string) {
	fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, str, 0x1B)
}

// 新增代办事项
func addTodo() {
	var Item ItemsList
	// 新增代办事项内容
	fmt.Printf("新增代办事项内容: ")
	fmt.Scanf("%s\n", &Item.Content)
	// 请输入自定义时间, 校验时间格式
	fmt.Println("时间格式 2021/07/07")
	fmt.Printf("输入代办事项时间: ")
	fmt.Scanf("%s\n", &Item.Time)

	// 读取json文件
	data, err := ioutil.ReadFile(JSONFILE)
	checkErr(err)
	newTodo := &Todo{}
	err = json.Unmarshal(data, &newTodo)
	checkErr(err)

	// 创建slice, 汇总新增代办事项
	item := []map[string]string{}

	// struct 转 map 类型
	for _, data := range newTodo.Items {
		itemMap := make(map[string]interface{})
		obj1 := reflect.TypeOf(data)
		obj2 := reflect.ValueOf(data)
		for i := 0; i < obj1.NumField(); i++ {
			itemMap[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
		// 拼接成单条代办事项
		subItemMap := map[string]string{}
		for k, v := range itemMap {
			// interface 类型转换 string类型
			ov, _ := v.(string)
			subItemMap[k] = ov
		}
		// 汇总本地json中代办事项
		item = append(item, subItemMap)
	}

	// 新增代办事项, 合并
	item = append(item, map[string]string{"Content": Item.Content, "Time": Item.Time})

	// 写入json文件
	writeJson(&item)
	formatInput(">> 新增成功")
}

// 查看代办事项
func queryTodo() {
	// 读取json文件
	data, err := ioutil.ReadFile(JSONFILE)
	checkErr(err)
	mapTodo := make(map[string]interface{})
	err = json.Unmarshal(data, &mapTodo)
	checkErr(err)
	// 显示所有代办列表
	fmt.Printf("序号\t\t名称\t\t日期\n")
	for index, v := range mapTodo["Items"].([]interface{}) {
		v := v.(map[string]interface{})
		fmt.Printf("%d\t\t%s\t\t%s\n", index, v["Content"], v["Time"])
	}
	formatInput(">> 查看成功")
}

// 查询当天代办事项
func queryDayTodo() {
	// 读取json
	mapTodo := make(map[string]interface{})
	readJson(&mapTodo)
	// 更新代办列表
	for index, v := range mapTodo["Items"].([]interface{}) {
		v := v.(map[string]interface{})
		if curTime == v["Time"].(string) {
			fmt.Printf("%d\t\t%s\t\t%s\n", index, v["Content"].(string), v["Time"].(string))
		}
	}
}

// 删除待办事项
func deleteTodo() {
	var no int
	fmt.Printf("输入删除代办事项的序号: ")
	fmt.Scanf("%d\n", &no)

	// 读取json文件
	mapTodo := make(map[string]interface{})
	readJson(&mapTodo)

	// 追加代办列表
	item := []map[string]string{}
	for index, v := range mapTodo["Items"].([]interface{}) {
		v := v.(map[string]interface{})
		if no != index {
			// 代办事项, 合并
			item = append(item, map[string]string{"Content": v["Content"].(string), "Time": v["Time"].(string)})
		}
	}

	// 写入json文件
	writeJson(&item)
	formatInput(">> 删除成功")
}

// 更新代办事项
func updateTodo() {
	var no int
	fmt.Printf("输入更新代办事项的序号: ")
	fmt.Scanf("%d\n", &no)

	// 读取json文件
	mapTodo := make(map[string]interface{})
	readJson(&mapTodo)

	// 更新代办列表
	item := []map[string]string{}
	for index, v := range mapTodo["Items"].([]interface{}) {
		v := v.(map[string]interface{})
		// 代办事项, 合并
		if no == index {
			// 代办日期是否过期
			t1, err := time.Parse("2006/01/02", v["Time"].(string))
			t2, err := time.Parse("2006/01/02", curTime)
			if err == nil && t1.Before(t2) {
				fmt.Println(">> 代办日期已过期, 无法更改!")
				return
			}
			// 更新代办事项内容
			var content, tTime string
			fmt.Printf("更新代办事项内容: ")
			fmt.Scanf("%s\n", &content)
			// 请输入自定义时间, 校验时间格式
			fmt.Println("时间格式 2021/07/07")
			fmt.Printf("更新代办事项时间: ")
			fmt.Scanf("%s\n", &tTime)
			item = append(item, map[string]string{"Content": content, "Time": tTime})
		} else {
			item = append(item, map[string]string{"Content": v["Content"].(string), "Time": v["Time"].(string)})
		}
	}

	// 写入json
	writeJson(&item)
	formatInput(">> 更新成功")
}
