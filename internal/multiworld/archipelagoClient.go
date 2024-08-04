package multiworld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/civilrights3/go-derek-go/internal/config"
	"net/url"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"syscall"
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
	}
}

func (a *ArchipelagoClient) Start(ctx context.Context, address string, port string, slot string) error {
	a.connection = connection{
		name:     slot,
		password: "", // TODO do i bother supporting passwords?
		address: url.URL{
			Scheme: "ws",
			Host:   fmt.Sprintf("%s:%s", address, port),
		},
	}

	for {
		if a.socket == nil {
			err := a.Connect(ctx)
			if err != nil {
				return err
			}
		}
	}
}

func (a *ArchipelagoClient) Connect(ctx context.Context) error {
	c, _, err := websocket.Dial(ctx, a.connection.address.String(), &websocket.DialOptions{
		CompressionMode: websocket.CompressionDisabled,
	})
	if err != nil && errors.Is(err, syscall.ECONNREFUSED) {
		fmt.Printf("%+v\n", err)
		return err
	}

	a.socket = c

	return nil
}

func (a *ArchipelagoClient) Disconnect(ctx context.Context) {
	sock := a.socket
	a.socket = nil

	err := sock.CloseNow()
	if err != nil && websocket.CloseStatus(err) != websocket.StatusNormalClosure {
		fmt.Printf("Error closing websocket: %s\n", err) // TODO move to a logger
	}
}

func (a *ArchipelagoClient) Read(ctx context.Context) ([]byte, error) {
	go a.startReadLoop(ctx)
	return nil, nil
}

func (a *ArchipelagoClient) startReadLoop(ctx context.Context) {
	dataChan := make(chan []byte)
	go func() {
		for {
			_, b, err := a.socket.Read(ctx)
			if err != nil && websocket.CloseStatus(err) != websocket.StatusNormalClosure {
				a.Disconnect(ctx)
				return
			}

			select {
			case <-ctx.Done():
				return
			case dataChan <- b:
				continue
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			err := a.socket.Close(websocket.StatusNormalClosure, "")
			if err != nil {
				fmt.Printf("error closing websocket: %s\n", err)
			}
			return
		case b := <-dataChan:
			err := a.handleMessage(ctx, b)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
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

func (a *ArchipelagoClient) handleMessage(ctx context.Context, msg []byte) error {
	fmt.Println("New messages, Sir!")
	fmt.Printf("%s\n", msg)
	msgs, err := a.parse(msg)
	if err != nil {
		return err
	}

	for _, m := range msgs {
		switch m.Type {
		case CmdRoomInfo:
			return a.handleRoomInfo(ctx, m.Payload)
		case CmdConnected:
		case CmdConnectionRefused:
		case CmdRoomUpdate:
		case CmdPrintJSON:

		default:
			fmt.Printf("Thats it! I've come up with a new type! \n%s\n", m.Payload)
			//return errors.New("STOP PLS")
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
