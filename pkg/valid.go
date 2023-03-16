package pkg

import (
	"fmt"
	"time"
)

func CheckValidData(check_in, check_out string) bool {
	in, err := time.Parse("2006-01-02", check_in)
	if err != nil {
		fmt.Println("error data:", check_in)
		return false
	}
	out, err := time.Parse("2006-01-02", check_out)
	if err != nil {
		fmt.Println("error data:", check_out)
		return false
	}
	if in.Before(time.Now()) || out.Before(in) {
		fmt.Println("error data in before time.now")
		return false
	}
	if in.Year() > time.Now().Year()+2 {
		fmt.Println("in a bigger than time now for 2 years")
		return false
	}
	return true
}
