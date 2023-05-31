package sysnc

import (
	model "ccccc/data/model/goods"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

var TotalLimit = 0
var fenweiLimit = 0

func Task_base() {
	cron := cronexpr.MustParse("* * */23 * * * *") //用cron库生成一个cronexpr.Expression对象
	next := cron.Next(time.Now())                  //计算下次触发时间的时间对象
	for {
		now := time.Now()                        //每次循环计算获取当前时间
		if next.Before(now) || next.Equal(now) { //下次触发时间与当前时间进行对比，等于或者时间已到 则进行任务触发
			StartSyncCp()
			next = cron.Next(now) //重新计算下次任务时间的时间对象
		}
		select {
		case <-time.NewTicker(1 * time.Hour).C: //每秒扫描一遍 循环频率设定
		}
	}
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
func StartSyncCp_base() {
	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
	fmt.Println(time.Now())
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
					log.Printf("sync GetListPage cp_info :%v \n", base)
					//不存在 直接插入
					//基础物料表
					goods.CpCode = base.CPBM
					goods.CpName = base.CPMC
					goods.CpDesc = base.CPJC
					goods.CpTypeCode = fmt.Sprintf("%d", base.CPFLID)
					goods.CpGuiGe = base.GG
					goods.CpMainUnit = fmt.Sprintf("%d", base.JLDWID_Z)
					goods.FuZhuUnit = fmt.Sprintf("%d", base.JLDWID_F)
					goods.MainXiShu = base.XS_ZJL
					goods.FuZhuXiShu = base.XS_FJL
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
					goods.SFTY = base.SFTY
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
				} else {
					if base.SFTY == 1 && goods.SFTY == 0 {
						goods.SFTY = 1
						goods.Update(nil)
						log.Printf("sync update cp_info :%v \n", base)
					}
					if base.SFTY == 0 && goods.SFTY == 1 {
						goods.SFTY = 0
						log.Printf("sync update cp_info :%v \n", base)
						goods.Update(nil)
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
