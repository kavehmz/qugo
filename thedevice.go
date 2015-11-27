package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"./src/queue"
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

func generateRandomEvents() {
	for id := 1; id <= 2; id++ {
		jsonVal, _ := json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Start", OrderID: id, ItemID: 0, Quantity: 0, Container: 0, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Pick", OrderID: id, ItemID: 1100, Quantity: 1, Container: 5, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Skip", OrderID: id, ItemID: 1100, Quantity: 1, Container: 5, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Stop", OrderID: id, ItemID: 0, Quantity: 0, Container: 0, PicklistID: id})
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Int()%2 != 0 {
			queue.AddTask(id, string(jsonVal))
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
			fmt.Println("Analyse", event.OrderID, event.Event)
			if event.Event == "Start" {
				fmt.Println("Starting", event.OrderID)
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

func main() {
	queue.Partitions(1)
	queue.QueuesInPartision(1)
	generateRandomEvents()
	analyse := analyse
	queue.AnalysePool(1, false, analyse)
}
