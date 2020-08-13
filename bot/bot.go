package bot

import (
	"fmt"
	"strings"

	"../config"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error)
		return
	}

	u, err := goBot.User("@me")

	goBot.AddHandler(channelHandler)

	goBot.AddHandler(messageHandler)

	if err != nil {
		fmt.Println(err.Error)
		return
	}
	BotID = u.ID

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error)
		return
	}

	fmt.Println("Bot is running")
}
func channelHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	if v.ChannelID == "743575002141687808" {
		channelData := discordgo.GuildChannelCreateData{Name: "test", Type: discordgo.ChannelTypeGuildVoice}
		_, _ = s.ChannelMessageSend("743455466213867540", "Kanal oluşturuluyor")
		chn, _ := s.GuildChannelCreateComplex(v.GuildID, channelData)

		s.GuildMemberMove(v.GuildID, v.UserID, &chn.ID)

	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {

		data := strings.Split(m.Content, " ")

		if m.Author.ID == BotID {
			return
		}
		if data[0] == "!ping" {
			_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
		}

	}

}
