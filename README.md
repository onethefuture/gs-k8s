
#                                               gs-k8s 服务

此服务的主要功能是将配置上传到 Apollo，并执行发布操作。此外，它还可以从 Kubernetes（k8s）集群中获取指定服务的版本信息。

## 目录

- 快速入门
  - [先决条件](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#先决条件)
  - [安装](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#安装)
  - [配置](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#配置)
- 使用
  - [上传配置](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#上传配置)
  - [发布配置](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#发布配置)
  - [获取服务版本](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#获取服务版本)
- [贡献](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#贡献)
- [许可证](https://chat.openai.com/c/12c09670-0c53-4758-bbfa-8b5ee1a29010#许可证)

## 快速入门

### 先决条件

开始之前，请确保已安装以下内容：

- Go（版本 1.21.4）

### 安装

克隆仓库并构建项目：

```shell
git clone https://github.com/onethefuture/gs-k8s.git
cd gs-k8s
go build
```

### 配置

通过更新 `config.yaml` 文件配置服务，其中包含 Apollo 凭据和其他必要的设置。

## 使用

### 上传配置

使用以下命令将配置上传到 Apollo：

```shell
curl --location --request GET '127.0.0.1/gstrain/getconf' \
--header 'Content-Type: application/json' \
--data '{}'
```

### 获取服务版本

从 Kubernetes 集群中检索指定服务的版本信息：

```shell
curl --location '127.0.0.1/gstrain/version' \
--header 'Content-Type: application/json' \
--data '{
    "servicename": "spider-man" 
}'
```

## 贡献

欢迎贡献！如果您发现了错误或有功能请求，请打开问题或提交拉取请求。

## 许可证

本项目根据 [MIT 许可证](https://chat.openai.com/c/LICENSE) 许可。
