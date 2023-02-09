package regex_util

import (
	"regexp"
)

//手机号
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^[1]([3-9])[0-9]{9}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

//电子邮件
func VerifyEmailFormat(email string) bool {
	regular := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(regular)
	return reg.MatchString(email)
}

// IP地址
func VerifyIPFormat(email string) bool {
	regular := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	reg := regexp.MustCompile(regular)
	return reg.MatchString(email)
}

//  身份证号码
func VerifyCardIDFormat(card string) bool {
	//身份证号码
	//rePersonCode = `<span>身份号码：[\s]?(([1-6]\d{5})(\d{4})(\d{2})(\d{2})(\d{4}))</span>`
	//rePersonCode = `<span>身份号码：[\s]?(([1-6]\d{5})((19\d{2})|(20\d{2}))(\d{2})(\d{2})(\d{4}))</span>`
	rePersonCode := `^?(([1-6]\d{5})(19\d{2}|20\d{2})(0[1-9]|1[012])(0[1-9]|[12]\d|3[01])(\d{3}[\dxX]))$`

	//身份证号码:  空格也需要表示出来, 直接打空格是不起作用的  \s表示空字符  [\s]任意空字符 ?出现0次或者次
	//([1-6]\d{5})(19\d{2}|20\d{2})(0[1-9]|1[012])(0[1-9]|[12]\d|3[01])(\d{3}[\dxX])
	// ([1-6]\d{5})  为一组   第一位表示地区  例如 华中  华南等  数字为1-6  [1-6]表示任意一个  \d{5} 数字共5位
	//(19\d{2}|20\d{2})  为一组 年的前两位要么为19  要么为20 后两位为\d{2}数字出现两次
	//(0[1-9]|1[012])  为一组  月第一位为0或者为1  为0的时候 后边为[1-9]任意一个数字  为1的时候  后边为[012]任意一个数字
	//(0[1-9]|[12]\d|3[01]) 为一组 日第一位要么为0要么为1或者2或者3 为0的时候第2位[1-9]任意一个数字  第1位为1或者2的时候 第2位为\d任意一个数字相当于[0-9]
	//第一位为3的时候 第2位为[0-1]任意一个数字  这里不处理2月28天或者29天的情况

	//regular := `^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{2}$`

	reg := regexp.MustCompile(rePersonCode)
	return reg.MatchString(card)
}
