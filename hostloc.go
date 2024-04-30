package main

import (
	"fmt"
	"hostloc/httputil"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Account struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"account"`
}

var BotConfig Config

func login(username, password string) bool {
	data := map[string]string{
		"fastloginfield": "username",
		"username":       username,
		"password":       password,
	}
	result, err := httputil.HttpPost("https://hostloc.com/member.php?mod=logging&action=login&loginsubmit=yes&infloat=yes&lssubmit=yes&inajax=1", data)
	if err != nil {
		log.Printf("login failed: %v\n", err)
		return false
	}
	if strings.Contains(result, "window.location.href") {
		return true
	} else {
		return false
	}
}

func signin(uid string) {
	url := "https://hostloc.com/space-uid-" + uid + ".html"
	httputil.HttpGet(url)
}

func main() {
        currentDir, err := os.Getwd()
        os.Chdir(currentDir)
        fmt.Println("Current directory:", currentDir)
	yamlFile, err := os.ReadFile("account.yaml")
	if err != nil {
		panic("read yaml fail: " + err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &BotConfig)
	if err != nil {
		panic("unmarshal yaml fail: " + err.Error())
	}
	if BotConfig.Account.Username == "" || BotConfig.Account.Password == "" {
		panic("username or password is empty")
	}

	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
	logFile, err := os.OpenFile("./logfile.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("无法创建日志文件:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	resultFile, _ := os.OpenFile("./result.txt", os.O_RDWR|os.O_CREATE, 0644)
	lastSigninDate := ""
	content, _ := io.ReadAll(resultFile)
	resultFile.Close()
	lastSigninDate = string(content)
        lastSigninDate = strings.TrimSpace(lastSigninDate) 
	currdate := time.Now().Format("20060102")
	if lastSigninDate == currdate {
		return
	}
	err = httputil.InitClient("")
	if err != nil {
		log.Fatal("init client failed")
	}
	loginresult := login(BotConfig.Account.Username, BotConfig.Account.Password)
	log.Printf("login result=%v\n", loginresult)
	if !loginresult {
		log.Fatal("login failed")
	}
	for i := 0; i < 15; i++ {
		time.Sleep(10 * time.Second)
		//生成一个随机uid,1-72991
		uid := fmt.Sprintf("%d", time.Now().UnixNano()%72991+1)
		log.Printf("%d.signin,uid=%s\n", i+1, uid)
		signin(uid)
	}
	resultFile, _ = os.OpenFile("./result.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	resultFile.WriteString(currdate)
	resultFile.Close()
}

