package controller

import (
	"ccccc/common"
	model "ccccc/data/model/goods"
	"ccccc/service"
	"ccccc/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"time"
)

func ChongSuanDraf(shafaCode string) {
	//根据沙发 找出所有成本
	allGongyi := model.GongYiDrafList{}

	err := allGongyi.GetBySoFaCode(nil, shafaCode)

	if err != nil {
		log.Printf("ChongSuan err :%v", err)
		return
	}

	for _, yi := range allGongyi {
		goods := model.Goods{}
		goods.CpCode = yi.CpCode
		err2 := goods.GetByCpCode(nil)
		if err2 != nil || goods.Id == 0 {
			continue
		}
		nums, _ := strconv.ParseFloat(yi.Nums, 64)
		ShunHao, _ := strconv.ParseFloat(goods.ShunHao, 64)
		newPeice := GetPrice(nums, ShunHao, goods.Price, goods.ChangeP, goods.FuZhuXiShu, goods.MainXiShu, yi.CpCode, goods.MainSize, yi.Size, yi.OwnerSize)

		yi.TotalPrice = fmt.Sprintf("%f", newPeice.TotalPrice)
		yi.JiJiaNum = fmt.Sprintf("%f", newPeice.JiJiaNums)
		errx := yi.Update(nil)
		if errx != nil {
			log.Printf("yi.Update err :%v", errx)

		}

	}

}

//重算沙发
func ReloadShaFaDraf(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	sf_code := params["cp_code"]

	ChongSuanDraf(sf_code)
	resp.Data = ""
	util.ReturnCompFunc(c, resp)
	return
}

//拷贝沙发
func CopyShaFaDraf(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	//被拷贝的
	copy_shafa_code := params["copy_shafa_code"]
	sf_code := params["cp_code"]

	d, err1 := service.GetDrafShaFaImportById(sf_code)

	if err1 != nil {
		log.Printf("GetDrafShaFaImportById.Delete err :%v", err1)
		resp.Status = 201
		resp.Desc = err1.Error()
		util.ReturnCompFunc(c, resp)
		return
	}
	transID := d.TransId
	if transID == "" {
		log.Printf("GetDrafShaFaImportById.  拷贝失败 没有选择版本")
		resp.Status = 201
		resp.Desc = "拷贝失败 没有选择版本"
		util.ReturnCompFunc(c, resp)
		return
	}

	allGongyi := model.GongYiDrafList{}

	//获取所有的配置
	err := allGongyi.GetBySoFaCodeDraf(nil, sf_code, transID)
	if err != nil || len(allGongyi) == 0 {
		log.Printf(" sf_code :%s, err:%v", sf_code, err)
		resp.Status = 201
		resp.Desc = "沙发成本为空"
		util.ReturnCompFunc(c, resp)
		return
	}

	shafa := model.ShaFaImportLog{}
	errme := shafa.GetByType(nil, copy_shafa_code)
	if shafa.Id == 0 {
		log.Printf(" 找不到沙发ID ：%s", copy_shafa_code)
		resp.Status = 201
		resp.Desc = "找不到该沙发"
		util.ReturnCompFunc(c, resp)
		return
	}

	/*	//先清空
		deleteInfo := model.GongYi{}
		deleteInfo.ShafaId = copy_shafa_code
		err = deleteInfo.Delete(nil)
		if err != nil {
			log.Printf("deleteInfo.Delete err :%v", err)
			resp.Status = 201
			resp.Desc = err.Error()
			util.ReturnCompFunc(c, resp)
			return
		}*/

	//要创建一个事务ID

	for _, yi := range allGongyi {
		yi.Id = 0
		yi.ShafaId = copy_shafa_code
		yi.Create(nil)
		if err != nil {
			log.Printf("insertInfo.Create err :%v", err)
			resp.Status = 201
			resp.Desc = err.Error()
			util.ReturnCompFunc(c, resp)
			return
		}
	}

	//修改 沙发表

	if errme == nil {
		shafa.IsSums = "是"
		shafa.Update(nil)
	}
	resp.Data = ""
	util.ReturnCompFunc(c, resp)
	return
}

