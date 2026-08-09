package model

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/internal/utils"
	"gorm.io/gorm"
)

type Note struct {
	Model
	NoteID      string `gorm:"unique;not null"` // 使用其他字段计算获得，可用于保证新增接口幂等性
	Writer      string `gorm:"not null"`        // user name
	WriterName  string `gorm:"not null"`        // user nickname(at that time)
	IsAnonymous bool   // 是否匿名，仅前端展示使用
	Title       string
	Content     string `gorm:"not null"`
}

func (n *Note) BeforeCreate(_ *gorm.DB) error {
	if n.NoteID == "" {
		payload := fmt.Sprintf(`"writer":%s,"is anonymous":%t,"title":%s,"content":%s`,
			n.Writer, n.IsAnonymous, n.Title, n.Content)

		n.NoteID = utils.UUIDv5(payload) // 保证新增接口幂等性
	}

	return nil
}
