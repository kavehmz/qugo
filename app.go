package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/kavehmz/queue"
)

// Event defines warehouse event
type Event struct {
	Username   string
	Timestamp  int64
	Event      string
	OrderID    int
	ItemID     int
	Quantity   int
	Container  int
	PicklistID int
}

func generateRandomEvents(q queue.Queue, n int) {
	for id := 1; id <= n; id++ {
		jsonVal, _ := json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Start", OrderID: id, ItemID: 0, Quantity: 0, Container: 0, PicklistID: id})
		q.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Pick", OrderID: id, ItemID: 1100, Quantity: 1, Container: 5, PicklistID: id})
		q.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Skip", OrderID: id, ItemID: 1101, Quantity: 1, Container: 5, PicklistID: id})
		q.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Stop", OrderID: id, ItemID: 0, Quantity: 0, Container: 0, PicklistID: id})
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Int()%2 != 4 {
			q.AddTask(id, string(jsonVal))
		}
	}

}
func parseEvent(msg string) Event {
	var event Event
	json.Unmarshal([]byte(msg), &event)
	return event
}

func analyse(id int, msg_channel chan string, success chan bool, next chan bool) {
	for {
		select {
		case msg := <-msg_channel:
			event := parseEvent(msg)
			if event.Event == "Start" {
				fmt.Println("Start", event.OrderID)
			} else if event.Event == "Pick" {
				fmt.Println("Pick", event.OrderID, event.ItemID)
			} else if event.Event == "Skip" {
				fmt.Println("Skip", event.OrderID, event.ItemID)
			} else if event.Event == "Stop" {
				fmt.Println("Stop", event.OrderID)
				<-next
				success <- true
				return
			}

		case <-time.After(2 * time.Second):
			fmt.Println("new event for 2 seconds for orderID", id)
			<-next
			success <- false
			return
		}
	}
}

var mode *string
var id, insert, workers *int

func init() {
	// Application parameters the Go way
	mode = flag.String("mode", "", "Specfies the mode for this application [device|analyser].")
	insert = flag.Int("insert", 10, "Number of inserts into queue. Only useful if mode is device.")
	id = flag.Int("id", 1, "Specfies the ID of analyser. This will set which redis and which queue this analyser will handle.")
	workers = flag.Int("workers", 0, "Specfies number of concurrent workers which each analysers will have.")
	flag.Parse()
}

func main() {
	var q queue.Queue
	q.Urls([]string{"redis://redisqueue.kaveh.me:6379"})
	if *mode == "device" {
		generateRandomEvents(q, *insert)
	} else {
		if *workers != 0 {
			q.AnalyzerBuff = *workers
		}
		wsStress := analyse
		exitOnEmpty := func() bool {
			return false
		}
		q.AnalysePool(*id, exitOnEmpty, wsStress)
	}
}
