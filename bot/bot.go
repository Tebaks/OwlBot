package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"../config"
	"github.com/bwmarrin/discordgo"
)

var BotID string
var goBot *discordgo.Session
var userChannel [][]string
var tempChannels []string

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

func deleteChannel(s *discordgo.Session, channelID string) {

	for x, y := range tempChannels {
		if y == channelID {
			tempChannels = append(tempChannels[:x], tempChannels[x+1:]...)
			s.ChannelDelete(y)
		}
	}
}

func deleteUser(array [][]string, data string) [][]string {

	for x, y := range array {
		if y[0] == data {
			result := append(array[:x], array[x+1:]...)
			return result
		}
	}

	return nil

}

func setUserChannel(array [][]string, userID string, channelID string) [][]string {

	for _, y := range array {
		if y[0] == userID {
			y[1] = channelID
			return array
		}
	}

	array = append(array, []string{userID, channelID})

	return array
}

func deleteEmptyChannels(s *discordgo.Session) {
	for _, y := range tempChannels {
		if checkChannelEmpty(y) {
			deleteChannel(s, y)
		}
	}

}

func checkChannelEmpty(tempChannel string) bool {
	for _, y := range userChannel {
		if y[1] == tempChannel {
			return false
		}
	}

	return true
}

func channelHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	if v.ChannelID == "743575002141687808" {
		usr, _ := s.User(v.UserID)

		channelData := discordgo.GuildChannelCreateData{Name: "Room " + usr.Username, Type: discordgo.ChannelTypeGuildVoice}
		_, _ = s.ChannelMessageSend("743455466213867540", "Kanal oluşturuluyor")
		chn, _ := s.GuildChannelCreateComplex(v.GuildID, channelData)

		tempChannels = append(tempChannels, chn.ID)

		s.GuildMemberMove(v.GuildID, v.UserID, &chn.ID)

	} else {
		if v.ChannelID == "" {
			userChannel = deleteUser(userChannel, v.UserID)
			fmt.Println(userChannel)
		} else {
			userChannel = setUserChannel(userChannel, v.UserID, v.ChannelID)
			fmt.Println(userChannel)
		}
		deleteEmptyChannels(s)
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
		if data[0] == "!delete" {
			if len(data) == 2 {
				var digitCheck = regexp.MustCompile(`^[0-9]+$`)

				if digitCheck.MatchString(data[1]) {

					count, _ := strconv.Atoi(data[1])

					msgs, _ := s.ChannelMessages(m.ChannelID, count, m.ID, "", "")

					msgsID := []string{}

					for _, x := range msgs {
						msgsID = append(msgsID, x.ID)
					}

					_ = s.ChannelMessagesBulkDelete(m.ChannelID, msgsID)

				} else {
					_, _ = s.ChannelMessageSend(m.ChannelID, "Yanlış format. (!delete count)")
				}

			} else {
				_, _ = s.ChannelMessageSend(m.ChannelID, "Yanlış format. (!delete count)")
			}

		}

	}

}
