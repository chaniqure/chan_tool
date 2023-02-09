package utils

import (
	"chan_tool/strs"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"time"
)

// 判断是否存在
func Exists(slice []interface{}, v interface{}) bool {
	for _, value := range slice {
		if reflect.DeepEqual(value, v) {
			return true
		}
	}

	return false
}

func ExistsInt64(slice []int64, v int64) bool {
	for _, value := range slice {
		if value == v {
			return true
		}
	}

	return false
}

func ExistsString(slice []string, v string) bool {
	for _, value := range slice {
		if value == v {
			return true
		}
	}

	return false
}

func GetIndex(slice []string, v string) int {
	for index, value := range slice {
		if value == v {
			return index
		}
	}

	return -1
}

// 移除全部int64
// [1,2,3,2,1] 移除 1 变为 [2,3,2]
func RemoveInt64(slice []int64, v int64) (newSlice []int64) {
	newSlice = make([]int64, 0)
	newSlice = append(newSlice, slice...)

	for i := 0; i < len(newSlice); {
		if newSlice[i] == v {
			tempSlice := make([]int64, 0)
			tempSlice = append(tempSlice, newSlice[:i]...)
			tempSlice = append(tempSlice, newSlice[i+1:]...)
			newSlice = tempSlice
		} else {
			i++
		}
	}

	return newSlice
}

// 移除全部string
func RemoveString(slice []string, v string) (newSlice []string) {
	newSlice = make([]string, 0)
	newSlice = append(newSlice, slice...)

	for i := 0; i < len(newSlice); {
		if newSlice[i] == v {
			tempSlice := make([]string, 0)
			tempSlice = append(tempSlice, newSlice[:i]...)
			tempSlice = append(tempSlice, newSlice[i+1:]...)
			newSlice = tempSlice
		} else {
			i++
		}
	}

	return newSlice
}

func ToSplitStringSlice(s []string, split string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(s), "[]"), " ", split, -1)
}

// 倒序排列
func RSortInt64(slice []int64) (newSlice []int64) {
	newSlice = append(make([]int64, 0), slice...)

	length := len(newSlice)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if newSlice[i] < newSlice[j] {
				newSlice[i], newSlice[j] = newSlice[j], newSlice[i]
			}
		}
	}

	return newSlice
}

func RSortString(slice []string) (newSlice []string) {
	return ReverseString(SortString(slice))
}

// 正序排列
func SortInt64(slice []int64) (newSlice []int64) {
	newSlice = append(make([]int64, 0), slice...)

	length := len(newSlice)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if newSlice[i] > newSlice[j] {
				newSlice[i], newSlice[j] = newSlice[j], newSlice[i]
			}
		}
	}

	return newSlice
}

func SortString(slice []string) (newSlice []string) {
	newSlice = append(make([]string, 0), slice...)

	sort.Strings(newSlice)

	return newSlice
}

// 打乱数组
func RandomInt64(slice []int64) (newSlice []int64) {
	newSlice = append(make([]int64, 0), slice...)

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(newSlice) - 1; i > 0; i-- {
		r := rand.Intn(i + 1)
		newSlice[r], newSlice[i] = newSlice[i], newSlice[r]
	}

	return newSlice
}

func RandomString(slice []string) (newSlice []string) {
	newSlice = append(make([]string, 0), slice...)

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(newSlice) - 1; i > 0; i-- {
		r := rand.Intn(i + 1)
		newSlice[r], newSlice[i] = newSlice[i], newSlice[r]
	}

	return newSlice
}

// 颠倒数组
// [1,2,3,5] 变为 [5,3,2,1]
func ReverseInt64(slice []int64) (newSlice []int64) {
	newSlice = make([]int64, len(slice))
	for i := 0; i < len(slice); i++ {
		newSlice[i] = slice[len(slice)-1-i]
	}

	return newSlice
}

func ReverseString(slice []string) (newSlice []string) {
	newSlice = make([]string, len(slice))
	for i := 0; i < len(slice); i++ {
		newSlice[i] = slice[len(slice)-1-i]
	}

	return newSlice
}

func StringToInt64Slice(s string, spliter string) []int64 {
	ids := strings.Split(s, spliter)
	idList := make([]int64, 0)
	for _, i := range ids {
		id := int64(strs.Str2Int(i))
		idList = append(idList, id)
	}
	return idList
}
func StringToJsonNumber(s string, spliter string) []json.Number {
	ids := strings.Split(s, spliter)
	idList := make([]json.Number, 0)
	for _, i := range ids {
		id := json.Number(rune(strs.Str2Int(i)))
		idList = append(idList, id)
	}
	return idList
}

func JsonNumbersToInt64Slice(numbers []json.Number) []int64 {
	list := make([]int64, 0)
	for _, number := range numbers {
		intNumber, _ := number.Int64()
		list = append(list, intNumber)
	}
	return list
}
