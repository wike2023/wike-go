package jwt

import (
	"fmt"
	"testing"
	"time"
)

type Info struct {
	Name string
}

func TestJWT(t *testing.T) {
	core := Info{
		Name: "123",
	}
	token, err := Create[Info](core, time.Second*0)
	fmt.Println(token, err)
	info, err := Parse[Info](token)
	fmt.Println(info, err)
}
