package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/greymd/ojichat/generator"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	<-ctx.Done()

	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := fmt.Sprintf("<@!%s>", s.State.User.ID)

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	msg := strings.TrimSpace(strings.TrimPrefix(m.Content, prefix))

	if msg == "" {
		msg = m.Author.Username
	}

	reply, err := generator.Start(generator.Config{
		EmojiNum:         4,
		TargetName:       msg,
		PunctuationLevel: 0,
	})

	if err != nil {
		reply = err.Error()
	}

	s.ChannelMessageSend(m.ChannelID, reply)
}
