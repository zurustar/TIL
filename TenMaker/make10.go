package main

import (
	"fmt"
	"go/constant"
	"go/token"
	"go/types"
	"os"
	"strconv"
)

func check10(f string) bool {
	fs := token.NewFileSet()
	tv, err := types.Eval(fs, nil, token.NoPos, f)
	if err == nil {
		if tv.Value != nil {
			v, ok := constant.Float64Val(tv.Value)
			if ok {
				if (9.999 < v) && (v < 10.001) {
					fmt.Println(f, "=", v)
					return true
				}
			}
		}
	}
	return false
}

func check(s string) {
	ls := []string{"", "(", "((", "((("}
	os := []string{"+", "-", "*", "/"}
	rs := []string{"", ")", "))", ")))"}
	for _, l0 := range ls {
		for _, o0 := range os {
			for _, l1 := range ls[:3] {
				for _, r0 := range rs[:2] {
					for _, o1 := range os {
						for _, l2 := range ls[:2] {
							for _, r1 := range rs[:3] {
								for _, o2 := range os {
									for _, r2 := range rs {
										f := l0 + s[:1] + ".0" + o0 + l1 + s[1:2] + ".0" + r0 + o1 + l2 + s[2:3] + ".0" + r1 + o2 + s[3:] + ".0" + r2
										if check10(f) {
											return
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	if len(os.Args) >= 2 {
		s := os.Args[1]
		check(s)
	} else {
		for i := 1000; i < 10000; i++ {
			check(strconv.Itoa(i))
		}
	}
}
