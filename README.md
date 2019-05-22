### 企业微信消息发送

    根据网上代码修改的加上了支持接收参数和配置文件功能,适用于zabbix监控使用
    相对于python版的不需要环境依赖。
    只有一个可执行程序和一个配置文件。
    作者：yang_youwei81@163.com
    
### useage    
    
    Usage: COMMAND args1 args2 args3
    args1 is usercount
    args2 is the mesages's title
    args3 is messages's content

### conf

配置文件和程序放到一个目录下。

    [main]
    #腾讯接口无变化的话，sendurl和get_token 不用修改。
    sendurl = https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=
    get_token = https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=
    
    #企业微信中建立的监控应用的AgentId
    agid = 1000001
    #企业微信中建立的监控应用的secret
    secret = uhTcjsdsdkfk34rg3gvbODU
    #企业ID
    corpid = wwsd2fwef4gahy500
    
### 编译

    centos
    yum install golang -y
    cd
    mkdir go{pkg,bin,src}
    cd go/src
    go get github.com/yangyouwei/wechat
    go get github.com/Unknwon/goconfig
    cd github.com/yangyouwei/wechat
    go build wechat.go
