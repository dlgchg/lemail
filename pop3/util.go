package pop3

import (
	"bytes"
	"encoding/base64"
	"github.com/axgle/mahonia"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"mime/quotedprintable"
	"regexp"
	"strings"
)

func IsOK(s string) bool {
	if len(s) > 0 {
		if strings.Fields(s)[0] != OK {
			return false
		}
	}
	return true
}

func IsErr(s string) bool {
	if strings.Fields(s)[0] != ERR {
		return false
	}
	return true
}

func From(s string) (from string) {
	if strings.Contains(s, "?B?") {
		split := strings.Split(s, " ")
		for _, v := range split {
			index := strings.Index(v, "?B?")
			decoderType := StringsReplaces(v[:index+3], []string{"=?", "?B?"}...)
			notQ := v[index+3:]
			replace := StringsReplaces(notQ, []string{"?=", "_", "\""}...)
			decodeBytes, _ := base64.StdEncoding.DecodeString(replace)
			if strings.ToUpper(decoderType) == "GB18030" {
				decoderType = "GBK"
			}
			encoder := mahonia.NewDecoder(decoderType)
			from = encoder.ConvertString(string(decodeBytes))
			if len(split) >= 2 {
				from = encoder.ConvertString(string(decodeBytes)) + split[1]
			}
			return
		}
	}
	return
}

func Subject(s string) (subject string) {
	if len(s) > 0 {
		if strings.Contains(s, "?B?") {
			split := strings.Split(s, " ")
			var list []string
			for _, v := range split {
				index := strings.Index(v, "?B?")
				decoderType := StringsReplaces(v[:index+3], []string{"=?", "?B?"}...)
				notQ := v[index+3:]
				replace := StringsReplaces(notQ, []string{"?=", "_"}...)
				decodeBytes, _ := base64.StdEncoding.DecodeString(replace)
				if strings.ToUpper(decoderType) == "GB2312" {
					decoderType = "GBK"
				}
				encoder := mahonia.NewDecoder(decoderType)
				list = append(list, encoder.ConvertString(string(decodeBytes)))
			}

			subject = strings.Join(list, "")
			return
		} else if strings.Contains(s, "?Q?") {
			split := strings.Split(s, " ")
			var list []string
			for _, v := range split {
				index := strings.Index(v, "?Q?")
				decoderType := StringsReplaces(v[:index+3], []string{"=?", "?Q?"}...)
				notQ := v[index+3:]
				replace := StringsReplaces(notQ, []string{"?=", "_"}...)
				newReader := quotedprintable.NewReader(strings.NewReader(replace))
				all, err := ioutil.ReadAll(newReader)
				if err != nil {
					return
				}
				if strings.ToUpper(decoderType) == "GB2312" {
					decoderType = "GBK"
				}
				encoder := mahonia.NewDecoder(decoderType)
				list = append(list, encoder.ConvertString(string(all)))
			}

			subject = strings.Join(list, "")
			return
		}
	}
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	subject = string(all)
	return
}

func StringsReplaces(s string, old ...string) (replaceStr string) {
	replaceStr = s
	var new = ""
	for _, v := range old {
		if v == "_" {
			new = " "
		} else {
			new = ""
		}
		if v == "+ " {
			new = "+"
		} else if v == "- " {
			new = "-"
		}
		replaceStr = strings.Replace(replaceStr, v, new, -1)
	}
	return
}

func Other(s string) (str string) {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GBK.NewDecoder())
	all, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}
	return string(all)
}

func RemoveHtml(src string) (text string) {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}
