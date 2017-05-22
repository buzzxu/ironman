package model

type (
	User struct {
		ID         int    `json:"id",gorm:"primary_key"`
		UserName   string `json:"userName",gorm:"column:username"`
		RealName   string `json:"realName",gorm:"column:realname"`
		Type       int8   `json:"type"`
		Sex        int8   `json:"sex"`
		Avatar     string `json:"avatar"`
		City       string `json:"city"`
		Province   string `json:"province"`
		Country    string `json:"country"`
		Email      string `json:"email"`
		Mobile     string `json:"mobile"`
		RegTime    string `json:"regTime"`
		UpgradedAt int    `json:"upgradedAt"`
		CreatedAt  int    `json:"createdAt"`
	}
	UserWx struct {
		ID      int    `json:"id",gorm:"primary_key"`
		UnionID string `json:"unionId"`
		OpenID  string `json:"openId"`
	}
)


func (u UserWx) TableName() string { return "user_wx" }
