package service

import (
	model "ccccc/data/model/goods"
	"fmt"
	"log"
	"strconv"
)

func GetGoodsListGroupByNameDraf(shafaId string, transId string) (list [][]string, err error) {
	//查询 沙发信息

	shafa := model.ShaFaImportLog{}
	err = shafa.Get("", shafaId)
	if err != nil {
		fmt.Printf("GetGoodsListGroupByName err %v\n", err)
		return
	}

	info := model.GongYiDrafList{}
	err = info.GetListPageWithTrans(shafaId, "", transId, nil)
	//有没有呢
	if err != nil {
		log.Printf("GetListPage err :%v", err)
	}
	fenweiMap := make(map[string][]model.GongYiDraf, 0)
	cailiaoMap := make(map[string][]model.GongYiDraf, 0)
	//先根据分位聚合 不需要人工和其他
	for _, yi := range info {
		if yi.CpCode != "" {
			//说明是材料
			_, ok := fenweiMap[yi.FenWeiName]
			if ok {
				fenweiMap[yi.FenWeiName] = append(fenweiMap[yi.FenWeiName], yi)
			} else {
				fenweiMap[yi.FenWeiName] = make([]model.GongYiDraf, 0)
				fenweiMap[yi.FenWeiName] = append(fenweiMap[yi.FenWeiName], yi)
			}

			_, ok1 := cailiaoMap[yi.CpCode]
			if ok1 {
				cailiaoMap[yi.CpCode] = append(cailiaoMap[yi.CpCode], yi)
			} else {
				cailiaoMap[yi.CpCode] = make([]model.GongYiDraf, 0)
				cailiaoMap[yi.CpCode] = append(cailiaoMap[yi.CpCode], yi)
			}
		}
	}
	//outinfo
	outInfo := make([][]string, len(cailiaoMap)+1)

	//第一行数据
	outInfo[0] = make([]string, 0)
	outInfo[0] = append(outInfo[0], "产品编码")
	outInfo[0] = append(outInfo[0], "产品名称")
	outInfo[0] = append(outInfo[0], "产品规格")
	outInfo[0] = append(outInfo[0], "材料编码")
	outInfo[0] = append(outInfo[0], "材料名称")
	outInfo[0] = append(outInfo[0], "单位")
	for s, _ := range fenweiMap {
		outInfo[0] = append(outInfo[0], s)
	}
	ii := 1
	for _, yis := range cailiaoMap {
		outInfo[ii] = make([]string, 0)
		outInfo[ii] = append(outInfo[ii], shafa.SfCode)
		outInfo[ii] = append(outInfo[ii], shafa.SfName)
		outInfo[ii] = append(outInfo[ii], shafa.GG)
		outInfo[ii] = append(outInfo[ii], yis[0].CpCode)

		//查询最新名称
		goods_info := model.Goods{}
		errxxx := goods_info.Get("", yis[0].CpCode)
		if errxxx == nil {
			outInfo[ii] = append(outInfo[ii], goods_info.CpName)
		} else {
			outInfo[ii] = append(outInfo[ii], yis[0].CLName)
		}

		//查询单位
		//这里去要区分单位和数量
		unit := ""

		if goods_info.MainSize != "" {
			unit = yis[0].JiJiaUnit
		} else {
			unit, _ = getGoodsNuit(goods_info)
		}

		outInfo[ii] = append(outInfo[ii], unit)
		tempMpa := make(map[string]float64, 0)
		for _, yi := range yis {
			_, okk := tempMpa[yi.FenWeiName]
			nums := ""
			if yi.Size != "" || goods_info.MainSize != "" {
				nums = yi.JiJiaNum
			} else {
				nums = yi.Nums
			}
			fll, _ := strconv.ParseFloat(nums, 64)
			if okk {
				tempMpa[yi.FenWeiName] = tempMpa[yi.FenWeiName] + fll
			} else {
				tempMpa[yi.FenWeiName] = fll
			}
		}

		for i := 6; i < len(outInfo[0]); i++ {
			llll, okkk := tempMpa[outInfo[0][i]]
			if !okkk {
				llll = 0
			}
			if llll != 0 {
				outInfo[ii] = append(outInfo[ii], fmt.Sprintf("%.4f", llll))

			} else {
				outInfo[ii] = append(outInfo[ii], "")
			}
		}
		ii++
	}
	list = outInfo
	return

}

func GetDrafShaFaImportList(page, size int, name, code string, issums string) (list []ShaFaInfoDto, total int, err error) {

	offset := (page - 1) * size

	data := model.ShaFaDrafImportLogList{}
	total, err = data.GetListPage(name, code, offset, size, issums, nil)
	if err != nil {
		log.Printf("GetShaFaImportList err :%v", err)
		return
	}

	for _, datum := range data {
		iiinfo := ShaFaInfoDto{
			Id:     datum.Id,
			SfName: datum.SfName,

			SfCode:     datum.SfCode,
			SDesc:      datum.SDesc,
			ImportUser: datum.ImportUser,
			CreateTime: datum.CreateTime,
			UpdateTime: datum.UpdateTime,
			GG:         datum.GG,
			IsSums:     datum.IsSums,
		}

		trans := model.Trans{}
		trans.GetByShafaIdAndTrans(nil, datum.SfCode, datum.TransId)
		if trans.IsSubmit == 1 {
			iiinfo.IsOnline = "是"
		} else {
			iiinfo.IsOnline = "否"

		}
		if trans.IsCheck == 1 {
			iiinfo.IsCheck = "是"
		} else {
			iiinfo.IsCheck = "否"
		}
		iiinfo.CheckUser = trans.CheckUser

		list = append(list, iiinfo)
	}
	return

}

