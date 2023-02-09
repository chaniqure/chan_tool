package utils

import (
	"bytes"
	"chan_tool/strs"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// 结构体转为json
func Struct2Json(obj interface{}) string {
	str, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Struct2Json]转换异常: %v", err))
	}
	return string(str)
}

// json转为结构体
func Json2Struct(str string, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal([]byte(str), obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}
func Byte2Struct(bytes []byte, obj interface{}) {
	// 将json转为结构体
	err := json.Unmarshal(bytes, obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// json interface转为结构体
func JsonI2Struct(str interface{}, obj interface{}) {
	// 将json interface转为string
	jsonStr, _ := str.(string)
	Json2Struct(jsonStr, obj)
}

// 结构体转结构体, json为中间桥梁, struct2必须以指针方式传递, 否则可能获取到空数据
func Struct2StructByJson(struct1 interface{}, struct2 interface{}) {
	// 转换为响应结构体, 隐藏部分字段
	jsonStr := Struct2Json(struct1)
	Json2Struct(jsonStr, struct2)
}

// 结构体转结构体，不会丢失精度的, json为中间桥梁, struct2必须以指针方式传递, 否则可能获取到空数据
func Struct2StructByJsonHighPrecision(struct1 interface{}, struct2 interface{}) {
	// 转换为响应结构体, 隐藏部分字段
	jsonStr := Struct2Json(struct1)
	Json2StructHighPrecision(jsonStr, struct2)
}

func Json2StructHighPrecision(str string, obj interface{}) {
	d := json.NewDecoder(bytes.NewBuffer([]byte(str)))
	d.UseNumber()
	err := d.Decode(obj)
	if err != nil {
		fmt.Println(fmt.Sprintf("[Json2Struct]转换异常: %v", err))
	}
}

// 两结构体比对不同的字段, 不同时将取newStruct中的字段返回, json为中间桥梁
func CompareDifferenceStructByJson(oldStruct interface{}, newStruct interface{}, update *map[string]interface{}) {
	// 通过json先将其转为map集合
	m1 := make(map[string]interface{}, 0)
	m2 := make(map[string]interface{}, 0)
	m3 := make(map[string]interface{}, 0)
	Struct2StructByJson(newStruct, &m1)
	Struct2StructByJson(oldStruct, &m2)
	for k1, v1 := range m1 {
		for k2, v2 := range m2 {
			switch v1.(type) {
			// 复杂结构不做对比
			case map[string]interface{}:
				continue
			}
			rv := reflect.ValueOf(v1)
			// 值类型必须有效
			if rv.Kind() != reflect.Invalid {
				// key相同, 值不同
				if k1 == k2 && v1 != v2 {
					t := reflect.TypeOf(oldStruct)
					key := strs.CamelCase(k1)
					var fieldType reflect.Type
					oldStructV := reflect.ValueOf(oldStruct)
					// map与结构体取值方式不同
					if oldStructV.Kind() == reflect.Map {
						mapV := oldStructV.MapIndex(reflect.ValueOf(k1))
						if !mapV.IsValid() {
							break
						}
						fieldType = mapV.Type()
					} else if oldStructV.Kind() == reflect.Struct {
						structField, ok := t.FieldByName(key)
						if !ok {
							break
						}
						fieldType = structField.Type
					} else {
						// oldStruct类型不对, 直接跳过不处理
						break
					}
					// 取到结构体对应字段
					realT := fieldType
					// 指针类型需要剥掉一层获取真实类型
					if fieldType.Kind() == reflect.Ptr {
						realT = fieldType.Elem()
					}
					// 获得元素
					e := reflect.New(realT).Elem()
					// 不同类型不一定可以强制转换
					switch e.Interface().(type) {
					default:
						// 强制转换rv赋值给e
						e.Set(rv.Convert(realT))
						m3[k1] = e.Interface()
					}
					break
				}
			}
		}
	}
	*update = m3
}

// 两结构体比对不同的字段, 将key转为蛇形
func CompareDifferenceStruct2SnakeKeyByJson(oldStruct interface{}, newStruct interface{}, update *map[string]interface{}) {
	m1 := make(map[string]interface{}, 0)
	m2 := make(map[string]interface{}, 0)
	CompareDifferenceStructByJson(oldStruct, newStruct, &m1)
	for key, item := range m1 {
		m2[strs.SnakeCase(key)] = item
	}
	*update = m2
}

var keyMatchRegex = regexp.MustCompile(`"(\w+)":`)
var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)

/*************************************** 下划线json ***************************************/

func JsonSnakeCase(value interface{}) ([]byte, error) {
	marshalled, err := json.Marshal(value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match,
				[]byte(`${1}_${2}`),
			))
		},
	)
	return converted, err
}

/*************************************** 驼峰json ***************************************/

func JsonCamelCase(value interface{}) ([]byte, error) {
	marshalled, err := json.Marshal(value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			matchStr := string(match)
			key := matchStr[1 : len(matchStr)-2]
			resKey := Lcfirst(Case2Camel(key))
			return []byte(`"` + resKey + `":`)
		},
	)
	return converted, err
}

/*************************************** 其他方法 ***************************************/

// 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}

// 下划线写法转为驼峰写法
func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

// 首字母大写
func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}
