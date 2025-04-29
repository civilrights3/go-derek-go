package multiworld

import "github.com/civilrights3/go-derek-go/internal/queue"

type RawMsg struct {
	Type    ServerMessageType
	Payload []byte
}

type ServerMessageType string

const (
	CmdRoomInfo          ServerMessageType = "RoomInfo"
	CmdDataPackage       ServerMessageType = "DataPackage"
	CmdConnectionRefused ServerMessageType = "ConnectionRefused"
	CmdConnected         ServerMessageType = "Connected"
	CmdReceivedItems     ServerMessageType = "ReceivedItems"
	CmdLocationInfo      ServerMessageType = "LocationInfo"
	CmdRoomUpdate        ServerMessageType = "RoomUpdate"
	CmdPrintJSON         ServerMessageType = "PrintJSON"
	CmdInvalidPacket     ServerMessageType = "InvalidPacket"
	CmdBounced           ServerMessageType = "Bounced"
	CmdSetReply          ServerMessageType = "SetReply"
)

func (s ServerMessageType) String() string {
	return string(s)
}

type ClientMessageType string

const (
	CmdConnect        ClientMessageType = "Connect"
	CmdSync           ClientMessageType = "Sync"
	CmdLocationChecks ClientMessageType = "LocationChecks"
	CmdLocationScouts ClientMessageType = "LocationScouts"
	CmdStatusUpdate   ClientMessageType = "StatusUpdate"
	CmdSay            ClientMessageType = "Say"
	CmdGetDataPackage ClientMessageType = "GetDataPackage"
	CmdBounce         ClientMessageType = "Bounce"
	CmdGet            ClientMessageType = "Get"
	CmdSet            ClientMessageType = "Set"
	CmdSetNotify      ClientMessageType = "SetNotify"
)

func (s ClientMessageType) String() string {
	return string(s)
}

type DataPacket []Message

type Message struct {
	Cmd ServerMessageType `json:"cmd"`
}

type RoomInfoMessage struct {
	Version             Version           `json:"version"`
	PasswordReqd        bool              `json:"password"`
	DataPackageChecksum map[string]string `json:"datapackage_checksums"`
	Games               []string          `json:"games"`
}

type ConnectMessage struct {
	Cmd           string   `json:"cmd"`
	Password      *string  `json:"password"`
	Name          string   `json:"name"`
	Version       Version  `json:"version"`
	Tags          []string `json:"tags"`
	ItemsHandling int      `json:"items_handling"`
	Uuid          string   `json:"uuid"`
	Game          string   `json:"game"`
}

type ConnectedMessage struct {
	Cmd      string              `json:"cmd"`
	Team     int                 `json:"team"`
	Slot     int                 `json:"slot"`
	Players  []Player            `json:"players"`
	SlotInfo map[string]SlotInfo `json:"slot_info"`
}

type GetDataPackageMessage struct {
	Cmd   string   `json:"cmd"`
	Games []string `json:"games"`
}

type DataPackageMessage struct {
	Data GameData `json:"data"`
}

type GameData struct {
	Games map[string]Game `json:"games"`
}

type Version struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Build uint64 `json:"build"`
	Class string `json:"class"`
}

type Player struct {
	Team  int    `json:"team"`
	Slot  int    `json:"slot"`
	Alias string `json:"alias"`
	Name  string `json:"name"`
	Class string `json:"class"`
}

type SlotInfo struct {
	Name  string `json:"name"`
	Game  string `json:"Game"`
	Type  int    `json:"type"`
	Class string `json:"class"`
}

const JSONDataTypeItemSend = "ItemSend"

type PrintJSONMessage struct {
	Cmd       string            `json:"cmd"`
	Data      []JSONDataElement `json:"data"`
	Type      string            `json:"type"`
	Item      JSONItem          `json:"item"`
	Receiving int               `json:"receiving"`
}

type JsonDataItemType string

const (
	DataTypePlayerID   JsonDataItemType = "player_id"
	DataTypeItemID     JsonDataItemType = "item_id"
	DataTypeLocationID JsonDataItemType = "location_id"
)

type JSONDataElement struct {
	Text   string           `json:"text"`
	Player int              `json:"Player"`
	Flags  int              `json:"flags"`
	Type   JsonDataItemType `json:"type"`
}

type JSONItem struct {
	Item     int                      `json:"item"`
	Location int                      `json:"location"`
	Player   int                      `json:"Player"`
	Flags    queue.ItemImportanceFlag `json:"flags"`
	Class    string                   `json:"class"`
}

type Game struct {
	LocationNameToId map[string]int `json:"location_name_to_id"`
	ItemNameToId     map[string]int `json:"item_name_to_id"`
	Checksum         string         `json:"checksum"`
}

//[
//  {
//    "cmd": "PrintJSONMessage",
//    "data": [
//      { "text": "1", "type": "player_id" },
//      { "text": " found their " },
//      { "text": "77771037", "Player": 1, "flags": 1, "type": "item_id" },
//      { "text": " (" },
//      { "text": "3790429", "Player": 1, "type": "location_id" },
//      { "text": ")" }
//    ],
//    "type": "ItemSend",
//    "receiving": 1,
//    "item": {
//      "item": 77771037,
//      "location": 3790429,
//      "Player": 1,
//      "flags": 1,
//      "class": "NetworkItem"
//    }
//  }
//]
