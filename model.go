package ironman

type (
	// User 用户基础信息
	User struct {
		ID         int    `json:"id" gorm:"primary_key"`
		Version    int    `json:"version,omitempty"gorm:"column:_version"`
		ParentID   int    `json:"parentId,omitempty" gorm:"column:parent_id"`
		UserName   string `json:"userName,omitempty" gorm:"column:user_name"`
		RealName   string `json:"realName,omitempty" gorm:"column:real_name"`
		Password   string `json:"password,omitempty" gorm:"column:password"`
		Type       int8   `json:"type,omitempty"`
		Sex        int8   `json:"sex,omitempty"`
		Avatar     string `json:"avatar,omitempty"`
		City       string `json:"city,omitempty"`
		Province   string `json:"province,omitempty"`
		Country    string `json:"country,omitempty"`
		Email      string `json:"email,omitempty"`
		Mobile     string `json:"mobile,omitempty"`
		RegTime    int64  `json:"regTime,omitempty"`
		UpgradedAt int64  `json:"upgradedAt,omitempty" gorm:"column:upgraded_at"`
		CreatedAt  int64  `json:"createdAt,omitempty" gorm:"column:created_at"`
		UpdatedAt  int64  `json:"updatedAt,omitempty" gorm:"column:updated_at"`
	}
	// UserWx 用户微信信息
	UserWx struct {
		ID      int    `json:"id" gorm:"primary_key"`
		UnionID string `json:"unionId"`
		OpenID  string `json:"openId"`
	}
	// Map 类型
	Map map[string]interface{}
)

// TableName user_wx
func (u UserWx) TableName() string { return "t_user_wx" }
