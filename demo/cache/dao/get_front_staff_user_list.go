package dao

import (
	"github.com/oaago/common/mobile"
	"github.com/oaago/component/mysql"
	"github.com/oaago/component/mysql/cache"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func InitFrontStaffUserListModel(where Where) *FrontStaffUserListModel {
	return &FrontStaffUserListModel{
		cache.NewDaoDataSource(cache.DBType{
			mysql.GetDBByName("scrm"),
		}),
		where,
		FrontStaffUserListResult{},
	}
}

func (frontStaffUserListModel *FrontStaffUserListModel) GetFrontStaffUserExcel() {
	UserListSearchResultList := frontStaffUserListModel.FrontStaffUserListResult.UserListSearchResultList
	db := frontStaffUserListModel.DB.DB
	sql := db.Debug().Table("user_base ub").
		Select(`ub.user_id AS userId,
						ub.user_name AS userName,
						ub.wx_num AS userWxNum,
						ub.wx_id AS wxId,
						tw.wx_num AS wxNum, 
						sum(tod.pay_amount) payAmount,
						count(tod.order_id) payCount,
						IFNULL(TIMESTAMPDIFF( YEAR, tu.user_age, NOW() ),'未知') AS userAge,
						IF(tua.birthday=-99,'未知',tua.birthday) birthday,
						tod.receiver_province receiverProvince, 
						tod.receiver_city receiverCity,
						tod.receiver_country receiverCountry,
						tu.goods_id AS goodsId,
						tc.category_name AS categoryName,
					CASE
							WHEN tu.user_sex = '' THEN
							'' 
							WHEN tu.user_sex = 0 THEN
							'未知' 
							WHEN tu.user_sex = 1 THEN
							'男' 
							WHEN tu.user_sex = 2 THEN
							'女' 
						END AS userSex,
						FROM_UNIXTIME( tu.create_time / 1000, '%Y-%m-%d %H:%i:%s' ) AS createTime,
						FROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d %H:%i:%s' ) AS intoTime,
						tw.staff_id staffId,
						FROM_UNIXTIME( tu.second_into_time / 1000, '%Y-%m-%d %H:%i:%s' ) AS secondTime,
						tw2.wx_num AS secondWx,
						ss.real_name AS realName,
						ss.department_no AS departmentNo,
						ss.work_no workNo,
						ss2.work_no secondWorkNo,
						ss2.real_name AS secondRealName,
						ss2.department_no AS secondDepartmentNo,
						ss2.work_no secondWorkNo,
					CASE
							
							WHEN tu.is_follow = 1 THEN
							'是' 
							WHEN tu.is_follow = 2 THEN
							'否' 
						END AS isFollow ,
					CASE
							WHEN tu.is_introduce = 1 THEN
							'是' 
							WHEN tu.is_introduce = 2 THEN
							'否' 
						END AS isIntroduce,
						tui.introduce_wx_no introduceWxNo`).
		Joins("LEFT JOIN t_user tu ON tu.id = ub.user_id").
		Joins("LEFT JOIN t_category tc ON tu.goods_id = tc.category_id").
		Joins("LEFT JOIN t_user_tag tut ON ub.user_id = tut.id").
		Joins("LEFT JOIN t_weixin_staff tw ON tw.id = ub.wx_id").
		Joins("LEFT JOIN t_weixin_staff tw2 ON tw2.id = ub.second_wx_id").
		Joins("LEFT JOIN sys_staff ss ON ub.current_staff_id =ss.id").
		Joins("LEFT JOIN sys_staff ss2 ON tw2.staff_id=ss2.id").
		Joins("LEFT JOIN t_user_introduce tui ON tui.user_id = ub.user_id").
		Joins("LEFT JOIN t_user_archives tua ON tua.id = tu.id")
	sql.Joins("LEFT JOIN t_order tod ON tu.id = tod.user_id  AND tod.order_status not IN (12,22,24) AND tod.pay_amount > 0")
	sql.Scopes(
		frontStaffUserListModel.ByStaffId,
		frontStaffUserListModel.ByDepartmentNo,
		frontStaffUserListModel.ByUserInfo,
		frontStaffUserListModel.ByCreateTime,
		frontStaffUserListModel.ByIntoTime,
		frontStaffUserListModel.ByWxInfo,
		frontStaffUserListModel.BySex,
		frontStaffUserListModel.ByIsSecond,
		frontStaffUserListModel.ByIsIntroduce,
		frontStaffUserListModel.ByIsFollow,
		frontStaffUserListModel.ByGoodsId,
		frontStaffUserListModel.ByIsRepeat,
		frontStaffUserListModel.ByPushSuccess,
		frontStaffUserListModel.ByUserSource,
		frontStaffUserListModel.ByQyUserType,
		frontStaffUserListModel.ByFirstBuyTime,
		frontStaffUserListModel.ByLastSuccessTime,
		frontStaffUserListModel.ByPushTime,
		frontStaffUserListModel.ByStaffWxNum,
		frontStaffUserListModel.ByBuyCate,
		frontStaffUserListModel.ByReceiverMobile,
		frontStaffUserListModel.ByBuyCount,
		frontStaffUserListModel.ByBuyTotalAmount,
		frontStaffUserListModel.ByFirstBuyAmount,
		frontStaffUserListModel.ByIntoFansNum,
		frontStaffUserListModel.ByIsBuy,
		frontStaffUserListModel.ByIsDelete,
	)
	sql = sql.Group("ub.user_id")

	sql.Order("ub.user_id DESC")

	//结果
	//sql.Scan(&UserListSearchResultList)
	frontStaffUserListModel.DB.DB.DB = sql
	//不查缓存的话
	//nativeDb := frontStaffUserListModel.DB.GetNativeDb()
	//nativeDb = sql
	//nativeDb.Scan(&UserListSearchResultList)
	//缓存
	frontStaffUserListModel.DB.DB.Scan(1000, frontStaffUserListModel.Where, &UserListSearchResultList)

	frontStaffUserListModel.FrontStaffUserListResult.UserListSearchResultList = UserListSearchResultList
}

func (v *FrontStaffUserListModel) GetColors() []Color {
	var result []Color
	v.DB.DB.Debug().Select("into_fans intoFans, color").Table("t_dict_fans").Where("`status` = 1 AND type = 1").Scan(&result)
	return result
}

//创建时间
func (m *FrontStaffUserListModel) ByCreateTime(db *gorm.DB) *gorm.DB {
	createTimeEnd := m.CreateTimeEnd
	createTimeStart := m.CreateTimeStart
	if createTimeStart != "" && createTimeEnd != "" {
		createTimeStartTimeStamp := GetTimeStampByTimeStr(createTimeStart) * 1000
		createTimeStartEndStamp := GetTimeStampByTimeStr(createTimeEnd) * 1000
		db = db.Where("tu.create_time >= ? AND tu.create_time <= ?", createTimeStartTimeStamp, createTimeStartEndStamp)
	}
	return db
}

//进粉时间
func (m *FrontStaffUserListModel) ByIntoTime(db *gorm.DB) *gorm.DB {
	//进粉时间
	intoTimeStart := m.IntoTimeStart
	intoTimeEnd := m.IntoTimeEnd
	if intoTimeStart != "" && intoTimeEnd != "" {
		intoTimeStartTimeStamp := GetTimeStampByTimeStr(intoTimeStart) * 1000
		intoTimeEndEndStamp := GetTimeStampByTimeStr(intoTimeEnd) * 1000
		db = db.Where("ub.into_time >= ? AND ub.into_time <= ?", intoTimeStartTimeStamp, intoTimeEndEndStamp)
	}
	return db
}

func (m *FrontStaffUserListModel) ByDepartmentNo(db *gorm.DB) *gorm.DB {
	if m.DepartmentID != "" {
		db = db.Where("ss.department_no like CONCAT((SELECT no FROM sys_department WHERE id = ?), '%') or ss3.department_no like CONCAT((SELECT no FROM sys_department WHERE id = ?), '%')", m.DepartmentID, m.DepartmentID)
	}
	return db
}

func (m *FrontStaffUserListModel) ByUserInfo(db *gorm.DB) *gorm.DB {
	if m.UserInfo != "" {
		db = db.Where("ub.user_id = ? OR ub.user_name = ?", m.UserInfo, m.UserInfo)
	}
	return db
}

func (m *FrontStaffUserListModel) ByStaffId(db *gorm.DB) *gorm.DB {
	if m.StaffID != "" {
		if m.WxIds != "" {
			split := strings.Split(m.WxIds, ",")
			db = db.Where("ub.current_staff_id = ? OR ub.wx_id IN (?)", m.StaffID, split)
		} else {
			db = db.Where("ub.current_staff_id = ?", m.StaffID)
		}
	}
	return db
}

func (m *FrontStaffUserListModel) ByBuyCate(db *gorm.DB) *gorm.DB {
	if m.BuyCate != "" {
		db = db.Where("tut.buy_cate = ?", m.BuyCate)
	}
	return db
}

//微信信息 微信号/微信昵称
func (m *FrontStaffUserListModel) ByWxInfo(db *gorm.DB) *gorm.DB {
	if m.WxInfo != "" {
		db = db.Where("(ub.nick_name = HEX(?) OR ub.user_name = ? OR ub.wx_num = ? OR ub.nick_name = ?)", m.WxInfo, m.WxInfo, m.WxInfo, m.WxInfo)
	}
	return db
}

//手机号
func (m *FrontStaffUserListModel) ByReceiverMobile(db *gorm.DB) *gorm.DB {
	if m.ReceiverMobile != "" {
		encMobile := mobile.EncMobile(m.ReceiverMobile)
		db = db.Where("tod.receiver_mobile = ?", encMobile)
	}
	return db
}

//性别
func (m *FrontStaffUserListModel) BySex(db *gorm.DB) *gorm.DB {
	if m.Sex != "" {
		db = db.Where("tu.user_sex = ?", m.Sex)
	}
	return db
}

//是否是VIP微信
func (m *FrontStaffUserListModel) ByIsSecond(db *gorm.DB) *gorm.DB {
	if m.IsSecond != "" {
		db = db.Where("tw.isSecond = ?", m.IsSecond)
	}
	return db
}

//是否是转介绍
func (m *FrontStaffUserListModel) ByIsIntroduce(db *gorm.DB) *gorm.DB {
	if m.IsIntroduce != "" {
		db = db.Where("(tu.is_introduce = ?)", m.IsIntroduce)
	}
	return db
}

//是否重粉
func (m *FrontStaffUserListModel) ByIsFollow(db *gorm.DB) *gorm.DB {
	if m.IsFollow != "" {
		db = db.Where("tu.is_follow = ?", m.IsFollow)
	}
	return db
}

//类目
func (m *FrontStaffUserListModel) ByGoodsId(db *gorm.DB) *gorm.DB {
	if m.GoodsID != "" {
		db = db.Where("tu.goods_id = ?", m.GoodsID)
	}
	return db
}

//是否复购
func (m *FrontStaffUserListModel) ByIsRepeat(db *gorm.DB) *gorm.DB {
	if m.IsRepeat != "" {
		db = db.Where("tut.is_repeat = ?", m.IsRepeat)
	}
	return db
}

//是否推送成功
func (m *FrontStaffUserListModel) ByPushSuccess(db *gorm.DB) *gorm.DB {
	if m.PushSuccess != "" {
		isPushSuccess, _ := strconv.Atoi(m.PushSuccess)
		if isPushSuccess == 1 {
			db = db.Where("ub.second_wx_id != -99")
		} else if isPushSuccess == 2 {
			db = db.Where("ub.second_wx_id = -99 OR ub.second_wx_id IS NULL")
		}
	}
	return db
}

//客户来源
func (m *FrontStaffUserListModel) ByUserSource(db *gorm.DB) *gorm.DB {
	if m.UserSource != "" {
		db = db.Where("ub.user_source = ?", m.UserSource)
	}
	return db
}

//是否是企业微信
func (m *FrontStaffUserListModel) ByQyUserType(db *gorm.DB) *gorm.DB {
	if m.QyUserType != "" {
		db = db.Where("tu.qy_type = ?", m.QyUserType)
	}
	return db
}

//首次购买时间
func (m *FrontStaffUserListModel) ByFirstBuyTime(db *gorm.DB) *gorm.DB {
	firstBuyTimeStart := m.FirstBuySTime
	firstBuyTimeEnd := m.FirstBuyETime
	if firstBuyTimeStart != "" && firstBuyTimeEnd != "" {
		firstBuyTimeStartStamp := GetTimeStampByTimeStrNoEnd(firstBuyTimeStart) * 1000
		firstBuyTimeEndStamp := GetTimeStampByTimeStrNoEnd(firstBuyTimeEnd) * 1000
		db = db.Where("tut.first_buy_time >= ? AND tut.first_buy_time <= ?", firstBuyTimeStartStamp, firstBuyTimeEndStamp)
	}
	return db
}

//最后一次完成时间
func (m *FrontStaffUserListModel) ByLastSuccessTime(db *gorm.DB) *gorm.DB {
	lastSuccessTimeStart := m.LastSuccessTimeStart
	lastSuccessTimeEnd := m.LastSuccessTimeEnd
	if lastSuccessTimeStart != "" && lastSuccessTimeEnd != "" {
		lastSuccessTimeStartStamp := GetTimeStampByTimeStrNoEnd(lastSuccessTimeStart) * 1000
		lastSuccessTimeEndStamp := GetTimeStampByTimeStrNoEnd(lastSuccessTimeEnd) * 1000
		db = db.Where("tut.last_receive_time >= ? AND tut.last_receive_time <= ?", lastSuccessTimeStartStamp, lastSuccessTimeEndStamp)
	}
	return db
}

//推送时间
func (m *FrontStaffUserListModel) ByPushTime(db *gorm.DB) *gorm.DB {
	PushSTimeStart := m.PushSTime
	PushSTimeEnd := m.PushETime
	if PushSTimeStart != "" && PushSTimeEnd != "" {
		PushSTimeStartStamp := GetTimeStampByTimeStr(PushSTimeStart) * 1000
		PushSTimeEndStamp := GetTimeStampByTimeStr(PushSTimeEnd) * 1000
		db = db.Where("tu.second_into_time >= ? AND tu.second_into_time <= ?", PushSTimeStartStamp, PushSTimeEndStamp)
	}
	return db
}

//工作微信号
func (m *FrontStaffUserListModel) ByStaffWxNum(db *gorm.DB) *gorm.DB {
	wxIds := m.WxID
	if wxIds != "" {
		replaceResult := strings.ReplaceAll(wxIds, "，", ",")
		split := strings.Split(replaceResult, ",")
		db = db.Where("tw.wx_num in (?) or tw2.wx_num in (?)", split, split)
	}
	return db
}

//购买总金额
func (m *FrontStaffUserListModel) ByBuyTotalAmount(db *gorm.DB) *gorm.DB {
	buyOrdersStart := m.BuySSumMoney
	buyOrdersEnd := m.BuyESumMoney
	if buyOrdersStart != "" {
		db = db.Where("tut.buy_total_amount >= ?", buyOrdersStart)
	}
	if buyOrdersEnd != "" {
		db = db.Where("tut.buy_total_amount <= ?", buyOrdersEnd)
	}
	return db
}

//购买总单数
func (m *FrontStaffUserListModel) ByBuyCount(db *gorm.DB) *gorm.DB {
	buyOrdersStart := m.BuySSumOrder
	buyOrdersEnd := m.BuyESumOrder
	if buyOrdersStart != "" {
		db = db.Where("tut.buy_total_count >= ?", buyOrdersStart)
	}
	if buyOrdersEnd != "" {
		db = db.Where("tut.buy_total_count <= ?", buyOrdersEnd)
	}
	return db
}

//首次购买金额
func (m *FrontStaffUserListModel) ByFirstBuyAmount(db *gorm.DB) *gorm.DB {
	buyOrdersStart := m.FirstBuySMoney
	buyOrdersEnd := m.FirstBuyEMoney
	if buyOrdersStart != "" {
		db = db.Where("tut.first_buy_amount >= ?", buyOrdersStart)
	}
	if buyOrdersEnd != "" {
		db = db.Where("tut.first_buy_amount <= ?", buyOrdersEnd)
	}
	return db
}

func (m *FrontStaffUserListModel) ByIsBuy(db *gorm.DB) *gorm.DB {
	isBuy := m.IsBuy
	if isBuy == "1" {
		db = db.Where("tu.level = 5")
	} else if isBuy == "2" {
		db = db.Where("tu.level = 1")
	} else if isBuy == "3" {
		db = db.Where("tu.level = 5 and tu.second_into_time != -99")
	} else if isBuy == "4" {
		db = db.Where("tu.level = 5 and tu.second_into_time = -99")
	}
	return db
}

func (m *FrontStaffUserListModel) ByIsDelete(db *gorm.DB) *gorm.DB {
	db = db.Where("ub.sys_delete_flag = 0")
	return db
}

func (m *FrontStaffUserListModel) ByIntoFansNum(db *gorm.DB) *gorm.DB {
	intoFansNum := m.IntoFansNum

	if intoFansNum == "300" {
		db.Where("tu.level = 1")
		db = db.Where("DATEDIFF(NOW(),\tFROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d' )) > 30")
	} else if intoFansNum != "" {
		db.Where("tu.level = 1")
		atoi, _ := strconv.Atoi(intoFansNum)
		if atoi <= 7 {
			db = db.Where("DATEDIFF(NOW(),\tFROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d' )) = ?", intoFansNum)
		} else if atoi > 7 && atoi <= 15 {
			db = db.Where("DATEDIFF(NOW(),\tFROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d' )) > 7 and DATEDIFF(NOW(),\tFROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d' )) <= 15")
		} else if atoi > 15 && atoi <= 30 {
			db = db.Where("DATEDIFF(NOW(),\tFROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d' )) > 15 and DATEDIFF(NOW(),\tFROM_UNIXTIME( ub.into_time / 1000, '%Y-%m-%d' )) <= 30")
		}
	}
	return db
}

func (m *FrontStaffUserListModel) GetStaffAllWxs() {
	var result string
	m.DB.DB.Debug().Select(`group_concat(id)`).Table("t_weixin_staff").Where("staff_id=?", m.StaffID).Scan(&result)
	m.WxIds = result
}

// 日期转化为时间戳
func GetTimeStampByTimeStr(timeStr string) int64 {
	//转化所需模板
	timeLayout := "2006-01-02 15:04:05"
	//获取时区
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	//转化为时间戳 类型是int64
	timestamp := tmp.Unix()
	return timestamp
}

// 日期转化为时间戳
func GetTimeStampByTimeStrNoEnd(timeStr string) int64 {
	//转化所需模板
	timeLayout := "2006-01-02"
	//获取时区
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, timeStr, loc)
	//转化为时间戳 类型是int64
	timestamp := tmp.Unix()
	return timestamp
}
