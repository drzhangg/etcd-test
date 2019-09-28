etcd用于服务注册和服务发现


## 1.服务注册
* 向etcd注册自己的信息，etcd通过kv保存信息
* 服务注册的话要创建租约，以及进行续租等


## 2.服务发现
* 通过watch etcd的client，来监听etcd的添加、删除、修改等




## 扩展
使用mongodb存一些日志文件到etcd中
etcd的val存放一个struct