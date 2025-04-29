package multiworld

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/civilrights3/go-derek-go/internal/queue"
)

func (a *ArchipelagoClient) handleRoomInfo(_ context.Context, b []byte) error {
	out := &RoomInfoMessage{}
	err := json.Unmarshal(b, &out)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("%+v\n", out)

	// determine data package updates needed
	updates := a.dataCache.getListOfUpdates(out.DataPackageChecksum)
	a.sendGetDataPackage(updates)

	// send slot info
	a.sendConnect()
	return nil
}

func (a *ArchipelagoClient) sendConnect() {
	body := ConnectMessage{
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

	a.messageChan <- body
	return
}

func (a *ArchipelagoClient) sendGetDataPackage(games []string) {
	if len(games) == 0 {
		return
	}

	body := GetDataPackageMessage{
		Cmd:   CmdGetDataPackage.String(),
		Games: games,
	}

	a.messageChan <- body
	return
}

func (a *ArchipelagoClient) handleDataPackage(_ context.Context, b []byte) error {
	out := &DataPackageMessage{}
	err := json.Unmarshal(b, &out)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = a.dataCache.updateCache(out)
	if err != nil {
		return err
	}

	return nil
}

func (a *ArchipelagoClient) handleConnected(_ context.Context, b []byte) error {
	out := &ConnectedMessage{}
	err := json.Unmarshal(b, &out)
	if err != nil {
		fmt.Println(err)
		return err
	}

	a.dataCache.setPlayers(out.Players, out.SlotInfo)
	return nil
}

func (a *ArchipelagoClient) handleConnectionRefused(_ context.Context, b []byte) error {
	return fmt.Errorf("connection refused: %s", b)
}

func (a *ArchipelagoClient) handleRoomUpdate(ctx context.Context, b []byte) error {
	return nil
}

func (a *ArchipelagoClient) handlePrintJSON(ctx context.Context, b []byte) error {
	out := &PrintJSONMessage{}
	err := json.Unmarshal(b, &out)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if out.Type == JSONDataTypeItemSend {
		transformed := queue.BroadcastMessage{
			Sender:     a.dataCache.GetPlayerNameForSlot(out.Item.Player),
			Receiver:   a.dataCache.GetPlayerNameForSlot(out.Receiving),
			Item:       a.dataCache.GetItemNameForIDForPlayer(out.Item.Item, out.Receiving),
			Location:   a.dataCache.GetLocationNameForIDForPlayer(out.Item.Location, out.Item.Player),
			Importance: out.Item.Flags,
		}

		queue.Queue.EnqueueMessage(transformed)
	}
	// ignore other message types
	return nil
}

func (a *ArchipelagoClient) handleInvalidPacket(ctx context.Context, b []byte) error {
	panic(fmt.Sprintf("%s\n", b))
}
