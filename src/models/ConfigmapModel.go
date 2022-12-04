package models

//提交 用的
type PostConfigMapModel struct {
	Name string
	NameSpace string
	Data map[string]string
}
//列表用
type ConfigMapModel struct {
	Name string
	NameSpace string
	CreateTime string
	Data map[string]string  // KV

}
