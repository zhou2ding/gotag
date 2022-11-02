package datafactory

import (
	"bufio"
	"gotag/pkg/l"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var typeList = []string{"int", "bool", "float32", "float64", "string", "struct"}

func GetArrayCnt(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		l.GetLogger().Error("open file failed", zap.Error(err))
		return 0, err
	}
	reader := bufio.NewReader(f)
	var cnt int
	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			l.GetLogger().Error("read line failed", zap.Error(err))
			return 0, err
		}
		line := string(data)
		if strings.Contains(line, "]struct") {
			left := strings.Index(line, "[")
			right := strings.Index(line, "]")
			cnt, err = strconv.Atoi(line[left+1 : right])
			if err != nil {
				return 0, err
			}
		}
	}
	return cnt, nil
}

func GetFiledTagAndFields(path string) (map[string]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		l.GetLogger().Error("open file failed", zap.Error(err))
		return nil, nil, err
	}
	reader := bufio.NewReader(f)
	fieldTags := make(map[string]string) // key是字段名，value是字段标签中的函数
	fields := make([]string, 0)          // 所有字段名从上往下排好序的数组
	for {
		data, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			l.GetLogger().Error("read line failed", zap.Error(err))
			return nil, nil, err
		}
		line := string(data)
		filed := getField(line)
		if strings.Contains(line, "`") {
			tag := line[strings.Index(line, "`"):]
			fieldTags[filed] = tag
		}
		if isField(line) {
			fields = append(fields, filed)
		}
	}

	return fieldTags, fields, nil
}

func getField(s string) string {
	var field []byte
	s = strings.TrimSpace(s)
	for _, c := range s {
		if !unicode.IsSpace(c) {
			field = append(field, byte(c))
		} else if len(field) > 0 {
			break
		}
	}
	return string(field)
}

func isField(s string) bool {
	for _, t := range typeList {
		if strings.Contains(s, t) {
			return true
		}
	}
	return false
}
