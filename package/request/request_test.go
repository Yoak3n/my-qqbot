package request

import "testing"

func TestGetArgs(t *testing.T) {
	err := FetchArgs("订阅动态 12312315 泛式")
	if err != nil {
		return
	}
}
