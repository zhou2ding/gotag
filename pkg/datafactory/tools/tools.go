package tools

import (
	"strings"
	"unicode"
)

// GetFuncName 提取函数名，按从内到外的顺序排序，级别相同的从左往右排，eg：`tag:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`，得到["funcB", "funcC", "funcA"]
func GetFuncName(tag string) []string {
	lefts := make([]int, 0) // '('的下标
	fns := make([]string, 0)
	for i, c := range tag {
		if c == '(' {
			lefts = append(lefts, i)
		} else if c == ')' {
			newestLeft := lefts[len(lefts)-1] // 最近一个'('的下标
			prevLeft := getFirstInNormalCharacterIndex(tag[:newestLeft])
			fns = append(fns, tag[prevLeft+1:newestLeft])
			lefts = lefts[:len(lefts)-1] // '('出栈
		}
	}
	return fns
}

// GetUnFuncParam 提取非函数执行结果的参数，eg：传参为`bor:funcA(funcB(b1,b2),funcC(c1,c2),a1,a2)`和funcB，则得到[b1,b2]
func GetUnFuncParam(tag, funcName string) []string {
	lefts := make([]int, 0) // '('的下标
	for i, c := range tag {
		if c == '(' {
			lefts = append(lefts, i)
		} else if c == ')' {
			newestLeft := lefts[len(lefts)-1] // 最近一个'('的下标
			prevLeft := getFirstInNormalCharacterIndex(tag[:newestLeft])
			if tag[prevLeft+1:newestLeft] == funcName {
				if subTag := tag[newestLeft+1 : i]; strings.Contains(subTag, ")") {
					if strings.LastIndex(subTag, ")") == len(subTag)-1 {
						return nil
					} else {
						return strings.Split(subTag[strings.LastIndex(subTag, ")")+2:], ",")
					}
				} else {
					return strings.Split(tag[newestLeft+1:i], ",")
				}
			}
			lefts = lefts[:len(lefts)-1] // '('出栈
		}
	}
	return nil
}

// GetForeignParam 提取调用链和从数组中查找的索引
// eg：`sum(foreign(config.value.measurePoints[A2].theoreticalX),this(dx),this(dy))`，得到 "config.value.measurePoints[pointName].theoreticalX" 和 "A2"
func GetForeignParam(tag string) ([]string, string) {
	lefts := make([]int, 0) // '('的下标
	for i, c := range tag {
		if c == '(' {
			lefts = append(lefts, i)
		} else if c == ')' {
			newestLeft := lefts[len(lefts)-1] // 最近一个'('的下标
			prevLeft := getFirstInNormalCharacterIndex(tag[:newestLeft])
			if tag[prevLeft+1:newestLeft] == "foreign" {
				subTag := tag[newestLeft+1 : i]
				var idxVal string
				if strings.Contains(subTag, "[") {
					idxVal = subTag[strings.Index(subTag, "[")+1 : strings.LastIndex(subTag, "]")]
				}
				return strings.Split(subTag, "."), idxVal
			}
			lefts = lefts[:len(lefts)-1] // '('出栈
		}
	}

	return nil, ""
}

// 获取所有目标字符在指定字符串中的索引
func getCharacterIndex(s string, target rune) []int {
	var idxes []int
	for i, v := range s {
		if v == target {
			idxes = append(idxes, i)
		}
	}
	return idxes
}

// 从右往左遍历获取第一个不是数字、字母、下划线的字符传索引
func getFirstInNormalCharacterIndex(s string) int {
	for i := len(s) - 1; i > 0; i-- {
		if s[i] != '_' && !unicode.IsDigit(rune(s[i])) && !unicode.IsLetter(rune(s[i])) {
			return i
		}
	}
	return -1
}

func Contains(s []string, target string) bool {
	for _, v := range s {
		if strings.Contains(v, target) {
			return true
		}
	}
	return false
}