func GetAllPriceDraf(c *gin.Context) {
	log.Printf("GetAllPriceDraf")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	sofa_code := params["sf_code"]
	trans_id := params["trans_id"]

	allGongyi := model.GongYiDrafList{}

	err := allGongyi.GetBySoFaCodeDraf(nil, sofa_code, trans_id)

	if err != nil {
		log.Printf("GetAllPrice err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}

	// 要根据全套 剔除掉其他的成本
	//生成一个key
	tttMp := make(map[string]bool, 0)

	tempList := make([]model.GongYiDraf, 0)

	for _, yi := range allGongyi {
		key := fmt.Sprintf("%s_%s_%s", yi.Types, yi.GongYiName, yi.CpCode)
		n1 := yi.FenWeiName
		n1 = strings.Replace(n1, " ", "", -1)
		n1 = strings.Replace(n1, "/r", "", -1)
		n1 = strings.Replace(n1, "\r", "", -1)
		n1 = strings.Replace(n1, "\n", "", -1)
		if n1 == "全套" {
			//说明有全套了
			tttMp[key] = true
		}
	}

	for _, yi := range allGongyi {
		key := fmt.Sprintf("%s_%s_%s", yi.Types, yi.GongYiName, yi.CpCode)
		n1 := yi.FenWeiName
		n1 = strings.Replace(n1, " ", "", -1)
		n1 = strings.Replace(n1, "/r", "", -1)
		n1 = strings.Replace(n1, "\r", "", -1)
		n1 = strings.Replace(n1, "\n", "", -1)
		if n1 == "全套" {
			tempList = append(tempList, yi)
		} else {
			//不是全套
			//判断有没有
			_, ok := tttMp[key]
			if ok {
				//如果有全套 就不添加了
			} else {
				tempList = append(tempList, yi)
			}
		}

	}

	//重新赋值
	allGongyi = tempList

	shafa := model.ShaFaImportLog{}
	err = shafa.Get("", sofa_code)

	if err != nil {
		log.Printf("hafa.Get err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}

	outInfo := Outt{
		SofaCode: sofa_code,
		SofaName: shafa.SfName,
	}

	//看是不是有重复的分位置
	FenWeiMap := make(map[string]int)
	ls := getshafaGuiGe(shafa.GG)

	for _, l := range ls {
		ll := strings.ReplaceAll(l, " ", "")
		nums, ok := FenWeiMap[ll]
		if ok {
			FenWeiMap[ll] = nums + 1
		} else {
			FenWeiMap[ll] = 1
		}

	}

	AllTotalPrice := 0.0
	var p1, p2, p3, p4, p5, p6, p7 float64
	var s1, s2, s3, s4, s5, s6, s7 float64
	yy1 := IInfo{}
	yy2 := IInfo{}
	yy3 := IInfo{}
	yy4 := IInfo{}
	yy5 := IInfo{}
	yy6 := IInfo{}
	yy7 := IInfo{}

	pp1 := make([]IIIIInfo, 0)
	pp2 := make([]IIIIInfo, 0)
	pp3 := make([]IIIIInfo, 0)
	pp4 := make([]IIIIInfo, 0)
	pp5 := make([]IIIIInfo, 0)
	pp6 := make([]IIIIInfo, 0)
	pp7 := make([]IIIIInfo, 0)

	// 要计算同名

	for _, yi := range allGongyi {

		//处理一下价格
		six, okkkk := FenWeiMap[yi.FenWeiName]

		six_fw := 1
		if okkkk {
			six_fw = six
		}

		if six_fw > 1 {
			yi.ShunHaoPrice = yi.ShunHaoPrice * float64(six_fw)
			yi.TotalPrice = covertPriceSix(yi.TotalPrice, six)
			yi.JiJiaNum = covertPriceSix(yi.JiJiaNum, six)
			yi.Nums = covertPriceSix(yi.Nums, six)

		}

		price, _ := strconv.ParseFloat(yi.TotalPrice, 64)

		if yi.Types == "裁工" {
			yy1.TypeName = "裁工"
			p1 += price
			p1, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p1), 64)
			yy1.TotalPrice = p1
			yy1.TotalSunhao += yi.ShunHaoPrice
			cccc1, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", yy1.TotalSunhao), 64)
			yy1.TotalSunhao = cccc1
			s1 = yy1.TotalSunhao

			pp1 = append(pp1, IIIIInfo{
				FenWeiName:  yi.FenWeiName,
				CLName:      GetCpName(yi.CpCode),
				CpCode:      yi.CpCode,
				Size:        yi.Size,
				Nums:        covertPrice(yi.Nums),
				Unit:        yi.JiJiaUnit,
				TotalPrice:  covertPrice(yi.TotalPrice),
				Price:       GetCpPrice(yi.CpCode),
				JiJiaNum:    covertPrice(yi.JiJiaNum),
				ShunHaoNums: fmt.Sprintf("%.4f", yi.ShunHaoPrice),
			})
		}
		if yi.Types == "车工" {
			yy2.TypeName = "车工"
			p2 += price
			p2, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p2), 64)
			yy2.TotalPrice = p2
			yy2.TotalSunhao += yi.ShunHaoPrice
			cccc, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", yy2.TotalSunhao), 64)
			yy2.TotalSunhao = cccc
			s2 = yy2.TotalSunhao
			pp2 = append(pp2, IIIIInfo{
				FenWeiName:  yi.FenWeiName,
				CpCode:      yi.CpCode,
				CLName:      GetCpName(yi.CpCode),
				Size:        yi.Size,
				Nums:        covertPrice(yi.Nums),
				Unit:        yi.JiJiaUnit,
				TotalPrice:  covertPrice(yi.TotalPrice),
				Price:       GetCpPrice(yi.CpCode),
				JiJiaNum:    covertPrice(yi.JiJiaNum),
				ShunHaoNums: fmt.Sprintf("%.4f", yi.ShunHaoPrice),
			})
		}
		if yi.Types == "海绵" {
			yy3.TypeName = "海绵"
			p3 += price
			p3, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p3), 64)
			yy3.TotalPrice = p3
			yy3.TotalSunhao += yi.ShunHaoPrice
			cccc, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", yy3.TotalSunhao), 64)
			yy3.TotalSunhao = cccc
			s3 = yy3.TotalSunhao
			pp3 = append(pp3, IIIIInfo{
				FenWeiName:  yi.FenWeiName,
				CpCode:      yi.CpCode,
				CLName:      GetCpName(yi.CpCode),
				Size:        yi.Size,
				Nums:        covertPrice(yi.Nums),
				Unit:        yi.JiJiaUnit,
				TotalPrice:  covertPrice(yi.TotalPrice),
				Price:       GetCpPrice(yi.CpCode),
				JiJiaNum:    covertPrice(yi.JiJiaNum),
				ShunHaoNums: fmt.Sprintf("%.4f", yi.ShunHaoPrice),
			})
		}
		if yi.Types == "扪工" {
			yy4.TypeName = "扪工"
			p4 += price
			p4, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p4), 64)
			yy4.TotalPrice = p4
			yy4.TotalSunhao += yi.ShunHaoPrice
			cccc, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", yy4.TotalSunhao), 64)
			yy4.TotalSunhao = cccc
			s4 = yy4.TotalSunhao
			pp4 = append(pp4, IIIIInfo{
				FenWeiName:  yi.FenWeiName,
				CLName:      GetCpName(yi.CpCode),
				CpCode:      yi.CpCode,
				Size:        yi.Size,
				Nums:        covertPrice(yi.Nums),
				Unit:        yi.JiJiaUnit,
				TotalPrice:  covertPrice(yi.TotalPrice),
				Price:       GetCpPrice(yi.CpCode),
				JiJiaNum:    covertPrice(yi.JiJiaNum),
				ShunHaoNums: fmt.Sprintf("%.4f", yi.ShunHaoPrice),
			})
		}
		if yi.Types == "木工" {
			yy5.TypeName = "木工"
			p5 += price
			p5, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p5), 64)
			yy5.TotalPrice = p5
			yy5.TotalSunhao += yi.ShunHaoPrice
			cccc, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", yy5.TotalSunhao), 64)
			yy5.TotalSunhao = cccc
			s5 = yy5.TotalSunhao
			pp5 = append(pp5, IIIIInfo{
				FenWeiName:  yi.FenWeiName,
				CLName:      GetCpName(yi.CpCode),
				CpCode:      yi.CpCode,
				Size:        yi.Size,
				Nums:        covertPrice(yi.Nums),
				Unit:        yi.JiJiaUnit,
				TotalPrice:  covertPrice(yi.TotalPrice),
				Price:       GetCpPrice(yi.CpCode),
				JiJiaNum:    covertPrice(yi.JiJiaNum),
				ShunHaoNums: fmt.Sprintf("%.4f", yi.ShunHaoPrice),
			})
		}
		if yi.Types == "人工" {
			yy6.TypeName = "人工"
			p6 += price
			p6, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p6), 64)
			yy6.TotalPrice = p6
			pp6 = append(pp6, IIIIInfo{
				FenWeiName: yi.FenWeiName,
				CLName:     yi.CLName,
				CpCode:     yi.CpCode,
				Size:       yi.Size,
				Nums:       covertPrice(yi.Nums),
				Unit:       yi.JiJiaUnit,
				TotalPrice: covertPrice(yi.TotalPrice),
				Price:      GetCpPrice(yi.CpCode),
				JiJiaNum:   covertPrice(yi.JiJiaNum),
			})
		}
		if yi.Types == "其他" {
			yy7.TypeName = "其他"
			p7 += price
			p7, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", p7), 64)
			yy7.TotalPrice = p7
			pp7 = append(pp7, IIIIInfo{
				FenWeiName: yi.FenWeiName,
				CLName:     yi.CLName,
				Size:       yi.Size,
				Nums:       covertPrice(yi.Nums),
				Unit:       yi.JiJiaUnit,
				TotalPrice: covertPrice(yi.TotalPrice),
				Price:      GetCpPrice(yi.CpCode),
				JiJiaNum:   covertPrice(yi.JiJiaNum),
			})
		}

	}
	AllTotalPrice = p1 + p2 + p3 + p4 + p5 + p6 + p7
	AllTotalPrice, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", AllTotalPrice), 64)

	if yy1.TypeName != "" {
		yy1.List = pp1
		outInfo.List = append(outInfo.List, yy1)
	}
	if yy2.TypeName != "" {
		yy2.List = pp2
		outInfo.List = append(outInfo.List, yy2)
	}
	if yy3.TypeName != "" {
		yy3.List = pp3
		outInfo.List = append(outInfo.List, yy3)
	}
	if yy4.TypeName != "" {
		yy4.List = pp4
		outInfo.List = append(outInfo.List, yy4)
	}
	if yy5.TypeName != "" {
		yy5.List = pp5
		outInfo.List = append(outInfo.List, yy5)
	}
	if yy6.TypeName != "" {
		yy6.List = pp6
		outInfo.List = append(outInfo.List, yy6)
	}
	if yy7.TypeName != "" {
		yy7.List = pp7
		outInfo.List = append(outInfo.List, yy7)
	}
	outInfo.TotalPrice = AllTotalPrice
	ppppppp1, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", s1+s2+s3+s4+s5+s6+s7), 64)
	outInfo.TotalShunHao = ppppppp1
	//处理下合并材料规则 todo

	for i, info := range outInfo.List {
		tempMap := make(map[int]IIIIInfo, 0)
		newi := make([]IIIIInfo, 0)
		for _, gongyiInfo := range info.List {
			if gongyiInfo.CpCode != "" {
				mergeDesc := model.GoodsMergeDesInfoDtoList{}
				errx1 := mergeDesc.GetListByCpCode(gongyiInfo.CpCode, nil)
				if errx1 == nil && len(mergeDesc) > 0 {

					lll, ok := tempMap[mergeDesc[0].MergeId]
					if ok {
						ppp1, _ := strconv.ParseFloat(lll.TotalPrice, 64)
						ppp2, _ := strconv.ParseFloat(gongyiInfo.TotalPrice, 64)
						lll.TotalPrice = fmt.Sprintf("%.4f", ppp1+ppp2)

						jjj1, _ := strconv.ParseFloat(lll.JiJiaNum, 64)
						jjj2, _ := strconv.ParseFloat(gongyiInfo.JiJiaNum, 64)
						lll.JiJiaNum = fmt.Sprintf("%.4f", jjj1+jjj2)

						// 设置损耗值
						xxx1, _ := strconv.ParseFloat(lll.ShunHaoNums, 64)
						xxx2, _ := strconv.ParseFloat(gongyiInfo.ShunHaoNums, 64)
						lll.ShunHaoNums = fmt.Sprintf("%.4f", xxx1+xxx2)

						tempMap[mergeDesc[0].MergeId] = lll

					} else {
						tempMap[mergeDesc[0].MergeId] = IIIIInfo{
							FenWeiName:  "",
							CpCode:      gongyiInfo.CpCode,
							CLName:      mergeDesc[0].CLName,
							Size:        "",
							Nums:        "",
							Unit:        mergeDesc[0].Unit,
							JiJiaNum:    gongyiInfo.JiJiaNum,
							Price:       mergeDesc[0].Price,
							TotalPrice:  gongyiInfo.TotalPrice,
							Descs:       "",
							ShunHaoNums: gongyiInfo.ShunHaoNums,
						}
					}

				} else {
					newi = append(newi, gongyiInfo)
				}
			} else {
				newi = append(newi, gongyiInfo)
			}
		} // 结束数据处理
		if len(tempMap) > 0 {
			for _, iiiiInfo := range tempMap {
				newi = append(newi, iiiiInfo)
			}
		}
		outInfo.List[i].List = newi

	}

	//相同材料合并
	for i, info := range outInfo.List {
		cpCodeMap := make(map[string]IIIIInfo, 0)
		newi := make([]IIIIInfo, 0)

		for _, gongyiInfo := range info.List {
			la, okkk := cpCodeMap[gongyiInfo.CpCode]
			if okkk {
				//如果存在 合并一下

				ppp1, _ := strconv.ParseFloat(la.TotalPrice, 64)
				ppp2, _ := strconv.ParseFloat(gongyiInfo.TotalPrice, 64)
				la.TotalPrice = fmt.Sprintf("%.4f", ppp1+ppp2)

				jjj1, _ := strconv.ParseFloat(la.JiJiaNum, 64)
				jjj2, _ := strconv.ParseFloat(gongyiInfo.JiJiaNum, 64)
				la.JiJiaNum = fmt.Sprintf("%.4f", jjj1+jjj2)

				// 设置损耗值
				xxx1, _ := strconv.ParseFloat(la.ShunHaoNums, 64)
				xxx2, _ := strconv.ParseFloat(gongyiInfo.ShunHaoNums, 64)
				la.ShunHaoNums = fmt.Sprintf("%.4f", xxx1+xxx2)
				cpCodeMap[gongyiInfo.CpCode] = la
			} else {
				cpCodeMap[gongyiInfo.CpCode] = gongyiInfo
			}

		}
		if len(cpCodeMap) > 0 {
			for _, iiiiInfo := range cpCodeMap {
				newi = append(newi, iiiiInfo)
			}
		}
		outInfo.List[i].List = newi

	}

	resp.Data = outInfo
	util.ReturnCompFunc(c, resp)
	return
}

