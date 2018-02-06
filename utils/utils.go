package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	"qt/host/compent"
	"qt/host/definition"

	"github.com/fsnotify/fsnotify"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

//SwitchDomainHost 切换生成host
func SwitchDomainHost(dns []*definition.Dns, btn *widgets.QPushButton) {
	for _, host := range dns {
		if host.Except {
			continue
		}
		if host.Btn == btn && !strings.Contains(btn.Text(), definition.SELECTEDSTR+" ") {
			host.Selected = true                                   //修改选中状态
			btn.SetText(definition.SELECTEDSTR + " " + btn.Text()) //设置btn内容
		} else {
			host.Selected = false
			host.Btn.SetText(strings.Replace(host.Btn.Text(), definition.SELECTEDSTR+" ", "", 1))
		}
	}
}

//SaveToHostFile 保存host文件
func SaveToHostFile(filePath string, dns *map[string]*definition.Domain) {
	definition.PreModifyTime = time.Now().In(definition.Location).Unix()
	content := ""
	for domain, hosts := range *dns {
		for _, host := range hosts.Dns {
			comment := ""
			if !host.Selected {
				comment = "#"
			}
			if domain != definition.ISNOTDOMAIN {
				content = content + comment + host.Ip + " " + host.Host + " " + host.Comment + definition.EOL
			} else {
				if host.Str != "" {
					content += host.Str + definition.EOL
				}
			}
		}
		content += strings.Repeat(definition.EOL, 1)
	}
	//fmt.Println(content) //打印最后的host内容
	// return
	mCBtn := widgets.QMessageBox__Close
	//备份host
	src, err := os.Open(filePath)
	if err != nil {
		HandleError(err)
		return
	}
	defer src.Close()
	if f, err := os.Open(definition.BackUpFile); err == nil { //存在文件不做任何操作
		f.Close()
	} else { //不存在备份文件则生成一个
		target, err := os.OpenFile(definition.BackUpFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			HandleError(err)
			return
		}
		defer target.Close()
		_, err = io.Copy(target, src)
		if err != nil {
			HandleError(err)
			return
		}
	}
	src.Close()

	err = ioutil.WriteFile(filePath, []byte(content), 0)
	if err != nil {
		HandleError(err)
		return
	}
	//https://zh.wikihow.com/%E5%88%B7%E6%96%B0-DNS 刷新dns
	go func() {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			//cmd := exec.Command("")
		case "linux":
			//cmd = exec.Command("")
		default:
			//start /b
			cmd = exec.Command("ipconfig", "/displaydns")
			cmd.Env = os.Environ()
		}
		if err := cmd.Run(); err != nil {
			fmt.Println(err.Error())
			return
		} else {
			fmt.Println("刷新缓存成功")
		}
	}()
	widgets.QMessageBox_Information(nil, "切换Host成功", "切换Host成功", mCBtn, mCBtn)
}

