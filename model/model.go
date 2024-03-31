package model

type Room struct {
	User          `json:"user"`
	ShortId       int    `json:"short_id"`
	LongId        int64  `json:"long_id"`
	FollowerCount int64  `json:"follower_count"`
	Cover         string `json:"cover"`
	Title         string `json:"title"`
}

type User struct {
	UID           int64  `json:"uid" gorm:"unique"`
	Name          string `json:"name"`
	Sex           string `json:"sex"`
	Avatar        string `json:"avatar"`
	FollowerCount int64  `json:"fans_count,omitempty"`
	*Medal        `json:"medal,omitempty" gorm:"embedded;embeddedPrefix:medal_"`
}

type Medal struct {
	Name     string `json:"name"`
	OwnerID  int64  `json:"owner_id"`
	Level    int    `json:"level,omitempty"`
	TargetID int64  `json:"target_id"`
}
