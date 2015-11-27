package queue

import (
	"testing"

	"github.com/garyburd/redigo/redis"
)

func TestSetRedisPool(t *testing.T) {
	Partitions(5)
	_, err := redisPool[4].Do("PING")
	if err != nil {
		t.Error("SetRedisPool items are not set correctly")
	}

}

func TestAddTask(t *testing.T) {
	QueuesInPartision(1)
	Partitions(1)
	redisdb := redisPool[0]
	redisdb.Do("DEL", "WAREHOUSE_0")
	AddTask(4, "test")
	r, e := redisdb.Do("RPOP", "WAREHOUSE_0")
	s, e := redis.String(r, e)
	if s != "4;test" {
		t.Error("Task is stored incorrectly: ", s)
	}
}

func TestQueuesInPartision(t *testing.T) {
	QueuesInPartision(5)
	Partitions(1)
	redisdb := redisPool[0]
	redisdb.Do("DEL", "WAREHOUSE_4")
	AddTask(4, "test")
	r, e := redisdb.Do("RPOP", "WAREHOUSE_4")
	s, e := redis.String(r, e)
	if s != "4;test" {
		t.Error("Task is stored incorrectly: ", s)
	}

}

func BenchmarkPrimes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AddTask(i, "test")
	}
}
