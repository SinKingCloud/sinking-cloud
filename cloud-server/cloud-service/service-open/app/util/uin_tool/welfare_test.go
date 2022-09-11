package uin_tool

import (
	"fmt"
	"testing"
)

func Test_Run(t *testing.T) {
	a := NewAutoWelfare(
		map[string]string{
			"uin":   "3214750144",
			"skey":  "MWjIZfXScw",
			"pskey": "xjTZZQTkg2eJSqVfefSMYl6aIBDdmKzPUMfHwq6rwvU_",
		},
		map[string]string{
			"uin":   "640170986",
			"skey":  "MQvAn9RQFG",
			"pskey": "pgFJzByZe6K7rZwEwmXrAFHWk392eWlPEFTZLFibS-8_",
		},
		map[string]string{
			"uin":   "109938129",
			"skey":  "MlZyDG5WTT",
			"pskey": "7i9EFcDqOp5ym-ivl0EHq10zYQjHkJhlDDiZNz*XXFo_",
		},
	)
	fmt.Println(a.Run())
}
