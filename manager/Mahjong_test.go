package manager

import (
	"fmt"
	"testing"
)

func Test_CheckHu(t *testing.T){
	fmt.Println(getCardSlice(&[]int{17,16,26,26,25}))
	fmt.Println(checkHu([]int{17,16,26,26,25},0))
}
