package str_util

import (
	"bytes"
)

//Input 必需。规定要填充的字符串。
//PadLength必需。规定新字符串的长度。如果该值小于原始字符串的长度，则不进行任何操作。
//PadString必选。规定供填充使用的字符串。
//PadType必选。0=填充字符串的左侧。1=填充字符串的右侧。2=填充字符串的两侧。
//StrPadPadType 0 = "left",1 = "right", 2 = both

func StrPad(Input string, PadLength int, PadString string, PadType int) string {
	var leftPad, rightPad = 0, 0
	numPadChars := PadLength - len(Input)
	if numPadChars <= 0 {
		return Input
	}
	var buffer bytes.Buffer
	buffer.WriteString(Input)
	switch PadType {
	case 0:
		leftPad = numPadChars
		rightPad = 0
	case 1:
		leftPad = 0
		rightPad = numPadChars
	case 2:
		rightPad = numPadChars / 2
		leftPad = numPadChars - rightPad
	}
	var leftBuffer bytes.Buffer
	/* 左填充：循环添加字符*/
	for i := 0; i < leftPad; i++ {
		leftBuffer.WriteString(PadString)
		if leftBuffer.Len() > leftPad {
			leftBuffer.Truncate(leftPad)
			break
		}
	}
	/* 右填充：循环添加字符串*/
	for i := 0; i < rightPad; i++ {
		buffer.WriteString(PadString)
		if buffer.Len() > PadLength {
			buffer.Truncate(PadLength)
			break
		}
	}
	leftBuffer.WriteString(buffer.String())
	return leftBuffer.String()
}

func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}

	return
}
