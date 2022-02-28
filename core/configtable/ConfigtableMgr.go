package configtable

var IsLoad bool 

func InitializeAllXlsxData() {
    IsLoad = true 
    AreaTableMgr_GetMe().LoadAreaTable("../conf/excel/AreaTable.xlsx")
}


func ReloadSingleXlsx(fileName string) {
    if !IsLoad {
        return
    }
    if fileName == "AreaTable.xlsx" { 
        AreaTableMgr_GetMe().LoadAreaTable("../conf/excel/AreaTable.xlsx")
    }
}


