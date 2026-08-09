package utils

import (
	"errors"
	"testing"
)

func TestLogStyle(t *testing.T) {
	e := ErrForTest().WithCause(errors.New("a new error")).
		WithParam("first param", "first value").
		WithParam("second param", 10000)

	t.Log(e.String())
}

func TestMultiUseOnOneInstance(t *testing.T) {
	e := ErrForTest().WithParam("key1", "value1")
	e = ErrForTest().WithParam("key2", "value2")

	t.Log(e.String())

	// 如果Err直接是*Error类型的实例，那么在多处使用该实例时，会继承历史数据
	// （即：假设数据库未启动，dbError一直触发，并且每次都会使用不同参数调用withParam()，
	// 那么一段时间以后，Err的map可能就会带满所有程序中写到的key；
	// 另外，多次触发同一位置的Err没有这个问题，因为error和相同的key会覆盖历史数据）
	//
	// 我们可以看到，实际场景中确实存在同一个Err多处使用的情况，所以我们将Err统一编辑成函数，
	// 这样在每次调用Err时，都会产生新的实例，不会存在继承历史数据的情况。
}