//RenderView 渲染展示
func RenderView(layout *widgets.QVBoxLayout) {
	file, err := os.Open(definition.HostFile)
	if err != nil {
		Error("打开host失败 原因: " + err.Error())
		return
	}
	reader := bufio.NewReader(file)
	reg := regexp.MustCompile(`\s+`)
	var str string
	var hostInfo []string
	for {
		str, err = reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		str = reg.ReplaceAllString(strings.Trim(strings.Trim(strings.Trim(str, "\n"), "\r"), " "), " ")
		str = strings.Replace(str, "# ", "#", 1) //将:# ->#
		hostInfo = append(hostInfo, str)
	}

	definition.Hosts[definition.ISNOTDOMAIN] = &definition.Domain{} //事先初始化非domain map的结构体
	//检查注释并且分组
	for _, v := range hostInfo {
		isStart, _ := regexp.Match(`^\d+`, []byte(v)) //以数字为开始的均视为被启用
		//先匹配标准的
		curDNS := &definition.Dns{
			Str:      v,
			Selected: isStart, // #前缀
		}
		//正则替换多个空格
		splits := strings.Split(v, " ")
		//只处理大于三个的
		if len(splits) >= 2 {
			curDNS.Ip = strings.Replace(splits[0], "#", "", 1)
			if !regexp.MustCompile(`^\d{1,3}\.`).Match([]byte(curDNS.Ip)) {
				curDNS.Selected = false
				curDNS.Ip = ""
				(*definition.Hosts[definition.ISNOTDOMAIN]).Dns = append((*definition.Hosts[definition.ISNOTDOMAIN]).Dns, curDNS)
				continue //非 ip 特征跳过
			}
			curDNS.Host = splits[1]
			if curDNS.Ip == "" || curDNS.Host == "" {
				curDNS.Selected = false
				continue
			}
			if len(splits) >= 3 {
				curDNS.Comment = strings.Join(splits[2:], " ")
				//处理不用被管理的host
				for _, flag := range definition.ExceptFlag {
					if strings.HasPrefix(curDNS.Comment, flag) {
						curDNS.Except = true //判定排除
						break
					}
				}
			}

			if _, ok := definition.Hosts[curDNS.Host]; !ok {
				definition.Hosts[curDNS.Host] = &definition.Domain{} //cannot assign to struct field xxx in map 处理这个问题
			}
			(*definition.Hosts[curDNS.Host]).Dns = append((*definition.Hosts[curDNS.Host]).Dns, curDNS) //域名分组
		} else {
			(*definition.Hosts[definition.ISNOTDOMAIN]).Dns = append((*definition.Hosts[definition.ISNOTDOMAIN]).Dns, curDNS)
		}
	}
	for k, v := range definition.Hosts {
		if k == definition.ISNOTDOMAIN {
			continue
		}
		groupLayout := widgets.NewQVBoxLayout()
		group := widgets.NewQGroupBox2(strings.ToUpper(k), nil)
		v.GroupBox = group
		v.Layout = groupLayout
		v.ChildBtnNum = 0
		v.GroupBox.SetFixedWidth(293) //设置按钮组的最大高度
		for _, v1 := range v.Dns {
			v1.Btn = widgets.NewQPushButton(nil)
			v1.Btn.SetToolTip(strings.TrimLeft(v1.Comment, "#"))
			v1.GroupBox = group
			v.ChildBtnNum++
			if v1.Selected {
				v1.Btn.SetText(definition.SELECTEDSTR + " " + v1.Ip)
			} else {
				v1.Btn.SetText(v1.Ip)
			}
			if v1.Except { //禁用点击按钮
				v1.Btn.SetDisabled(true)
			} else {
				RegisterBtnEvent(v1.Btn, v, k)
			}
			groupLayout.AddWidget(v1.Btn, 1, 0)
		}
		v.GroupBox.SetLayout(groupLayout)
		groupLayout.SetSpacing(5)
		layout.AddWidget(v.GroupBox, 1, 0)
	}

}

//SearchDomain 查询过滤
func SearchDomain(val string) {
	for domain, hosts := range definition.Hosts {
		if domain == definition.ISNOTDOMAIN {
			continue
		}
		if val == "" { //为空全部显示
			hosts.GroupBox.Show()
			continue
		}
		if val != "" && !strings.Contains(strings.ToUpper(domain), strings.ToUpper(val)) { //匹配不到关键字的隐藏
			hosts.GroupBox.Hide()
		} else {
			hosts.GroupBox.Show()
		}
	}
}

//HandleError 处理error
func HandleError(err error) {
	if os.IsPermission(err) {
		Error("没有操作权限")
	} else if os.IsNotExist(err) {
		Error("文件不存在")
	} else if os.IsExist(err) {
		Error("文件已存在")
	} else {
		Error(err.Error())
	}
}

//Error 打印error
func Error(msg string) {
	fmt.Println(msg)
	widgets.QMessageBox_Critical(nil, "错误提醒", msg, widgets.QMessageBox__Yes, widgets.QMessageBox__Yes)
}

//ListenFileModifyTime 监听文件修改
func ListenFileModifyTime(layout *widgets.QVBoxLayout) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				definition.Locker.Lock()
				if event.Op&fsnotify.Write == fsnotify.Write {
					f, err := os.Stat(definition.HostFile)
					if err == nil {
						preFileModifyTime := f.ModTime().In(definition.Location).Unix()
						if preFileModifyTime-definition.PreModifyTime > 5 {
							fmt.Println("reloading...")
							Reload(layout)
							definition.PreModifyTime = preFileModifyTime
							fmt.Println("reloaded")
						}
					}
				}
				definition.Locker.Unlock()
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}

		}
	}()
	err = watcher.Add(definition.HostFile)
	if err != nil {
		log.Fatal(err)
	}
	<-done //死锁, 处于等待
}

//Reload todo:目前动态添加组件不会显示, 还不知道怎么搞
//Reload 监听Host被外部修改重载
func Reload(layout *widgets.QVBoxLayout) {
	if len(definition.Hosts) > 0 {
		//for _, item := range dns {
		//	if item.GroupBox.Pointer() != nil {
		//		layout.RemoveWidget(item.GroupBox) //移除这些组件即可
		//	}
		//}
		//dns = map[string]*Domain{}
		time.Sleep(500 * time.Millisecond) //防止再次读取host为空内容, 休眠0.5秒
		RenderView(layout)
	}
}

//FindCurrentDNSIndexByBtn 通过btn查找当前的对应的dns
func FindCurrentDNSIndexByBtn(button *widgets.QPushButton, dns *[]*definition.Dns) (int, error) {
	for k, v := range *dns {
		if v.Btn == button {
			return k, nil
		}
	}
	return -1, errors.New("无法找到指定的button")
}

