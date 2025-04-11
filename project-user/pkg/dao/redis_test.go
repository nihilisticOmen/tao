package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRedisCache_Put(t *testing.T) {
	ctx := context.Background()

	t.Run("正常设置缓存", func(t *testing.T) {
		// 验证 Redis 连接
		pong, err := Rc.rdb.Ping(ctx).Result()
		assert.NoError(t, err, "Redis 连接失败")
		assert.Equal(t, "PONG", pong, "Redis 未返回 PONG")

		// 测试 Put 方法
		err = Rc.Put(ctx, "test_key", "test_value", time.Hour)
		assert.NoError(t, err)

		// 验证数据存在
		val, _ := Rc.rdb.Get(ctx, "test_key").Result()
		assert.Equal(t, "test_value", val)
		log.Println("测试数据已设置:-------as", val)
		// 清理测试数据
		//Rc.rdb.Del(ctx, "test_key")
	})

	t.Run("空键处理", func(t *testing.T) {
		err := Rc.Put(ctx, "", "value", time.Minute)
		assert.Error(t, err)
	})

	t.Run("过期时间测试", func(t *testing.T) {
		key := "expire_test"
		err := Rc.Put(ctx, key, "data", time.Second)
		assert.NoError(t, err)

		// 立即获取应存在
		exists, _ := Rc.rdb.Exists(ctx, key).Result()
		assert.Equal(t, int64(1), exists)

		// 等待过期
		time.Sleep(2 * time.Second)
		exists, _ = Rc.rdb.Exists(ctx, key).Result()
		assert.Equal(t, int64(0), exists)
	})
}

func TestRedisCache_Get(t *testing.T) {
	ctx := context.Background()

	t.Run("获取存在的键", func(t *testing.T) {
		Rc.rdb.Set(ctx, "exist_key", "exist_value", 0)
		val, err := Rc.Get(ctx, "exist_key")
		if err != nil {
			log.Println("获取键值失败:", err)
			return
		}
		log.Println("------!!!" + val + "--")
		assert.NoError(t, err)
		assert.Equal(t, "exist_value", val)
	})

	t.Run("获取不存在的键", func(t *testing.T) {
		_, err := Rc.Get(ctx, "non_exist_key")
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("错误类型键", func(t *testing.T) {
		// 设置哈希类型测试错误处理
		Rc.rdb.HSet(ctx, "wrong_type", "field", "value")
		_, err := Rc.Get(ctx, "wrong_type")
		assert.Error(t, err)
	})
}

func TestRedisConnection(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t.Run("正常连接测试", func(t *testing.T) {
		// 使用Ping命令验证连接
		pong, err := Rc.rdb.Ping(ctx).Result()
		assert.NoError(t, err, "连接应成功")
		assert.Equal(t, "PONG", pong, "应返回PONG响应")

		// 验证连接配置参数
		opt := Rc.rdb.Options()
		assert.Equal(t, "localhost:6379", opt.Addr, "地址配置验证")
		assert.Equal(t, "", opt.Password, "密码配置验证")
		assert.Equal(t, 0, opt.DB, "数据库编号验证")
	})

	t.Run("异常连接测试", func(t *testing.T) {
		// 创建错误配置的客户端
		badClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6380", // 错误端口
			Password: "wrongpass",
			DB:       1,
		})

		// 测试连接超时
		_, err := badClient.Ping(ctx).Result()
		assert.Error(t, err, "错误配置应返回连接失败")
		assert.Contains(t, err.Error(), "connect: connection refused", "错误类型验证")
	})

	t.Run("连接稳定性测试", func(t *testing.T) {
		// 模拟网络波动
		for i := 0; i < 3; i++ {
			_, err := Rc.rdb.Ping(ctx).Result()
			assert.NoError(t, err, "第%d次重连失败", i+1)
			time.Sleep(100 * time.Millisecond)
		}
	})
}
func TestMain(m *testing.M) {
	// 测试前初始化连接
	ctx := context.Background()

	// 清空测试数据库
	Rc.rdb.FlushDB(ctx)

	// 执行测试
	m.Run()

	// 测试后清理
	Rc.rdb.FlushDB(ctx)
}
