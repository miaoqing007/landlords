package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/tealeg/xlsx"
)

var (
	tableList = flag.String("tableList", "./conf/excel/server.xml", "配表列表")

	conf = flag.String("conf", "./conf/excel", "game config path")

	pkgName = flag.String("pkgName", "configtable", "包名")

	genFile = flag.String("genFile", "./src/configtable", "config table gen go code file")

	mgrFile = flag.String("mgrFile", "./src/configtable/ConfigtableMgr.go", "config table gen mgr go code file")

	typesMap = map[string]string{
		"bool":        "Bool",
		"float32":     "Float",
		"float":       "Float",
		"int32":       "Int",
		"uint32":      "Int",
		"int64":       "Int64",
		"uint64":      "Int64",
		"string":      "String",
		"[]string":    "String",
		"[][]string":  "String",
		"[]int32":     "String",
		"[][]int32":   "String",
		"[]uint32":    "String",
		"[][]uint32":  "String",
		"[]float32":   "String",
		"[][]float32": "String",
		"time":        "String",
	}

	validateTypes = []string{
		"bool",
		"float32",
		"int32",
		"uint32",
		"uint64",
		"string",
		"[]string",
		"[][]string",
		"[]int32",
		"[][]int32",
		"[]uint32",
		"[][]uint32",
		"[]float32",
		"[][]float32",
		"time",
	}

	fieldsVarName = make([]string, 0, 128)

	fieldsName = make([]string, 0, 128)

	fieldsType = make([]string, 0, 128)

	testMap = make(map[uint32]string, 10000)
)

type configTable struct {
	XMLName   xml.Name `xml:"table"`
	TableName string   `xml:"tableName,attr"`
}

type configTables struct {
	XMLName xml.Name      `xml:"root"`
	Tables  []configTable `xml:"table"`
}

func upperfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func lowerfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 检查一下字段变量类型定义是否正确
func checkFieldType(fieldType string) bool {
	for _, fType := range validateTypes {
		if fieldType == fType {
			return true
		}
	}
	return false
}

func getXlsxType(t string) string {
	fType, ok := typesMap[t]
	if ok {
		return fType
	}
	return "Bool"
}

func checkFieldsNumber() bool {
	if len(fieldsName) == len(fieldsType) && len(fieldsName) == len(fieldsVarName) {
		return true
	}
	return false
}

// isParseHeader:是否是解析头(前三行表示该表的定义)
func parseExcelFile(fileName string) {
	excelFileName := filepath.FromSlash(fileName)
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("open excel file %v error %v", excelFileName, err)
	}
	if len(xlFile.Sheets) > 0 {
		parseOneSheet(xlFile.Sheets[0])
		autoGenLoadEachTableFuncCode(excelFileName, xlFile.Sheets[0].Name)
	}
}

