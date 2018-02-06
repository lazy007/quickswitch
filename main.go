package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"qt/host/definition"
	"qt/host/utils"

	"qt/host/compent"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var mw *widgets.QMainWindow

const VERSION = "0.0.1"

const TITLE = "QuickSwitch"

func init() {
	definition.Location, _ = time.LoadLocation("Asia/Shanghai")
	switch runtime.GOOS {
	case "windows":
		definition.HostFile = "C:\\Windows\\System32\\drivers\\etc\\hosts" // window host
	default:
		definition.HostFile = "/etc/hosts" // linux darwin host
	}
	definition.BackUpFile = definition.HostFile + "_" + TITLE
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	//读取host文件
	mw = widgets.NewQMainWindow(nil, core.Qt__CustomizeWindowHint|core.Qt__WindowCloseButtonHint) //core.Qt__FramelessWindowHint | core.Qt__WindowStaysOnTopHint
	mw.SetWindowTitle(TITLE + " - V" + VERSION)
	mw.SetFixedSize2(305, 680)

	//mw.SetFont(gui.NewQFont2("宋体", 8, 0, false))
	opBtnOk := widgets.NewQPushButton2("确定", mw)
	opBtnOk.SetGeometry2(10, 645, 90, 25)
	opBtnOk.ConnectClicked(func(_ bool) {
		utils.SaveToHostFile(definition.HostFile, &definition.Hosts)
	})

	addBtn := widgets.NewQPushButton2("添加", mw)
	addBtn.SetGeometry2(108, 645, 90, 25)
	addBtn.ConnectClicked(func(_ bool) {
		dns := &definition.Dns{
			Btn: widgets.NewQPushButton(nil),
		}
		compent.NewDialog(dns, true, "添加 host")

		if dns.Host == "" { //没有域名数据 返回
			return
		}
		_, ok := definition.Hosts[dns.Host]
		//如果已经存在了
		if ok {
			definition.Hosts[dns.Host].ChildBtnNum++
			(*definition.Hosts[dns.Host]).Dns = append((*definition.Hosts[dns.Host]).Dns, dns)
			dns.GroupBox = definition.Hosts[dns.Host].GroupBox
			definition.Hosts[dns.Host].Layout.AddWidget(dns.Btn, 0, 0)
			dns.Btn.Show()
		} else { //节点不存在, 创建一个节点
			groupLayout := widgets.NewQVBoxLayout()
			group := widgets.NewQGroupBox2(strings.ToUpper(dns.Host), nil)
			definition.Hosts[dns.Host] = &definition.Domain{
				ChildBtnNum: 1,
				Dns:         []*definition.Dns{},
				GroupBox:    group,
				Layout:      groupLayout,
			}
			(*definition.Hosts[dns.Host]).Dns = append((*definition.Hosts[dns.Host]).Dns, dns)
			definition.Hosts[dns.Host].GroupBox.SetFixedWidth(293)
			groupLayout.AddWidget(dns.Btn, 1, 0)
			group.SetLayout(groupLayout)
			groupLayout.SetSpacing(5)
			definition.MainLayout.AddWidget(group, 1, 0)
		}

		utils.RegisterBtnEvent(dns.Btn, definition.Hosts[dns.Host], dns.Host)
		fmt.Println(dns)

	})

	delBakBtn := widgets.NewQPushButton2("删除备份", mw)
	delBakBtn.SetToolTip("不删除 缓存软件不会再次生成备份文件" + definition.BackUpFile)
	delBakBtn.SetGeometry2(206, 645, 90, 25)
	delBakBtn.ConnectClicked(func(_ bool) {
		if widgets.QMessageBox_Question(mw,
			"删除提醒", "是否要删除"+definition.BackUpFile,
			widgets.QMessageBox__Yes|widgets.QMessageBox__Cancel,
			widgets.QMessageBox__Cancel) == widgets.QMessageBox__Yes {
			err := os.Remove(definition.BackUpFile)
			if err != nil {
				utils.HandleError(err)
				return
			}
			widgets.QMessageBox_Information(nil, "删除成功", "切换Host成功", widgets.QMessageBox__Yes, widgets.QMessageBox__Yes)
		}
	})

	//==========\\ INPUT ==============
	searchInput := widgets.NewQLineEdit(mw)
	searchInput.SetGeometry2(10, 8, 280, 25)
	searchInput.SetPlaceholderText("查询域名")
	searchInput.ConnectKeyReleaseEvent(func(event *gui.QKeyEvent) {
		utils.SearchDomain(searchInput.Text())
	})
	//===========\\ END ===============

	layout := widgets.NewQVBoxLayout()
	layout.SetSizeConstraint(widgets.QLayout__SetFixedSize) //设置布局跟随组件的宽高
	layout.AddSpacing(10)
	layout.SetContentsMargins(5, 5, 5, 5)

	QScrollArea := widgets.NewQScrollArea(mw)
	qWidget := widgets.NewQWidget(QScrollArea, 0) //父组件是滚动条
	QScrollArea.SetWidget(qWidget)                //设置要监听的控件
	QScrollArea.SetGeometry2(0, 40, 325, 600)     //设置滚动条区域把滚动区域放到主窗体外面
	qWidget.SetFixedSize2(295, 600)
	qWidget.SetLayout(layout)
	qWidget.SetObjectName("mainWidget")
	mw.SetStyleSheet(`
#mainWidget QPushButton {text-align:left;}
QMainWindow,#mainWidget {background-color:#fdfdfd;}
`)
	definition.MainLayout = layout
	utils.RenderView(layout)
	go utils.ListenFileModifyTime(layout)
	utils.TrayMenu(app, mw, TITLE)
	mw.Show()
	app.Exec()
}
