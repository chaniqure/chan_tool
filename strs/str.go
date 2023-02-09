package strs

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

// 是否空字符串
func StrIsEmpty(str string) bool {
	return str == "null" || strings.TrimSpace(str) == ""
}

func Int64SliceToStrSlice(arr []int64) []string {
	strArr := make([]string, 0)
	for _, n := range arr {
		strArr = append(strArr, Int64ToStr(n))
	}
	return strArr
}

// 字符串转uint数组, 默认逗号分割
func Str2UintArr(str string) (ids []uint) {
	idArr := strings.Split(str, ",")
	for _, v := range idArr {
		ids = append(ids, Str2Uint(v))
	}
	return
}

// 字符串转int
func Str2Int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return num
}

func Str2Int64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func StrListToInt64List(strList []string) []int64 {
	result := make([]int64, len(strList))
	for _, s := range strList {
		result = append(result, Str2Int64(s))
	}
	return result
}

// 字符串转uint
func Str2Uint(str string) uint {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint(num)
}

// 字符串转uint
func Str2Uint32(str string) uint32 {
	num, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(num)
}

// 字符串转uint
func Str2Bool(str string) bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}
	return b
}

// 字符串转float64
func Str2Float64(str string) float64 {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return num
}

// 字符串转uint数组, 默认逗号分割
func UintArr2IntArr(arr []uint) (newArr []int) {
	for _, v := range arr {
		newArr = append(newArr, int(v))
	}
	return
}

