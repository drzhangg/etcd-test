package common

import "errors"

var (
	ERR_LOCAL_ALREADY_REQUIRED = errors.New("锁🔐已被占用")

	ERR_NO_LOCAL_IP_FOUND = errors.New("没有找到网卡IP")
)
