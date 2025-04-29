package chat

import (
	"fmt"
	"github.com/civilrights3/go-derek-go/internal/queue"
)

type textHandler func(message queue.BroadcastMessage, isSelfFind bool) string

func formatPlainMessage(msg queue.BroadcastMessage, isSelfFind bool) string {
	if isSelfFind {
		return fmt.Sprintf("[%s] found their <%s> (%s)", msg.Receiver, msg.Item, msg.Location)
	}

	return fmt.Sprintf("[%s] sent <%s> to {%s} (%s)", msg.Sender, msg.Item, msg.Receiver, msg.Location)
}

func formatMonospacedMessage(msg queue.BroadcastMessage, isSelfFind bool) string {
	if isSelfFind {
		return fmt.Sprintf("`[%s] found their <%s> (%s)`", msg.Receiver, msg.Item, msg.Location)
	}

	return fmt.Sprintf("`[%s] sent <%s> to {%s} (%s)`", msg.Sender, msg.Item, msg.Receiver, msg.Location)
}

const (
	ColorNeutral = `[0m`
	ColorGold    = `[3;33m`
	ColorWhite   = `[3;37m`
	ColorMagenta = `[3;35m`
	ColorBlue    = `[3;34m`
	ColorRed     = `[3;31m`
	ColorTeal    = `[3;36m`
)

var (
	importanceToColor = map[queue.ItemImportanceFlag]string{
		queue.ItemNormal:      ColorWhite,
		queue.ItemHelpful:     ColorBlue,
		queue.ItemProgression: ColorMagenta,
		queue.ItemTrap:        ColorRed,
	}
)

func formatColorMessage(msg queue.BroadcastMessage, isSelfFind bool) string {
	if isSelfFind {
		return fmt.Sprintf("```ansi\n%s[%s]%s found their %s<%s> %s(%s)\n```", ColorGold, msg.Receiver, ColorNeutral, importanceToColor[msg.Importance], msg.Item, ColorTeal, msg.Location)
	}

	return fmt.Sprintf("```ansi\n%s[%s]%s sent %s<%s>%s to %s{%s} %s(%s)\n```", ColorGold, msg.Sender, ColorNeutral, importanceToColor[msg.Importance], msg.Item, ColorNeutral, ColorGold, msg.Receiver, ColorTeal, msg.Location)
}
