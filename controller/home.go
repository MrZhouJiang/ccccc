package controller

import (
	"ccccc/common"
	model "ccccc/data/model/goods"
	model3 "ccccc/data/model/role"
	model2 "ccccc/data/model/user"
	"ccccc/service"
	"ccccc/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Data   interface{} `json:"data"`
	Status int         `json:"status"`
	Desc   string      `json:"desc"`
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func GetById(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 2001,
	}
	list := make([]Person, 0)

	list = append(list, Person{
		Name: params["p1"],
		Age:  12,
	})
	list = append(list, Person{
		Name: "罗慧1",
		Age:  122,
	})
	resp.Data = list

	util.ReturnCompFunc(c, resp)
	return
}

func Post1(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 2001,
	}
	list := make([]Person, 0)

	list = append(list, Person{
		Name: params["p1"],
		Age:  12,
	})
	list = append(list, Person{
		Name: "罗慧1",
		Age:  122,
	})
	resp.Data = list

	util.ReturnCompFunc(c, resp)
	return
}

func UserLogin(c *gin.Context) {
	log.Printf("user login")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}
	name := params["username"]
	word := params["password"]

	log.Printf("username:%s word: %s", name, word)
	//查找用户名

	user := model2.UserInfo{
		Name: name,
	}
	err := user.Get("", name)
	if err != nil {
		resp.Status = 201
		resp.Desc = "用户不存在"

	} else {
		//查看秘密
		password_ssa, err1 := common.Encrypt(word)
		if err1 != nil {
			resp.Status = 202
			resp.Desc = "密码错误"
		} else {
			if password_ssa != user.PassWord {
				resp.Status = 203
				resp.Desc = "密码错误"
			}
		}
	}

	r1 := model3.UserRole{}
	roleId := 0
	err1 := r1.GetUserRole(user.UserId, nil)
	if err1 != nil {
		log.Printf("GetUserRole err :%v", err1)
		resp.Status = 204
		resp.Desc = "无权限"
	}
	roleId = r1.RoleId
	login := LoginUser{
		Name:   user.Name,
		UserId: user.UserId,
		RoleId: roleId,
	}
	log.Printf("login sucess  userInfo :%v", login)

	resp.Data = login
	util.ReturnCompFunc(c, resp)
	return
}

type LoginUser struct {
	//用户唯一ID
	UserId string `json:"user_id"`
	//昵称
	Name   string `json:"name"`
	RoleId int    `json:"role_id"`
}

func CreateUser(c *gin.Context) {
	log.Printf("user login")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}
	name := params["username"]
	create_password := params["create_password"]
	create_name := params["create_user_name"]
	role_id := params["role_id"]
	phone := params["phone"]
	desc := params["desc"]
	log.Printf("username:%s create_password: %s create_name: %s", name, create_password, create_name)
	//查找当前用户有没有创建用户的权限
	user := model2.UserInfo{
		Name: name,
	}
	err := user.Get("", name)
	if err != nil {
		resp.Status = 201
		resp.Desc = "当前用户无权限201"
		goto Out
	} else {
		//查找用户权限 创建用户权限是 1
		b, err1 := CheckUserPermission(user.UserId, 1)
		if err1 != nil || !b {
			resp.Status = 202
			resp.Desc = "当前用户无权限202"
			goto Out
		}
		//有权限  开始创建用户

		//加密密码
		password_ssa, err1 := common.Encrypt(create_password)
		if err1 != nil {
			resp.Status = 203
			resp.Desc = "密码生成失败"
			goto Out
		}
		createUser := model2.UserInfo{
			CreateUser: name,
			Name:       create_name,
			UserId:     common.RandStr(12),
			PassWord:   password_ssa,
			Phone:      phone,
			DescInfo:   desc,
		}
		err1 = createUser.Create(nil)
		if err1 != nil {
			resp.Status = 204
			resp.Desc = "创建用户失败 ，请重试"
			goto Out
		}
		//开始创建角色
		intRole, _ := strconv.Atoi(role_id)
		userRole := model3.UserRole{
			UserId:   createUser.UserId,
			RoleId:   intRole,
			AddUser:  name,
			AddTime:  time.Now(),
			DescInfo: desc,
			RoleName: GetRoleNameDesc(intRole),
		}
		userRole.Create(nil)
		if err1 != nil {
			resp.Status = 205
			resp.Desc = "创建用户权限失败 ，请重试"
			goto Out
		}
	}
