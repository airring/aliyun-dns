# aliyun-dns

通过阿里云公共api调用自动修改dns解析信息,实现ddns功能

## 配置修改
修改[config.yaml](./config.yaml)文件中

Domain为需要修改的二级域名

Dnsdomain为需要修改的dns的三级域名

accessKeyId 和 accessKeySecret需要从阿里云 用户AccessKey下啦菜单中申请,详细可以进入

[如何申请AccessKey](https://help.aliyun.com/knowledge_detail/63482.html)

## 快速运行

```
git clone https://github.com/airring/aliyun-dns
cd aliyun-dns
# 修改 config.yaml (pass)
nohup main &
```
