// +build windows

package utils

import (
	"github.com/CodyGuo/win"
	"fmt"
)

func FlashDns(callback func(string)) {
	cmd := "ipconfig /flushdns"
	lpCmdLine := win.StringToBytePtr(cmd)
	ret := win.WinExec(lpCmdLine, win.SW_HIDE)
	if ret <= 31 {
		winExecError := map[uint32]string{
			0:  "The system is out of memory or resources.",
			2:  "The .exe file is invalid.",
			3:  "The specified file was not found.",
			11: "The specified path was not found.",
		}
		callback(winExecError[ret])
	} else {
		fmt.Println("刷新缓存成功")
	}
}