var (
	camelRe = regexp.MustCompile("(_)([a-zA-Z]+)")
	snakeRe = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// 字符串转为驼峰
func CamelCase(str string) string {
	camel := camelRe.ReplaceAllString(str, " $2")
	camel = strings.Title(camel)
	camel = strings.Replace(camel, " ", "", -1)
	return camel
}

// 字符串转为驼峰(首字母小写)
func CamelCaseLowerFirst(str string) string {
	camel := CamelCase(str)
	for i, v := range camel {
		return string(unicode.ToLower(v)) + camel[i+1:]
	}
	return camel
}

// 驼峰式写法转为下划线蛇形写法
func SnakeCase(str string) string {
	snake := snakeRe.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

// 加密base64字符串
func EncodeStr2Base64(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// 解密base64字符串
func DecodeStrFromBase64(str string) string {
	decodeBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(decodeBytes)
}
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// convert string to specify type

type StrTo string

func (f *StrTo) Set(v string) {
	if v != "" {
		*f = StrTo(v)
	} else {
		f.Clear()
	}
}

func (f *StrTo) Clear() {
	*f = StrTo(0x1E)
}

func (f StrTo) Exist() bool {
	return string(f) != string(0x1E)
}

func (f StrTo) Bool() (bool, error) {
	if f == "on" {
		return true, nil
	}
	return strconv.ParseBool(f.String())
}

func (f StrTo) Float32() (float32, error) {
	v, err := strconv.ParseFloat(f.String(), 32)
	return float32(v), err
}

func (f StrTo) Float64() (float64, error) {
	return strconv.ParseFloat(f.String(), 64)
}

func (f StrTo) Int() (int, error) {
	v, err := strconv.ParseInt(f.String(), 10, 32)
	return int(v), err
}

func (f StrTo) Int8() (int8, error) {
	v, err := strconv.ParseInt(f.String(), 10, 8)
	return int8(v), err
}

func (f StrTo) Int16() (int16, error) {
	v, err := strconv.ParseInt(f.String(), 10, 16)
	return int16(v), err
}

func (f StrTo) Int32() (int32, error) {
	v, err := strconv.ParseInt(f.String(), 10, 32)
	return int32(v), err
}

func (f StrTo) Int64() (int64, error) {
	v, err := strconv.ParseInt(f.String(), 10, 64)
	return int64(v), err
}

func (f StrTo) Uint() (uint, error) {
	v, err := strconv.ParseUint(f.String(), 10, 32)
	return uint(v), err
}

func (f StrTo) Uint8() (uint8, error) {
	v, err := strconv.ParseUint(f.String(), 10, 8)
	return uint8(v), err
}

func (f StrTo) Uint16() (uint16, error) {
	v, err := strconv.ParseUint(f.String(), 10, 16)
	return uint16(v), err
}

func (f StrTo) Uint32() (uint32, error) {
	v, err := strconv.ParseUint(f.String(), 10, 32)
	return uint32(v), err
}

func (f StrTo) Uint64() (uint64, error) {
	v, err := strconv.ParseUint(f.String(), 10, 64)
	return uint64(v), err
}

func (f StrTo) String() string {
	if f.Exist() {
		return string(f)
	}
	return ""
}

// convert any type to string
func ToStr(value interface{}, args ...int) (s string) {
	switch v := value.(type) {
	case bool:
		s = strconv.FormatBool(v)
	case float32:
		s = strconv.FormatFloat(float64(v), 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 32))
	case float64:
		s = strconv.FormatFloat(v, 'f', argInt(args).Get(0, -1), argInt(args).Get(1, 64))
	case int:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int8:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int16:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int32:
		s = strconv.FormatInt(int64(v), argInt(args).Get(0, 10))
	case int64:
		s = strconv.FormatInt(v, argInt(args).Get(0, 10))
	case uint:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint8:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint16:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint32:
		s = strconv.FormatUint(uint64(v), argInt(args).Get(0, 10))
	case uint64:
		s = strconv.FormatUint(v, argInt(args).Get(0, 10))
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		s = fmt.Sprintf("%v", v)
	}
	return s
}

func ToInt(value interface{}) (d int64, err error) {
	s := ToStr(value)
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

// convert any numeric value to int64
func ToInt64(value interface{}) (d int64, err error) {
	val := reflect.ValueOf(value)
	switch value.(type) {
	case int, int8, int16, int32, int64:
		d = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		d = int64(val.Uint())
	default:
		err = fmt.Errorf("ToInt64 need numeric not `%T`", value)
	}
	return
}

type argString []string

func (a argString) Get(i int, args ...string) (r string) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

type argInt []int

func (a argInt) Get(i int, args ...int) (r int) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

type argAny []interface{}

func (a argAny) Get(i int, args ...interface{}) (r interface{}) {
	if i >= 0 && i < len(a) {
		r = a[i]
	}
	if len(args) > 0 {
		r = args[0]
	}
	return
}

func formatMapToXML(req map[string]string) (buf []byte, err error) {
	bodyBuf := textBufferPool.Get().(*bytes.Buffer)
	bodyBuf.Reset()
	defer textBufferPool.Put(bodyBuf)

	if bodyBuf == nil {
		return []byte{}, errors.New("nil xmlWriter")
	}

	if _, err = io.WriteString(bodyBuf, "<xml>"); err != nil {
		return
	}

	for k, v := range req {
		if _, err = io.WriteString(bodyBuf, "<"+k+">"); err != nil {
			return
		}
		if err = xml.EscapeText(bodyBuf, []byte(v)); err != nil {
			return
		}
		if _, err = io.WriteString(bodyBuf, "</"+k+">"); err != nil {
			return
		}
	}

	if _, err = io.WriteString(bodyBuf, "</xml>"); err != nil {
		return
	}

	return bodyBuf.Bytes(), nil
}

var textBufferPool = sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(make([]byte, 0, 16<<10)) // 16KB
	},
}

func GetMD5(args ...string) string {
	var str string
	for _, s := range args {
		str += s
	}
	value := md5.Sum([]byte(str))
	rs := []rune(fmt.Sprintf("%x", value))
	return string(rs)
}

func GetSign(args ...string) string {
	salt := "Yexhj8agldf3yaexuda7da"
	var str string
	for _, s := range args {
		str += s
	}
	str += salt
	value := md5.Sum([]byte(str))
	rs := []rune(fmt.Sprintf("%x", value))
	return string(rs)
}

// res:原字符串，sep替换的字符串，idx开始替换的位置
// 如果替换的字符串超出，就加在原字符串后面
// idx从0开始
func ReplaceString(res, sep string, idx int) string {
	sepLen := len(sep)
	if sepLen == 0 {
		return res
	}

	resLen := len(res)
	if idx > resLen-1 {
		return res + sep
	}

	allLen := resLen
	if sepLen > resLen-idx {
		allLen = idx + sepLen
	}

	buf := bytes.Buffer{}
	sepIdx := 0
	for i := 0; i < allLen; i++ {
		if i < idx {
			buf.WriteByte(res[i])
		} else {
			if sepIdx < sepLen {
				buf.WriteByte(sep[sepIdx])
			} else {
				buf.WriteByte(res[i])
			}
			sepIdx++
		}
	}
	return buf.String()
}

func ConvFormMapToString(mData map[string]string) string {
	formBuf := bytes.Buffer{}
	l := len(mData)
	i := 0
	for k, v := range mData {
		formBuf.WriteString(k)
		formBuf.WriteString("=")
		formBuf.WriteString(v)
		if i < l {
			formBuf.WriteString("&")
			i++
		}
	}
	return string(formBuf.Bytes())
}

// 检查是否包含,必须全部包含
func CheckContains(s string, subArr ...string) bool {
	for _, sub := range subArr {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// 有任何一个包含就返回true
func CheckContainsAny(s string, subArr ...string) bool {
	for _, sub := range subArr {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
func StrToUInt(num string) uint {
	return uint(StrToInt(num))
}

func StrToUInt8(num string) uint8 {
	return uint8(StrToInt(num))
}

func StrToUInt16(num string) uint16 {
	return uint16(StrToInt(num))
}

func StrToUInt32(num string) uint32 {
	return uint32(StrToInt(num))
}

func StrToUInt64(num string) uint64 {
	return uint64(StrToInt(num))
}

// 字符串转int
func StrToInt(num string) int {
	result, err := strconv.Atoi(num)

	if err != nil {
		float, err := strconv.ParseFloat(num, 64)
		if err == nil {
			return int(float)
		}

		return 0
	}

	return result
}

// 字符串转int32
func StrToInt8(num string) int8 {
	result, err := strconv.ParseInt(num, 10, 8)
	if err != nil {
		float, err := strconv.ParseFloat(num, 8)
		if err == nil {
			return int8(float)
		}

		return 0
	}

	return int8(result)
}

func StrToInt16(num string) int16 {
	result, err := strconv.ParseInt(num, 10, 16)
	if err != nil {
		float, err := strconv.ParseFloat(num, 16)
		if err == nil {
			return int16(float)
		}

		return 0
	}

	return int16(result)
}

func StrToInt32(num string) int32 {
	result, err := strconv.ParseInt(num, 10, 32)
	if err != nil {
		float, err := strconv.ParseFloat(num, 32)
		if err == nil {
			return int32(float)
		}

		return 0
	}

	return int32(result)
}

// 字符串转int64
func StrToInt64(num string) int64 {
	result, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		float, err := strconv.ParseFloat(num, 64)
		if err == nil {
			return int64(float)
		}

		return 0
	}

	return result
}

// int到字符串
func IntToStr(num int) string {
	return strconv.Itoa(num)
}

// int64到字符串
func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

// 任意类型转字符串
func InterfaceToByte(v interface{}) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

// 任意类型转字符串
func InterfaceToStr(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

func InterfacesToStrs(v []interface{}) []string {
	result := make([]string, 0)
	for _, i := range v {
		result = append(result, InterfaceToStr(i))
	}
	return result
}

func InterfacesToInt64s(v []interface{}) []int64 {
	result := make([]int64, 0)
	for _, i := range v {
		result = append(result, InterfaceToInt64(i))
	}
	return result
}

// 给list里面的str数组都拼接前后缀
func BatchAppend(prefix, postfix string, strs []string) []string {
	resultList := make([]string, 0)
	for _, str := range strs {
		resultList = append(resultList, prefix+str+postfix)
	}
	return resultList
}

//任意类型转数字

func InterfaceToUInt(v interface{}) uint {
	return StrToUInt(fmt.Sprintf("%v", v))
}

func InterfaceToUInt8(v interface{}) uint8 {
	return StrToUInt8(fmt.Sprintf("%v", v))
}

func InterfaceToUInt16(v interface{}) uint16 {
	return StrToUInt16(fmt.Sprintf("%v", v))
}

func InterfaceToUInt32(v interface{}) uint32 {
	return StrToUInt32(fmt.Sprintf("%v", v))
}

func InterfaceToUInt64(v interface{}) uint64 {
	return StrToUInt64(fmt.Sprintf("%v", v))
}

func InterfaceToInt(v interface{}) int {
	return StrToInt(fmt.Sprintf("%v", v))
}

func InterfaceToInt8(v interface{}) int8 {
	return StrToInt8(fmt.Sprintf("%v", v))
}

func InterfaceToInt16(v interface{}) int16 {
	return StrToInt16(fmt.Sprintf("%v", v))
}

func InterfaceToInt32(v interface{}) int32 {
	return StrToInt32(fmt.Sprintf("%v", v))
}

func InterfaceToInt64(v interface{}) int64 {
	return StrToInt64(fmt.Sprintf("%v", v))
}

// 转float64
func InterfaceToFloat32(v interface{}) float32 {
	f, e := strconv.ParseFloat(fmt.Sprintf("%v", v), 32)

	if e != nil {
		return 0
	}

	return float32(f)
}

func InterfaceToFloat64(v interface{}) float64 {
	f, e := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)

	if e != nil {
		return 0
	}

	return f
}

// 转big.int
func InterfaceToBigInt(v interface{}) *big.Int {
	switch v.(type) {
	case float32, float64:
		return big.NewInt(InterfaceToInt64(v))
	case int, int32, int64:
		return big.NewInt(InterfaceToInt64(v))
	default:
		num := new(big.Int)
		num.SetString(fmt.Sprintf("%v", v), 10)
		return num
	}

	return nil
}

// 16进制字符串转int64
func HexToInt64(hex string) int64 {
	if strings.HasPrefix(hex, "0x") {
		hex = strings.TrimLeft(hex, "0x")
	}

	result, err := strconv.ParseInt(hex, 16, 64)

	if err != nil {
		return 0
	}

	return result
}

// 16进制字符串转BigInt
func HexToBigInt(hex string) *big.Int {
	if strings.HasPrefix(hex, "0x") {
		hex = strings.TrimLeft(hex, "0x")
	}

	num := new(big.Int)
	num.SetString(hex, 16)

	return num
}

// 产生指定长度的数字随机数
func RandomNumStr(len int) (randomNum string) {
	var buffer bytes.Buffer
	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len; i++ {
		buffer.WriteString(fmt.Sprintf("%d", rand1.Int()%10))
	}

	return buffer.String()
}

// 产生指定长度的16进制字符串
func RandomHexStr(len int) (str string) {
	var hexArr = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	var buffer bytes.Buffer

	rand1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < len; i++ {
		buffer.WriteString(hexArr[rand1.Int()%16])
	}

	return buffer.String()
}

// 将字符串填充到指定位数
func FillToLen(str string, length int) string {
	return FillToLenByChar(str, length, "0")
}

func FillToLenByChar(str string, length int, char string) string {
	fillStr := str

	for i := 0; i < length-len(str); i++ {
		fillStr = fmt.Sprintf("%s%s", char, fillStr)
	}

	return fillStr
}

// 隐藏字符串以*代替
// 包含fromIndex，但是不包含toIndex 和 slice的规则一致
func HidStr(str string, prefixLen int, suffixLen int) string {
	var hidStrCount = 0
	var buffer bytes.Buffer
	for i := 0; i < len(str); i++ {
		if i >= prefixLen && i < len(str)-suffixLen {
			if hidStrCount < 4 {
				hidStrCount++
				buffer.WriteString("*")
			}
		} else {
			buffer.WriteString(str[i : i+1])
		}
	}

	return buffer.String()
}

// 下划线转大驼峰
func UnderscoreToCamel(name string) string {

	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
				vv[0] -= 32
			}
			s += string(vv)
		}
	}

	return s
}

/*
*
struct 转 map
*/
func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		data[name] = v.Field(i).Interface()
	}
	return data
}

var hzRegexp = regexp.MustCompile("^[\u4e00-\u9fa5]$")

// 去除字符串中的中文字符
func StrFilterChinese(src string) string {

	str := ""
	for _, c := range src {
		if !hzRegexp.MatchString(string(c)) {
			str += string(c)
		}
	}

	return str
}

// 提取字符串中的数字 支持提取 浮点数
func StrFilterNum(src string) string {

	str := ""
	j := 0
	a := ""
	for i, c := range src {
		_, err := strconv.Atoi(string(c))
		// 是数字
		if err == nil {
			str += string(c)
			j = i
		} else {
			//非数字，则判断是否是. 小数点，如果j 即上一行是数字，则将.加入
			if string(c) == "." && a != "." {
				if i-j == 1 && a != "." {
					str += string(c)
					j = i
					a = string(c)
				}
			}
		}
	}

	return str
}

func StrListIntoStrPar(strList []string, sep string) string {
	return "(" + strings.Join(strList, sep) + ")"
}
