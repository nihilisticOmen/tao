package discovery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

// Register for grpc server
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	leasesID    clientv3.LeaseID
	keepAliveCh <-chan *clientv3.LeaseKeepAliveResponse

	srvInfo Server
	srvTTL  int64
	cli     *clientv3.Client
	logger  *zap.Logger
}

// NewRegister create a register base on etcd
func NewRegister(etcdAddrs []string, logger *zap.Logger) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
		logger:      logger,
	}
}

// Register a service
func (r *Register) Register(srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error
	//校验服务地址合法性（如IP格式）。
	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}
	//创建etcd客户端连接。
	if r.cli, err = clientv3.New(clientv3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl
	//调用register()将服务信息写入etcd（带租约）。
	if err = r.register(); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})
	//启动keepAlive()协程维持租约。
	go r.keepAlive()
	//返回一个closeCh通道，用于外部控制注销。
	return r.closeCh, nil
}

// Stop stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 申请租约并写入服务信息
func (r *Register) register() error {
	leaseCtx, cancel := context.WithTimeout(context.Background(), time.Duration(r.DialTimeout)*time.Second)
	defer cancel()
	//通过cli.Grant方法获取租约响应
	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}
	//保存租约ID
	r.leasesID = leaseResp.ID
	zap.L().Info("租约获取成功", zap.Int64("lease_id", int64(leaseResp.ID)), zap.Int64("lease_ttl", leaseResp.TTL))
	//向etcd申请租约（TTL）。
	if r.keepAliveCh, err = r.cli.KeepAlive(context.Background(), leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	//将服务信息（JSON格式）写入etcd，绑定租约。
	//路径格式：/<服务名>/<版本>/<地址>（通过BuildRegPath生成）。
	_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister() error {
	_, err := r.cli.Delete(context.Background(), BuildRegPath(r.srvInfo))
	return err
}

// keepAlive监听两个事件：
// 租约续约响应,若失败则重新注册。
// 定时器：定期检查租约状态，异常时重新注册。
// 收到closeCh信号时，调用unregister()删除etcd中的服务信息。
func (r *Register) keepAlive() {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(); err != nil {
				r.logger.Error("unregister failed", zap.Error(err))
			}
			if _, err := r.cli.Revoke(context.Background(), r.leasesID); err != nil {
				r.logger.Error("revoke failed", zap.Error(err))
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(); err != nil {
					r.logger.Error("register failed", zap.Error(err))
				}
			}
		}
	}
}

// UpdateHandler return http handler
// 提供HTTP接口动态更新服务权重（如负载均衡调整）。
func (r *Register) UpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			_, err = r.cli.Put(context.Background(), BuildRegPath(r.srvInfo), string(data), clientv3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("update server weight success"))
	}
}

// GetServerInfo 从etcd查询当前服务的注册信息。
func (r *Register) GetServerInfo() (Server, error) {
	resp, err := r.cli.Get(context.Background(), BuildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}
