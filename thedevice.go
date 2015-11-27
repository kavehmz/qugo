package main

import (
	"encoding/json"
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
	for id := 1; id <= 1000; id++ {
		jsonVal, _ := json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Start", OrderID: id, ItemID: 0, Quantity: 0, Container: 0, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Pick", OrderID: id, ItemID: 1100, Quantity: 1, Container: 5, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Skip", OrderID: id, ItemID: 1100, Quantity: 1, Container: 5, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
		jsonVal, _ = json.Marshal(Event{Username: "system", Timestamp: time.Now().Unix(), Event: "Stop", OrderID: id, ItemID: 0, Quantity: 0, Container: 0, PicklistID: id})
		queue.AddTask(id, string(jsonVal))
	}
}

func analyse(id string) bool {
	return true
}

func main() {
	queue.Partitions(1)
	queue.QueuesInPartision(1)
	generateRandomEvents()
	queue.FetchPool(1, analyse("d"))
}
