package testdata

import (
	"fmt"

	"github.com/mats0319/unnamed_plan/server/internal/db/model"
)

// TestNotes 提供大量数据给前端测试展示效果，建议前端在测试阶段把每页尺寸调整为3
func TestNotes() []*model.Note {
	res := make([]*model.Note, 0, 10)
	for i := range 10 {
		res = append(res, &model.Note{
			Writer:      "mats0319",
			WriterName:  "Mario",
			IsAnonymous: false,
			Title:       fmt.Sprintf("test title %d", i+1),
			Content:     fmt.Sprintf("test content %d", i+1),
		})
	}

	return res
}
