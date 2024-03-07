package test

import (
	"fmt"
	"testing"
)

func TestStringSlices(t *testing.T) {
	str := "[1,2]"
	fmt.Println(str, str[0])
	//strings.TrimLeft(strings.TrimRight(str[0], "]"), "[")
	//s := str[1 : len(str)-1]
	//slices := strings.Split(s, ",")
	//for _, slice := range slices {
	//	fmt.Println(slice, reflect.TypeOf(slice))
	//}

}
