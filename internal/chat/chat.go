package chat

import "context"

type ChatClient interface {
	SendMessage(context.Context, []byte) error
}
