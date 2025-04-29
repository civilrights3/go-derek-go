package queue

import (
	"fmt"
	"sync"
	"time"
)

const (
	delay = 250 * time.Millisecond
)

var (
	Queue *messageQueue
)

type messageQueue struct {
	queue            []BroadcastMessage
	messageListeners []func(message BroadcastMessage) error
	lock             sync.RWMutex
}

func StartMessageQueue() {
	q := &messageQueue{
		queue:            make([]BroadcastMessage, 0),
		messageListeners: make([]func(message BroadcastMessage) error, 0),
		lock:             sync.RWMutex{},
	}

	Queue = q

	go Queue.run()
}

func (m *messageQueue) run() {
	for {
		select {
		case <-time.After(delay):
			if len(m.messageListeners) > 0 && len(m.queue) > 0 {
				msg := m.GetNext()
				for _, l := range m.messageListeners {
					err := l(msg)
					if err != nil {
						fmt.Printf("error sending message %s\n", err.Error())
					}
				}
				m.AckNext()
			}
		}
	}
}

func (m *messageQueue) RegisterMessageListener(f func(message BroadcastMessage) error) {
	m.messageListeners = append(m.messageListeners, f)
}

func (m *messageQueue) EnqueueMessage(message BroadcastMessage) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.queue = append(m.queue, message)
}

func (m *messageQueue) GetNext() BroadcastMessage {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.queue[0]
}

func (m *messageQueue) AckNext() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.queue = m.queue[1:]
}

type ItemImportanceFlag int

const (
	ItemNormal      ItemImportanceFlag = 0
	ItemProgression ItemImportanceFlag = 0b001
	ItemHelpful     ItemImportanceFlag = 0b010
	ItemTrap        ItemImportanceFlag = 0b100
)

type BroadcastMessage struct {
	Sender     string
	Receiver   string
	Item       string
	Location   string
	Importance ItemImportanceFlag
}

func (m *messageQueue) TestHandler(message BroadcastMessage) error {
	fmt.Println("---------------------------")
	//fmt.Printf("%s %s %s %s\n", message.Sender, message.Receiver, message.Item, message.Location)
	fmt.Printf("%d\n", len(m.queue))
	fmt.Println("---------------------------")
	return nil
}