func RegisterBtnEvent(button *widgets.QPushButton, domain *definition.Domain, k string) {
	button.ConnectClicked(func(_ bool) {
		SwitchDomainHost(definition.Hosts[k].Dns, button)
	})
	button.ConnectContextMenuEvent(func(event *gui.QContextMenuEvent) {
		menu := widgets.NewQMenu(nil)
		editAction := menu.AddAction("修改此项")
		exceptAction := menu.AddAction("禁止管理") // 设置为不管理
		delAction := menu.AddAction("删除此项")

		delAction.ConnectTriggered(func(_ bool) {
			if widgets.QMessageBox_Information(
				nil,
				"删除提醒",
				"确定要删除该host信息吗?",
				widgets.QMessageBox__Yes|widgets.QMessageBox__Close,
				widgets.QMessageBox__Close,
			) == widgets.QMessageBox__Close {
				return
			}

			domain.ChildBtnNum--
			if domain.ChildBtnNum <= 0 {
				delete(definition.Hosts, k)        //移除这个节点
				domain.GroupBox.DestroyQGroupBox() //销毁这个Box
			} else {
				var index int
				for k, btn := range definition.Hosts[k].Dns {
					if btn.Btn == button {
						index = k
						break
					}
				}
				//剔除button
				definition.Hosts[k].Dns = append(definition.Hosts[k].Dns[0:index], definition.Hosts[k].Dns[index+1:]...)
				button.DestroyQPushButton() //销毁Btn
			}
		})

		editAction.ConnectTriggered(func(_ bool) {
			index, err := FindCurrentDNSIndexByBtn(button, &definition.Hosts[k].Dns)
			if err == nil {
				compent.NewDialog(definition.Hosts[k].Dns[index], false, "修改host")
			}
		})

		exceptAction.ConnectTriggered(func(_ bool) {
			if widgets.QMessageBox_Question(nil,
				"修改提醒", "是否要排除管理此项?",
				widgets.QMessageBox__Yes|widgets.QMessageBox__Cancel,
				widgets.QMessageBox__Cancel) == widgets.QMessageBox__Yes {
				index, err := FindCurrentDNSIndexByBtn(button, &definition.Hosts[k].Dns)
				if err == nil {
					definition.Hosts[k].Dns[index].Except = true
					definition.Hosts[k].Dns[index].Btn.SetDisabled(true)
					definition.Hosts[k].Dns[index].Comment = definition.ExceptFlag[0] + " " + strings.TrimLeft(definition.Hosts[k].Dns[index].Comment, "#")
				}
				return
			}
		})

		menu.Exec2(gui.QCursor_Pos(), nil)
	})
}

func TrayMenu(app *widgets.QApplication, mw *widgets.QMainWindow, tips string) {
	//检查是否支持托盘
	if widgets.QSystemTrayIcon_IsSystemTrayAvailable() {
		trayMenu := widgets.NewQMenu(widgets.QApplication_Desktop())
		restoreWinAction := widgets.NewQAction2("显示(&R)", mw)
		quitAction := widgets.NewQAction2("退出(&Q)", mw)
		//添加菜单
		trayMenu.AddActions([]*widgets.QAction{restoreWinAction, quitAction})
		//恢复
		restoreWinAction.ConnectTriggered(func(checked bool) {
			mw.Show()
		})
		//隐藏
		quitAction.ConnectTriggered(func(checked bool) {
			app.Quit()
		})
		//实例化QSystemTryIcon
		myTrayIcon := widgets.NewQSystemTrayIcon(mw)
		//设置图标
		myTrayIcon.SetIcon(gui.NewQIcon5(":/qrc/app.ico"))
		//鼠标放托盘图标上提示信息
		myTrayIcon.SetToolTip(tips)
		myTrayIcon.SetContextMenu(trayMenu) //托盘菜单
		myTrayIcon.Show()
		//拦截关闭事件
		mw.ConnectCloseEvent(func(event *gui.QCloseEvent) {
			mw.Hide() //隐藏主窗体
			myTrayIcon.ShowMessage(tips, "程序已最小化", widgets.QSystemTrayIcon__Information, 2)
			event.Ignore() //忽略信号
		})

		myTrayIcon.ConnectActivated(func(reason widgets.QSystemTrayIcon__ActivationReason) {
			switch reason {
			case widgets.QSystemTrayIcon__Trigger: //触发
				if mw.IsHidden() { //隐藏则显示
					mw.Show()
				} else {
					myTrayIcon.ShowMessage(tips, "程序已最小化", widgets.QSystemTrayIcon__Information, 2)
					mw.Hide()
				}
			}
		})
	}
}
