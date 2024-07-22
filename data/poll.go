package data

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Poll struct {
	ID          uint64       `json:"id" gorm:"primaryKey"`
	Question    string       `json:"question"`
	Options     []Option     `json:"options" gorm:"foreignKey:PollID"`
	CategoryID  uint64       `json:"category_id"`
	Category    PollCategory `json:"category" gorm:"foreignKey:CategoryID"`
	CreatedById uint64       `json:"created_by_id"`
	CreatedBy   User         `json:"created_by"`
	CreatedAt   *time.Time   `json:"created_at"`
	UpdatedAt   *time.Time   `json:"updated_at"`
}

type Option struct {
	ID         uint64     `json:"id" gorm:"primaryKey"`
	Title      string     `json:"title"`
	TotalVotes int        `json:"total_votes"`
	PollID     uint64     `json:"poll_id"`
	Poll       Poll       `json:"poll"`
	Votes      []Vote     `json:"votes" gorm:"foreignKey:OptionID"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type Vote struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	OptionID uint64 `json:"option_id"`
	Option   Option `json:"option"`
	UserID   uint64 `json:"user_id"`
}

type PollCategory struct {
	ID    uint64 `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

func (p *PollCategory) BeforeCreate(*gorm.DB) error {
	p.Label = strings.ToLower(p.Name)
	return nil
}