func parseOneSheet(sheet *xlsx.Sheet) {
	createGolangCodeFile(sheet.Name)
	for i, row := range sheet.Rows {
		// 前三行分别为字段名、变量名、数据类型定义
		if i == 0 {
			parseFieldsName(row)
		} else if i == 1 {
			parseFieldsType(row)
		} else if i == 2 {
			parseFieldsVariableName(row)
			fvNameLen, fNameLen, fTypeLen := len(fieldsVarName), len(fieldsName), len(fieldsType)
			if fvNameLen != fNameLen || fvNameLen != fTypeLen {
				panic(fmt.Errorf("%s表 字段长度错误 %v %v %v ", sheet.Name, fvNameLen, fNameLen, fTypeLen))
			}
			if strings.HasSuffix(sheet.Name, "Table") {
				autoGenGolangCode(sheet.Name)
				autoGenLoadXlsxDataCode(sheet.Name)
			}
		} else { // 数据行
			cellsLen := len(row.Cells)
			if cellsLen > 0 && cellsLen < len(fieldsVarName) {
				if len(row.Cells[0].String()) == 0 {
					continue
				}
				panic(fmt.Errorf("%s表 数据配置错误,第%v行可能有空白配置 请检查配表", sheet.Name, i+1))
			}

			for col, cell := range row.Cells {
				if col >= len(fieldsType) || col >= len(fieldsName) {
					continue
				}
				if col == 0 && len(cell.String()) == 0 {
					continue
				}
				ftype := fieldsType[col]
				if ftype == "uint32" || ftype == "uint64" {
					value, _ := cell.Int()
					if uint32(value) > math.MaxInt32 {
						errStr := fmt.Sprintf("数值配置错误 表:%v 第%v行 字段:%v 配置值:%v 解析值:%v \n", sheet.Name, i+1, fieldsName[col], cell.String(), uint32(value))
						panic(errStr)
					}
				} else if ftype == "time" {
					value, err := time.ParseInLocation("2006/1/2 15:04:05", cell.String(), time.Local)
					if err != nil {
						errStr := fmt.Sprintf("数值配置错误 表:%v 第%v行 字段:%v 配置值:%v 解析值:%v 正确格式:2006/01/02 15:04:05 \n", sheet.Name, i+1, fieldsName[col], cell.String(), value)
						panic(errStr)
					}
				}
				if sheet.Name == "SeverListTable" || sheet.Name == "UnitTable" ||
					sheet.Name == "NameTable" {
					continue
				}
				/*if col == 0 {
					value, _ := cell.Int()
					if table, ok := testMap[uint32(value)]; ok {
						fmt.Println("配表载入错误  表名1:"+sheet.Name+" 表名2:"+table+" id重复:", value)
					} else {
						testMap[uint32(value)] = sheet.Name
					}
				}*/
			}
		}
	}
	fieldsVarName = fieldsVarName[:0]
	fieldsName = fieldsName[:0]
	fieldsType = fieldsType[:0]
}

// 解析字段名称
func parseFieldsName(row *xlsx.Row) {
	for _, cell := range row.Cells {
		fieldName := cell.String()
		if len(fieldName) == 0 {
			continue
		}
		fieldsName = append(fieldsName, fieldName)
	}
}

// 解析变量名定义
func parseFieldsVariableName(row *xlsx.Row) {
	for _, cell := range row.Cells {
		fieldVarName := cell.String()
		if fieldVarName == "" {
			continue
		}
		fieldVarName = upperfirst(fieldVarName)
		fieldsVarName = append(fieldsVarName, fieldVarName)
	}
}

func typeTotype(str string) string {
	//if str == "time" {
	//	return "int64"
	//}
	if str == "time" {
		return "time.Time"
	}
	return str
}

// 解析字段类型
func parseFieldsType(row *xlsx.Row) {
	for _, cell := range row.Cells {
		fieldType := cell.String()
		if fieldType == "" {
			continue
		}
		switch fieldType {
		case "float":
			fieldType += "32"
		case "string[]":
			fieldType = "[]string"
		case "string[][]":
			fieldType = "[][]string"
		case "int32[]":
			fieldType = "[]int32"
		case "int32[][]":
			fieldType = "[][]int32"
		case "uint32[]":
			fieldType = "[]uint32"
		case "uint32[][]":
			fieldType = "[][]uint32"
		case "float[]":
			fieldType = "[]float32"
		case "float[][]":
			fieldType = "[][]float32"
		}
		fieldsType = append(fieldsType, fieldType)
	}
}

