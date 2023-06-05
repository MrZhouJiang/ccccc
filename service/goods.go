package service

import (
	"ccccc/common"
	model "ccccc/data/model/goods"
	model2 "ccccc/data/model/user"
	"fmt"
	"log"
	"strconv"
	"time"
)

func GetGoodsList(page, size int, name, goodsCode, typeName, shunhao string) (gs []GoodsDto, total int, err error) {
	offset := (page - 1) * size

	data := model.GoodsList{}

	//根据 GoodsTypeMap 获取 GoodsTypeCOde
	typeInfo, ok := common.GoodsTypeMap[typeName]
	typeCode := ""
	if ok {
		typeCode = typeInfo.GoodsTypeId
	}

	total, err = data.GetListPage(name, goodsCode, typeCode, shunhao, offset, size, nil)
	if err != nil {
		log.Printf("GetShaFaImportList err :%v", err)
		return
	}

	//设置产品类型
	for i, datum := range data {
		cd, okk := common.GoodsTypeCodeMap[datum.CpTypeCode]
		if okk {
			data[i].CpType = cd.GoodsTypeName
		}

		//获取 损耗率 和bom转换 损耗率 以数据库为准
		changeDesc := model.GoodsChangeDesInfoDtoList{}
		errx := changeDesc.GetListByCpCode(datum.CpCode, nil)
		if errx == nil {
			for _, info := range changeDesc {
				/*				if info.ChangeType == "损耗" && (datum.ShunHao == "" || datum.ShunHao == "0") {
								data[i].ShunHao = fmt.Sprintf("%.3f", info.ValuesL)
							}*/
				if info.ChangeType == "换算" {

					data[i].DomH = fmt.Sprintf("%s%.7f", info.Types, info.ValuesL)
					data[i].DomId = info.Id
				}
			}
		}
		//获取单位
		unit := model.UnitDesc{}

		intv1, err1 := strconv.Atoi(datum.CpMainUnit)
		unit.GetById(nil, intv1)
		//获取
		if err1 != nil {
			log.Printf("GetShaFaImportList err :%v", err)
			return
		}
		unit2 := model.UnitDesc{}
		intv2, err1 := strconv.Atoi(datum.FuZhuUnit)
		unit2.GetById(nil, intv2)
		//获取
		if err1 != nil {
			log.Printf("GetShaFaImportList err :%v", err)
			return
		}

		dto := GoodsDto{
			Id:          data[i].Id,
			CpCode:      data[i].CpCode,
			CpName:      data[i].CpName,
			CpDesc:      data[i].CpDesc,
			CpType:      data[i].CpType,
			CpTypeCode:  data[i].CpTypeCode,
			CpGuiGe:     data[i].CpGuiGe,
			CpMainUnit:  unit.Name,
			FuZhuUnit:   unit2.Name,
			MainSize:    data[i].MainSize,
			MainXiShu:   data[i].MainXiShu,
			FuZhuXiShu:  data[i].FuZhuXiShu,
			Price:       data[i].Price,
			ChangeP:     data[i].ChangeP,
			LoadTime:    data[i].LoadTime,
			CreateTime:  data[i].CreateTime,
			DomH:        data[i].DomH,
			DomId:       data[i].DomId,
			TyName:      data[i].TyName,
			ShunHao:     data[i].ShunHao,
			GuDingPrice: data[i].GuDingPrice,
		} /**/
		// 查看是否有 合并价格配置
		mergeDesc := model.GoodsMergeDesInfoDtoList{}
		errx1 := mergeDesc.GetListByCpCode(datum.CpCode, nil)
		if errx1 == nil {
			for _, info := range mergeDesc {
				dto.MergeId = info.Id
				dto.MergeName = info.Name

			}
		}
		gs = append(gs, dto)

	}

	return
}

