package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	v, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		ans := []int{}
		for i := 0; i < v; i++ {
			ans = append(ans, i+1)
		}
		for len(ans) > 0 {
			pos := rand.Intn(len(ans))
			fmt.Println(ans[pos])
			newans := []int{}
			for i := 0; i < len(ans); i++ {
				if pos != i {
					newans = append(newans, ans[i])
				}
			}
			ans = newans
		}
	}
}
