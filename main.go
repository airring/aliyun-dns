package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"net/http"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	domain20180129 "github.com/alibabacloud-go/domain-20180129/client"
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
			fmt.Println("获取到的外网ip为:", string(content))
			return string(content)
		}
/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *domain20180129.Client, _err error) {
	config := &openapi.Config{}
	// 您的AccessKey ID
	config.AccessKeyId = accessKeyId
	// 您的AccessKey Secret
	config.AccessKeySecret = accessKeySecret
	// 访问的域名
	config.Endpoint = tea.String("domain.aliyuncs.com")
	_result = &domain20180129.Client{}
	_result, _err = domain20180129.NewClient(config)
	return _result, _err
}

func _main(args []*string) (_err error) {
	var InstanceId string
	var access conf
	access.getConf()
	// fmt.Println("accessKeyId:", access)
	client, _err := CreateClient(tea.String(access.AccessKeyId), tea.String(access.AccessKeySecret))
	if _err != nil {
		return _err
	}

	queryDomainListRequest := &domain20180129.QueryDomainListRequest{
		PageNum:  tea.Int32(1),
		PageSize: tea.Int32(5),
	}
	urlrequest, _err := client.QueryDomainList(queryDomainListRequest)
	if _err != nil {
		return _err
	}
	// 获取域名id
	domainlist := urlrequest.Body.Data.Domain
	for _, n := range domainlist {

		if *n.DomainName == access.Domain {
			InstanceId = *n.InstanceId
		}
	}
	fmt.Println(InstanceId)
	// 获取ip地址
	clientip := get_external()

	// 修改dns地址
	saveSingleTaskForModifyingDnsHostRequest := &domain20180129.SaveSingleTaskForModifyingDnsHostRequest{
		InstanceId: tea.String(InstanceId),
		DnsName: tea.String(access.Dnsdomain),
		Ip: []*string{tea.String(clientip)},
	  }
	  _, _err = client.SaveSingleTaskForModifyingDnsHost(saveSingleTaskForModifyingDnsHostRequest)
	  if _err != nil {
		return _err
	  }

	return _err
}

func main() {
	err := _main(tea.StringSlice(os.Args[1:]))
	if err != nil {
		panic(err)
	}

}
