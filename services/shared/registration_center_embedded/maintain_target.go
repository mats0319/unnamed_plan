package rc_embedded

import (
	"github.com/pkg/errors"
)

type rpcTarget struct {
	list  []string
	index int // polling index
	// 后续拟添加“响应时间”等字段，以支持更多选择实例的方式，默认采用轮询
}

func newTarget(target []string) *rpcTarget {
	return &rpcTarget{
		list:  target,
		index: 0,
	}
}

func (t *rpcTarget) getTarget() (string, error) {
	if len(t.list) < 1 {
		return "", errors.New("no valid target")
	}

	res := t.list[t.index]

	t.index = (t.index + 1) % len(t.list)

	return res, nil
}
