package multiworld

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/civilrights3/go-derek-go/internal/config"
	"net/url"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
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
}

type connection struct {
	name     string
	password string
	address  url.URL
}

func NewArchipelagoClient(cfg config.Multiworld) *ArchipelagoClient {
	return &ArchipelagoClient{
		clientID:      cfg.ClientID,
		clientVersion: semver.MustParse(cfg.ClientVersion),
		maxRetry:      time.Duration(cfg.MaxConnectionRetry) * time.Second,
		minRetry:      1 * time.Second,
	}
}

func (a *ArchipelagoClient) Start(ctx context.Context, address string, port string, slot string) {
	a.connection = connection{
		name:     slot,
		password: "", // TODO do i bother supporting passwords?
		address: url.URL{
			Scheme: "ws",
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
			if a.socket == nil {
				a.connect(ctx)
			}
			a.readLoop(ctx)
			a.disconnect(ctx)
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

func (a *ArchipelagoClient) connect(ctx context.Context) {
	currentRetry := a.minRetry

	for {
		timer := time.NewTimer(currentRetry)
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			c, _, err := websocket.Dial(ctx, a.connection.address.String(), &websocket.DialOptions{
				CompressionMode: websocket.CompressionDisabled,
			})
			if err == nil {
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
	msgs, err := a.parse(msg)
	if err != nil {
		return err
	}

	for _, m := range msgs {
		fmt.Printf("received command %s\n", m.Type)
		switch m.Type {
		case CmdRoomInfo:
			return a.handleRoomInfo(ctx, m.Payload)
		case CmdConnected:
		case CmdConnectionRefused:
		case CmdRoomUpdate:
		case CmdPrintJSON:

		default:
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

func (a *ArchipelagoClient) handleRoomInfo(ctx context.Context, b []byte) error {
	out := &RoomInfoMessage{}
	err := json.Unmarshal(b, &out)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("%+v\n", out)

	// TODO determine data package updates needed

	// send slot info
	return a.sendConnect(ctx) // TODO should sends be in a separate worker? is that overkill?
}

func (a *ArchipelagoClient) sendConnect(ctx context.Context) error {
	body := ConnectPacket{
		Cmd:  CmdConnect.String(),
		Name: a.connection.name,
		Version: Version{
			Major: a.clientVersion.Major(),
			Minor: a.clientVersion.Minor(),
			Build: a.clientVersion.Patch(),
			Class: "Version",
		},
		Uuid:          a.clientID,
		ItemsHandling: 0b011,
		Tags:          []string{"TextOnly", "IgnoreGame", "AP", "Derek"},
	}

	return a.send(ctx, body)
}

func (a *ArchipelagoClient) send(ctx context.Context, message interface{}) error {
	msgs := []interface{}{message}
	fmt.Println("Sending message")
	b, _ := json.Marshal(msgs)
	fmt.Printf("%s\n", b)
	err := wsjson.Write(ctx, a.socket, msgs)
	if err != nil {
		return err
	}
	fmt.Println("Sent message")
	return nil
}
