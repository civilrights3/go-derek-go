package chat

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/civilrights3/go-derek-go/internal/config"
	"github.com/civilrights3/go-derek-go/internal/queue"
)

type DiscordClient struct {
	discord   *discordgo.Session
	channelID string
	guildID   string
}

func NewDiscordClient(cfg config.Chat) (*DiscordClient, error) {
	c := &DiscordClient{
		channelID: cfg.ChannelID,
		guildID:   cfg.GuildID,
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
	// TODO how do i do @?
	_, err := d.discord.ChannelMessageSend(d.channelID, formatMessage(msg))
	if err != nil {
		return fmt.Errorf("unable to send message: %w", err)
	}

	return nil
}

func formatMessage(msg queue.BroadcastMessage) string {
	selfFind := msg.Sender == msg.Receiver

	if selfFind {
		return fmt.Sprintf("[%s] found their <%s> (%s)", msg.Receiver, msg.Item, msg.Location)
	}

	return fmt.Sprintf("[%s] sent <%s> to {%s} (%s)", msg.Sender, msg.Item, msg.Receiver, msg.Location)
}
