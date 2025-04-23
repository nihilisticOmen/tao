package test

import (
	"context"
	"testing"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestEtcdConnection(t *testing.T) {
	// 从配置文件中获取的地址可能有误，这里直接使用正确地址进行测试
	etcdAddrs := []string{"127.0.0.1:2379"}
	dialTimeout := 5 * time.Second

	t.Log("开始测试连接 etcd...")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdAddrs,
		DialTimeout: dialTimeout,
	})
	if err != nil {
		t.Fatalf("连接 etcd 失败: %v", err)
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			t.Fatalf("关闭 etcd 连接失败: %v", err)
		} else {
			t.Log("成功关闭 etcd 连接")
		}
	}(cli)

	// 使用超时上下文确保测试不会无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 测试写入操作
	testKey := "/test/connection"
	testValue := "test-value"
	t.Logf("尝试写入测试数据: %s -> %s", testKey, testValue)
	_, err = cli.Put(ctx, testKey, testValue)
	if err != nil {
		t.Fatalf("写入测试数据失败: %v", err)
	}

	// 测试读取操作
	t.Logf("尝试读取测试数据: %s", testKey)
	resp, err := cli.Get(ctx, testKey)
	if err != nil {
		t.Fatalf("读取测试数据失败: %v", err)
	}

	if len(resp.Kvs) > 0 {
		actual := string(resp.Kvs[0].Value)
		if actual != testValue {
			t.Errorf("读取数据不匹配, 期望: %s, 实际: %s", testValue, actual)
		} else {
			t.Logf("成功读取到数据: %s", actual)
		}
	} else {
		t.Error("未找到测试数据")
	}

	// 清理测试数据
	_, err = cli.Delete(ctx, testKey)
	if err != nil {
		t.Logf("清理测试数据失败: %v", err)
	}

	t.Log("etcd 连接测试成功！")
}
