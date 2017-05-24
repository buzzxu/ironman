package ironman

type (


	User struct {
		ID         int    `json:"id",gorm:"primary_key"`
		UserName   string `json:"userName",gorm:"column:username"`
		RealName   string `json:"realName",gorm:"column:realname"`
		Password   string `json:"password",gorm:"column:password"`
		Type       int8   `json:"type"`
		Sex        int8   `json:"sex"`
		Avatar     string `json:"avatar"`
		City       string `json:"city"`
		Province   string `json:"province"`
		Country    string `json:"country"`
		Email      string `json:"email"`
		Mobile     string `json:"mobile"`
		RegTime    string `json:"regTime"`
		UpgradedAt int    `json:"upgradedAt",gorm:"column:upgraded_at"`
		CreatedAt  int    `json:"createdAt",gorm:"column:created_at"`
		UpdatedAt  int    `json:"updatedAt",gorm:"column:updated_at"`
	}
	UserWx struct {
		ID      int    `json:"id",gorm:"primary_key"`
		UnionID string `json:"unionId"`
		OpenID  string `json:"openId"`
	}
)

func (u UserWx) TableName() string { return "t_user_wx" }


