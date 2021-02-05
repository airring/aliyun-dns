package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	"strings"
	"time"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	alidns20150109  "github.com/alibabacloud-go/alidns-20150109/client"
	"github.com/alibabacloud-go/tea/tea"
	"gopkg.in/yaml.v2"
)

type conf struct {
	AccessKeyId     string  `yaml:"accessKeyId"`
	AccessKeySecret string  `yaml:"accessKeySecret"`
	Domain          string `yaml:"domain"`
	Dnsdomain          string `yaml:"dnsdomain"`
}

func (c *conf) getConf() *conf {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println("Unmarshal: %v", err)
	}
	return c
}


func get_external() string {
		resp, err := http.Get("http://ip.cip.cc")
		    if err != nil {
				return "接口访问失败"    
				}
			defer resp.Body.Close()
			content, _ := ioutil.ReadAll(resp.Body)  
			// fmt.Println("获取到的外网ip为:", string(content))
			return strings.Replace(string(content),"\n","",-1)
		}
/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
 func CreateClient (accessKeyId *string, accessKeySecret *string) (_result *alidns20150109.Client, _err error) {
	config := &openapi.Config{}
	// 您的AccessKey ID
	config.AccessKeyId = accessKeyId
	// 您的AccessKey Secret
	config.AccessKeySecret = accessKeySecret
	// 访问的域名
	config.Endpoint = tea.String("dns.aliyuncs.com")
	_result = &alidns20150109.Client{}
	_result, _err = alidns20150109.NewClient(config)
	return _result, _err
  }
  
func _main(args []*string) (_err error) {
	var requsetmap alidns20150109.DescribeSubDomainRecordsResponseBodyDomainRecordsRecord
	var access conf
	access.getConf()
	// fmt.Println("accessKeyId:", access)
	client, _err := CreateClient(tea.String(access.AccessKeyId), tea.String(access.AccessKeySecret))
	if _err != nil {
		return _err
	}

	describeSubDomainRecordsRequest := &alidns20150109.DescribeSubDomainRecordsRequest{
		SubDomain: tea.String(access.Dnsdomain),
	  }
	  urlrequest, _err := client.DescribeSubDomainRecords(describeSubDomainRecordsRequest)
	  if _err != nil {
		return _err
	  }
	// 获取域名id

	domainlist := urlrequest.Body.DomainRecords.Record
	for _, n := range domainlist {
		requsetmap = *n
	}

	// 获取ip地址
	clientip := get_external()
	fmt.Println("获取到外网ip为:", clientip)
	// 判断dns地址是否与当前地址相等
	if *requsetmap.Value == clientip {
		fmt.Printf("当前ip相等,不做处理")
		return _err
	}
	// 修改dns地址
	updateDomainRecordRequest := &alidns20150109.UpdateDomainRecordRequest{
		RecordId: tea.String(*requsetmap.RecordId),
		RR: tea.String(*requsetmap.RR),
		Type: tea.String(*requsetmap.Type),
		Value: tea.String(clientip),
	  }
	  _, _err = client.UpdateDomainRecord(updateDomainRecordRequest)
	  if _err != nil {
		return _err
	  }
	fmt.Printf("dns修改完成")
	return _err
}

func main() {
	LABLE:
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Duration(5)*time.Minute)
	goto LABLE
}
