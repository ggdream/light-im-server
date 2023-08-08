package typec

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"unsafe"
)

func NumToStr(num any) string {
	return fmt.Sprintf("%v", num)
}

func StrToInt(str string) int {
	i, _ := strconv.ParseInt(str, 10, 64)
	return int(i)
}

func StrToInt64(str string) int64 {
	i, _ := strconv.ParseInt(str, 10, 64)
	return i
}

func StrToUint(str string) uint {
	u, _ := strconv.ParseUint(str, 10, 64)
	return uint(u)
}

func StrToUint8(str string) uint8 {
	u, _ := strconv.ParseUint(str, 10, 64)
	return uint8(u)
}

func FloatToStr(num float64, prec int) string {
	return strconv.FormatFloat(num, 'f', prec, 64)
}
func StrToFloat(str string) float64 {
	f, _ := strconv.ParseFloat(str, 64)
	return f
}

func MapToUrlValuesToJson(data map[string]any) string {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, fmt.Sprintf("%v", v))
	}

	return values.Encode()
}

func MapToJson(data map[string]any) string {
	res, _ := json.Marshal(data)

	return unsafe.String(&res[0], len(res))
}

func JsonToMap(data string) map[string]any {
	var ret map[string]any
	_ = json.Unmarshal(unsafe.Slice(unsafe.StringData(data), len(data)), &ret)

	return ret
}

func Three[T any](fn func() bool, a, b T) T {
	if fn() {
		return a
	} else {
		return b
	}
}

func BoolSliceToUint8Slice(data []bool) []uint8 {
	ret := make([]uint8, 0, len(data))
	for _, v := range data {
		if v {
			ret = append(ret, 1)
		} else {
			ret = append(ret, 0)
		}
	}

	return ret
}

func LastBytes(id uint, num int) string {
	str := NumToStr(id)
	return str[len(str)-num:]
}
