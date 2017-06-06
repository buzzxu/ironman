package ironman

import (
	"fmt"
	"testing"
)

func TestComparePasswordAndStr(t *testing.T) {

	//token := GeneratePassword("111111")
	//fmt.Printf("%s \n",token)
	//$2a$15$NflzmYbtZtRiJaIR0pnY.uq
	fmt.Printf("%t \n", ComparePasswordAndStr("111111", "$2a$15$u7OC5fL8xNQQSWP7MoJCyerLqx.H73G1JFUuzw6ztRELiwAGDVKtu"))
	fmt.Printf("%d \n", len("$2a$10$/QnFAsA7.9LrdgL.cjlscOi1eeJ4ODsB9l7..Ooy4ErLS1oOgRubS"))
}

func TestStruct2Map(t *testing.T) {
	wx := UserWx{1, "2323", "323"}
	fmt.Println(Struct2Map(wx))
}
