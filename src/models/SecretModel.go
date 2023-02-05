package models

var SecretType map[string]string

func init() {
	SecretType = map[string]string{
		"Opaque":"自定义类型",
		"kubernetes.io/service-account-token": "服务账号令牌",
		"kubernetes.io/dockercfg": "docker配置",
		"kubernetes.io/dockerconfigjson": "docker配置(JSON)",
		"kubernetes.io/basic-auth": "Basic认证凭据",
		"kubernetes.io/ssh-auth": "SSH凭据",
		"kubernetes.io/tls": "TLS凭据",
		"bootstrap.kubernetes.io/token": "启动引导令牌数据",
	}
}

type SecretModel struct {
	Name string
	NameSpace string
	CreateTime string
	Type string  //类型
	Data map[string][]byte  // KV
}

// 提交用
type PostSecretModel struct {
	Name string
	NameSpace string
	Type string  //类型
	Data map[string]string
}
