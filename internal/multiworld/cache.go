package multiworld

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type dataCache struct {
	fileRoot     string
	playersByID  map[int]Player
	playerToGame map[int]string
	games        map[string]saneGame
	checksums    map[string]string
}

// TODO consider a lock between read and write?

func newDataCache(fileRoot string) *dataCache {
	return &dataCache{
		playersByID:  make(map[int]Player),
		playerToGame: make(map[int]string),
		games:        make(map[string]saneGame),
		checksums:    make(map[string]string),
		fileRoot:     fileRoot,
	}
}

func (c *dataCache) loadCacheFromFS() error {
	gameList, err := os.ReadDir(c.fileRoot)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(c.fileRoot, fs.ModeDir)
			if err != nil {
				return fmt.Errorf("cannot create cache directory: %w", err)
			}
		}

		return fmt.Errorf("unable to read cache dir: %w", err)
	}

	for _, f := range gameList {
		gameName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		contents, err := os.ReadFile(filepath.Join(c.fileRoot, f.Name()))
		if err != nil {
			return fmt.Errorf("unable to read file %s: %w", f.Name(), err)
		}

		g := saneGame{}
		err = json.Unmarshal(contents, &g)
		if err != nil {
			return fmt.Errorf("unable to unmarshal contents of file %s: %w", f.Name(), err)
		}

		c.games[gameName] = g
		c.checksums[gameName] = g.Checksum
	}

	return nil
}

func (c *dataCache) setPlayers(players []Player, info map[string]SlotInfo) {
	c.playersByID = make(map[int]Player)
	for _, p := range players {
		c.playersByID[p.Slot] = p
	}

	for _, i := range info {
		for _, p := range players {
			if i.Name == p.Name {
				c.playerToGame[p.Slot] = i.Game
			}
		}
	}
}

func (c *dataCache) getListOfUpdates(games map[string]string) []string {
	var updates []string
	for name, check := range games {
		cs, ok := c.checksums[name]
		if !ok || cs != check {
			updates = append(updates, name)
		}
	}

	return updates
}

func (c *dataCache) updateCache(updates *DataPackageMessage) error {
	for name, g := range updates.Data.Games {
		c.games[name] = saneitizeGame(g)
		c.checksums[name] = g.Checksum
	}

	return c.saveCacheToFS()
}

func (c *dataCache) saveCacheToFS() error {
	for name, g := range c.games {
		b, err := json.Marshal(g)
		if err != nil {
			return fmt.Errorf("unable to marshal Game %s: %w", name, err)
		}

		err = os.WriteFile(filepath.Join(c.fileRoot, fmt.Sprintf("%s.json", name)), b, fs.ModePerm)
		if err != nil {
			return fmt.Errorf("unable to write file %s: %w", name, err)
		}
	}

	return nil
}

func (c *dataCache) GetPlayerNameForSlotStr(slot string) string {
	slotNum, _ := strconv.Atoi(slot)
	return c.GetPlayerNameForSlot(slotNum)
}

func (c *dataCache) GetPlayerNameForSlot(slot int) string {
	player, ok := c.playersByID[slot]
	if !ok {
		return fmt.Sprintf("%d", slot)
	}

	return player.Name
}

func (c *dataCache) GetLocationNameForIDForPlayer(locationID int, playerID int) string {
	gameName := c.playerToGame[playerID]
	gameDetails := c.games[gameName]
	return gameDetails.LocationIDToName[locationID]
}

func (c *dataCache) GetItemNameForIDForPlayer(itemID int, playerID int) string {
	gameName := c.playerToGame[playerID]
	gameDetails := c.games[gameName]
	return gameDetails.ItemNameToId[itemID]
}

type saneGame struct {
	LocationIDToName map[int]string `json:"location_id_to_name"`
	ItemNameToId     map[int]string `json:"item_name_to_id"`
	Checksum         string         `json:"checksum"`
}

func saneitizeGame(input Game) saneGame {
	locationIDToName := make(map[int]string)
	itemNameToId := make(map[int]string)

	for l, i := range input.LocationNameToId {
		locationIDToName[i] = l
	}

	for n, i := range input.ItemNameToId {
		itemNameToId[i] = n
	}

	return saneGame{
		LocationIDToName: locationIDToName,
		ItemNameToId:     itemNameToId,
		Checksum:         input.Checksum,
	}
}
