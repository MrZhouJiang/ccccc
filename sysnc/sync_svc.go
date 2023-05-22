package sysnc

import (
	model "ccccc/data/model/goods"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

var TotalLimit = 0
var fenweiLimit = 0

func Task() {
	/*s := gocron.NewScheduler()

	s.Every(1).Second().Do(Test)

	<-s.Start()*/

	/*	s1 := gocron.NewScheduler()

		s1.Every(10).Second().Do(Test)

		<-s1.Start()*/

}
func Test() {
	fmt.Println(time.Now())
	base := SqlCpBase{}
	err := base.GetListOne(nil)
	if err != nil {
		log.Printf("Test err :%v", err)
	} else {
		log.Printf("Test success :%v", base)
	}
}

func StartSyncCp() {
	log.Printf("sync start time:%v", time.Now())

	start := TotalLimit
	size := 100

	//

	for {
		list := CpList{}

		err := list.GetListPage(start, size, nil)
		if err != nil {
			log.Printf("sync GetListPage err :%v \n", err)
		} else {
			log.Printf("sync GetListPage success : \n")

		}

		log.Printf("sync GetListPage len :%d  \n", len(list))

		for _, base := range list {
			log.Printf("sync GetListPage cp_code :%s \n", base.CPBM)
			log.Printf("sync GetListPage cp_code :%v \n", base)

			timeLayout := "2006-01-02 15:04:05"

			//查询是否存在
			goods := model.Goods{}
			if base.CPBM != "" {
				err = goods.Get("", base.CPBM)
				if err != nil && err != gorm.ErrRecordNotFound {
					log.Printf("sync goods.Get err :%v", err)
					continue
				}
				//不存在或者要更新
				if goods.Id == 0 {
					if base.SFTY != 1 {
						//不存在 直接插入
						//基础物料表
						goods.CpCode = base.CPBM
						goods.CpName = base.CPMC
						goods.CpDesc = base.CPJC
						goods.CpTypeCode = fmt.Sprintf("%d", base.CPFLID)
						goods.CpGuiGe = base.GG
						goods.CpMainUnit = fmt.Sprintf("%d", base.JLDWID_Z)
						goods.FuZhuUnit = fmt.Sprintf("%d", base.JLDWID_F)
						goods.MainXiShu = int(base.XS_ZJL)
						goods.FuZhuXiShu = int(base.XS_FJL)
						goods.Price = base.CBJ
						goods.ChangeP, _ = strconv.ParseFloat(base.ZHL, 64)
						goods.TyName = base.TYMC
						goods.DomH = ""
						goods.DomId = 0
						goods.ShunHao = "0"
						if strings.Contains(base.CJSJ, "T") {
							timeLayout = "2006-01-02T15:04:05Z"
						}
						goods.CreateTime, _ = time.Parse(timeLayout, base.CJSJ)
						//goods.CreateTime = time.Now()
						goods.LoadTime = time.Now()
						goods.Create(nil)
						// 沙发表

						if goods.CpTypeCode == "1019" || goods.CpTypeCode == "1020" || goods.CpTypeCode == "1018" || goods.CpTypeCode == "1029" {
							shafaimport := model.ShaFaImportLog{}
							shafaimport.SfName = base.CPMC
							shafaimport.SfCode = base.CPBM
							shafaimport.SDesc = base.CPJC
							shafaimport.GG = base.GG
							shafaimport.IsSums = "否"
							if strings.Contains(base.CJSJ, "T") {
								timeLayout = "2006-01-02T15:04:05Z"
							}
							shafaimport.CreateTime, _ = time.Parse(timeLayout, base.CJSJ)
							//shafaimport.CreateTime = time.Now()

							shafaimport.Create(nil)

						}

					}

				} else {
					//判断是够需要更新
					if base.SFTY == 1 {
						//停用了 直接删除
						goods.Delete(nil)
					} else {

						if strings.Contains(base.XGSJ, "T") {
							timeLayout = "2006-01-02T15:04:05Z"
						}
						cl1, _ := time.Parse(timeLayout, base.XGSJ)
						if cl1.After(goods.LoadTime) {
							goods.CpName = base.CPMC
							goods.Update(nil)
							shafaimport := model.ShaFaImportLog{}
							err = shafaimport.Get("", base.CPBM)
							shafaimport.SfName = base.CPMC
							shafaimport.Update(nil)
						}

					}

				}

			}

		}

		if len(list) < size {
			break
		} else {
			start += size
			TotalLimit = start
		}

	}
	log.Printf("sync end time:%v", time.Now())

}

func StartSyncFenWei() {
	log.Printf("sync start StartSyncFenWei time:%v", time.Now())

	start := fenweiLimit
	size := 100

	for {
		list := CpFenWeiList{}

		err := list.GetListPage(start, size, nil)
		if err != nil {
			log.Printf("sync GetListPage err :%v", err)
		}

		log.Printf("sync GetListPage len :%d", len(list))

		for _, base := range list {
			//查询是否存在
			goods := model.FenWei{}
			if base.FWBM != "" {
				goods.FWBM = base.FWBM
				err = goods.GetByCode(nil)
				if err != nil && err != gorm.ErrRecordNotFound {
					log.Printf("sync goods.Get err :%v", err)
					continue
				}
				goods.FWMC = base.FWMC
				goods.FWJC = base.FWJC
				goods.FWXS = base.FWXS
				goods.Create(nil)
			}

		}

		if len(list) < size {
			break
		} else {
			start += size
			fenweiLimit = start
		}

	}
	log.Printf("sync fenweiLimit end time:%v", time.Now())

}
