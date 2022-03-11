package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"go-nacos-samples/src/main/go/golog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/**
文件
path 文件路径
name 文件名
*/
type FileInf struct {
	path string
	name string
}

/**
配置文件属性
path 文件路径
name 文件名
md5 md5值
content 文件内容
statusCode 状态值  200存在 404表示不存在
*/
type ConfigFileInfo struct {
	path       string
	name       string
	md5        string
	content    string
	statusCode int
}

var propertiesPath string //nacos链接信息
var configPath string     //应用配置信息
var logFile = "info.log"

func main() {
	//本地测试用
	//propertiesPath = "绝对路径\\go-nacos\\src\\main\\go\\bootstrap.properties"
	//configPath = "绝对路径\\go-nacos\\src\\main\\config\\paas"
	propertiesPath = "D:\\ideawsoy\\go-nacos-samples\\src\\main\\resources\\bootstrap.properties"
	configPath = "D:\\ideawsoy\\go-nacos-samples\\src\\main\\config\\paas"
	if propertiesPath == "" {
		propertiesPath = "bootstrap.properties"
	}
	if configPath == "" {
		configPath = "config/paas"
	}
	config := InitConfig(propertiesPath)
	ip := config["spring.cloud.nacos.config.server-addr"]
	tenant := config["spring.cloud.nacos.config.namespace"]
	group := "paas" // 配置分组

	fromSlash := filepath.FromSlash(configPath)
	files := getFiles(fromSlash)

	for _, fi := range files {
		dataId := fi.name
		Log(fmt.Sprintf("ip:%s,tenant:%s,dataId:%s,group:%s", ip, tenant, dataId, group))
		remoteFileInfo := get(ip, tenant, dataId, group)
		//不存在直接上传本地
		if remoteFileInfo.statusCode != 200 {
			localPath := configPath + string(os.PathSeparator) + dataId
			content, err := ioutil.ReadFile(localPath)
			if err != nil {
				Log(fmt.Sprintf("error: read %s file fail. %s", localPath, err))
				os.Exit(0)
			}
			Log(fmt.Sprintf("dataId=%s 发布", dataId))
			post(ip, tenant, dataId, group, string(content))
		} else {
			Log(fmt.Sprintf("dataId=%s 已存在，不发布", dataId))
		}
	}
}

func post(ip, tenant, dataId, group, content string) {
	urlValues := url.Values{}
	urlValues.Add("tenant", tenant)
	urlValues.Add("dataId", dataId)
	urlValues.Add("group", group)
	urlValues.Add("type", "yaml")
	urlValues.Add("content", string(content))

	// 发布配置
	url := "http://" + ip + "/nacos/v1/cs/configs"
	resp, err := http.PostForm(url, urlValues)
	if err != nil {
		Log(fmt.Sprintf("error: post %s fail. %s", url, err))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	Log(fmt.Sprintf("post result: %s", string(body)))
}

func get(ip, tenant, dataId, group string) ConfigFileInfo {
	url := build(ip, tenant, dataId, group)
	resp, err := http.Get(url)
	if err != nil {
		Log(fmt.Sprintf("error: get %s fail. %s", url, err))
	}
	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		s := fmt.Sprintf("%x", md5.Sum(body))
		return ConfigFileInfo{name: dataId, content: string(body), md5: s, statusCode: resp.StatusCode}
	} else {
		return ConfigFileInfo{statusCode: resp.StatusCode}
	}
}
func build(ip, tenant, dataId, group string) string {
	u := url.URL{}
	u.Scheme = "http"
	u.Host = ip
	u.Path = "/nacos/v1/cs/configs"

	values := url.Values{} //拼接query参数
	values.Add("tenant", tenant)
	values.Add("dataId", dataId)
	values.Add("group", group)

	u.RawQuery = values.Encode()

	Log(fmt.Sprintf("URL:%s", u.String()))
	return u.String()
}

//写日志
func Log(msg string) {
	fmt.Println(msg)
	golog.WriteLog(logFile, msg)
}

//读取key=value类型的配置文件
func InitConfig(path string) map[string]string {
	config := make(map[string]string)

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		config[key] = value
	}
	return config
}

//propertiesPath := "bootstrap.properties"
func getFiles(configPath string) []FileInf {
	var files []FileInf

	_, err1 := os.Lstat(configPath)
	if err1 != nil {
		Log(fmt.Sprintf("error: configPath=%s fail. %s", configPath, err1))
		os.Exit(1)
	}

	err := filepath.Walk(configPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			//path是全路径
			files = append(files, FileInf{path: path, name: f.Name()})
		}
		return nil
	})

	if err != nil {
		panic(err)
	}
	/*for _, file := range files {
		Log(fmt.Sprintf("%s", file))
	}*/
	return files
}
