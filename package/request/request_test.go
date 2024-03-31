package request

import "testing"

func TestGetArgs(t *testing.T) {
	err := GetArgs("订阅动态 12312315 泛式")
	if err != nil {
		return
	}
}
