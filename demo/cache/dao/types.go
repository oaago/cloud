package dao

import (
	"github.com/oaago/component/mysql/cache"
)

type FrontStaffUserListModel struct {
	DB *cache.DaoDataSource
	Where
	FrontStaffUserListResult
}

type FrontStaffUserListResult struct {
	UserListSearchResultList []UserListSearchResult
}

type Where struct {
	UserInfo             string `json:"userInfo"`
	IsIntroduce          string `json:"isIntroduce"`
	StaffID              string `json:"staffId"`
	WxInfo               string `json:"wxInfo"` // 微信账号/昵称/备注
	WxID                 string `json:"wxId"`
	BuyCate              string `json:"buyCate"`
	GoodsID              string `json:"goodsId"`
	ReceiverMobile       string `json:"receiverMobile"`
	DepartmentID         string `json:"departmentId"`
	IntoTimeStart        string `json:"intoTimeStart"`
	IntoTimeEnd          string `json:"intoTimeEnd"`
	CreateTimeStart      string `json:"createTimeStart"`
	CreateTimeEnd        string `json:"createTimeEnd"`
	IsRepeat             string `json:"isRepeat"`
	PushSuccess          string `json:"pushSuccess"`
	QyUserType           string `json:"userType"`
	FirstBuySTime        string `json:"firstBuySTime"`
	FirstBuyETime        string `json:"firstBuyETime"`
	FirstBuySMoney       string `json:"firstBuySMoney"`
	FirstBuyEMoney       string `json:"firstBuyEMoney"`
	BuySSumMoney         string `json:"buySSumMoney"`
	BuyESumMoney         string `json:"buyESumMoney"`
	BuySSumOrder         string `json:"buySSumOrder"`
	BuyESumOrder         string `json:"buyESumOrder"`
	LastSuccessTimeStart string `json:"lastSuccessTimeStart"`
	LastSuccessTimeEnd   string `json:"lastSuccessTimeEnd"`
	IsBuy                string `json:"isBuy"`
	PushSTime            string `json:"pushSTime"`
	PushETime            string `json:"pushETime"`
	IsSecond             string `json:"isSecond"`
	Sex                  string `json:"sex"`
	All                  string `json:"all"`
	IsFollow             string `json:"isFollow"`
	UserSource           string `json:"userSource"`
	IntoFansNum          string `json:"intoFansNum" form:"intoFansNum"`
	WxIds                string `json:"wxIds" form:"wxIds"`
}

type UserListSearchResult struct {
	Avatar             string  `json:"avatar" gorm:"column:avatar"`
	CategoryName       string  `json:"categoryName" gorm:"column:categoryName"`
	CreateTime         string  `json:"createTime" gorm:"column:createTime"`
	DataDiff           int     `json:"dataDiff" gorm:"column:dataDiff"`
	DepartmentNo       string  `json:"departmentNo" gorm:"column:departmentNo"`
	DeptName           string  `json:"deptName" gorm:"column:deptName"`
	EncryptWxNo        string  `json:"encryptWxNo" gorm:"column:encryptWxNo"`
	GoodsID            int     `json:"goodsId" gorm:"column:goodsId"`
	ID                 int     `json:"id" gorm:"column:userId"`
	Improve            string  `json:"improve" gorm:"column:improve"`
	IntoTime           string  `json:"intoTime" gorm:"column:intoTime"`
	IsFollow           string  `json:"isFollow" gorm:"column:isFollow"`
	IsIntroduce        string  `json:"isIntroduce" gorm:"column:isIntroduce"`
	IntroduceWxNo      string  `json:"introduceWxNo" gorm:"column:introduceWxNo"`
	Level              int     `json:"level" gorm:"column:level"`
	NickName           string  `json:"nickName" gorm:"column:nickName"`
	RealName           string  `json:"realName" gorm:"column:realName"`
	SecondTime         string  `json:"secondTime" gorm:"column:secondTime"`
	Special            string  `json:"special" gorm:"column:special"`
	StaffID            int     `json:"staffId" gorm:"column:staffId"`
	TmNickname         string  `json:"tmNickname" gorm:"column:tmNickname"`
	UserName           string  `json:"userName" gorm:"column:userName"`
	UserSex            string  `json:"userSex" gorm:"column:userSex"`
	VirtualAccount     string  `json:"virtualAccount" gorm:"column:virtualAccount"`
	WxComm             string  `json:"wxComm" gorm:"column:wxComm"`
	WxID               int     `json:"wxId" gorm:"column:wxId"`
	UserWxNum          string  `json:"userWxNum" gorm:"column:userWxNum"`
	WxNum              string  `json:"wxNum" gorm:"column:wxNum"`
	SecondWx           string  `json:"secondWx" gorm:"column:secondWx"`
	WorkNo             string  `json:"workNo" gorm:"column:workNo"`
	SecondWorkNo       string  `json:"secondWorkNo" gorm:"column:secondWorkNo"`
	PayAmount          float64 `json:"payAmount" gorm:"column:payAmount"`
	PayCount           int     `json:"payCount" gorm:"column:payCount"`
	UserAge            string  `json:"userAge" gorm:"column:userAge"`
	Birthday           string  `json:"birthday" gorm:"column:birthday"`
	ReceiverProvince   string  `json:"receiverProvince" gorm:"column:receiverProvince"`
	ReceiverCity       string  `json:"receiverCity" gorm:"column:receiverCity"`
	ReceiverCountry    string  `json:"receiverCountry" gorm:"column:receiverCountry"`
	SecondRealName     string  `json:"secondRealName" gorm:"column:secondRealName"`
	SecondDepartmentNo string  `json:"secondDepartmentNo" gorm:"column:secondDepartmentNo"`
	SecondDeptName     string  `json:"secondDeptName" gorm:"column:secondDeptName"`
	LastChatTime       string  `json:"lastChatTime" gorm:"column:lastChatTime"`
	Color              string  `json:"color" gorm:"column:color"`
}

type Color struct {
	IntoFans int    `json:"intoFans" gorm:"column:intoFans"`
	Color    string `json:"color" gorm:"column:color"`
}
