package main

import (
    "github.com/Jeffail/gabs"
    "github.com/ugjka/cleverbot-go"
    "github.com/bwmarrin/discordgo"
    "fmt"
)

var (
    cleverbotSessions map[string]*cleverbot.Session
)

func GetConfig(path string) *gabs.Container {
    json, err := gabs.ParseJSONFile(path)

    if err != nil {
        panic(err)
    }

    return json
}

func CleverbotSend(channel string, message string) {
    var msg string

    if cleverbotSessions[channel] == nil {
        if len(cleverbotSessions) == 0 {
            cleverbotSessions = make(map[string]*cleverbot.Session)
        }

        cleverbotSessions[channel] = cleverbot.New()
    }

    response, err := cleverbotSessions[channel].Ask(message)
    if err != nil {
        msg = "Error :frowning:\n```\n" + err.Error() + "\n```"
    } else {
        msg = response
    }

    discordSession.ChannelMessageSend(channel, msg)
}

func CleverbotRefreshSession(channel string) {
    cleverbotSessions[channel] = cleverbot.New()
}

func CCTV(message *discordgo.Message) {
    var (
        channelName string = "?"
        channelID string = "?"
        serverName string = "?"
        serverID string = "?"
    )

    channel, err := discordSession.Channel(message.ChannelID)
    if err == nil {
        channelName = channel.Name
        channelID = channel.ID

        server, err := discordSession.Guild(channel.ID)
        if err == nil {
            serverName = server.Name
            serverID = server.ID
        }
    }

    template := `
At:      %s
Origin:  #%s (%s) in %s (%s)
Author:  %s#%s (%s)
Message:
%s
`

    msg :=
        "```\n" +
            fmt.Sprintf(
                template,
                message.Timestamp,
                channelName,
                channelID,
                serverName,
                serverID,
                message.Author.Username,
                message.Author.Discriminator,
                message.Author.ID,
                message.Content,
            ) +
            "\n```"

    discordSession.ChannelMessageSend(
        config.Path("cctv").Data().(string),
        msg,
    )
}

func GetPrefixForServer(guild string) (string, error) {
    return GuildSettingGet(guild, "prefix")
}

func SetPrefixForServer(guild string, prefix string) error {
    return GuildSettingSet(guild, "prefix", prefix)
}

func SendError(channel string, err error) {
    discordSession.ChannelMessageSend(
        channel,
        "Error :frowning:\n```\n" + err.Error() + "\n```",
    )
}

type Callback func()

func WhileTypingIn(channel string, cb Callback) {
    discordSession.ChannelTyping(channel)
    cb()
}