func GetDrafShaFaImportById(code string) (data model.ShaFaDrafImportLog, err error) {

	data = model.ShaFaDrafImportLog{}
	err = data.GetOne(code)
	if err != nil {
		log.Printf("GetShaFaImportList err :%v", err)
		return
	}

	return data, nil

}

func GetFenWeiListGroupByNameDraf(shafaId, types, transeId string) (list []GetFenWeiListGroupByNameDto, err error) {

	info := model.GongYiDrafList{}
	err = info.GetListPageWithTrans(shafaId, types, transeId, nil)
	//有没有呢
	if err != nil {
		log.Printf("GetFenWeiList err :%v", err)
	}
	// 要按照 名称分类
	typeMaps := make(map[string][]GetFenWeiListGroupByNameInfo, 0)

	for _, yi := range info {
		l3 := GetFenWeiListGroupByNameInfo{
			Types:      yi.Types,
			FenWeiName: yi.FenWeiName,
			CLName:     yi.CLName,
			Size:       yi.Size,
			Nums:       yi.Nums,
			Unit:       yi.Unit,
			Descs:      yi.Descs,
			TotalPrice: yi.TotalPrice,
			CpCode:     yi.CpCode,
			GoodsPoint: yi.GoodsPoint,
		}
		//获取单位
		goods_info := model.Goods{}
		errxxx := goods_info.Get("", yi.CpCode)
		if errxxx == nil {
			l3.CLName = goods_info.CpName
			//
			if goods_info.CpMainUnit == goods_info.FuZhuUnit {
				//说明是一样的 可以都不做
			} else if goods_info.MainXiShu > 0 {
				intv1 := goods_info.MainXiShu
				intv2 := goods_info.FuZhuXiShu

				tmp := 0
				if intv1 > intv2 {
					tmp, _ = strconv.Atoi(goods_info.CpMainUnit)

				} else {
					tmp, _ = strconv.Atoi(goods_info.FuZhuUnit)
				}
				unit := model.UnitDesc{}
				err2 := unit.GetById(nil, tmp)
				//获取
				if err2 != nil {
					log.Printf("unit.GetById(nil, intv1) err :%v", err)
				} else {
					l3.Unit = unit.Name
				}

			}

		}

		if len(l3.Size) > 0 {
			l3.Unit = "件"
		}

		ll, ok := typeMaps[yi.GongYiName]
		if ok {
			ll = append(ll, l3)
			typeMaps[yi.GongYiName] = ll
		} else {
			l1 := make([]GetFenWeiListGroupByNameInfo, 0)
			l1 = append(l1, l3)
			typeMaps[yi.GongYiName] = l1
		}

	}

	for gongyiName, infos := range typeMaps {

		// 这里 还要处理下数据
		outList := GetFenWeiListGroupByNameDto{
			GongYiName: gongyiName,
			Size:       len(infos),
			List:       infos,
		}
		fenWeiString := make(map[string]bool, 0)
		for _, nameInfo := range infos {
			if nameInfo.FenWeiName != "全套" && nameInfo.FenWeiName != "子件" && nameInfo.FenWeiName != "外框架" {
				fenWeiString[nameInfo.FenWeiName] = true
			}
		}
		sssl := ""
		for s, _ := range fenWeiString {
			sssl += s
			sssl += "+"
		}
		if len(sssl) > 0 {
			sssl = sssl[0 : len(sssl)-1]
		}
		outList.FenWeiList = sssl
		list = append(list, outList)
	}
	return

}

type TranseDto struct {
	Id   string `json:"id"`
	User string `json:"user"`
	Time string `json:"time"`
	//是否审核过
	Issubmit   string `json:"issubmit"`
	Ischeck    string `json:"ischeck"`
	SfCode     string `json:"sf_code"`
	OnlineUser string `json:"online_user"`
	OnlineTime string `json:"online_time"`
	//是否提交
	IsCheck   string `json:"is_check"`
	CheckUser string `json:"check_user"`
}

func GetTransList(shafaId string) (out []TranseDto) {
	out = make([]TranseDto, 0)
	list := model.TransList{}
	list.GetByShafaId(nil, shafaId)

	if list != nil && len(list) != 0 {
		for _, trans := range list {
			dto := TranseDto{}
			dto.Id = trans.TransId
			dto.SfCode = trans.ShafaCode
			dto.Time = trans.CreateTime.Format("2006-01-02 15:04:05")
			dto.User = trans.CreateUser
			dto.OnlineUser = trans.OnlineUser
			dto.OnlineTime = trans.OnlineTime.Format("2006-01-02 15:04:05")
			if trans.IsSubmit == 1 {
				dto.Issubmit = "是"
			} else {
				dto.Issubmit = "否"
			}
			if trans.IsCheck == 1 {
				dto.IsCheck = "是"
			} else {
				dto.IsCheck = "否"
			}
			dto.CheckUser = trans.CheckUser

			// 查询当前的 事务
			imports := model.ShaFaDrafImportLog{}
			imports.GetOne(shafaId)
			if imports.TransId == trans.TransId {
				dto.Ischeck = "是"
			} else {
				dto.Ischeck = "否"
			}
			out = append(out, dto)
		}
	}
	return

}
