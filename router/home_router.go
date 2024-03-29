package router

import (
	"ccccc/controller"
	"github.com/gin-gonic/gin"
)

func HomeRouter(engine *gin.Engine) {
	eg_v2 := engine.Group("/v2")
	{
		//用户相关接口
		eg_v2.GET("/get1", controller.GetById)
		eg_v2.POST("/post1", controller.Post1)
		eg_v2.POST("/user_login", controller.UserLogin)
		eg_v2.POST("/create_user", controller.CreateUser)
		eg_v2.GET("/get_user_list", controller.GeUserList)
		eg_v2.GET("/delete_user", controller.DeleteUser)

		//物料
		eg_v2.GET("/get_goods_list", controller.GetGoodsList)
		eg_v2.GET("/get_goods_by_id", controller.GetGoodsById)
		eg_v2.POST("/update_goods", controller.UpdatGoods)
		eg_v2.GET("/get_all_goods_desc", controller.GetAllGoodsListDesc)

		// 沙发
		eg_v2.GET("/get_shafa_import_list", controller.GetShaFaImportList)
		eg_v2.POST("/post_feng_wei", controller.PostFengWei)
		//拷贝一个沙发配置
		eg_v2.POST("/cpoy_shafa", controller.CopyShaFa)
		//拷贝一个沙发配置
		eg_v2.POST("/reload_shafa", controller.ReloadShaFa)

		eg_v2.GET("/get_feng_wei", controller.GetFinWei)
		eg_v2.GET("/get_all_fen_wei_list", controller.GetAllFenWeiDesc)
		//根据工艺名称分组获取 工艺成本项目 用户导出成本清单给财务
		eg_v2.GET("/get_feng_wei_group_by_name", controller.GetFinWeiGroupByName)
		//获取材料表
		eg_v2.GET("/get_export_goods_group_by_fen_wei", controller.GetExportGoodsGroupByFenWei)

		//成本计价
		eg_v2.GET("/sums", controller.Sums)

		//
		eg_v2.GET("/get_uint", controller.GetUnit)
		eg_v2.GET("/get_goods_type", controller.GetAllGoodsTypeListDesc)

		//换算 损耗
		eg_v2.GET("/get_goods_changeList", controller.GetGoodsChangeList)
		eg_v2.GET("/get_goods_change_byid", controller.GetGoodsChangeById)
		eg_v2.POST("/post_goods_change_desc", controller.PostGoodsChangeById)
		eg_v2.POST("/create_change", controller.CreateChange)
		eg_v2.POST("/update_change", controller.UpdateChange)
		eg_v2.GET("/delete_change", controller.DeleteCharge)

		// 沙发总价计算
		eg_v2.GET("/get_all_price", controller.GetAllPrice)

		// 成本merge相关
		eg_v2.GET("/get_merge_list", controller.GetMergeList)
		eg_v2.GET("/get_merge_by_id", controller.GetMergeById)
		eg_v2.POST("/update_merge", controller.UpdateMerge)
		eg_v2.POST("/create_merge", controller.CreateMerge)
		eg_v2.POST("/post_goods_merge_desc", controller.PostGoodsMergeById)

		eg_v2.GET("/delete_merge", controller.DeleteMergeById)

		//沙发草稿 开始********************************************************************************************************************************************

		//上传excel
		eg_v2.POST("/uewiohrewr", controller.ImportFenweiInfo)

		// 沙发
		eg_v2.GET("/get_draf_shafa_import_list", controller.GetDrafShaFaImportList)

		eg_v2.GET("/get_draf_shafa_import_byid", controller.GetDrafShaFaImportById)
		eg_v2.POST("/post_draf_feng_wei", controller.PostDrafFengWei)
		eg_v2.GET("/get_draf_feng_wei", controller.GetDrafFinWei)
		eg_v2.GET("/get_feng_wei_group_by_name_draf", controller.GetFinWeiGroupByNameDraf)

		//获取材料表
		eg_v2.GET("/get_export_goods_group_by_fen_wei_draf", controller.GetExportGoodsGroupByFenWeiDraf)

		eg_v2.GET("/get_all_price_draf", controller.GetAllPriceDraf)
		//拷贝一个沙发配置
		eg_v2.POST("/cpoy_shafa_draf", controller.CopyShaFaDraf)
		eg_v2.POST("/reload_shafa_draf", controller.ReloadShaFaDraf)
		eg_v2.GET("/get_shafa_trans", controller.GetShafaTransList)
		eg_v2.GET("/check_trans", controller.CheckTrans)
		eg_v2.GET("/online_trans", controller.OnlineTrans)
		eg_v2.GET("/check_gongyi", controller.CheckGongyi)

		//拷贝一个沙发配置
		//eg_v2.POST("/cpoy_shafa", controller.CopyShaFa)
		//拷贝一个沙发配置
		//沙发草稿 结束********************************************************************************************************************************************
	}

}
