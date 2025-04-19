package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// 从配置文件中获取的地址可能有误，这里直接使用正确地址进行测试
	etcdAddrs := []string{"127.0.0.1:2379"}
	dialTimeout := 5 * time.Second

	fmt.Println("开始测试连接 etcd...")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddrs,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatalf("连接 etcd 失败: %v", err)
	}
	defer cli.Close()

	// 使用超时上下文确保测试不会无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 测试写入操作
	testKey := "/test/connection"
	testValue := "test-value"
	fmt.Printf("尝试写入测试数据: %s -> %s\n", testKey, testValue)
	_, err = cli.Put(ctx, testKey, testValue)
	if err != nil {
		log.Fatalf("写入测试数据失败: %v", err)
	}

	// 测试读取操作
	fmt.Printf("尝试读取测试数据: %s\n", testKey)
	resp, err := cli.Get(ctx, testKey)
	if err != nil {
		log.Fatalf("读取测试数据失败: %v", err)
	}

	if len(resp.Kvs) > 0 {
		fmt.Printf("成功读取到数据: %s\n", string(resp.Kvs[0].Value))
	} else {
		fmt.Println("未找到测试数据")
	}

	// 清理测试数据
	_, err = cli.Delete(ctx, testKey)
	if err != nil {
		log.Printf("清理测试数据失败: %v", err)
	}

	fmt.Println("etcd 连接测试成功！")
}
