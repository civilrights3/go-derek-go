package multiworld

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
	Version      Version `json:"version"`
	PasswordReqd bool    `json:"password"`
}

type Version struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Build uint64 `json:"build"`
	Class string `json:"class"`
}

type ConnectPacket struct {
	Cmd           string   `json:"cmd"`
	Password      *string  `json:"password"`
	Name          string   `json:"name"`
	Version       Version  `json:"version"`
	Tags          []string `json:"tags"`
	ItemsHandling int      `json:"items_handling"`
	Uuid          string   `json:"uuid"`
	Game          string   `json:"game"`
}
