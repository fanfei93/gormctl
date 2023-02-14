gormctl 是一款可根据mysql表结构自动生成gorm模型的工具。
构建：
1. 编译
在根目录下运行
```
go build .
```
2. 为生成的gormctl可执行文件增加执行权限
```
chmod +x gormctl
```
3. 将gormctl所在目录添加至环境变量中
根据系统自行设置

使用：
```
gormctl -H dbHost -U dbUsername -P password -D dbName -T "tableName" -O "modelDir"
```
dbHost: 数据库host
dbUsername: 用户名
password: 密码
dbName: 数据库名
tableName: 数据表名
modelDir: 生成的模型文件存放目录