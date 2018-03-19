package compent

import (
	"strings"

	"quickswitch/definition"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func NewDialog(dns *definition.Dns, isNewDNS bool, title string) {
	widget := widgets.NewQDialog(nil, core.Qt__WindowCloseButtonHint)
	widget.SetWindowTitle(title)
	layout := widgets.NewQGridLayout2()
	widget.SetLayout(layout)
	widget.SetFixedSize2(250, 200)
	layout.SetContentsMargins(7, 7, 7, 7)
	layout.SetSpacing(10)

	domainLabel := widgets.NewQLabel2("域名:", nil, 0)
	ipLabel := widgets.NewQLabel2("  IP:", nil, 0)
	commentLabel := widgets.NewQLabel2("注释:", nil, 0)
	selectLabel := widgets.NewQLabel2("启用:", nil, 0)

	layout.AddWidget(domainLabel, 0, 0, 0)
	layout.AddWidget(ipLabel, 1, 0, 0)
	layout.AddWidget(commentLabel, 2, 0, 0)
	layout.AddWidget(selectLabel, 3, 0, 0)

	domainInput := widgets.NewQLineEdit(nil)
	ipInput := widgets.NewQLineEdit(nil)
	commentInput := widgets.NewQLineEdit(nil)
	selectComboBox := widgets.NewQComboBox(nil)
	selectComboBox.AddItems([]string{"是", "否"})

	layout.AddWidget(domainInput, 0, 1, 0)
	layout.AddWidget(ipInput, 1, 1, 0)
	layout.AddWidget(commentInput, 2, 1, 0)
	layout.AddWidget(selectComboBox, 3, 1, 0)

	if !isNewDNS {
		domainInput.SetDisabled(true) //禁用,修改时期不允许修改域名信息
		domainInput.SetText(dns.Host)
		ipInput.SetText(dns.Ip)
		commentInput.SetText(strings.TrimLeft(dns.Comment, "#"))
		if dns.Selected {
			selectComboBox.SetCurrentText("是")
		} else {
			selectComboBox.SetCurrentText("否")
		}
	}
	okBtn := widgets.NewQPushButton2("确定", nil)
	cancelBtn := widgets.NewQPushButton2("取消", nil)

	cancelBtn.ConnectClicked(func(_ bool) {
		if widget != nil {
			widget.Close()
			widget.Destroy(true, true)
		}
	})

	okBtn.ConnectClicked(func(_ bool) {
		dns.Ip = ipInput.Text()
		commentPrefix := "#"
		dns.Except = false
		se := selectComboBox.CurrentIndex()
		if se == 0 {
			dns.Selected = true
			dns.Btn.SetText(definition.SELECTEDSTR + " " + dns.Ip)
		} else {
			dns.Btn.SetText(dns.Ip)
			dns.Selected = false
		}
		dns.Comment = commentPrefix + commentInput.Text()
		//刷新按钮信息
		if !isNewDNS {
			dns.Btn.SetToolTip(strings.TrimLeft(dns.Comment, "#"))
		}
		dns.Host = domainInput.Text()
		if widget != nil {
			widget.Close()
			widget.Destroy(true, true)
		}
	})

	btnLayout := widgets.NewQHBoxLayout2(nil)

	btnLayout.AddStretch(2)
	btnLayout.AddWidget(okBtn, 1, 0)
	btnLayout.AddWidget(cancelBtn, 1, 0)
	layout.AddLayout2(btnLayout, 5, 0, 1, 2, 0)
	widget.Exec() //qDialog 特性, 阻塞进程 , 必须调用Close方法才可以.真正返回给函数结果
}
