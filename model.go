package ironman

type (
	User struct {
		ID         int    `json:"id",gorm:"primary_key"`
		Version    int    `gorm:"column:_version"`
		ParentId   int    `json:"parentId",gorm:"column:parent_id"`
		UserName   string `json:"userName",gorm:"column:user_name"`
		RealName   string `json:"realName",gorm:"column:real_name"`
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
		CreatedAt  int64  `json:"createdAt",gorm:"column:created_at"`
		UpdatedAt  int64  `json:"updatedAt",gorm:"column:updated_at"`
	}
	UserWx struct {
		ID      int    `json:"id",gorm:"primary_key"`
		UnionID string `json:"unionId"`
		OpenID  string `json:"openId"`
	}

	Map map[string]interface{}
)

func (u UserWx) TableName() string { return "user_wx" }
