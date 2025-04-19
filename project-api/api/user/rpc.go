package user

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"project-api/config"
	"project-common/discovery"
	"project-common/logs"
	loginServiceV1 "project-user/pkg/service/login.service.v1"
)

var LoginServiceClient loginServiceV1.LoginServiceClient

func InitRpcUserClient() {
	zap.L().Info("初始化rpc客户端")
	etcdRegister := discovery.NewResolver(config.AppConf.EtcdConfig.Addrs, logs.LG)
	resolver.Register(etcdRegister)

	// 使用NewClient方法替代DialContext
	conn, err := grpc.NewClient(
		"etcd:///user",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)

	if err != nil {
		log.Fatalf("连接服务失败: %v", err)
	}

	LoginServiceClient = loginServiceV1.NewLoginServiceClient(conn)
	zap.L().Info("rpc客户端初始化成功")
}
