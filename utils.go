package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var (
	chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
)

func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func InArrayString(val string, array []string) (bool, int) {
	for i := range array {
		if array[i] == val {
			return true, i
		}
	}
	return false, -1
}

func InArrayInt(val int, array []int) (bool, int) {
	for i := range array {
		if array[i] == val {
			return true, i
		}
	}
	return false, -1
}

func InArrayInt64(val int64, array []int64) (bool, int) {
	for i := range array {
		if array[i] == val {
			return true, i
		}
	}
	return false, -1
}

func Md5Str(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func StrToInt(str string) int {
	if val, err := strconv.Atoi(str); err != nil {
		return 0
	} else {
		return val
	}
}

func StrToInt64(str string) int64 {
	if val, err := strconv.ParseInt(str, 10, 64); err != nil {
		return 0
	} else {
		return val
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ByteToString(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfInt(condition bool, trueVal, falseVal int) int {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfInt64(condition bool, trueVal, falseVal int64) int64 {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfString(condition bool, trueVal, falseVal string) string {
	if condition {
		return trueVal
	}
	return falseVal
}

func DeepCopyByGob(dst, src interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(&buffer).Decode(dst)
}
