package chat

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/civilrights3/go-derek-go/internal/config"
	"github.com/civilrights3/go-derek-go/internal/queue"
)

type DiscordClient struct {
	discord          *discordgo.Session
	channelID        string
	guildID          string
	messageFormatter textHandler
}

var (
	formattingFuncs = map[config.DisplayMode]textHandler{
		config.DisplayPlain:      formatPlainMessage,
		config.DisplayMonospaced: formatMonospacedMessage,
		config.DisplayColor:      formatColorMessage,
	}
)

func NewDiscordClient(cfg config.Chat) (*DiscordClient, error) {
	c := &DiscordClient{
		channelID:        cfg.ChannelID,
		guildID:          cfg.GuildID,
		messageFormatter: formattingFuncs[cfg.DisplayMode],
	}

	discord, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.Key))
	if err != nil {
		return nil, fmt.Errorf("unable to create Discord session: %w", err)
	}
	discord.ShouldRetryOnRateLimit = true
	discord.ShouldReconnectOnError = true

	discord.Identify.Intents = discordgo.IntentGuildMessages
	discord.AddHandler(c.HandleOnReady)

	c.discord = discord
	return c, nil
}

func (d *DiscordClient) Connect() error {
	return d.discord.Open()
}

func (d *DiscordClient) Disconnect() error {
	return d.discord.Close()
}

func (d *DiscordClient) HandleOnReady(s *discordgo.Session, m *discordgo.Ready) {
	_, err := d.discord.ChannelMessageSend(d.channelID, "Engaging Maximum Derek!")
	if err != nil {
		fmt.Println(fmt.Errorf("unable to send message: %w", err))
	}
}

func (d *DiscordClient) SendMessage(msg queue.BroadcastMessage) error {
	selfFind := msg.Sender == msg.Receiver

	// TODO how do i do @?
	_, err := d.discord.ChannelMessageSend(d.channelID, d.messageFormatter(msg, selfFind))
	if err != nil {
		return fmt.Errorf("unable to send message: %w", err)
	}

	return nil
}