type GoodsDto struct {
	Id         int    `json:"id"`
	CpCode     string `json:"cp_code"`
	CpName     string `json:"cp_name"`
	CpDesc     string `json:"cp_desc"`
	CpType     string `json:"cp_type"`
	CpTypeCode string `json:"cp_type_code"`
	//规格(成品有)
	CpGuiGe string `json:"cp_gui_ge"`
	//主计量单位
	CpMainUnit string `json:"cp_main_unit"`
	//辅计量单位
	FuZhuUnit string `json:"fu_zhu_unit"`
	//产品尺寸 3*5 （cm） 注意 都是里面单位 3 表示3 cm   （3*5*2）
	MainSize string `json:"main_size"`
	//主系数
	MainXiShu float64 `json:"main_xi_shu"`
	//辅助系数
	FuZhuXiShu float64   `json:"fu_zhu_xi_shu"`
	Price      float64   `json:"price"`
	ShunHao    string    `json:"shun_hao"`
	ChangeP    float64   `json:"change_p"`
	LoadTime   time.Time `json:"load_time"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	DomH       string    `json:"dom_h"`
	DomId      int       `json:"dom_id"`
	//通用名称
	TyName    string `json:"ty_name"`
	MergeName string `json:"merge_name"`
	MergeId   int    `json:"merge_id"`
	//固定价格
	GuDingPrice float64 `json:"gu_ding_price"`
}

func GeAllGoodsList() (gs []model.AllGoodsDesc, err error) {

	//需要去除一点 goods_type ：

	notInid := make([]int, 0)
	notInid = append(notInid, 1015)
	notInid = append(notInid, 1018)
	notInid = append(notInid, 1019)
	notInid = append(notInid, 1020)
	notInid = append(notInid, 1021)
	notInid = append(notInid, 1004)
	notInid = append(notInid, 1026)
	notInid = append(notInid, 1027)
	notInid = append(notInid, 1012)
	notInid = append(notInid, 1029)
	notInid = append(notInid, 1030)

	data := model.GoodsList{}
	gs, err = data.GetAllGoodsList(nil, notInid)
	if err != nil {
		log.Printf("GetShaFaImportList err :%v", err)
		return
	}
	return
}

func GeGoodsById(id string) (gs model.AllGoodsDesc, err error) {

	data := model.Goods{}
	gs, err = data.GetGoodsById(id, nil)
	//获取单位
	unit := model.UnitDesc{}
	unit2 := model.UnitDesc{}
	intv1, err := strconv.Atoi(gs.CpMainUnit)
	intv2, err := strconv.Atoi(gs.FuZhuUnit)
	unit.GetById(nil, intv1)
	gs.CpMainUnit = unit.Name
	gs.CpMainUnitId = intv1
	unit2.GetById(nil, intv2)
	gs.FuZhuUnit = unit2.Name
	gs.FuZhuUnitId = intv2

	//获取
	if err != nil {
		log.Printf("GetShaFaImportList err :%v", err)
		return
	}

	cd, okk := common.GoodsTypeCodeMap[gs.CpTypeCode]
	if okk {
		gs.CpType = cd.GoodsTypeName
	}

	return
}

func GetShaFaImportList(page, size int, name, code string, issums string) (list []ShaFaInfoDto, total int, err error) {

	offset := (page - 1) * size

	data := model.ShaFaImportLogList{}
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
		list = append(list, iiinfo)
	}
	return

}

type ShaFaInfoDto struct {
	Id         int       `json:"id"`
	SfName     string    `json:"sf_name"`
	SfCode     string    `json:"sf_code"`
	SDesc      string    `json:"s_desc"`
	ImportUser string    `json:"import_user"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	//沙发规格（可以解析分位名称）
	GG     string `json:"gg"`
	IsSums string `json:"is_sums"`
}

func GetFenWeiList(shafaId, types string) (list []model.GongYi, err error) {

	info := model.GongYiList{}
	err = info.GetListPage(shafaId, types, nil)
	//有没有呢
	if err != nil {
		log.Printf("GetFenWeiList err :%v", err)
	}
	list = info
	return

}

func GetFenWeiListGroupByName(shafaId, types string) (list []GetFenWeiListGroupByNameDto, err error) {

	info := model.GongYiList{}
	err = info.GetListPage(shafaId, types, nil)
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
				intv1, _ := strconv.Atoi(goods_info.CpMainUnit)
				intv2, _ := strconv.Atoi(goods_info.FuZhuUnit)

				tmp := 0
				if intv1 > intv2 {
					tmp = intv1
				} else {
					tmp = intv2
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
		list = append(list, outList)
	}
	return

}

