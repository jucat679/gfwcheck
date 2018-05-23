## GFW Check

> 这是一个使用 go 编写的代理健康检测工具，支持多服务器配置，支持标准的 cron 定时执行与邮件告警；
该工具使用指定的代理地址尝试连接 google，如果不成功则认为代理失效，从而触发远程与本地的命令执行；
当多次检测失败并且执行命令后代理未恢复正常，则触发告警

### 安装

自行编译出二进制文件，然后执行 `./gfwcheck install` 即可，如下所示

``` sh
mritd.docker ➜  ~ ./gfwcheck install
2018/05/23 22:55:57 Stop gfwcheck
2018/05/23 22:55:57 Clean files
2018/05/23 22:55:57 Systemd reload
2018/05/23 22:55:57 Create config dir /etc/gfwcheck
2018/05/23 22:55:57 Copy file to /usr/bin
2018/05/23 22:55:57 Create config file /etc/gfwcheck/config.yaml
2018/05/23 22:55:57 Create systemd config file /lib/systemd/system/gfwcheck.service
```

安装成功后会创建 systemd service 配置，同时会将样例配置写入到 `/etc/gfwcheck/config.yaml`
修改配置文件后执行 `systemctl start gfwcheck` 即可启动，执行 `journalctl -fu gfwcheck` 可实时观察日志

### 配置解释

``` yaml
# 告警配置
alarms:
# 目前仅支持 smtp 邮件告警，并且仅支持 TLS smtp 服务器
- type: smtp
  # 指定接收者，可以有多个
  targets:
  - mritd1234@gmail.com

# 服务器配置
servers:
  # 指定名称(可以任意填写，用于 log 显示与告警)
- name: test1.com
  # 远程主机地址
  host: test1.com
  # ssh 端口
  port: 22
  # ssh 用户
  user: root
  # 如果采用 ssh key 登录密码留空，否则填写 ssh 密码
  password: ""
  # 登录方式 (password/pem)
  method: pem
  # 如果采用 ssh key 方式登录，则此处填写 ssh key 绝对路径
  key: "/root/.ssh/id_rsa"
  # 当检测失败时 ssh 连接超时时间
  timeout: 10s
  # 需要检测的代理地址
  proxy: socks5://192.168.1.10:2018
  # 检测失败后本地执行的命令
  localcmd: "docker restart proxy"
  # 检测失败后远程需要执行的命令
  remotecmd: "docker restart proxy"
  # 检测时间间隔，支持标准 cron 表达式，同时支持快捷命令(以下为 每 30s 检测一次)
  # 检测时间应当尽量 >= timeout*2
  cron: '@every 30s'
  # 最大失败次数，超过该次数后将触发告警
  maxfailed: 5
- name: test2.com
  host: test2.com
  port: 22
  user: root
  password: ""
  method: pem
  key: "/root/.ssh/id_rsa"
  timeout: 10s
  proxy: socks5://192.168.1.10:9999
  localcmd: "docker restart proxy_bak"
  remotecmd: "systemctl reboot"
  cron: '@every 30s'
  maxfailed: 5

# 邮件告警设置
smtp:
  # 邮箱地址
  username: mritd@mritd.me
  # 密码
  password: "xxxxxxx"
  # 发送者，建议与邮箱地址保持一致
  from: mritd@mritd.me
  # 邮箱发件服务器地址(仅支持 TLS 端口)
  server: "smtp.xxx.com:465"
```

### 卸载

卸载直接执行 `./gfwcheck uninstall` 即可