Out:
	resp.Data = user
	util.ReturnCompFunc(c, resp)
	return
}

func CheckUserPermission(userId string, permissionId int) (is bool, err error) {
	role := model3.UserRole{
		UserId: userId,
	}
	perList, err := role.GetUserPermission()
	if err != nil || perList == nil || len(perList) == 0 {
		return
	}
	hasPer := false
	for _, permission := range perList {
		if permission.PermissionId == permissionId {
			hasPer = true
		}
	}
	return hasPer, nil
}

func GetGoodsList(c *gin.Context) {
	log.Printf("GetGoodsList")
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
	types := params["goods_type"]
	shunhao := params["shunhao"]
	log.Printf("goods_name:%s goods_code: %s goods_type:%s", name, code, types)

	d, total, err := service.GetGoodsList(page, size, name, code, types, shunhao)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	dto := GetGoodsListDto{
		page,
		size,
		total,
		d,
	}
	resp.Data = dto

	util.ReturnCompFunc(c, resp)
	return
}

type GetGoodsListDto struct {
	Page  int                `json:"page"`
	Size  int                `json:"size"`
	Total int                `json:"total"`
	List  []service.GoodsDto `json:"list"`
}

func GetShaFaImportList(c *gin.Context) {
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

	d, total, err := service.GetShaFaImportList(page, size, name, code, isSums)
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

type GetShaFaImportListDto struct {
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
	Total int                    `json:"total"`
	List  []service.ShaFaInfoDto `json:"list"`
}

type PostGoodsInfo struct {
	List []PostGongYi `json:"list"`
}
type PostGoodsInfoInDTO struct {
	Details []model.GongYi `json:"details"`
	// 物料的类型 车工 裁工
	Types   string `json:"types"`
	ShafaId string `json:"shafa_id"`
}

type PostGongYi struct {
	FenWeiName string         `json:"fen_wei_name"`
	Details    []model.GongYi `json:"details"`
}

func PostFengWei(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	data_ := params["data"]

	dto := PostGoodsInfoInDTO{}

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

	//先清空
	deleteInfo := model.GongYi{}
	deleteInfo.Types = dto.Types
	deleteInfo.ShafaId = dto.ShafaId
	err = deleteInfo.DeleteByTypes(nil)
	if err != nil {
		log.Printf("deleteInfo.Delete err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
		return
	}
	// 要将 全套用量 替换掉其他掉

	//循环插入
	for _, detail := range dto.Details {
		insertInfo := model.GongYi{}
		insertInfo.CreateTime = time.Now()
		insertInfo.Types = dto.Types
		insertInfo.ShafaId = dto.ShafaId
		insertInfo.FenWeiName = detail.FenWeiName
		insertInfo.CLName = detail.CLName
		insertInfo.Size = detail.Size
		insertInfo.Nums = detail.Nums
		insertInfo.Unit = detail.Unit
		insertInfo.Descs = detail.Descs
		insertInfo.TotalPrice = detail.TotalPrice
		insertInfo.CpCode = detail.CpCode
		insertInfo.GongYiName = detail.GongYiName
		insertInfo.JiJiaNum = detail.JiJiaNum
		insertInfo.GoodsPoint = detail.GoodsPoint
		insertInfo.JiJiaUnit = detail.Unit

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
	//修改 沙发表
	shafa := model.ShaFaImportLog{}
	errme := shafa.GetByType(nil, dto.ShafaId)
	if errme == nil {
		size := len(dto.Details)
		if size == 0 {
			shafa.IsSums = "否"
		} else {
			shafa.IsSums = "是"
		}
		shafa.Update(nil)
	}

	resp.Data = dto
	util.ReturnCompFunc(c, resp)
	return
}

func GetFinWei(c *gin.Context) {
	log.Printf("GetFenWei")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	types := params["types"]

	log.Printf("shafa_id:%s types: %s ", shafaId, types)

	d, err := service.GetFenWeiList(shafaId, types)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

func GeUserList(c *gin.Context) {
	log.Printf("getUserList")
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
	user_name := params["user_name"]
	log.Printf("user_name:%s ", user_name)

	d, total, err := service.GetUserListD(page, size, user_name)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	dto := GetUserListDto{
		page,
		size,
		total,
		make([]UserInfo, 0),
	}
	for _, info := range d {
		bb := UserInfo{
			Name:       info.Name,
			UserId:     info.UserId,
			Phone:      info.Phone,
			CreateTime: info.CreateTime,
			DescInfo:   info.DescInfo,
			CreateUser: info.CreateUser,
			RoleName:   getRoleName(info.UserId),
		}
		dto.List = append(dto.List, bb)
	}

	resp.Data = dto

	util.ReturnCompFunc(c, resp)
	return
}

type GetUserListDto struct {
	Page  int        `json:"page"`
	Size  int        `json:"size"`
	Total int        `json:"total"`
	List  []UserInfo `json:"list"`
}
type UserInfo struct {

	//用户唯一ID
	UserId string `json:"user_id"`
	//昵称
	Name string `json:"name"`
	//Phone
	Phone string `json:"phone"`
	//创建时间
	CreateTime time.Time `json:"create_time"`
	//备注
	DescInfo   string `json:"desc_info"`
	CreateUser string `json:"create_user"`
	RoleName   string `json:"role_name"`
}

func getRoleName(userId string) string {

	userR := model3.UserRole{}
	err := userR.GetUserRole(userId, nil)
	if err != nil {
		return ""
	} else {
		return userR.RoleName
	}
}

var RoleMap = map[int]string{}

func init() {
	RoleMap = make(map[int]string, 0)
	RoleMap[1] = "超级管理员"
	RoleMap[2] = "工艺部管理员"
	RoleMap[3] = "财务部门管理员"
	RoleMap[4] = "生产部门管理员"
}

func GetRoleNameDesc(roleId int) string {
	val, ok := RoleMap[roleId]

	if ok {
		return val
	} else {
		return "未知"
	}
}

func GetUnit(c *gin.Context) {
	log.Printf("GetUnit")

	resp := Response{
		Status: 200,
	}
	resp.Data = common.UnitList

	util.ReturnCompFunc(c, resp)
	return
}

func GetAllFenWeiDesc(c *gin.Context) {
	log.Printf("GetAllFenWeiDesc")

	resp := Response{
		Status: 200,
	}
	resp.Data = common.FenWeiListDesc

	util.ReturnCompFunc(c, resp)
	return
}

func GetAllGoodsTypeListDesc(c *gin.Context) {
	log.Printf("GetAllGoodsTypeListDesc")

	resp := Response{
		Status: 200,
	}
	resp.Data = common.GoodsTypeListDesc

	util.ReturnCompFunc(c, resp)
	return
}

func GetAllGoodsListDesc(c *gin.Context) {
	log.Printf("GetGoodsList")

	resp := Response{
		Status: 200,
	}

	d, err := service.GeAllGoodsList()
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return
}

func GetGoodsById(c *gin.Context) {
	log.Printf("GetGoodsById")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}

	cp_cpde := params["cp_code"]

	d, err := service.GeGoodsById(cp_cpde)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return
}

func Sums(c *gin.Context) {
	log.Printf("GetFenWei")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	num := params["num"]
	num_i, e1 := strconv.ParseFloat(num, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}
	cpCode := params["cp_code"]

	shunhao := params["shun_hao"]
	shunhao_i, e1 := strconv.ParseFloat(shunhao, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}

	price := params["price"]
	price_i, e1 := strconv.ParseFloat(price, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}
	huansuan := params["huan_suan"]
	huansuan_i, e1 := strconv.ParseFloat(huansuan, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}
	//输入的尺寸
	size := params["size"]

	main_xishu := params["main_xishu"]
	main_xishu_i, e1 := strconv.ParseFloat(main_xishu, 64)
	if e1 != nil {
		log.Printf("Sums:err :%v", e1)

	}

	main_size := params["main_size"]

	fuzhu_xishu := params["fuzhu_xishu"]
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
		log.Printf("GetGoodsById: err :%v", err)

		resp.Status = 201
		resp.Desc = "没有找到该物料"
		resp.Data = totalPrice

		util.ReturnCompFunc(c, resp)
		return
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

	if l1 != "" {
		if strings.Contains(l1, "*") {
			s_list := strings.Split(l1, "*")
			a1_s := s_list[0]
			a1, err = strconv.ParseFloat(a1_s, 64)
			if err != nil {
				resp.Status = 201
				resp.Desc = "尺寸格式不合法"
				goto out
			}
			a2_s := s_list[1]
			a2, err = strconv.ParseFloat(a2_s, 64)
			if err != nil {
				resp.Status = 201
				resp.Desc = "尺寸格式不合法"
				goto out
			}
		} else {
			a1, err = strconv.ParseFloat(l1, 64)
			if err != nil {
				resp.Status = 201
				resp.Desc = "尺寸格式不合法"
				goto out
			}
		}
	}

	if size != "" {
		if strings.Contains(size, "*") {
			s_list := strings.Split(size, "*")
			a3_s := s_list[0]
			a3, err = strconv.ParseFloat(a3_s, 64)
			if err != nil {
				resp.Status = 201
				resp.Desc = "尺寸格式不合法"
				goto out
			}
			a4_s := s_list[1]
			a4, err = strconv.ParseFloat(a4_s, 64)
			if err != nil {
				resp.Status = 201
				resp.Desc = "尺寸格式不合法"
				goto out
			}
		} else {
			a3, err = strconv.ParseFloat(size, 64)
			if err != nil {
				resp.Status = 201
				resp.Desc = "尺寸格式不合法"
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
	jijiaNums = ((a1 * a2 * a3 * a4) * num_i / huansuan_i) / float64(main_xishu_i) / float64(fuzhu_xishu_i)
	//
	//没有合并规则下的价格
	totalPrice_t = (((shunhao_i/100 + 1) * (a1 * a2 * a3 * a4) * num_i) / huansuan_i) / float64(main_xishu_i) / float64(fuzhu_xishu_i)
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

			resp.Status = 201
			resp.Desc = "未找到产品单位基本信息"
		}
		unit := model.UnitDesc{}
		intv1, err1 := strconv.Atoi(goodsMerge.CpMainUnit)

		if err1 != nil {
			log.Printf("GetGoodsChangeById err :%v", err1)

			resp.Status = 201
			resp.Desc = "未找到产品单位信息"
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
	riceInfo := PriceInfo{
		TotalPrice: totalPrice,
		JiJiaNums:  jijiaNums,
	}

	resp.Data = riceInfo

	util.ReturnCompFunc(c, resp)
	return

}

type PriceInfo struct {
	TotalPrice float64 `json:"total_price"`
	JiJiaNums  float64 `json:"ji_jia_nums"`
}

func GetGoodsChangeList(c *gin.Context) {
	log.Printf("GetGoodsList")
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
		size = 20
	}
	change_type := params["change_type"]
	name := params["name"]
	log.Printf("change_type:%s", change_type)

	d, total, err := service.GetGoodsChangeList(page, size, change_type, name)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	dto := GetGoodsChangeListDto{
		page,
		size,
		total,
		d,
	}
	resp.Data = dto

	util.ReturnCompFunc(c, resp)
	return
}

type GetGoodsChangeListDto struct {
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
	Total int                     `json:"total"`
	List  []model.GoodsChangeDesc `json:"list"`
}

func GetGoodsChangeById(c *gin.Context) {
	log.Printf("GetGoodsChangeById")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}
	id := params["goods_dom_id"]

	info := model.GoodsChangeDesc{}
	intl, err := strconv.Atoi(id)
	err = info.GetById(nil, intl)
	if err != nil {
		log.Printf("GetGoodsChangeById err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	util.ReturnCompFunc(c, resp)
	return
}

type GetGoodsChangeByIdDto struct {
	Id    int    `json:"id"`
	CName string `json:"c_name"`
	// 1 损耗 2 换算
	ChangeType string `json:"change_type"`
	// * ？、/
	Types    string               `json:"types"`
	ValuesL  float64              `json:"values_l"`
	DataList []GoodsChangeDescDto `json:"data_list"`
}

type GoodsChangeDescDto struct {
	CpCode string `json:"cp_code"`
	CpName string `json:"cp_name"`
	CpType string `json:"cp_type"`
}

func PostGoodsChangeById(c *gin.Context) {
	log.Printf("PostGoodsChangeById")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	id := params["id"]
	info := model.GoodsChangeDesc{}
	intl, err := strconv.Atoi(id)
	err = info.GetById(nil, intl)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
		util.ReturnCompFunc(c, resp)
		return
	}
	goods_ids := params["goods_ids"]

	//根据ID 获取配置

	goodsCodeList := strings.Split(goods_ids, ",")

	go InsertGoods(goodsCodeList, info)

	util.ReturnCompFunc(c, resp)
	return
}

func InsertGoods(codes []string, change model.GoodsChangeDesc) {
	//获取一下
	for _, code := range codes {
		if code != "" {
			code = strings.Replace(code, " ", "", -1)
			code = strings.Replace(code, "\n", "", -1)
			change_des := model.GoodsChangeDesInfo{}
			change_des.ChangeType = change.ChangeType
			change_des.CpCode = code
			change_des.GetByTypCode(nil)
			if change_des.Id != 0 {
				//先删除
				change_des.Delete(nil)
			}
			change_des.Id = 0
			change_des.ChangeId = change.Id
			change_des.Create(nil)

			if change.ChangeType == "损耗" {
				//
				goods := model.Goods{}
				goods.Get("", code)
				goods.ShunHao = fmt.Sprintf("%f", change.ValuesL)
				goods.Update(nil)

			}
		}

	}

}

func CreateChange(c *gin.Context) {
	log.Printf("PostGoodsChangeById")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	name := params["name"]
	types := params["types"]
	types_x := params["types_x"]
	values := params["values"]

	info := model.GoodsChangeDesc{}
	info.CName = name
	info.ChangeType = types
	info.Types = types_x
	intl, err := strconv.ParseFloat(values, 64)
	info.ValuesL = intl

	err = info.Create(nil)

	if err != nil {
		log.Printf("PostGoodsChangeById err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	util.ReturnCompFunc(c, resp)
	return
}

func UpdateChange(c *gin.Context) {
	log.Printf("UpdateChange")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	id := params["id"]

	types_x := params["types_x"]
	values := params["values"]

	intl, err := strconv.Atoi(id)
	if intl == 0 {
		log.Printf("PostGoodsChangeById err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}
	info := model.GoodsChangeDesc{}
	info.Id = intl
	//先查找
	err = info.GetById(nil, intl)

	if err != nil {
		log.Printf("PostGoodsChangeById err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}
	info.Types = types_x
	floatl, err := strconv.ParseFloat(values, 64)
	info.ValuesL = floatl

	err = info.Save(nil)
	if err != nil {
		log.Printf("PostGoodsChangeById err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	util.ReturnCompFunc(c, resp)
	return
}

func UpdatGoods(c *gin.Context) {
	log.Printf("UpdateChange")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	goods_size := params["goods_size"]
	goods_shunhao := params["goods_shunhao"]
	goods_xishu := params["goods_xishu"]
	goods_dom := params["goods_dom"]
	cp_code := params["cp_code"]
	goods_merge := params["goods_merge"]
	cp_name := params["cp_name"]
	main_xi_shu := params["main_xi_shu"]
	fu_zhu_xi_shu := params["fu_zhu_xi_shu"]
	main_unit := params["main_unit"]
	fu_unit := params["fu_unit"]
	price := params["price"]
	cp_gui_ge := params["cp_gui_ge"]
	cp_type_name := params["cp_type_code"]
	cp_desc := params["cp_desc"]
	gu_ding_price := params["gu_ding_price"]

	typeInfo, ok := common.GoodsTypeMap[cp_type_name]
	typeCode := ""
	if ok {
		typeCode = typeInfo.GoodsTypeId
	}
	//处理合并的
	GoodsChangeDesInfo := model.GoodsChangeDesInfo{}
	GoodsChangeDesInfo.CpCode = cp_code
	GoodsChangeDesInfo.ChangeType = "换算"
	GoodsChangeDesInfo.DeleteByType(nil)

	//再新增
	if goods_dom != "0" && goods_dom != "" {
		info := model.GoodsChangeDesc{}
		c123, _ := strconv.Atoi(goods_dom)
		err := info.GetById(nil, c123)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("未找到DOM记录 err :%v", err)
			resp.Status = 201
			resp.Desc = "未找到DOM记录"
			util.ReturnCompFunc(c, resp)
			return
		}
		//
		GoodsChangeDesInfo.ChangeId = info.Id
		GoodsChangeDesInfo.Save(nil)
	}
	// 处理材料合并的
	GoodsMerge := model.GoodsMergeDesInfo{}
	GoodsMerge.CpCode = cp_code
	GoodsMerge.DeleteByType(nil)
	//再新增
	if goods_merge != "0" && goods_merge != "" {
		info := model.GoodMergeDesc{}
		c1234, _ := strconv.Atoi(goods_merge)
		err := info.GetById(nil, c1234)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("未找到Merge记录 err :%v", err)
			resp.Status = 201
			resp.Desc = "未找到Merge记录"
			util.ReturnCompFunc(c, resp)
			return
		}
		//
		GoodsMerge.CreateTime = time.Now()
		GoodsMerge.MergeId = info.Id
		GoodsMerge.Save(nil)
	}

	//修改goods
	goods := model.Goods{
		CpCode: cp_code,
	}
	err := goods.GetByCpCode(nil)
	if err != nil {
		log.Printf("未找到物料记录 err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到物料记录"
		util.ReturnCompFunc(c, resp)
		return
	}

	//设置数据
	goods.ShunHao = goods_shunhao
	goods.MainSize = goods_size
	cccccl, _ := strconv.ParseFloat(goods_xishu, 64)
	goods.ChangeP = cccccl

	x1, _ := strconv.ParseFloat(main_xi_shu, 64)
	f1, _ := strconv.ParseFloat(fu_zhu_xi_shu, 64)

	p1, _ := strconv.ParseFloat(price, 64)

	p2, _ := strconv.ParseFloat(gu_ding_price, 64)
	updateData := map[string]interface{}{
		"shun_hao":      goods_shunhao,
		"main_size":     goods_size,
		"change_p":      cccccl,
		"cp_name":       cp_name,
		"main_xi_shu":   x1,
		"fu_zhu_xi_shu": f1,
		"cp_main_unit":  main_unit,
		"fu_zhu_unit":   fu_unit,
		"price":         p1,
		"cp_gui_ge":     cp_gui_ge,
		"cp_type_code":  typeCode,
		"cp_desc":       cp_desc,
		"gu_ding_price": p2,
	}

	err = goods.Save(updateData, nil)
	if err != nil {
		log.Printf("PostGoodsChangeById err :%v", err)
		resp.Status = 201
		resp.Desc = "修改物料失败"
		util.ReturnCompFunc(c, resp)
		return
	}
	if err != nil {
		log.Printf("PostGoodsChangeById err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	//如果是沙发 还需要改沙发名称

	if typeCode == "1019" || typeCode == "1020" || typeCode == "1018" || typeCode == "1029" {
		shafaLogs := model.ShaFaImportLog{}
		shafaLogs.GetByType(nil, goods.CpCode)
		shafaLogs.SfName = cp_name
		shafaLogs.Update(nil)

	}

	util.ReturnCompFunc(c, resp)
	return
}

func covertPrice(string2 string) string {
	if string2 != "" {
		ppp1, _ := strconv.ParseFloat(string2, 64)
		return fmt.Sprintf("%.4f", ppp1)
	} else {
		return "0"
	}

}

func GetAllPrice(c *gin.Context) {
	log.Printf("GetAllPrice")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	sofa_code := params["sf_code"]

	allGongyi := model.GongYiList{}

	err := allGongyi.GetBySoFaCode(nil, sofa_code)

	if err != nil {
		log.Printf("GetAllPrice err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}

	// 要根据全套 剔除掉其他的成本
	//生成一个key
	tttMp := make(map[string]bool, 0)

	tempList := make([]model.GongYi, 0)

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
	for _, yi := range allGongyi {
		price, _ := strconv.ParseFloat(yi.TotalPrice, 64)

		if yi.Types == "裁工" {
			yy1.TypeName = "裁工"
			p1 += price
			yy1.TotalPrice = p1
			yy1.TotalSunhao += yi.ShunHaoPrice
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
			yy2.TotalPrice = p2
			yy2.TotalSunhao += yi.ShunHaoPrice
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
			yy3.TotalPrice = p3
			yy3.TotalSunhao += yi.ShunHaoPrice
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
			yy4.TotalPrice = p4
			yy4.TotalSunhao += yi.ShunHaoPrice
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
			yy5.TotalPrice = p5
			yy5.TotalSunhao += yi.ShunHaoPrice
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
	outInfo.TotalShunHao = s1 + s2 + s3 + s4 + s5 + s6 + s7
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

func GetCpPrice(cp_code string) float64 {

	goods := model.Goods{}
	d, _ := goods.GetGoodsById(cp_code, nil)
	return d.Price

}
func GetCpName(cp_code string) string {

	goods := model.Goods{}
	d, _ := goods.GetGoodsById(cp_code, nil)
	return d.CpName

}

type Outt struct {
	SofaName     string  `json:"sofa_name"`
	SofaCode     string  `json:"sofa_code"`
	TotalPrice   float64 `json:"total_price"`
	TotalShunHao float64 `json:"total_shun_hao"`
	List         []IInfo `json:"list"`
}

type IInfo struct {
	//车工 裁工
	TypeName    string     `json:"type_name"`
	TotalPrice  float64    `json:"total_price"`
	TotalSunhao float64    `json:"total_sunhao"`
	List        []IIIIInfo `json:"list"`
}

type IIIIInfo struct {
	CpCode      string  `json:"cp_code"`
	FenWeiName  string  `json:"fen_wei_name"`
	CLName      string  `json:"cl_name"`
	Size        string  `json:"size"`
	Nums        string  `json:"nums"`
	Unit        string  `json:"unit"`
	Descs       string  `json:"descs"`
	TotalPrice  string  `json:"total_price"`
	Price       float64 `json:"price"`
	JiJiaNum    string  `json:"ji_jia_num"`
	ShunHaoNums string  `json:"shun_hao_nums"`
}

func GetFinWeiGroupByName(c *gin.Context) {
	log.Printf("GetFinWeiGroupByName")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	types := params["types"]

	log.Printf("shafa_id:%s types: %s ", shafaId, types)

	d, err := service.GetFenWeiListGroupByName(shafaId, types)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

type GetMergeListDto struct {
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
	Total int                   `json:"total"`
	List  []model.GoodMergeDesc `json:"list"`
}

func GetMergeList(c *gin.Context) {
	log.Printf("GetMergeList")
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
	merge_name := params["merge_name"]
	log.Printf("merge_name:%s", merge_name)

	d, total, err := service.GetGoodsMergeList(page, size, merge_name)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	dto := GetMergeListDto{
		page,
		size,
		total,
		d,
	}
	resp.Data = dto

	util.ReturnCompFunc(c, resp)
	return
}

func GetMergeById(c *gin.Context) {
	log.Printf("GetMergeById")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}
	id := params["id"]

	info := model.GoodMergeDesc{}
	intl, err := strconv.Atoi(id)
	err = info.GetById(nil, intl)
	if err != nil {
		log.Printf("GetGoodsChangeById err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	resp.Data = info
	util.ReturnCompFunc(c, resp)
	return
}

func UpdateMerge(c *gin.Context) {
	log.Printf("UpdateMerge")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	id := params["id"]

	merge_name := params["merge_name"]
	cl_name := params["cl_name"]
	unit := params["unit"]
	price := params["price"]

	intl, err := strconv.Atoi(id)
	if intl == 0 {
		log.Printf("UpdateMerge err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}
	info := model.GoodMergeDesc{}
	info.Id = intl
	//先查找
	err = info.GetById(nil, intl)

	if err != nil {
		log.Printf("UpdateMerge err :%v", err)
		resp.Status = 201
		resp.Desc = "未找到该记录"
		return
	}
	info.Name = merge_name
	info.CLName = cl_name
	info.Unit = GetUnitById(unit)
	floatl, err := strconv.ParseFloat(price, 64)
	info.Price = floatl

	err = info.Save(nil)
	if err != nil {
		log.Printf("UpdateMerge err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	util.ReturnCompFunc(c, resp)
	return
}

func GetUnitById(id string) string {
	unit := model.UnitDesc{}
	intv1, err1 := strconv.Atoi(id)
	unit.GetById(nil, intv1)
	//获取
	if err1 != nil {
		log.Printf("GetShaFaImportList err :%v", err1)
		return "未知"
	}
	return unit.Name
}

func CreateMerge(c *gin.Context) {
	log.Printf("CreateMerge")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	merge_name := params["merge_name"]
	cl_name := params["cl_name"]
	unit := params["unit"]
	price := params["price"]

	info := model.GoodMergeDesc{}
	info.Name = merge_name
	info.CLName = cl_name
	info.Unit = GetUnitById(unit)
	floatl, err := strconv.ParseFloat(price, 64)
	info.Price = floatl

	err = info.Create(nil)
	if err != nil {
		log.Printf("UpdateMerge err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}

	util.ReturnCompFunc(c, resp)
	return
}

func PostGoodsMergeById(c *gin.Context) {
	log.Printf("PostGoodsMergeById")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	id := params["id"]
	info := model.GoodMergeDesc{}
	intl, err := strconv.Atoi(id)
	err = info.GetById(nil, intl)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
		util.ReturnCompFunc(c, resp)
		return
	}
	goods_ids := params["goods_ids"]

	//根据ID 获取配置

	goodsCodeList := strings.Split(goods_ids, ",")

	go InsertMerge(goodsCodeList, info)

	util.ReturnCompFunc(c, resp)
	return
}

func InsertMerge(codes []string, change model.GoodMergeDesc) {
	//获取一下
	for _, code := range codes {
		if code != "" {
			code = strings.Replace(code, " ", "", -1)
			code = strings.Replace(code, "\n", "", -1)
			change_des := model.GoodsMergeDesInfo{}
			change_des.CpCode = code
			change_des.GetByTypCode(nil)
			if change_des.Id != 0 {
				//先删除
				change_des.Delete(nil)
			}
			change_des.Id = 0
			change_des.MergeId = change.Id
			change_des.Create(nil)
		}

	}

}

func GetExportGoodsGroupByFenWei(c *gin.Context) {
	log.Printf("GetExportGoodsGroupByFenWei")
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	shafaId := params["shafa_id"]
	types := params["types"]

	log.Printf("shafa_id:%s types: %s ", shafaId, types)

	d, err := service.GetGoodsListGroupByName(shafaId)
	if err != nil {
		log.Printf("GetGoodsList err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}

	resp.Data = d

	util.ReturnCompFunc(c, resp)
	return

}

func DeleteMergeById(c *gin.Context) {
	log.Printf("DeleteMergeById")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}
	id := params["id"]

	info := model.GoodMergeDesc{}
	intl, err := strconv.Atoi(id)
	info.Id = intl
	err = info.Delete(nil)
	if err != nil {
		log.Printf("info.Delete err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	//删除desc

	change_des := model.GoodsMergeDesInfo{}
	change_des.MergeId = intl
	if change_des.Id != 0 {
		//先删除
		change_des.DeleteByMerGetId(nil)
	}

	resp.Data = info
	util.ReturnCompFunc(c, resp)
	return
}

func DeleteCharge(c *gin.Context) {
	log.Printf("DeleteCharge")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}
	id := params["id"]

	info := model.GoodsChangeDesc{}
	intl, err := strconv.Atoi(id)
	info.Id = intl
	err = info.Delete(nil)
	if err != nil {
		log.Printf("info.Delete err :%v", err)

		resp.Status = 201
		resp.Desc = err.Error()
	}
	//删除desc

	GoodsChangeDesInfo := model.GoodsChangeDesInfo{}
	GoodsChangeDesInfo.ChangeId = intl
	GoodsChangeDesInfo.DeleteByChange_id(nil)

	resp.Data = info
	util.ReturnCompFunc(c, resp)
	return
}

func DeleteUser(c *gin.Context) {
	log.Printf("DeleteUser")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}
	id := params["user_id"]

	user := model2.UserInfo{
		UserId: id,
	}
	err := user.Delete(nil)
	if err != nil {
		log.Printf("user.Delete err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
	}
	//删除desc
	resp.Data = ""
	util.ReturnCompFunc(c, resp)
	return
}
