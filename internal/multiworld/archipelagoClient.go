package multiworld

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/civilrights3/go-derek-go/internal/config"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type MultiworldClient interface {
	Connect(ctx context.Context, address string, port string, slot string) error
	Disconnect(ctx context.Context) error
	Read(ctx context.Context) ([]byte, error)
}

type ArchipelagoClient struct {
	socket        *websocket.Conn
	clientID      string
	clientVersion *semver.Version
	connection    connection
	maxRetry      time.Duration
	minRetry      time.Duration
	dataCache     *dataCache
	messageChan   chan any
}

type connection struct {
	name     string
	password string
	address  url.URL
}

func NewArchipelagoClient(cfg config.Multiworld) (*ArchipelagoClient, error) {
	cache := newDataCache(cfg.Cache.Filepath)
	err := cache.loadCacheFromFS()
	if err != nil {
		// Yes this will crash the start of the application. If not it'll just have a crash loop later when saving caches
		return nil, fmt.Errorf("error loading cache from FS: %w", err)
	}

	return &ArchipelagoClient{
		clientID:      cfg.ClientID,
		clientVersion: semver.MustParse(cfg.ClientVersion),
		maxRetry:      time.Duration(cfg.MaxConnectionRetry) * time.Second,
		minRetry:      1 * time.Second,
		dataCache:     cache,
	}, nil
}

func (a *ArchipelagoClient) Start(ctx context.Context, address string, port string, slot string) {
	a.connection = connection{
		name:     slot,
		password: "", // TODO do i bother supporting passwords?
		address: url.URL{
			Scheme: "wss", // TODO make this smart enough to determine based on URL
			Host:   fmt.Sprintf("%s:%s", address, port),
		},
	}

	go a.startReadLoop(ctx)

	return
}

func (a *ArchipelagoClient) startReadLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.messageChan = make(chan any, 10)
			if a.socket == nil {
				a.connect(ctx)
			}

			go a.writeLoop(ctx)
			a.readLoop(ctx)

			a.disconnect(ctx)
			close(a.messageChan)
			a.messageChan = nil
		}
	}
}

func (a *ArchipelagoClient) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			a.disconnect(ctx)
			return
		default:
			b, err := a.readSock(ctx)
			if err != nil {
				fmt.Printf("error reading socket: %s\n", err)
				return
			}

			err = a.handleMessage(ctx, b)
			if err != nil {
				fmt.Printf("unable to handle message: %s\n", err)
				return
			}
		}
	}
}

func (a *ArchipelagoClient) writeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, open := <-a.messageChan:
			if !open || a.socket == nil {
				return
			}

			fmt.Println("Sending message")
			msgs := []interface{}{msg}
			b, err := json.Marshal(msgs)
			if err != nil {
				fmt.Printf("error marshalling message: %s\n", err)
			}
			fmt.Printf("%s\n", b)

			err = wsjson.Write(ctx, a.socket, msgs)
			if err != nil {
				fmt.Printf("error occured when sending message to server: %s\n", err)
				// TODO how do we retry sends? should we?
				return
			}
			fmt.Println("Sent message")
		}
	}
}

func (a *ArchipelagoClient) connect(ctx context.Context) {
	currentRetry := a.minRetry

	for {
		timer := time.NewTimer(currentRetry)
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			fmt.Println(a.connection.address.String())
			c, _, err := websocket.Dial(ctx, a.connection.address.String(), &websocket.DialOptions{
				CompressionMode: websocket.CompressionDisabled,
			})
			if err == nil {
				c.SetReadLimit(-1)

				a.socket = c
				return
			}

			currentRetry = currentRetry * 2
			if currentRetry > a.maxRetry {
				currentRetry = a.maxRetry
			}
			fmt.Printf("Failed to connect %s\n", err)
			fmt.Printf("Retry in %s seconds\n", currentRetry)
		}
	}
}

func (a *ArchipelagoClient) disconnect(ctx context.Context) {
	sock := a.socket
	a.socket = nil

	err := sock.CloseNow()
	if err != nil && websocket.CloseStatus(err) != websocket.StatusNormalClosure {
		fmt.Printf("Error closing websocket: %s\n", err) // TODO move to a logger
	}
}

func (a *ArchipelagoClient) readSock(ctx context.Context) (b []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("remote server suddenly closed: %v", r)
		}
	}()

	_, b, err = a.socket.Read(ctx)
	return
}

func (a *ArchipelagoClient) handleMessage(ctx context.Context, msg []byte) error {
	fmt.Printf("%s\n", msg)
	msgs, err := a.parse(msg)
	if err != nil {
		return err
	}

	for _, m := range msgs {
		fmt.Printf("received command %s\n", m.Type)
		switch m.Type {
		case CmdRoomInfo:
			return a.handleRoomInfo(ctx, m.Payload)
		case CmdDataPackage:
			return a.handleDataPackage(ctx, m.Payload)
		case CmdConnected:
			return a.handleConnected(ctx, m.Payload)
		case CmdConnectionRefused:
			return a.handleConnectionRefused(ctx, m.Payload)
		case CmdRoomUpdate:
			return a.handleRoomUpdate(ctx, m.Payload)
		case CmdPrintJSON:
			return a.handlePrintJSON(ctx, m.Payload)
		case CmdInvalidPacket:
			return a.handleInvalidPacket(ctx, m.Payload)
		default:
			fmt.Printf("unknown command: %s\n", m.Type)
			continue
		}
	}

	return nil
}

func (a *ArchipelagoClient) parse(msg []byte) ([]RawMsg, error) {
	out := make([]map[string]interface{}, 0)
	err := json.Unmarshal(msg, &out)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	parsed := []RawMsg{}
	for _, m := range out {
		raw, _ := json.Marshal(m)
		t := m["cmd"].(string)
		msg := RawMsg{
			Type:    ServerMessageType(t),
			Payload: raw,
		}
		parsed = append(parsed, msg)
	}

	return parsed, nil
}

// keeping just in case we need to debug the send loop later
//func (a *ArchipelagoClient) send(ctx context.Context, message interface{}) error {
//	msgs := []interface{}{message}
//	fmt.Println("Sending message")
//	b, _ := json.Marshal(msgs)
//	fmt.Printf("%s\n", b)
//	err := wsjson.Write(ctx, a.socket, msgs)
//	if err != nil {
//		return err
//	}
//	fmt.Println("Sent message")
//	return nil
//}
