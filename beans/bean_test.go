package beans

import (
	"fmt"
	"middle-center/models/entity"
	"testing"
)

func TestGetBeansProperties(t *testing.T) {
	var roles = []entity.Role{
		{
			Code: "111",
			Name: "一一一",
		},
		{
			Code: "222",
			Name: "二二二",
		},
	}
	properties, err := GetBeansProperties(roles[0], "code", "name")
	if err != nil {
		return
	}
	fmt.Println(properties)
}
