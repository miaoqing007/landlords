package configtable

import (
    "sync"
    "fmt"
    "path/filepath"
    "strings"
    "time"
    "github.com/tealeg/xlsx"
)

var (
    areaTableMgrIns *AreaTableMgr
    areaTableMgrOnce sync.Once
)

type AreaTableMgr struct {
    Datas *sync.Map
}

func AreaTableMgr_GetMe() *AreaTableMgr {
    areaTableMgrOnce.Do( func () {
        areaTableMgrIns = &AreaTableMgr {
            Datas: &sync.Map{},
        }
    })
    return areaTableMgrIns
}

func AreaTableMgr_GetSize() int {
    var size int
    areaTableMgrIns.Datas.Range(func(key, value interface{}) bool {
        size++
        return true
    })
    return size
}

type AreaTable struct {
    AreaId uint32 // 场景ID
    Desc string // 备注
    Name string // 场景名称
    MainScene uint32 // 所属主场景
    BornPoint uint32 // 出生点
    AreaNPC string // 区域NPC
    AreaIcon string // 大地图区块
    AreaScene string // 区域室内场景
    AreaSequenceId uint32 // 大地图显示排序
    Sound uint32 // 背景音乐
    AreaBuilding string // 区域功能物件
    Elficon string // 精灵图标
    ElfCoord string // 精灵传送坐标
    ElfList string // 怪物列表
    FishCoord string // 钓鱼点坐标
}

func (mgr *AreaTableMgr) Get(areaId uint32) *AreaTable {
    areaTable , ok := areaTableMgrIns.Datas.Load(uint32(areaId))
    if !ok {
        return nil
    }
    return areaTable.(*AreaTable)
}

func (mgr *AreaTableMgr) LoadAreaTable(fileName string) {
    time.Now() // 不要删除这段代码, 这段代码是为了防止time包导入未使用的错误
    xlFile, err := xlsx.OpenFile(filepath.FromSlash(strings.TrimSpace(fileName)))
    if err != nil {
        panic(err)
    }
    mgr.Datas = &sync.Map{} 
    sheet := xlFile.Sheets[0]
    var cell *xlsx.Cell
    for i, row := range sheet.Rows {
        // 第三行开始才是配置数据
        if i <= 2  || len(row.Cells) == 0  { 
            continue
        }
        if row.Cells[0].String() == "" { 
            continue
        }
        if len(row.Cells) < 15 { 
            panic(fmt.Sprintf("配表载入错误 表名:AreaTable 行号:%v", i))
        }
        item := &AreaTable{}
        cell = row.Cells[0]
        dataAreaId, _ := cell.Int()
        item.AreaId = uint32(dataAreaId)
        typeId := uint32(dataAreaId)
        cell = row.Cells[1]
        item.Desc = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[2]
        item.Name = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[3]
        dataMainScene, _ := cell.Int()
        item.MainScene = uint32(dataMainScene)
        cell = row.Cells[4]
        dataBornPoint, _ := cell.Int()
        item.BornPoint = uint32(dataBornPoint)
        cell = row.Cells[5]
        item.AreaNPC = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[6]
        item.AreaIcon = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[7]
        item.AreaScene = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[8]
        dataAreaSequenceId, _ := cell.Int()
        item.AreaSequenceId = uint32(dataAreaSequenceId)
        cell = row.Cells[9]
        dataSound, _ := cell.Int()
        item.Sound = uint32(dataSound)
        cell = row.Cells[10]
        item.AreaBuilding = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[11]
        item.Elficon = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[12]
        item.ElfCoord = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[13]
        item.ElfList = string(strings.TrimSpace(cell.String()))
        cell = row.Cells[14]
        item.FishCoord = string(strings.TrimSpace(cell.String()))
        mgr.Datas.Store(typeId, item)
    }
}
