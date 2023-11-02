package controller

import (
	"ccccc/common"
	model "ccccc/data/model/goods"
	model3 "ccccc/data/model/role"
	model2 "ccccc/data/model/user"
	"ccccc/service"
	"ccccc/util"
	"encoding/json"
	"errors"
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
	log.Printf("Sums")
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
	a5 := 1.0
	a6 := 1.0

	if l1 != "" {
		if strings.Contains(l1, "*") {
			s_list := strings.Split(l1, "*")
			if len(s_list) >= 1 {
				a1_s := s_list[0]
				a1, err = strconv.ParseFloat(a1_s, 64)
				if err != nil {
					resp.Status = 201
					resp.Desc = "尺寸格式不合法"
					goto out
				}
			}
			if len(s_list) >= 2 {
				a2_s := s_list[1]
				a2, err = strconv.ParseFloat(a2_s, 64)
				if err != nil {
					resp.Status = 201
					resp.Desc = "尺寸格式不合法"
					goto out
				}
			}
			if len(s_list) >= 3 {
				a3_s := s_list[2]
				a3, err = strconv.ParseFloat(a3_s, 64)
				if err != nil {
					resp.Status = 201
					resp.Desc = "尺寸格式不合法"
					goto out
				}
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
			if len(s_list) >= 1 {
				a4_s := s_list[0]
				a4, err = strconv.ParseFloat(a4_s, 64)
				if err != nil {
					resp.Status = 201
					resp.Desc = "尺寸格式不合法"
					goto out
				}
			}
			if len(s_list) >= 2 {
				a5_s := s_list[1]
				a5, err = strconv.ParseFloat(a5_s, 64)
				if err != nil {
					resp.Status = 201
					resp.Desc = "尺寸格式不合法"
					goto out
				}
			}

		} else {
			a4, err = strconv.ParseFloat(size, 64)
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

func ChongSuan(shafaCode string) {
	//根据沙发 找出所有成本
	allGongyi := model.GongYiList{}

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

func GetPrice(num_i, shunhao_i, price_i, huansuan_i, fuzhu_xishu_i, main_xishu_i float64, cpCode, main_size, size, owner_size string) PriceInfo {

	totalPrice := 0.0
	totalPrice_t := 0.0
	//计价数量
	jijiaNums := 0.0

	//根据产品ID 获取产品

	goods := model.Goods{}
	info, err := goods.GetGoodsById(cpCode, nil)
	if err != nil {
		log.Printf("GetGoodsById: err :%v", err)
		return PriceInfo{
			0, 0,
		}
	}
	//判断长宽高
	l1 := ""
	if main_size == "" {
		l1 = info.MainSize
	} else {
		l1 = main_size
	}
	if owner_size != "" {
		l1 = owner_size
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
			log.Printf("%s GetGoodsChangeById err :%v", "未找到产品单位基本信息", err3)
		}
		unit := model.UnitDesc{}
		intv1, err1 := strconv.Atoi(goodsMerge.CpMainUnit)

		if err1 != nil {
			log.Printf("%s GetGoodsChangeById err :%v", "未找到产品单位基本信息", err1)
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

	return riceInfo
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
	cp_gui_ge := strings.TrimSpace(params["cp_gui_ge"])
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
		shafaLogs.GG = cp_gui_ge
		shafaLogs.Update(nil)

		// 改draf表

		shafaLogsDraf := model.ShaFaDrafImportLog{}
		shafaLogsDraf.GetByType(nil, goods.CpCode)
		shafaLogsDraf.SfName = cp_name
		shafaLogsDraf.GG = cp_gui_ge
		shafaLogsDraf.Update(nil)

		// 看一下规格是不是和原来一样
		if goods.CpGuiGe != cp_gui_ge {
			//修改了 规格 先把原来的 放到map中

			guigeMap := make(map[string]bool, 0)

			for _, s := range getshafaGuiGe(goods.CpGuiGe) {
				guigeMap[s] = false
			}
			for _, s := range getshafaGuiGe(cp_gui_ge) {
				_, Okk := guigeMap[s]
				if Okk {
					//如果存在 就弄成 true
					guigeMap[s] = true
				}
			}
			for s, b := range guigeMap {
				if !b {
					//说明这次没有这个了 要把成本表 对应分位的成本删除。
					deleteInfo := model.GongYi{}
					deleteInfo.FenWeiName = s
					deleteInfo.ShafaId = goods.CpCode
					err = deleteInfo.DeleteByFenweiName(nil)
					if err != nil {
						log.Printf("deleteInfo.Delete err :%v", err)
					}
					//删除草稿的
					deleteInfoDraf := model.GongYiDraf{}
					deleteInfoDraf.FenWeiName = s
					deleteInfoDraf.ShafaId = goods.CpCode
					err = deleteInfoDraf.DeleteByFenweiName(nil)
					if err != nil {
						log.Printf("deleteInfo.Delete err :%v", err)
					}

				}
			}
		}
	}

	util.ReturnCompFunc(c, resp)
	return
}

func getshafaGuiGe(guige string) []string {

	asc := strings.Split(guige, "+")
	return asc

}

func covertPrice(string2 string) string {
	if string2 != "" {
		ppp1, _ := strconv.ParseFloat(string2, 64)
		return fmt.Sprintf("%.4f", ppp1)
	} else {
		return "0"
	}

}

func covertPriceSix(string2 string, sex int) string {
	if string2 != "" {
		ppp1, _ := strconv.ParseFloat(string2, 64)

		return fmt.Sprintf("%.4f", ppp1*float64(sex))
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

func CopyShaFa(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	copy_shafa_code := params["copy_shafa_code"]
	sf_code := params["cp_code"]

	allGongyi := model.GongYiList{}

	//获取所有的配置
	err := allGongyi.GetBySoFaCode(nil, sf_code)
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

	//先清空
	deleteInfo := model.GongYi{}
	deleteInfo.ShafaId = copy_shafa_code
	err = deleteInfo.Delete(nil)
	if err != nil {
		log.Printf("deleteInfo.Delete err :%v", err)
		resp.Status = 201
		resp.Desc = err.Error()
		util.ReturnCompFunc(c, resp)
		return
	}

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

func ReloadShaFa(c *gin.Context) {
	params := common.GetUrlParams(c.Request)

	resp := Response{
		Status: 200,
	}

	sf_code := params["cp_code"]

	ChongSuan(sf_code)
	resp.Data = ""
	util.ReturnCompFunc(c, resp)
	return
}

type AllSheetData struct {
	List []SheetData `json:"sheet_name" form:"sheet_name"`
}

type SheetData struct {
	SheetName string     `json:"sheet_name" form:"sheet_name"`
	Data      [][]string `json:"data" form:"data"`
}

func ImportFenweiInfo(c *gin.Context) {
	log.Printf("DeleteMergeById")
	params := common.GetUrlParams(c.Request)
	resp := Response{
		Status: 200,
	}

	name := params["sheet_name"]
	britData := params["data"]
	user := params["user"]

	sheetNames := strings.Split(name, "@")

	sheetDatas := strings.Split(britData, "@")

	if len(sheetNames) != len(sheetDatas) {
		log.Printf("文件格式错误，sheet数量和内容不匹配 sheetNames ： %s", sheetNames)
		resp.Status = 201
		resp.Desc = "文件格式错误，sheet数量和内容不匹配"
		util.ReturnCompFunc(c, resp)
		return
	}

	allSheetData := AllSheetData{}
	allSheetData.List = make([]SheetData, 0)

	for i, sheetName := range sheetNames {
		sheetData := SheetData{}
		if strings.TrimSpace(sheetName) == "" {
			continue
		}
		sheetData.SheetName = sheetName
		ll_data := sheetDatas[i]
		//拆分每行的数据
		allrow := strings.Split(ll_data, "$")
		//小于两行不处理
		if len(allrow) <= 3 {
			continue
		}
		//拆分列
		for _, s := range allrow {
			d1 := strings.Split(s, ",")
			sheetData.Data = append(sheetData.Data, d1)
		}
		allSheetData.List = append(allSheetData.List, sheetData)
	}

	if len(allSheetData.List) == 0 {
		log.Printf("文件格式错误， 没有识别到有效的数据 sheetNames ： %s", sheetNames)
		resp.Status = 201
		resp.Desc = "文件格式错误，没有识别到有效的sheet数据 "
		util.ReturnCompFunc(c, resp)
		return
	}

	//校验沙发名称

	for _, data := range allSheetData.List {
		if len(data.Data[0]) == 0 {
			log.Printf("文件格式错误， sheet : %s  没有配置沙发名称", data.SheetName)
			resp.Status = 201
			resp.Desc = fmt.Sprintf("文件格式错误， sheet : %s  没有配置沙发名称 ", data.SheetName)
			util.ReturnCompFunc(c, resp)
			return
		}
	}
	//获取 沙发名称

	shafaName := strings.TrimSpace(allSheetData.List[0].Data[3][1])
	if shafaName == "" {
		log.Printf(" 找不到沙发 : %s", shafaName)
		resp.Status = 201
		resp.Desc = fmt.Sprintf("找不到沙发 : %s ", shafaName)
		util.ReturnCompFunc(c, resp)
		return
	}
	shafaId := strings.TrimSpace(allSheetData.List[0].Data[3][0])
	if shafaId == "" {
		log.Printf(" 找不到沙发 : %s", shafaId)
		resp.Status = 201
		resp.Desc = fmt.Sprintf("找不到沙发 : %s ", shafaId)
		util.ReturnCompFunc(c, resp)
		return
	}
	// 获取沙发 信息
	shafa := model.ShaFaImportLog{}
	errme := shafa.Get("", shafaId)
	if shafa.Id == 0 || errme != nil {
		log.Printf(" 找不到沙发 : %s", shafaId)
		resp.Status = 201
		resp.Desc = fmt.Sprintf("找不到沙发 : %s ", shafaId)
		util.ReturnCompFunc(c, resp)
		return
	}

	//找到沙发了 开始校验

	outInfo := ""

	goodsMap := make(map[string]model.Goods, 0)

	for _, data := range allSheetData.List {
		//每个sheet
		maxCell := 0
		if len(data.Data) <= 2 {
			//只有两行
			outInfo = fmt.Sprintf("%s \n 表格: %s  少于三行; ", outInfo, data.SheetName)
			continue
		}
		for iji := 0; iji < len(data.Data[2]); iji++ {
			if strings.TrimSpace(data.Data[2][iji]) == "单位" {
				maxCell = iji
			}
		}
		if maxCell == 0 {
			//说明没有单位字段  报错？
			outInfo = fmt.Sprintf("%s \n 表格: %s  没有单位数据; ", outInfo, data.SheetName)
			continue
		}
		if (maxCell+1)%2 != 0 {
			//如果是偶数  不合法 因为单位肯定在偶数列上
			//加1 是因为 下标第一位是 0
			outInfo = fmt.Sprintf("%s \n 表格: %s  没有单位列位置不合法; ", outInfo, data.SheetName)
			continue
		}

		//从第四行开始
		for ii := 3; ii < len(data.Data); ii++ {

			if len(data.Data[ii]) <= 7 {
				continue
			}
			if len(data.Data[ii]) <= 10 {
				outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，数据无效 ", outInfo, data.SheetName, ii+1)
				continue
			}
			//这里是每一行数据 。
			//校验每个单元格的数据
			//校验材料是否存在 //第三个单元格是材料
			cp_code := data.Data[ii][2]

			if strings.TrimSpace(cp_code) == "" {
				outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 材料编码为空; ", outInfo, data.SheetName, ii+1, 3)
			}
			//判断材料是否存在
			_, ok := goodsMap[cp_code]
			if !ok {
				goods := model.Goods{}
				goods.CpCode = strings.TrimSpace(cp_code)
				err := goods.GetByCpCode(nil)
				if err != nil || goods.Id == 0 {
					log.Printf("GetGoodsById: err :%v", err)
					outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的材料编码; ", outInfo, data.SheetName, ii+1, 3)
				} else {
					goodsMap[cp_code] = goods
				}
			}
			//第三行 是表头  且包含了 分位置的信息。

			//校验尺寸
			for jjj := 7; jjj < maxCell; jjj = jjj + 3 {
				//单元格从低6个开始
				//分两种情况
				//情况 1 两两一组
				//情况 2 处在最后 且包含单位。
				// 只需要校验是不是数字就行
				base_str := strings.TrimSpace(data.Data[ii][jjj])

				if base_str == "" {
					continue
				}
				if strings.Contains(base_str, "*") {
					//说明有*链接
					ssLLL := strings.Split(base_str, "*")
					if len(ssLLL) == 0 {
						outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的尺寸; ", outInfo, data.SheetName, ii+1, jjj+1)
					} else {
						for _, sssss := range ssLLL {
							_, erq1 := strconv.ParseFloat(sssss, 64)
							if erq1 != nil {
								//不是数值
								outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的尺寸; ", outInfo, data.SheetName, ii+1, jjj+1)
							}
						}
					}

				} else {
					_, erq := strconv.ParseFloat(base_str, 64)
					if erq != nil {
						//不是数值
						outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的尺寸; ", outInfo, data.SheetName, ii+1, jjj+1)
					}
				}

			}
			//校验规格
			for jjj := 8; jjj < maxCell; jjj = jjj + 3 {
				//单元格从低6个开始
				//分两种情况
				//情况 1 两两一组
				//情况 2 处在最后 且包含单位。
				// 只需要校验是不是数字就行
				base_str := strings.TrimSpace(data.Data[ii][jjj])

				if base_str == "" {
					continue
				}
				if strings.Contains(base_str, "*") {
					//说明有*链接
					ssLLL := strings.Split(base_str, "*")
					if len(ssLLL) == 0 {
						outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的规格; ", outInfo, data.SheetName, ii+1, jjj+1)
					} else {
						for _, sssss := range ssLLL {
							_, erq1 := strconv.ParseFloat(sssss, 64)
							if erq1 != nil {
								//不是数值
								outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的规格; ", outInfo, data.SheetName, ii+1, jjj+1)
							}
						}
					}

				} else {
					_, erq := strconv.ParseFloat(base_str, 64)
					if erq != nil {
						//不是数值
						outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的规格; ", outInfo, data.SheetName, ii+1, jjj+1)
					}
				}

			}

			//校验数量
			for jjj := 9; jjj < maxCell; jjj = jjj + 3 {
				//单元格从低6个开始
				//分两种情况
				//情况 1 两两一组
				//情况 2 处在最后 且包含单位。
				// 只需要校验是不是数字就行
				base_str := strings.TrimSpace(data.Data[ii][jjj])

				if base_str == "" {
					continue
				}

				_, erq := strconv.ParseFloat(base_str, 64)
				if erq != nil {
					//不是数值
					outInfo = fmt.Sprintf("%s \n 表格: %s  第%d行，第%d列 无效的数量; ", outInfo, data.SheetName, ii+1, jjj+1)

				}
			}

		}
	}
	//所有的表格都校验过了

	if outInfo != "" {
		//说明没有校验通过
		log.Printf(" 导入失败  沙发名称: %s", shafaName)
		resp.Status = 201
		resp.Desc = fmt.Sprintf(" 导入失败  具体原因为  :\n %s ", outInfo)
		util.ReturnCompFunc(c, resp)
		return
	}

	//校验通过 开始组装数据了；

	//获取transeID

	for _, data := range allSheetData.List {
		//每个sheet
		maxCell := 0
		if len(data.Data) <= 2 {
			//只有两行
			outInfo = fmt.Sprintf("%s \n 表格: %s  少于三行; ", outInfo, data.SheetName)
			continue
		}
		for iji := 0; iji < len(data.Data[2]); iji++ {
			if strings.TrimSpace(data.Data[2][iji]) == "单位" {
				maxCell = iji
			}
		}
		transID, errxas := GetTranseId(strings.TrimSpace(data.Data[3][0]), user)
		if transID == "" || errxas != nil {
			//说明没有校验通过
			log.Printf(" 导入失败  沙发名称: %s", shafaName)
			resp.Status = 201
			resp.Desc = fmt.Sprintf(" 导入失败  具体原因为  :\n %s ", "创建事务失败")
			util.ReturnCompFunc(c, resp)
			return
		}

		ConvertPostInfo(user, maxCell, goodsMap, transID, getGongyiNamesa(data.SheetName), data.Data[3][0], data.Data)
	}

	resp.Data = name
	resp.Data = britData

	util.ReturnCompFunc(c, resp)
	return
}

func GetTranseId(shafa_code, user string) (string, error) {
	//

	oldTransList := model.TransList{}

	errme := oldTransList.GetByShafaId(nil, shafa_code)
	if errme != nil && errme != gorm.ErrRecordNotFound {
		return "", errors.New("导入失败 获取导出记录 数据库错误")
	} else if errme == gorm.ErrRecordNotFound || len(oldTransList) == 0 {
		//创建一个新的
		return RomId(shafa_code, user)
	}
	if oldTransList[0].IsSubmit == 1 {
		//如果最近的一个已经上线了 创建一个新的
		return RomId(shafa_code, user)
	} else {
		//最近的一个还没上线  返回最新的一个ID
		return oldTransList[0].TransId, nil
	}

}
func RomId(code string, user string) (string, error) {
	////生成一个随机ID

	nix := time.Now().Unix()

	s := fmt.Sprintf("%d_%s", nix, code)

	// 要记录一下

	trans := model.Trans{
		TransId:    s,
		ShafaCode:  code,
		CreateTime: time.Now(),
		CreateUser: user,
	}
	er := trans.Create(nil)
	if er != nil {
		log.Printf(" 导入失败 创建事务入库失败 沙发code: %s,err:%v", code, er)
		return s, er
	}

	return s, nil

}

func getGongyiNamesa(base string) string {

	if strings.Contains(base, "海绵") {
		return "海绵"
	}
	if strings.Contains(base, "裁工") {
		return "裁工"
	}
	if strings.Contains(base, "车工") {
		return "车工"
	}
	if strings.Contains(base, "木工") {
		return "木工"
	}
	if strings.Contains(base, "扪工") {
		return "扪工"
	}

	return ""
}
