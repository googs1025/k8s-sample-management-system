## k8s 资源对象的简易后台管理系统
### 项目思路与功能
功能：预计提供k8s资源的增删改查。
![](https://github.com/googs1025/k8s-sample-management-system/blob/main/images/%E6%B5%81%E7%A8%8B%E5%9B%BE.jpg?raw=true)
#### 查询功能：
**workload**
1. pod
2. deployment
3. statefulset
4. job
5. cronjob

**命名空间**
1. namespace

**服务发现**
1. service
2. ingress

**配置文件**
1. configmap
2. secret
### 项目启动
1. 进入把目标集群的.kube/config文件放入项目根目录
```
➜  k8s-Management-System git:(main) ls -a | grep config
config
```
config文件示例(部分...)
```
apiVersion: v1
clusters:
- cluster:
    server: https://xxxxxxxxx:6443
  name: kubernetes
contexts:
.......
```
2. 加入远程node连接的配置文件app.yaml
* 根目录下加入app.yaml，内容如下，分别填入节点名、ip、用户与密码
```
k8s:
  nodes:
    - name: xxxxxxx
      ip: xxxxxxx
      user: xxxxx
      pass: xxxxxx
    - name: xxxxxxx
      ip: xxxxxxx
      user: xxxxxxx
      pass: xxxxxxx
    - name: xxxxxxx
      ip: xxxxxxx
      user: xxxxxxx
      pass: xxxxxxx
```


3. 启动server
```bigquery
➜  k8s-Management-System git:(main) ✗  go run main.go
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

2022/12/03 20:27:16 open /Users/zhenyu.jiang/go/src/golanglearning/new_project/k8s-Management-System/application.yaml: no such file or directory
[GIN-debug] GET    /deployments              --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /pods                     --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /jobs                     --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /services                 --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /namespaces               --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] POST   /vue-admin-template/user/login --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] POST   /vue-admin-template/user/logout --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /vue-admin-template/user/info --> github.com/shenyisyn/goft-gin/goft.StringResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /ws                       --> github.com/shenyisyn/goft-gin/goft.StringResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /statefulsets             --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] GET    /cronjobs                 --> github.com/shenyisyn/goft-gin/goft.JsonResponder.RespondTo.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080

```
4. 使用postman调用接口
调用接口方法
```
查询列表接口
http://localhost:8080/deployments
http://localhost:8080/pods
http://localhost:8080/services
http://localhost:8080/jobs
http://localhost:8080/cronjobs
http://localhost:8080/configmaps
http://localhost:8080/statefulsets
...

```