func GetExportGoodsGroupByFenWeiDraf(c *gin.Context) {
	log.Printf("GetExportGoodsGroupByFenWei")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	types := params["types"]
	trans_id := params["trans_id"]
	log.Printf("shafa_id:%s types: %s ", shafaId, types)

	d, err := service.GetGoodsListGroupByNameDraf(shafaId, trans_id)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

func GetFinWeiGroupByNameDraf(c *gin.Context) {
	log.Printf("GetFinWeiGroupByName")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	types := params["types"]
	trans_id := params["trans_id"]

	log.Printf("shafa_id:%s types: %s ", shafaId, types)

	d, err := service.GetFenWeiListGroupByNameDraf(shafaId, types, trans_id)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

func GetDrafFinWei(c *gin.Context) {
	log.Printf("GetFenWei")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	types := params["types"]
	transId := params["trans_id"]

	log.Printf("shafa_id:%s types: %s ;transId:%s ", shafaId, types, transId)

	d, err := GetDrafFenWeiList(shafaId, types, transId)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

func GetDrafFenWeiList(shafaId, types, transId string) (list []model.GongYiDraf, err error) {

	info := model.GongYiDrafList{}
	err = info.GetListPageWithTrans(shafaId, types, transId, nil)
	//有没有呢
	if err != nil {
		log.Printf("GetFenWeiList err :%v", err)
	}
	list = info
	return

}

func GetDrafShaFaImportList(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	page_ := params["page"]
	size_ := params["size"]
	page, _ := strconv.Atoi(page_)
	size, _ := strconv.Atoi(size_)
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	name := params["goods_name"]
	code := params["goods_code"]
	isSums := params["is_sums"]
	log.Printf("http GetGoodsList : goods_name:%s goods_code: %s ", name, code)

	d, total, err := service.GetDrafShaFaImportList(page, size, name, code, isSums)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()

	}

	dto := GetShaFaImportListDto{
		page,
		size,
		total,
		d,
	}
	resp.Data = dto

	util.ReturnCompFunc(c, resp)
	return
}

func GetDrafShaFaImportById(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	id := params["id"]

	log.Printf("http GetDrafShaFaImportById : shafa_id:%s  ", id)

	d, err := service.GetDrafShaFaImportById(id)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()

	}

	dto := GetDrafShaFaImportByIdDto{

		Data: d,
	}
	resp.Data = dto

	util.ReturnCompFunc(c, resp)
	return
}

type GetDrafShaFaImportByIdDto struct {
	Data model.ShaFaDrafImportLog `json:"data"`
}

type PostDrafFengWeiInDTO struct {
	Details []model.GongYiDraf `json:"details"`
	// 物料的类型 车工 裁工
	Types   string `json:"types"`
	ShafaId string `json:"shafa_id"`
}

func PostDrafFengWei(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	data_ := params["data"]

	user := params["user"]

	dto := PostDrafFengWeiInDTO{}

	err := json.Unmarshal([]byte(data_), &dto)
	if err != nil {
		log.Printf("json.Unmarshal  err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
		return
	}

	log.Printf("PostFengWei req :%+v", dto)
	//保存物料成本到每个具体的沙发
	//保存都是直接覆盖的

	// 先获取 trans_id

	drafLog := model.ShaFaDrafImportLog{}
	errad := drafLog.GetByType(nil, dto.ShafaId)
	if errad != nil || drafLog.Id == 0 {
		log.Printf("jdrafLog.GetByType  err :%v", err)
		resp.Status = 201
		resp.Desc = "获取导入表失败  获取版本失败"
		return
	}

	//先清空
	deleteInfo := model.GongYiDraf{}
	deleteInfo.Types = dto.Types
	deleteInfo.ShafaId = dto.ShafaId
	deleteInfo.TransId = drafLog.TransId
	err = deleteInfo.DeleteWithTrans(nil)
	if err != nil {
		log.Printf("deleteInfo.Delete err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
		return
	}
	// 要将 全套用量 替换掉其他掉

	//循环插入
	for _, detail := range dto.Details {
		insertInfo := model.GongYiDraf{}
		insertInfo.CreateTime = time.Now()
		insertInfo.Types = dto.Types
		insertInfo.ShafaId = dto.ShafaId
		insertInfo.FenWeiName = strings.TrimSpace(detail.FenWeiName)
		insertInfo.CLName = detail.CLName
		insertInfo.Size = detail.Size
		insertInfo.Nums = detail.Nums
		insertInfo.Unit = detail.Unit
		insertInfo.Descs = detail.Descs
		insertInfo.TotalPrice = detail.TotalPrice
		insertInfo.CpCode = detail.CpCode
		insertInfo.GongYiName = strings.TrimSpace(detail.GongYiName)
		insertInfo.JiJiaNum = detail.JiJiaNum
		insertInfo.GoodsPoint = detail.GoodsPoint
		insertInfo.JiJiaUnit = detail.Unit
		insertInfo.OwnerSize = detail.OwnerSize
		insertInfo.ImportUser = user
		insertInfo.TransId = drafLog.TransId

		//获取物料信息
		goods := model.Goods{}
		err1 := goods.Get("", detail.CpCode)
		if err1 != nil {
			log.Printf("goods.Get err :%v", err)
			resp.Status = 201
			resp.Desc = err.Error()
			return
		}
		if goods.FuZhuXiShu != 1 {
			//用辅助计量单位
			nn := GetUnitById(goods.FuZhuUnit)
			insertInfo.JiJiaUnit = nn
		}
		//
		if goods.ShunHao != "" {
			price_f, _ := strconv.ParseFloat(detail.TotalPrice, 64)
			shunhao_f, _ := strconv.ParseFloat(goods.ShunHao, 64)

			xx := shunhao_f / 100
			insertInfo.ShunHaoPrice = (xx / (xx + 1)) * price_f
		}
		//判断单位
		mergeDesc := model.GoodsMergeDesInfoDtoList{}
		errx1 := mergeDesc.GetListByCpCode(detail.CpCode, nil)
		if errx1 == nil && len(mergeDesc) > 0 {
			insertInfo.JiJiaUnit = mergeDesc[0].Unit
		}

		err = insertInfo.Create(nil)
		if err != nil {
			log.Printf("insertInfo.Create err :%v", err)
			resp.Status = 201
			resp.Desc = err.Error()
			return
		}
	}
	//修改事务表

	transe := model.Trans{}
	transe.GetByShafaIdAndTrans(nil, dto.ShafaId, drafLog.TransId)
	transe.IsCheck = 0
	transe.CheckUser = user
	transe.Update(nil)

	//修改 沙发表
	/*	shafa := model.ShaFaImportLog{}
		errme := shafa.GetByType(nil, dto.ShafaId)
		if errme == nil {
			size := len(dto.Details)
			if size == 0 {
				shafa.IsSums = "否"
			} else {
				shafa.IsSums = "是"
			}
			shafa.Update(nil)
		}*/

	resp.Data = dto
	util.ReturnCompFunc(c, resp)
	return
}

func ConvertPostInfo(user string, maxCell int, goodsMap map[string]model.Goods, transID string, gongyiTypes string, shafaId string, rowData [][]string) {

	//组装数据
	dto := PostGoodsInfoDrafInDTO{}
	dto.ShafaId = shafaId
	dto.Types = gongyiTypes

	dto.Details = make([]model.GongYiDraf, 0)

	unitMap := make(map[string]model.AllGoodsDesc, 0)

	//要知道有几个分位
	fenWeiMap := make(map[string]int, 0)
	for iiiiii := 9; iiiiii < maxCell; iiiiii += 3 {
		fenWeiMap[strings.TrimSpace(rowData[2][iiiiii])] = iiiiii
	}

	for i := 3; i < len(rowData); i++ {

		if len(rowData[i]) <= 10 {
			continue
		}
		//每一行表示一个成本
		cp_code := rowData[i][2]
		goods, ok := goodsMap[cp_code]
		if !ok {
			continue
		}
		//获取 desc

		//GeGoodsById
		desc, ok := unitMap[cp_code]
		if !ok {
			ll, errr := service.GeGoodsById(cp_code)
			if errr != nil {
				log.Printf("获取单位失败了 . err :%v", errr)
				continue
			} else {
				unitMap[cp_code] = ll
				desc = ll
			}
		}
		//这一行有多长
		leng := len(rowData[i])
		for fenwei_name, fenwei_index := range fenWeiMap {
			//如果数量为空 直接返回
			if leng < fenwei_index {
				continue
			}
			numsdd := strings.TrimSpace(rowData[i][fenwei_index])
			if numsdd == "" || numsdd == "0" {
				continue
			}

			draf := model.GongYiDraf{}
			draf.CLName = goods.CpName
			draf.CpCode = goods.CpCode
			draf.GoodsPoint = rowData[i][5]
			draf.Descs = rowData[i][6]
			draf.Unit = desc.CpMainUnit
			draf.JiJiaUnit = desc.CpMainUnit
			draf.Size = rowData[i][fenwei_index-2]
			draf.Nums = rowData[i][fenwei_index]
			draf.ImportUser = user
			draf.GongYiName = rowData[i][4]
			draf.FenWeiName = fenwei_name
			draf.ShafaId = shafaId
			draf.Types = gongyiTypes
			// 获取价格和计价数量

			//获取 规格  注意 规格 看是否要以 用户熟肉的为准。
			mainSize := goods.MainSize
			if strings.TrimSpace(rowData[i][fenwei_index-1]) != "" {
				mainSize = strings.TrimSpace(rowData[i][fenwei_index-1])
				draf.OwnerSize = mainSize
			}
			sumsDraf("0.5", "10210077", "2", "40.2000", "1.000", "", "1.000", "", "1..000")
			sum_price := sumsDraf(draf.Nums, draf.CpCode, goods.ShunHao, fmt.Sprintf("%f", goods.Price), fmt.Sprintf("%f", goods.ChangeP), draf.Size, fmt.Sprintf("%f", goods.MainXiShu), mainSize, fmt.Sprintf("%f", goods.FuZhuXiShu))
			draf.TotalPrice = fmt.Sprintf("%f", sum_price.TotalPrice)
			draf.JiJiaNum = fmt.Sprintf("%f", sum_price.JiJiaNums)
			draf.TransId = transID
			dto.Details = append(dto.Details, draf)
		}

	}

	postFenweiFraf(transID, user, dto)

}

type PostGoodsInfoDrafInDTO struct {
	Details []model.GongYiDraf `json:"details"`
	// 物料的类型 车工 裁工
	Types   string `json:"types"`
	ShafaId string `json:"shafa_id"`
}

//一次上传一种工艺。
func postFenweiFraf(transeId, user string, dto PostGoodsInfoDrafInDTO) error {

	log.Printf("导入 插入数据开始 postFenweiFraf req :%+v", dto)
	//先判断有没有这个沙发。
	//修改 沙发表
	/*	shafa := model.ShaFaDrafImportLog{}
		errme := shafa.GetByType(nil, dto.ShafaId)

		if errme == gorm.ErrRecordNotFound || shafa.Id == 0 {
			//需要创建一个。
			copy1 := model.ShaFaImportLog{}
			errme1 := copy1.GetByType(nil, dto.ShafaId)
			if errme1 != nil {
				log.Printf("导入 失败 没有找到沙发 shafa_id:%s", dto.ShafaId)
				return errme1
			}
			shafa.GG = copy1.GG
			shafa.ImportUser = user
			shafa.CreateTime = time.Now()
			shafa.UpdateTime = time.Now()
			shafa.SfName = copy1.SfName
			shafa.SfCode = copy1.SfCode
			shafa.SDesc = copy1.SDesc
			shafa.TransId = transeId
			errr := shafa.Create(nil)
			if errr != nil {
				log.Printf("导入 失败 创建失败 shafa_id:%s；err:%v", dto.ShafaId, errr)
				return errr
			}

		}*/

	//先执行删除操作；

	//先清空
	deleteInfo := model.GongYiDraf{}
	deleteInfo.Types = dto.Types
	deleteInfo.ShafaId = dto.ShafaId
	deleteInfo.TransId = transeId
	err1 := deleteInfo.DeleteWithTrans(nil)
	if err1 != nil {
		log.Printf("deleteInfo.Delete err :%v", err1)

		return errors.New(fmt.Sprintf("删除原有成本草稿失败 sf_code：%s  transId:%s", dto.ShafaId, transeId))
	}

	// 要将 全套用量 替换掉其他掉
	//循环插入
	for _, detail := range dto.Details {
		insertInfo := model.GongYiDraf{}
		insertInfo.CreateTime = time.Now()
		insertInfo.Types = dto.Types
		insertInfo.ShafaId = dto.ShafaId
		insertInfo.FenWeiName = strings.TrimSpace(detail.FenWeiName)
		insertInfo.CLName = detail.CLName
		insertInfo.Size = detail.Size
		insertInfo.Nums = detail.Nums
		insertInfo.Unit = detail.Unit
		insertInfo.Descs = detail.Descs
		insertInfo.TotalPrice = detail.TotalPrice
		insertInfo.CpCode = detail.CpCode
		insertInfo.GongYiName = strings.TrimSpace(detail.GongYiName)
		insertInfo.JiJiaNum = detail.JiJiaNum
		insertInfo.GoodsPoint = detail.GoodsPoint
		insertInfo.JiJiaUnit = detail.Unit
		insertInfo.ImportUser = user
		insertInfo.OwnerSize = detail.OwnerSize
		insertInfo.TransId = transeId
		//获取物料信息
		goods := model.Goods{}
		err1 := goods.Get("", detail.CpCode)
		if err1 != nil {
			log.Printf("导入失败 ， goods.Get err :%v", err1)
			return err1
		}
		if goods.FuZhuXiShu != 1 {
			//用辅助计量单位
			nn := GetUnitById(goods.FuZhuUnit)
			insertInfo.JiJiaUnit = nn
		}
		//
		if goods.ShunHao != "" {
			price_f, _ := strconv.ParseFloat(detail.TotalPrice, 64)
			shunhao_f, _ := strconv.ParseFloat(goods.ShunHao, 64)

			xx := shunhao_f / 100
			insertInfo.ShunHaoPrice = (xx / (xx + 1)) * price_f
		}
		//判断单位
		mergeDesc := model.GoodsMergeDesInfoDtoList{}
		errx1 := mergeDesc.GetListByCpCode(detail.CpCode, nil)
		if errx1 == nil && len(mergeDesc) > 0 {
			insertInfo.JiJiaUnit = mergeDesc[0].Unit
		}

		err := insertInfo.Create(nil)
		if err != nil {
			log.Printf("导入失败 ， insertInfo.Create err :%v", err)

			return err
		}
	}
	return nil
}

func sumsDraf(num, cpCode, shunhao, price, huansuan, size, main_xishu, main_size, fuzhu_xishu string) PriceInfo {

	cdl := fmt.Sprintf("mus: %s;cpCOde:%s;shunhao:%s;price %s; huansuan:%s;size:%s manxishu:%s;main_size:%s;fuzhuxushu:%s",
		num, cpCode, shunhao, price, huansuan, size, main_xishu, main_size, fuzhu_xishu)
	log.Printf("Sums:param  :%s", cdl)

	out := PriceInfo{}
	num_i, e1 := strconv.ParseFloat(num, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}
	//cpCode := params["cp_code"]

	//shunhao := params["shun_hao"]
	shunhao_i, e1 := strconv.ParseFloat(shunhao, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}

	//price := params["price"]
	price_i, e1 := strconv.ParseFloat(price, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}
	//huansuan := params["huan_suan"]
	huansuan_i, e1 := strconv.ParseFloat(huansuan, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}
	//输入的尺寸
	//size := params["size"]

	//main_xishu := params["main_xishu"]
	main_xishu_i, e1 := strconv.ParseFloat(main_xishu, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}

	//main_size := params["main_size"]

	//fuzhu_xishu := params["fuzhu_xishu"]
	fuzhu_xishu_i, e1 := strconv.ParseFloat(fuzhu_xishu, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}

	totalPrice := 0.0
	totalPrice_t := 0.0
	//计价数量
	jijiaNums := 0.0

	//根据产品ID 获取产品

	goods := model.Goods{}
	info, err := goods.GetGoodsById(cpCode, nil)
	if err != nil {
		return out
	}
	//判断长宽高
	l1 := ""
	if main_size == "" {
		l1 = info.MainSize
	} else {
		l1 = main_size
	}

	a1 := 1.0
	a2 := 1.0
	a3 := 1.0
	a4 := 1.0
	a5 := 1.0
	a6 := 1.0

	if l1 != "" {
		if strings.Contains(l1, "*") {
			s_list := strings.Split(l1, "*")
			if len(s_list) >= 1 {
				a1_s := s_list[0]
				a1, err = strconv.ParseFloat(a1_s, 64)
				if err != nil {
					goto out
				}
			}
			if len(s_list) >= 2 {
				a2_s := s_list[1]
				a2, err = strconv.ParseFloat(a2_s, 64)
				if err != nil {
					goto out
				}
			}
			if len(s_list) >= 3 {
				a3_s := s_list[2]
				a3, err = strconv.ParseFloat(a3_s, 64)
				if err != nil {
					goto out
				}
			}

		} else {
			a1, err = strconv.ParseFloat(l1, 64)
			if err != nil {

				goto out
			}
		}
	}
	if size != "" {
		if strings.Contains(size, "*") {
			s_list := strings.Split(size, "*")
			if len(s_list) >= 1 {
				a4_s := s_list[0]
				a4, err = strconv.ParseFloat(a4_s, 64)
				if err != nil {
					goto out
				}
			}
			if len(s_list) >= 2 {
				a5_s := s_list[1]
				a5, err = strconv.ParseFloat(a5_s, 64)
				if err != nil {
					goto out
				}
			}

		} else {
			a4, err = strconv.ParseFloat(size, 64)
			if err != nil {
				goto out
			}
		}
	}

	if huansuan_i == 0 {
		huansuan_i = 1
	}
	if fuzhu_xishu_i == 0 {
		fuzhu_xishu_i = 1
	}
	if main_xishu_i == 0 {
		main_xishu_i = 1
	}
	main_xishu_i = 1
	//	算了损耗的 数量
	//jijiaNums = (((shunhao_i/100 + 1) * (a1 * a2 * a3 * a4) * num_i) / huansuan_i) * float64(main_xishu_i) / float64(fuzhu_xishu_i)
	//没有算损耗的数量
	jijiaNums = ((a1 * a2 * a3 * a4 * a5 * a6) * num_i / huansuan_i) / float64(main_xishu_i) / float64(fuzhu_xishu_i)
	//
	//没有合并规则下的价格
	totalPrice_t = (((shunhao_i/100 + 1) * (a1 * a2 * a3 * a4 * a5 * a6) * num_i) / huansuan_i) / float64(main_xishu_i) / float64(fuzhu_xishu_i)
	//
	//有合并规则下的价格【【

out:

	////////是否有材料合并规则
	mergeDesc := model.GoodsMergeDesInfoDtoList{}
	errx1 := mergeDesc.GetListByCpCode(cpCode, nil)
	if errx1 == nil && len(mergeDesc) > 0 && mergeDesc[0].Price > 0 {
		//要判断一下 合并的的单位是否和主计量单位一致
		goodsMerge := model.Goods{}
		err3 := goodsMerge.Get("", cpCode)
		if err3 != nil {
			log.Printf("GetGoodsChangeById err :%v", err3)

		}
		unit := model.UnitDesc{}
		intv1, err1 := strconv.Atoi(goodsMerge.CpMainUnit)

		if err1 != nil {
			log.Printf("GetGoodsChangeById err :%v", err1)

		}
		unit.GetById(nil, intv1)
		unit_str := strings.Replace(unit.Name, "\n", "", -1)
		unit_str = strings.Replace(unit_str, " ", "", -1)
		unit_str = strings.Replace(unit_str, "\r", "", -1)
		unit_merge_str := strings.Replace(mergeDesc[0].Unit, "\n", "", -1)
		unit_merge_str = strings.Replace(unit_merge_str, " ", "", -1)
		unit_merge_str = strings.Replace(unit_merge_str, "\r", "", -1)

		//单位一样 直接取值
		if unit_str == unit_merge_str {
			newPric := mergeDesc[0].Price
			totalPrice = newPric * totalPrice_t
		} else {
			xishu := goodsMerge.MainXiShu
			if xishu == 0 {
				xishu = 1
			}
			newPric := mergeDesc[0].Price / xishu
			jijiaNums = jijiaNums / xishu
			//newPric := mergeDesc[0].Price
			totalPrice = newPric * totalPrice_t
		}

	} else {
		//查询是否有固定价格
		if info.GuDingPrice != 0 {
			totalPrice = info.GuDingPrice
		} else {
			totalPrice = totalPrice_t * price_i
		}
	}

	// 查找下是否有换算比例
	change := model.GoodsChangeDesInfoDtoList{}
	err24 := change.GetListByCpCode(cpCode, nil)

	if err24 == nil && len(change) > 0 {

		for _, dto := range change {
			if dto.ChangeType == "换算" {
				if dto.Types == "/" {
					totalPrice = totalPrice / dto.ValuesL
					jijiaNums = jijiaNums / dto.ValuesL
				} else {
					totalPrice = totalPrice * dto.ValuesL
					jijiaNums = jijiaNums * dto.ValuesL
				}
				break
			}

		}

	}
	log.Printf("sums :  a1:%f; a2:%f", totalPrice, jijiaNums)
	riceInfo := PriceInfo{
		TotalPrice: totalPrice,
		JiJiaNums:  jijiaNums,
	}

	return riceInfo

}

func GetShafaTransList(c *gin.Context) {
	log.Printf("GetFinWeiGroupByName")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["id"]

	log.Printf("shafa_id:%s", shafaId)

	d := service.GetTransList(shafaId)
	/*	if err != nil {
			log.Printf("GetGoodsList err :%v", err)

			resp.Status = 201
			resp.Desc = err.Error()
		}
	*/
	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

//选择版本
func CheckTrans(c *gin.Context) {
	log.Printf("CheckTrans")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["sf_code"]

	trans_code := params["trans_id"]

	log.Printf("CheckTrans shafa_id:%s trans_code:%s ", shafaId, trans_code)

	imports := model.ShaFaDrafImportLog{}
	err := imports.GetOne(shafaId)

	if err != nil || imports.Id == 0 {
		resp.Data = "切换版本失败 没有找到该沙发"
		resp.Status = 201
		util.ReturnCompFunc(c, resp)
	}

	imports.TransId = trans_code

	err = imports.Update(nil)
	if err != nil {
		resp.Data = "切换版本失败 没有找到该版本"
		resp.Status = 201
		util.ReturnCompFunc(c, resp)
	}

	resp.Data = "成功"

	util.ReturnCompFunc(c, resp)
	return

}

//选择版本
func OnlineTrans(c *gin.Context) {
	log.Printf("OnlineTrans")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["sf_code"]
	user := params["user"]
	log.Printf("CheckTrans shafa_id:%s ", shafaId)

	imports := model.ShaFaDrafImportLog{}
	err := imports.GetOne(shafaId)

	if err != nil || imports.Id == 0 {
		resp.Data = "通过失败 没有找到该沙发"
		resp.Status = 201
		util.ReturnCompFunc(c, resp)
	}

	transId := imports.TransId

	gongyiDrafS := model.GongYiDrafList{}
	err = gongyiDrafS.GetBySoFaCodeDraf(nil, shafaId, transId)

	if err != nil || len(gongyiDrafS) == 0 {
		resp.Data = "通过失败 该版本没有成本"
		resp.Status = 201
		util.ReturnCompFunc(c, resp)
	}

	// 将原来的主表成本全部删除。
	//先清空
	deleteInfo := model.GongYi{}
	deleteInfo.ShafaId = shafaId
	err = deleteInfo.Delete(nil)
	if err != nil {
		log.Printf("deleteInfo.Delete err :%v", err)
		resp.Status = 201
		resp.Desc = "通过失败 删除原来成本失败"
		util.ReturnCompFunc(c, resp)
		return
	}
	for _, yi := range gongyiDrafS {
		gy := model.GongYi{}
		gy.Id = 0
		gy.ShafaId = shafaId
		gy.CreateTime = time.Now()
		gy.Types = yi.Types
		gy.JiJiaNum = yi.JiJiaNum
		gy.Nums = yi.Nums
		gy.GongYiName = yi.GongYiName
		gy.OwnerSize = yi.OwnerSize
		gy.CpCode = yi.CpCode
		gy.CLName = yi.CLName
		gy.FenWeiName = yi.FenWeiName
		gy.ShunHaoPrice = yi.ShunHaoPrice
		gy.GoodsPoint = yi.GoodsPoint
		gy.JiJiaUnit = yi.JiJiaUnit
		gy.TotalPrice = yi.TotalPrice
		gy.Descs = yi.Descs
		gy.Unit = yi.Unit
		gy.Size = yi.Size
		gy.Types = yi.Types

		gy.Create(nil)
		if err != nil {
			log.Printf("insertInfo.Create err :%v", err)
			resp.Status = 201
			resp.Desc = "创建成本失败"
			util.ReturnCompFunc(c, resp)
			return
		}
	}
	//更新事务表 记录为成功。

	transInfo := model.Trans{}
	err = transInfo.GetByShafaIdAndTrans(nil, shafaId, transId)
	if err != nil {
		log.Printf("GetByShafaIdAndTransr :%v", err)
		resp.Status = 201
		resp.Desc = "修改事务失败"
		util.ReturnCompFunc(c, resp)
		return
	}
	transInfo.IsSubmit = 1
	transInfo.OnlineUser = user
	transInfo.OnlineTime = time.Now()
	err = transInfo.Update(nil)
	if err != nil {
		log.Printf("transInfo.Update(nil) :%v", err)
		resp.Status = 201
		resp.Desc = "修改事务失败"
		util.ReturnCompFunc(c, resp)
		return
	}

	base_shafa := model.ShaFaImportLog{}
	errme := base_shafa.GetByType(nil, shafaId)
	if errme != nil || base_shafa.Id == 0 {
		log.Printf(" 找不到沙发ID ：%s", shafaId)
		resp.Status = 201
		resp.Desc = "找不到该沙发"
		util.ReturnCompFunc(c, resp)
		return
	}

	base_shafa.IsSums = "是"
	base_shafa.Update(nil)
	resp.Data = "成功"
	util.ReturnCompFunc(c, resp)
	return

}

func CheckGongyi(c *gin.Context) {
	log.Printf("CheckGongyi")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	user := params["user"]

	log.Printf("CheckGongyi  shafa_id:%s user: %s ;", shafaId, user)

	shafaDrafInfo := model.ShaFaDrafImportLog{}
	errme := shafaDrafInfo.GetByType(nil, shafaId)
	if errme != nil && shafaDrafInfo.Id == 0 {
		log.Printf("oldTransList.GetByShafaId err :%v", errme)

		resp.Status = 201
		resp.Desc = "未找到沙发"
		util.ReturnCompFunc(c, resp)
		return
	}

	transe := model.Trans{}
	transe.GetByShafaIdAndTrans(nil, shafaId, shafaDrafInfo.TransId)
	transe.IsCheck = 1
	transe.CheckUser = user
	er := transe.Update(nil)
	if er != nil {
		log.Printf("GetGoodsList err :%v", er)

		resp.Status = 201
		resp.Desc = "更改失败"
		util.ReturnCompFunc(c, resp)
		return
	}
	resp.Data = ""

	util.ReturnCompFunc(c, resp)
	return

}
