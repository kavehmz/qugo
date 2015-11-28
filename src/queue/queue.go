// Queue is a simple Queue system written in Go that will use Redis underneath.
// Focus of this design is mainly horisontal scalability via concurrency, paritioning and fault-detection
// For scalability we can separately add more Redis instances and more Analysers.
//
// Partitioning factor to identify where each event must save is the OrderID.
// RedisParition => OrderId % redisParitions Queue => (OrderID/redisParitions) % queuePartitions
// AnalysePool will send all events of one OrderID to the same AnalyserWorker.
package queue

import (
	"regexp"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

var queuePartitions = 1

var redisParitions = 1

//Pool of redis connections
type redisStruct struct {
	conn redis.Conn
	url  string
}

var redisPool []redisStruct

//QueuesInPartision set number of queue in each partition. Each analyser will work on one queue in one partition and start its workers.
// Settign number of queues. This paritioning is for safely distributing related task into one queue
// If there is 2 partitions all events related to task 1,3,5.. will go to WAREHOUSE_1 and all events related to tasks 2,4,6,.. will go to WAREHOUSE_0
func QueuesInPartision(n int) {
	queuePartitions = n
}

// Partitions will define the number of redis partitions required for queue
// Redis has a singe processor implementations. If one redis can't handle a load of events we can use more than one
func Partitions(urls []string) {
	redisParitions = len(urls)
	redisPool = redisPool[:0]
	for _, v := range urls {
		r, _ := redis.DialURL(v)
		redisPool = append(redisPool, redisStruct{r, v})
	}
}

// AddTask will add a task event to the queue of tasks
func AddTask(id int, task string) {
	task = strconv.Itoa(id) + ";" + task
	_, e := redisPool[id%redisParitions].conn.Do("RPUSH", "WAREHOUSE_"+strconv.Itoa((id/redisParitions)%queuePartitions), task)
	checkErr(e)
}

func waitforSuccess(n int, id int, success chan bool, pool map[int]chan string) {
	redisdb, _ := redis.DialURL(redisPool[id%redisParitions].url)
	redisdb.Do("SET", "PENDING::"+strconv.Itoa(id), 1)
	r := <-success
	if r {
		delete(pool, id)
		redisdb.Do("DEL", "PENDING::"+strconv.Itoa(id))
	}
}

func removeTask(redisdb redis.Conn, queue string) (int, string) {
	r, e := redisdb.Do("LPOP", queue)
	checkErr(e)
	if r != nil {
		s, _ := redis.String(r, e)
		m := regexp.MustCompile(`(\d+);(.*)$`).FindStringSubmatch(s)
		id, _ := strconv.Atoi(m[1])
		redisdb.Do("SET", "PENDING::"+strconv.Itoa(id), 1)
		return id, m[2]
	}
	return 0, ""
}

//AnalysePool accepts an analyser function and empty the pool
func AnalysePool(n int, poolSize int, exitOnEmpy bool, f func(int, chan string, chan bool, chan bool)) {
	redisdb := redisPool[n%redisParitions].conn
	queue := "WAREHOUSE_" + strconv.Itoa((n/redisParitions)%queuePartitions)
	next := make(chan bool, poolSize)
	pool := make(map[int]chan string)
	for {

		id, task := removeTask(redisdb, queue)

		if task == "" {
			if exitOnEmpy {
				break
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			if pool[id] == nil {
				pool[id] = make(chan string)
				success := make(chan bool)
				go f(id, pool[id], success, next)
				go waitforSuccess(n, id, success, pool)
				pool[id] <- task
				next <- true
			} else {
				pool[id] <- task
			}
		}
	}

	for i := 0; i < poolSize; i++ {
		next <- true
	}
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
