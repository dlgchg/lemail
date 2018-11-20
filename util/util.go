package util

import (
	"os"
	"io"
	"runtime"
	"path"
	"github.com/urfave/cli"
	"fmt"
)

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func CopyFileToProjectAttach(attach string) (string, error) {
	if len(attach) > 0 {
		file, err := os.Open(attach)
		defer file.Close()
		if err == nil {
			_, f, _, ok := runtime.Caller(1)
			if ok {
				f = path.Dir(f)
				CopyFile(f+"/attach/"+path.Base(attach), attach)
				return "./attach/" + path.Base(attach), nil
			} else {
				return "", nil
			}
		} else {
			return "", err
		}
	}
	return "", nil
}

func CheckIsEmpty(c *cli.Context, value ...string) error {
	for _, v := range value {
		if len(v) == 0 {
			cli.ShowCommandHelp(c, c.Command.Name)
			return fmt.Errorf("参数缺失")
		}
	}
	return nil
}

//截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
