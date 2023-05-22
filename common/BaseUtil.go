package common

import (
	model "ccccc/data/model/goods"
	"log"
	"strings"
)

//单位集合
var UnitList []Unit

//分位描述集合
var FenWeiListDesc []FenWeiDesc

// 产品分类集合
var GoodsTypeListDesc []GoodsTypeDesc
var GoodsTypeMap map[string]GoodsTypeDesc
var GoodsTypeCodeMap map[string]GoodsTypeDesc
var UnitMap map[int]Unit

func InitBaseData() {
	initUnit()
	initFenWei()
	initGoodsType()

}

type Unit struct {
	Id int `json:"id"`
	//单位名称
	Name string `json:"name"`
	//换算比例
	Hy int `json:"hy"`
}

//
type FenWeiDesc struct {
	//编码
	FWBM string `json:"fwbm"`
	//名称
	FWMC string `json:"fwmc"`
	//分位系数
	FWXS float64 `json:"fwxs"`
	//简介
	FWJC string `json:"fwjc"`
}

//产品分类
type GoodsTypeDesc struct {
	//编码
	GoodsTypeId   string `json:"goods_type_id"`
	GoodsTypeCode string `json:"goods_type_code"`
	GoodsTypeName string `json:"goods_type_name"`
}

func initGoodsType() {
	dbDataList := model.GoodsTypeList{}
	GoodsTypeMap = make(map[string]GoodsTypeDesc, 0)
	GoodsTypeCodeMap = make(map[string]GoodsTypeDesc, 0)
	err := dbDataList.GetAll(nil)
	if err != nil {
		log.Fatalf("initGoodsType faile err:  %v\n", err)
		return
	}
	//
	GoodsTypeListDesc = make([]GoodsTypeDesc, 0)
	for _, wei := range dbDataList {
		info := GoodsTypeDesc{
			GoodsTypeId:   wei.GoodsTypeId,
			GoodsTypeName: strings.TrimSpace(wei.GoodsTypeName),
			GoodsTypeCode: strings.TrimSpace(wei.GoodsTypeCode),
		}
		GoodsTypeListDesc = append(GoodsTypeListDesc, info)
		GoodsTypeMap[strings.TrimSpace(wei.GoodsTypeName)] = info
		GoodsTypeCodeMap[strings.TrimSpace(wei.GoodsTypeId)] = info
	}
}

//从数据库读取 分位描述
func initFenWei() {
	dbDataList := model.FenWeiList{}
	err := dbDataList.GetAll(nil)
	if err != nil {
		log.Fatalf("initFenWei faile err:  %v\n", err)
		return
	}
	//
	FenWeiListDesc = make([]FenWeiDesc, 0)
	for _, wei := range dbDataList {
		FenWeiListDesc = append(FenWeiListDesc, FenWeiDesc{
			FWMC: wei.FWMC,
			FWBM: wei.FWBM,
			FWXS: wei.FWXS,
			FWJC: wei.FWJC,
		})
	}

}

//初始化单位 直接写死了
func initUnit() {
	UnitList = make([]Unit, 0)
	/*	UnitList = append(UnitList, Unit{"套", 1})
		UnitList = append(UnitList, Unit{"公斤", 1})
		UnitList = append(UnitList, Unit{"件", 1})
		UnitList = append(UnitList, Unit{"根", 1})
		UnitList = append(UnitList, Unit{"个", 1})
		UnitList = append(UnitList, Unit{"米", 1})
		UnitList = append(UnitList, Unit{"立方米", 1})
		UnitList = append(UnitList, Unit{"张", 1})
		UnitList = append(UnitList, Unit{"CM", 1})
		UnitList = append(UnitList, Unit{"片", 1})
		UnitList = append(UnitList, Unit{"板", 1})
		UnitList = append(UnitList, Unit{"把", 1})
		UnitList = append(UnitList, Unit{"条", 1})
		UnitList = append(UnitList, Unit{"对", 1})
		UnitList = append(UnitList, Unit{"把(根)", 1})
		UnitList = append(UnitList, Unit{"只", 1})
		UnitList = append(UnitList, Unit{"盘", 1})
		UnitList = append(UnitList, Unit{"KG", 1})
		UnitList = append(UnitList, Unit{"桶", 1})
		UnitList = append(UnitList, Unit{"盒", 1})
		UnitList = append(UnitList, Unit{"码", 1})
		UnitList = append(UnitList, Unit{"扎", 1})
		UnitList = append(UnitList, Unit{"本", 1})
		UnitList = append(UnitList, Unit{"支", 1})
		UnitList = append(UnitList, Unit{"瓶", 1})
		UnitList = append(UnitList, Unit{"卷", 1})
		UnitList = append(UnitList, Unit{"包", 1})
		UnitList = append(UnitList, Unit{"台", 1})
		UnitList = append(UnitList, Unit{"颗", 1})
		UnitList = append(UnitList, Unit{"双", 1})
		UnitList = append(UnitList, Unit{"圈", 1})
		UnitList = append(UnitList, Unit{"副", 1})
		UnitList = append(UnitList, Unit{"米", 1})
		UnitList = append(UnitList, Unit{"排", 1})
		UnitList = append(UnitList, Unit{"箱", 1})
		UnitList = append(UnitList, Unit{"英尺", 1})
		UnitList = append(UnitList, Unit{"份", 1})
		UnitList = append(UnitList, Unit{"立方", 1})
		UnitList = append(UnitList, Unit{"平方米", 1})
		UnitList = append(UnitList, Unit{"套数", 1})
	*/

	uuList := model.UnitDescList{}
	err := uuList.GetAll(nil)

	mapp := make(map[string]bool, 0)

	if err == nil {
		for _, desc := range uuList {

			name := desc.Name
			name = strings.Replace(name, " ", "", -1)
			name = strings.Replace(name, "/r", "", -1)
			name = strings.Replace(name, "\r", "", -1)
			name = strings.Replace(name, "\n", "", -1)
			_, ok := mapp[desc.Name]
			if !ok {
				UnitList = append(UnitList, Unit{
					desc.Id,
					desc.Name,
					desc.Types})
				mapp[desc.Name] = true
			}

		}
	} else {
		log.Printf("init uuList error : %v", err)
		panic("init uuList error")
	}

}
