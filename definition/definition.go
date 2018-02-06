package definition

import (
	"sync"
	"time"

	"github.com/therecipe/qt/widgets"
)

var Locker sync.Mutex

//一行host映射结果
type Dns struct {
	Selected bool                 //是否非注释
	Ip       string               //ip字符串
	Host     string               //域名
	Comment  string               //注释
	Str      string               //原始数据
	Btn      *widgets.QPushButton //关联的Btn
	Except   bool                 //排除不管理
	GroupBox *widgets.QGroupBox   //父级group
}

//一组domain的集合
type Domain struct {
	Dns          []*Dns             //集合
	ChildBtnNum  int64              //子节点数量
	GroupBox     *widgets.QGroupBox //对应的组件
	Layout       *widgets.QVBoxLayout
}

const (
	ISNOTDOMAIN = "isNotDomain" //非正式的host集合KEY
	EOL         = "\r\n"        //换行符
	SELECTEDSTR = "√"           //→
)

var (
	PreModifyTime int64                  //上次修改时间
	BackUpFile    string                 //备份文件名
	HostFile      string                 //文件名
	Location      *time.Location         //本地时区
	Hosts         = map[string]*Domain{} //全局HOST信息对象
	ExceptFlag    = [...]string{
		"#except",
		"#EXCEPT",
		"#不管理",
		"#排除",
		"#hide",
		"#HIDE",
	}                                    //排序字段flag
)

var MainLayout *widgets.QVBoxLayout