func createGolangCodeFile(sheetName string) {
	defFileName := filepath.FromSlash(*genFile + "/" + sheetName + ".go")
	defFile, err := os.OpenFile(defFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	mgrName := sheetName + "Mgr"
	mgrInstance := lowerfirst(mgrName) + "Ins"
	w := bufio.NewWriter(defFile)
	pkgName := fmt.Sprintf("package %v\n\n", *pkgName)
	w.WriteString(pkgName)
	w.WriteString("import (\n")
	w.WriteString(`    "sync"`)
	w.WriteString("\n")
	w.WriteString(`    "fmt"`)
	w.WriteString("\n")
	w.WriteString(`    "path/filepath"`)
	w.WriteString("\n")
	w.WriteString(`    "strings"`)
	//w.WriteString("\n")
	//w.WriteString(`    "component/function"`)
	w.WriteString("\n")
	w.WriteString(`    "time"`)
	w.WriteString("\n")
	w.WriteString(`    "github.com/tealeg/xlsx"`)
	w.WriteString("\n")
	w.WriteString(")\n\n")
	w.WriteString("var (\n")
	w.WriteString("    " + mgrInstance + " *" + mgrName + "\n")
	mgrOnceVar := lowerfirst(mgrName) + "Once"
	w.WriteString("    " + mgrOnceVar + " sync.Once\n")
	w.WriteString(")\n\n")
	w.WriteString("type " + mgrName + " struct {\n")
	w.WriteString("    Datas *sync.Map\n")
	w.WriteString("}\n\n")
	w.WriteString("func " + mgrName + "_GetMe() " + "*" + mgrName + " {\n")
	w.WriteString("    " + mgrOnceVar + ".Do( func () {\n")
	w.WriteString("        " + mgrInstance + " = &" + mgrName + " {\n")
	w.WriteString("            Datas: &sync.Map{},\n")
	w.WriteString("        }\n")
	w.WriteString("    })\n")
	w.WriteString("    return " + mgrInstance + "\n")
	w.WriteString("}\n\n")

	w.WriteString("func " + mgrName + "_GetSize() int " + "{\n")
	w.WriteString("    var size int\n")
	w.WriteString("    " + mgrInstance + ".Datas.Range(func(key, value interface{}) bool {\n")
	w.WriteString("        size++\n")
	w.WriteString("        return true\n")
	w.WriteString("    })\n")
	w.WriteString("    return size\n")
	w.WriteString("}\n\n")
	w.Flush()

}

func autoGenGolangCode(sheetName string) {
	defFileName := filepath.FromSlash(*genFile + "/" + sheetName + ".go")
	defFile, err := os.OpenFile(defFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(defFile)

	w.WriteString("type " + sheetName + " struct {\n")

	index := 1
	for i, fieldVarName := range fieldsVarName {
		if !checkFieldType(fieldsType[i]) {
			panic(fmt.Errorf("%s表 字段类型配置错误 字段名:%s 字段类型:%s", sheetName, fieldsName[i], fieldsType[i]))
		}
		w.WriteString("    " + fieldVarName + " " + typeTotype(fieldsType[i]) + " // " + fieldsName[i] + "\n")
		index = index + 1
	}

	w.WriteString("}\n\n")

	w.Flush()
}

func autoGenLoadXlsxDataCode(sheetName string) {

	defFileName := filepath.FromSlash(*genFile + "/" + sheetName + ".go")
	defFile, err := os.OpenFile(defFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(defFile)
	mgrName := sheetName + "Mgr"
	mgrInstance := lowerfirst(mgrName) + "Ins"

	w.WriteString("func (mgr *" + mgrName + ")" + " Get" + "(")
	lowfieldsVarName := lowerfirst(fieldsVarName[0])
	if !checkFieldType(fieldsType[0]) {
		panic(fmt.Errorf("%s表 字段类型配置错误 字段名:%s 字段类型:%s", sheetName, fieldsName[0], fieldsType[0]))
	}
	lowsheetName := lowerfirst(sheetName)
	w.WriteString(lowfieldsVarName + " " + fieldsType[0] + ")" + " *" + sheetName + " {\n")
	w.WriteString("    " + lowsheetName + " , ok" + " := " + mgrInstance + ".Datas.Load(uint32(" + lowfieldsVarName + "))\n")
	w.WriteString("    if !ok {\n")
	w.WriteString("        return nil\n")
	w.WriteString("    }\n")
	w.WriteString("    return " + lowsheetName + ".(*" + sheetName + ")" + "\n")
	w.WriteString("}\n\n")

	w.WriteString("func (mgr *" + mgrName + ") " + "Load" + sheetName + "(fileName string)" + " {\n")
	//w.WriteString(`    function.DummyFunc() // 不要删除这段代码, 这段代码是为了防止function包导入未使用的错误`)
	//w.WriteString("\n")
	w.WriteString(`    time.Now() // 不要删除这段代码, 这段代码是为了防止time包导入未使用的错误`)
	w.WriteString("\n")
	w.WriteString("    xlFile, err := xlsx.OpenFile(filepath.FromSlash(strings.TrimSpace(fileName)))\n")
	w.WriteString("    if err != nil {\n")
	w.WriteString("        panic(err)\n")
	w.WriteString("    }\n")
	w.WriteString("    mgr.Datas = &sync.Map{} \n")
	w.WriteString("    sheet := xlFile.Sheets[0]\n")
	w.WriteString("    var cell *xlsx.Cell\n")
	w.WriteString("    for i, row := range sheet.Rows {\n")
	w.WriteString("        // 第三行开始才是配置数据\n")
	lenFieldsStr := strconv.FormatInt(int64(len(fieldsVarName)), 10)
	w.WriteString("        if i <= 2  || len(row.Cells) == 0  { \n")
	w.WriteString("            continue\n")
	w.WriteString("        }\n")
	w.WriteString(`        if row.Cells[0].String() == "" { ` + "\n")
	w.WriteString("            continue\n")
	w.WriteString("        }\n")
	w.WriteString("        if len(row.Cells) < " + lenFieldsStr + " { \n")
	w.WriteString(`            panic(fmt.Sprintf("配表载入错误 ` + "表名:" + sheetName + ` 行号:%v", i))` + "\n")
	w.WriteString("        }\n")
	w.WriteString("        item := &" + sheetName + "{}\n")
	for i, fieldVarName := range fieldsVarName {
		w.WriteString("        cell = row.Cells[" + strconv.Itoa(i) + "]\n")
		if !checkFieldType(fieldsType[i]) {
			panic(fmt.Errorf("%s表 字段类型配置错误 字段名:%s 字段类型:%s", sheetName, fieldsName[i], fieldsType[i]))
		}
		xlsxType := getXlsxType(fieldsType[i])
		if xlsxType == "String" {
			switch fieldsType[i] {
			case "[]string":
				w.WriteString("        item." + fieldVarName + " = " + "strings.Split" + "(cell." + xlsxType + "()")
				w.WriteString(`, "|"`)
				w.WriteString(")\n")
			case "[][]string":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringTo2dStringSlice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|", ","`)
				w.WriteString(")\n")
			case "[]int32":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringToInt32Slice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|"`)
				w.WriteString(")\n")
			case "[][]int32":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringTo2dInt32Slice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|", ","`)
				w.WriteString(")\n")
			case "[]uint32":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringToUint32Slice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|"`)
				w.WriteString(")\n")
			case "[][]uint32":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringTo2dUint32Slice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|", ","`)
				w.WriteString(")\n")
			case "[]float32":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringToFloat32Slice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|"`)
				w.WriteString(")\n")
			case "[][]float32":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.SplitStringTo2dFloat32Slice(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(`, "|", ","`)
				w.WriteString(")\n")
			case "time":
				w.WriteString("        item." + fieldVarName + " = ")
				w.WriteString("function.ParseTime2(strings.TrimSpace(cell." + xlsxType + "())")
				w.WriteString(")\n")
			default:
				w.WriteString("        item." + fieldVarName + " = " + fieldsType[i] + "(strings.TrimSpace(cell." + xlsxType + "()))\n")
			}
		} else if xlsxType == "Bool" {
			w.WriteString("        item." + fieldVarName + " = " + fieldsType[i] + "(cell." + xlsxType + "())\n")
		} else {
			w.WriteString("        data" + fieldVarName + ", _ := cell." + xlsxType + "()\n")
			//w.WriteString("        if ok != nil {  panic(ok) } \n")
			w.WriteString("        item." + fieldVarName + " = " + fieldsType[i] + "(data" + fieldVarName + ")\n")
			if i == 0 {
				w.WriteString("        typeId := " + "uint32(data" + fieldVarName + ")\n")
			}
		}
	}
	w.WriteString("        mgr.Datas.Store(typeId, item)\n")
	w.WriteString("    }\n")
	w.WriteString("}\n")

	w.Flush()

}

func createLoadAllXlsxFileCode() {
	mgrFileName := filepath.FromSlash(*mgrFile)
	mgrFile, err := os.OpenFile(mgrFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer mgrFile.Close()
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(mgrFile)

	pkgName := fmt.Sprintf("package %v\n\n", *pkgName)

	w.WriteString(pkgName)

	w.WriteString("var IsLoad bool \n\n")

	w.WriteString("func InitializeAllXlsxData() {\n")

	w.WriteString("    IsLoad = true \n")

	w.Flush()
}

func createLoadSingleXlsxFileCode() {
	mgrFileName := filepath.FromSlash(*mgrFile)
	mgrFile, err := os.OpenFile(mgrFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(mgrFile)
	w.WriteString("func ReloadSingleXlsx(fileName string) {\n")
	w.WriteString("    if !IsLoad {\n")
	w.WriteString("        return\n")
	w.WriteString("    }\n")
	w.Flush()
}

func autoGenReloadEachTableFuncCode(fileName string) {
	excelFileName := filepath.FromSlash(fileName)
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("open excel file %v error %v", excelFileName, err)
	}
	if len(xlFile.Sheets) > 0 {
		sheetName := xlFile.Sheets[0].Name
		mgrFileName := filepath.FromSlash(*mgrFile)
		mgrFile, err := os.OpenFile(mgrFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}
		w := bufio.NewWriter(mgrFile)
		w.WriteString("    if fileName == ")
		w.WriteString(`"`)
		w.WriteString(path.Base(fileName))
		w.WriteString(`" {`)
		w.WriteString(" \n")
		w.WriteString("        " + sheetName + "Mgr_GetMe().Load" + sheetName + "(")
		w.WriteString(`"`)
		w.WriteString(filepath.ToSlash(fileName))
		w.WriteString(`"`)
		w.WriteString(")\n")
		w.WriteString("    }\n")
		w.Flush()
	}
}

func autoGenLoadEachTableFuncCode(fileName, sheetName string) {
	mgrFileName := filepath.FromSlash(*mgrFile)
	mgrFile, err := os.OpenFile(mgrFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(mgrFile)
	w.WriteString("    " + sheetName + "Mgr_GetMe().Load" + sheetName + "(")
	w.WriteString(`"`)
	w.WriteString(filepath.ToSlash(fileName))
	w.WriteString(`"`)
	w.WriteString(")\n")
	w.Flush()
}

func autoGenFuncCodeEnd() {
	mgrFileName := filepath.FromSlash(*mgrFile)
	mgrFile, err := os.OpenFile(mgrFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(mgrFile)
	w.WriteString("}\n\n\n")
	w.Flush()
}

func main() {
	flag.Parse()

	xmlFile, err := os.Open(filepath.FromSlash(*tableList)) // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer xmlFile.Close()
	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	tables := configTables{}
	err = xml.Unmarshal(data, &tables)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	createLoadAllXlsxFileCode()
	for _, table := range tables.Tables {
		if path.Ext(table.TableName) != ".xlsx" && path.Ext(table.TableName) != "xls" {
			continue
		} else {
			fileName := *conf + "/" + table.TableName
			parseExcelFile(fileName)
		}
	}
	autoGenFuncCodeEnd()
	createLoadSingleXlsxFileCode()
	for _, table := range tables.Tables {
		if path.Ext(table.TableName) != ".xlsx" && path.Ext(table.TableName) != "xls" {
			continue
		} else {
			fileName := *conf + "/" + table.TableName
			autoGenReloadEachTableFuncCode(fileName)
		}
	}
	autoGenFuncCodeEnd()
	testMap = make(map[uint32]string, 0)
	fmt.Println("xlsx2pbdata gen code file file finish...")
}
