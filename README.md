# QuickSwitch 快速切换本地HOST

平时在开发的时候经常遇到本地环境,测试环境,生产环境来回切换的问题,本人平时也是用谷歌浏览器的插件`hostAdmin`,但是平时得翻墙太麻烦了. 不如自己实现一个简单的. 参考功能也是依照hostadmin实现出来的.

# 功能
- [x] 备份文件/删除备份文件
- [x] 域名自动分组
- [x] 快速切换host
- [x] 查询domain
- [x] 删除指定的host/domain
# feature
- [x] 编辑host
- [ ] 外部文件修改flash程序
- [x] 添加host
- [x] flash dns(存在问题,执行命令的时候会弹出黑窗口windows)
- [x] 托盘管理
- [ ] 优化UI界面


# 使用 / 安装

1. 下载

```bash
go get -u -v github.com/lazy007/quickswitch
```

2. 部署app 

代码依赖`github.com/therecipe/qt`这个类库, 使用提供的工具生成

```bash
qtdeloy build desktop
```

使用:
1. 设置软件使用管理员权限运行
2. all done!
