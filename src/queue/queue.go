package queue

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

// Settign number of queues. This paritioning is for safely distributing related task into one queue
// If there is 2 partitions all events related to task 1,3,5.. will go to WAREHOUSE_1 and all events related to tasks 2,4,6,.. will go to WAREHOUSE_0
var queuePartitions = 1

// StorageParitions will define the number of redis partitions required for queue
// Redis has a singe processor implementations. If one redis can't handle a load of events we can use more than one
var redisParitions = 1

//Pool of redis connections
var redisPool []redis.Conn

//QueuesInPartision set number of queue in each partition. Each analyser will work on one queue in one partition and start its workers
func QueuesInPartision(n int) {
	queuePartitions = n
}

func Partitions(n int) {
	redisPool = redisPool[:0]
	for i := 0; i < n; i++ {
		r, _ := redis.DialURL(defaultRedis)
		redisPool = append(redisPool, r)
	}

}

// AddTask will add a task event to the queue of tasks
func AddTask(id int, task string) {
	task = strconv.Itoa(id) + ";" + task
	_, e := redisPool[id%redisParitions].Do("RPUSH", "WAREHOUSE_"+strconv.Itoa((id/redisParitions)%queuePartitions), task)
	checkErr(e)
}

func FetchPool(n int, f func(string) string) {
	//	redisdb := redisPool[n%redisParitions]
	//	queue := "WAREHOUSE_" + strconv.Itoa((n/redisParitions)%queuePartitions)

}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

var defaultRedis = "redis://redisqueue.kaveh.me:6379"
