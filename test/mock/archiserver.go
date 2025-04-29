package mock

import (
	"context"
	"github.com/civilrights3/go-derek-go/internal/queue"
)

type ArchiServer struct {
}

var (
	testMessages = []queue.BroadcastMessage{
		{
			Sender:     "Civil",
			Receiver:   "Tea",
			Item:       "A bag full of math rocks",
			Location:   "Under the couch",
			Importance: queue.ItemNormal,
		},
		{
			Sender:     "Tea",
			Receiver:   "Nintendale",
			Item:       "Way too many checks",
			Location:   "Somewhere in Canada",
			Importance: queue.ItemProgression,
		},
		{
			Sender:     "Salty",
			Receiver:   "EOG",
			Item:       "Turkey sandwich",
			Location:   "The kitchen",
			Importance: queue.ItemHelpful,
		},
		{
			Sender:     "Iruga",
			Receiver:   "Iruga",
			Item:       "A backflip into the void",
			Location:   "The Navel",
			Importance: queue.ItemTrap,
		},
	}
)

func SendTestMessages(_ context.Context) {
	for _, msg := range testMessages {
		queue.Queue.EnqueueMessage(msg)
	}
}