type GetFenWeiListGroupByNameDto struct {
	GongYiName  string                         `json:"gong_yi_name"`
	Size        int                            `json:"size"`
	List        []GetFenWeiListGroupByNameInfo `json:"list"`
	ShunhaoList []ShunHaoInfo                  `json:"shunhao_list"`
}
type ShunHaoInfo struct {
	Name    string  `json:"name"`
	ClList  string  `json:"cl_list"`
	Shunhao float64 `json:"shunhao"`
	Price   float64 `json:"price"`
}
type GetFenWeiListGroupByNameInfo struct {

	//裁工
	Types      string  `json:"types"`
	FenWeiName string  `json:"fen_wei_name"`
	CLName     string  `json:"cl_name"`
	Size       string  `json:"size"`
	Nums       string  `json:"nums"`
	Unit       string  `json:"unit"`
	Descs      string  `json:"descs"`
	Price      float64 `json:"price"`
	TotalPrice string  `json:"total_price"`
	CpCode     string  `json:"cp_code"`
	GoodsPoint string  `json:"goods_point"`
}

func GetUserListD(page, size int, username string) (gs []model2.UserInfo, total int, err error) {
	offset := (page - 1) * size

	data := model2.UserInfoList{}
	total, err = data.GetListPage(username, offset, size, nil)
	if err != nil {
		log.Printf("GetShaFaImportList err :%v", err)
		return
	}
	gs = data
	return
}

func GetGoodsChangeList(page, size int, changeType, name string) (gs []model.GoodsChangeDesc, total int, err error) {
	offset := (page - 1) * size

	data := model.GoodsChangeDescList{}

	total, err = data.GetListPage(changeType, name, offset, size, nil)
	if err != nil {
		log.Printf("GetGoodsChangeList err :%v", err)
		return
	}
	gs = data
	return

}

func GetGoodsMergeList(page, size int, name string) (gs []model.GoodMergeDesc, total int, err error) {
	offset := (page - 1) * size

	data := model.GoodMergeDescList{}

	total, err = data.GetListPage(name, offset, size, nil)
	if err != nil {
		log.Printf("GetGoodsMergeList err :%v", err)
		return
	}
	gs = data
	return

}

func getGoodsNuit(goods model.Goods) (string, int) {
	if goods.MainXiShu != 1 {
		unit := model.UnitDesc{}
		intv1, _ := strconv.Atoi(goods.CpMainUnit)
		unit.GetById(nil, intv1)
		return unit.Name, 1
	} else if goods.FuZhuXiShu != 1 {
		unit := model.UnitDesc{}
		intv1, _ := strconv.Atoi(goods.FuZhuUnit)
		unit.GetById(nil, intv1)
		return unit.Name, 2
	} else {
		unit := model.UnitDesc{}
		intv1, _ := strconv.Atoi(goods.CpMainUnit)
		unit.GetById(nil, intv1)
		return unit.Name, 1
	}
	return "", 1

}

func GetGoodsListGroupByName(shafaId string) (list [][]string, err error) {
	//查询 沙发信息

	shafa := model.ShaFaImportLog{}
	err = shafa.Get("", shafaId)
	if err != nil {
		fmt.Printf("GetGoodsListGroupByName err %v\n", err)
		return
	}

	info := model.GongYiList{}
	err = info.GetListPage(shafaId, "", nil)
	//有没有呢
	if err != nil {
		log.Printf("GetListPage err :%v", err)
	}
	fenweiMap := make(map[string][]model.GongYi, 0)
	cailiaoMap := make(map[string][]model.GongYi, 0)
	//先根据分位聚合 不需要人工和其他
	for _, yi := range info {
		if yi.CpCode != "" {
			//说明是材料
			_, ok := fenweiMap[yi.FenWeiName]
			if ok {
				fenweiMap[yi.FenWeiName] = append(fenweiMap[yi.FenWeiName], yi)
			} else {
				fenweiMap[yi.FenWeiName] = make([]model.GongYi, 0)
				fenweiMap[yi.FenWeiName] = append(fenweiMap[yi.FenWeiName], yi)
			}

			_, ok1 := cailiaoMap[yi.CpCode]
			if ok1 {
				cailiaoMap[yi.CpCode] = append(cailiaoMap[yi.CpCode], yi)
			} else {
				cailiaoMap[yi.CpCode] = make([]model.GongYi, 0)
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
