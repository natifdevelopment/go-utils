package utils

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/net/html"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func StringToInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	num, err := strconv.Atoi(fmt.Sprintf("%v", value))
	if err != nil {
		return 0
	}
	return num
}

func Float64ToString(value float64) string {
	s := strconv.FormatFloat(value, 'f', 3, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func Float64ToStringMoney(value float64) string {
	s := strconv.FormatFloat(value, 'f', 2, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func Float64FormatUS(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	parts := strings.Split(s, ".")
	intPart := parts[0]

	var out []byte
	n := len(intPart)
	for i, d := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(d))
	}

	if len(parts) == 2 {
		return string(out) + "." + parts[1]
	}
	return string(out)
}

func Float64ToStringTrunc(value float64) string {
	s := strconv.FormatFloat(value, 'f', 0, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func StringWithDefault(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func StringToBool(value string, defaultValue bool) bool {
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(fmt.Sprintf("%v", value))
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func SnakeToCamel(s string) string {
	words := strings.Split(s, "_")
	for i := range words {
		words[i] = cases.Title(language.English).String(words[i])
	}
	return strings.Join(words, "")
}

func SetFieldValue(model interface{}, values map[string]string) error {
	ref := reflect.ValueOf(model)

	for key, code := range values {
		field := ref.Elem().FieldByName(key)

		if !field.IsValid() {
			return fmt.Errorf("field %s not found in struct", key)
		}

		if !field.CanSet() {
			return fmt.Errorf("cannot set field %s value", key)
		}

		fieldValueReflect := reflect.ValueOf(code)

		if fieldValueReflect.Type() != field.Type() {
			return fmt.Errorf("field %s type mismatch", key)
		}

		field.Set(fieldValueReflect)
	}

	return nil
}

func Includes(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func SliceStrToAny(slice []string) []any {
	var anySlice []any
	for _, id := range slice {
		anySlice = append(anySlice, id)
	}
	return anySlice
}

func Contains(slice []any, value any) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func PrintStructName(data interface{}) {
	t := reflect.TypeOf(data)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Struct {
		fmt.Println("Struct name:", t.Name())
	} else {
		fmt.Println("Not a struct")
	}
}

func FormatNumber(number float64) string {
	floatStr := strconv.FormatFloat(number, 'f', 3, 64)
	parts := splitFloatString(floatStr)
	integerPartWithCommas := addCommas(parts[0])
	formattedNumber := integerPartWithCommas + "." + parts[1]
	return formattedNumber
}

func splitFloatString(floatStr string) [2]string {
	var result [2]string
	dotIndex := strings.Index(floatStr, ".")
	if dotIndex == -1 {
		result[0] = floatStr
		result[1] = ""
	} else {
		result[0] = floatStr[:dotIndex]
		result[1] = floatStr[dotIndex+1:]
	}
	return result
}

func addCommas(integerPart string) string {
	var result strings.Builder
	length := len(integerPart)
	commaIndex := length % 3

	for i := 0; i < length; i++ {
		if i == commaIndex && i != 0 {
			result.WriteByte(',')
			commaIndex += 3
		}
		result.WriteByte(integerPart[i])
	}

	return result.String()
}

func Abbreviate(input string) string {
	parts := strings.Split(input, "_")
	var abbreviation string
	for _, part := range parts {
		if part != "" {
			abbreviation += string(part[0])
		}
	}
	return abbreviation
}

func Truncate(f float64, precision int) float64 {
	pow10 := math.Pow(10, float64(precision))
	scaled := f * pow10
	truncatedScaled := math.Trunc(scaled)
	return truncatedScaled / pow10
}

func FormatNumberIntl(f float64, decimals ...int) string {
	decimalPlaces := 3
	if len(decimals) > 0 {
		decimalPlaces = decimals[0]
	}

	s := fmt.Sprintf("%.*f", decimalPlaces, f)
	parts := strings.Split(s, ".")
	integerPart := parts[0]
	fractionalPart := ""
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	var formattedIntegerPart strings.Builder
	length := len(integerPart)
	for i, c := range integerPart {
		if (length-i)%3 == 0 && i != 0 {
			formattedIntegerPart.WriteString(",")
		}
		formattedIntegerPart.WriteRune(c)
	}

	if fractionalPart != "" {
		return formattedIntegerPart.String() + "." + fractionalPart
	}
	return formattedIntegerPart.String()
}

func FormatNumberIntlTrimZero(f float64, decimals ...int) string {
	decimalPlaces := 3
	if len(decimals) > 0 {
		decimalPlaces = decimals[0]
	}

	s := fmt.Sprintf("%.*f", decimalPlaces, f)

	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	parts := strings.Split(s, ".")
	integerPart := parts[0]
	fractionalPart := ""
	if len(parts) > 1 {
		fractionalPart = parts[1]
	}

	var formattedIntegerPart strings.Builder
	length := len(integerPart)
	for i, c := range integerPart {
		if (length-i)%3 == 0 && i != 0 {
			formattedIntegerPart.WriteString(",")
		}
		formattedIntegerPart.WriteRune(c)
	}

	if fractionalPart != "" {
		return formattedIntegerPart.String() + "." + fractionalPart
	}
	return formattedIntegerPart.String()
}

func IntToRoman(num int) string {
	romanNumerals := []struct {
		value  int
		symbol string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	var result strings.Builder
	for _, numeral := range romanNumerals {
		for num >= numeral.value {
			result.WriteString(numeral.symbol)
			num -= numeral.value
		}
	}

	return result.String()
}

func ParsePercentage(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")

	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return num / 100.0, nil
}

func StripHtmlTags(input string) string {
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return ""
	}
	var b strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			b.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	out := html.UnescapeString(b.String())
	return strings.Join(strings.Fields(out), " ")
}

func TruncStr(s string, truncLen int) string {
	runes := []rune(s)
	if len(runes) > truncLen {
		return string(runes[:truncLen])
	}
	return s
}

func UniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, v := range input {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	return result
}

func AnyToFloat64(val any) float64 {
	var f float64

	switch v := val.(type) {
	case int:
		f = float64(v)
	case int64:
		f = float64(v)
	case float32:
		f = float64(v)
	case float64:
		f = v
	default:
		return 0
	}

	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0
	}

	return f
}

var angka = []string{
	"", "satu", "dua", "tiga", "empat", "lima",
	"enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas",
}

func TerbilangFloat(n float64) string {
	if n == 0 {
		return "nol"
	}

	result := ""

	if n < 0 {
		result = "minus "
		n = math.Abs(n)
	}

	intPart := int64(n)
	decPart := int64(math.Round((n - float64(intPart)) * 100))

	result += strings.TrimSpace(terbilangInt(intPart))

	if decPart > 0 {
		result += " koma"
		for _, d := range fmt.Sprintf("%02d", decPart) {
			result += " " + angka[int(d-'0')]
		}
	}

	result = strings.TrimSpace(result)
	return ucfirst(result)
}

func ucfirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func terbilangInt(n int64) string {
	switch {
	case n < 12:
		return angka[n]
	case n < 20:
		return angka[n-10] + " belas"
	case n < 100:
		return terbilangInt(n/10) + " puluh " + terbilangInt(n%10)
	case n < 200:
		return "seratus " + terbilangInt(n-100)
	case n < 1000:
		return terbilangInt(n/100) + " ratus " + terbilangInt(n%100)
	case n < 2000:
		return "seribu " + terbilangInt(n-1000)
	case n < 1000000:
		return terbilangInt(n/1000) + " ribu " + terbilangInt(n%1000)
	case n < 1000000000:
		return terbilangInt(n/1000000) + " juta " + terbilangInt(n%1000000)
	case n < 1000000000000:
		return terbilangInt(n/1000000000) + " miliar " + terbilangInt(n%1000000000)
	default:
		return terbilangInt(n/1000000000000) + " triliun " + terbilangInt(n%1000000000000)
	}
}

var _ = errors.New
