package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Unknwon/goconfig"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var requestError = errors.New("request error,check url or network")

var (
	corpid    string
	agid      int
	secret    string
	sendurl   string
	get_token string
)

type access_token struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
}

//定义一个简单的文本消息格式
type send_msg struct {
	Touser  string            `json:"touser"`
	Toparty string            `json:"toparty"`
	Totag   string            `json:"totag"`
	Msgtype string            `json:"msgtype"`
	Agentid int               `json:"agentid"`
	Text    map[string]string `json:"text"`
	Safe    int               `json:"safe"`
}

type send_msg_error struct {
	Errcode int    `json:"errcode`
	Errmsg  string `json:"errmsg"`
}

var Usage = func() {
	fmt.Println("Usage: COMMAND args1 args2 args3")
	fmt.Println("args1 is usercount")
	fmt.Println("args2 is the mesages's title")
	fmt.Println("args3 is messages's content")
}

func init() {
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
		return
	}
	//
	sendurl, err = cfg.GetValue("main", "sendurl")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "sendurl", err)
	}

	//
	get_token, err = cfg.GetValue("main", "get_token")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "get_token", err)
	}

	//
	corpid, err = cfg.GetValue("main", "corpid")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "corpid", err)
	}

	//
	agid, err = cfg.Int("main", "agid")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "agid", err)
	}

	//
	secret, err = cfg.GetValue("main", "secret")
	if err != nil {
		log.Fatalf("无法获取键值（%s）：%s", "secret", err)
	}


}

func main() {

	args := os.Args

	if args == nil || len(args) < 2 {
		Usage() //如果用户没有输入,或参数个数不够,则调用该函数提示用户
		return
	}
	touser := &args[1]
	agentid := &agid
	h := args[2]
	head := &h
	txt := args[3]
	content := &txt
	c := &corpid
	corpsecret := &secret

	var m send_msg = send_msg{Touser: *touser, Msgtype: "text", Agentid: *agentid, Text: map[string]string{"content": *head + "\n" + *content}}

	///-p "wx246" -s "JbjkM"
	token, err := Get_token(*c, *corpsecret)
	if err != nil {
		println(err.Error())
		return
	}
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = Send_msg(token.Access_token, buf)
	if err != nil {
		println(err.Error())
	}
}

//发送消息.msgbody 必须是 API支持的类型
func Send_msg(Access_token string, msgbody []byte) error {
	body := bytes.NewBuffer(msgbody)
	resp, err := http.Post(sendurl+Access_token, "application/json", body)
	if resp.StatusCode != 200 {
		return requestError
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var e send_msg_error
	err = json.Unmarshal(buf, &e)
	if err != nil {
		return err
	}
	if e.Errcode != 0 && e.Errmsg != "ok" {
		return errors.New(string(buf))
	}
	return nil
}

//通过corpid 和 corpsecret 获取token
func Get_token(corpid, corpsecret string) (at access_token, err error) {
	resp, err := http.Get(get_token + corpid + "&corpsecret=" + corpsecret)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = requestError
		return
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(buf, &at)
	if at.Access_token == "" {
		err = errors.New("corpid or corpsecret error.")
	}
	return
}

func Parse(jsonpath string) ([]byte, error) {
	var zs = []byte("//")
	File, err := os.Open(jsonpath)
	if err != nil {
		return nil, err
	}
	defer File.Close()
	var buf []byte
	b := bufio.NewReader(File)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		line = bytes.TrimSpace(line)
		if len(line) <= 0 {
			continue
		}
		index := bytes.Index(line, zs)
		if index == 0 {
			continue
		}
		if index > 0 {
			line = line[:index]
		}
		buf = append(buf, line...)
	}
	return buf, nil
}

