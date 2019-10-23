package master

import (
	"encoding/json"
	"github.com/drzhangg/etcd-test/prepare/crontab/common"
	"net"
	"net/http"
	"strconv"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

// 保存任务接口
// POST job={"name": "job1", "command": "echo hello", "cronExpr": "* * * * *"}
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	var (
		err     error
		bytes   []byte
		postJob string
		job     common.Job
		oldJob  *common.Job
	)
	//解析post表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	//获取表单中的job字段
	postJob = req.PostForm.Get("job")

	//反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	//将job信息保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	//返回正常应答
	if bytes, err = common.BuildResponse(0, "succ", oldJob); err == nil {
		resp.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		httpServer *http.Server
		listener   net.Listener
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	httpServer = &http.Server{
		Handler: mux,
	}

	G_apiServer = &ApiServer{httpServer: httpServer}

	go httpServer.Serve(listener)
	return
}
