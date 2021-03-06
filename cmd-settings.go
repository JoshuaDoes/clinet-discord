package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"4d63.com/tz"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/structs"
)

var (
	// regexpSwitchFC contains a regular expression for matching valid Nintendo Switch friend codes
	regexpSwitchFC = regexp.MustCompile(`SW-[0-9]{4}-[0-9]{4}-[0-9]{4}`)
)

var (
	// LogEventsRecommended contains pre-enabled recommended logging events
	LogEventsRecommended = LogEvents{
		ChannelCreate:     true,
		ChannelDelete:     true,
		GuildBanAdd:       true,
		GuildBanRemove:    true,
		GuildMemberAdd:    true,
		GuildMemberRemove: true,
		GuildRoleCreate:   true,
		GuildRoleDelete:   true,
		GuildRoleUpdate:   true,
		GuildUpdate:       true,
		SwearDetect:       true,
		UserModlog:        true,
		VoiceStateUpdate:  true,
	}
)

// GuildSettings holds settings specific to a guild
type GuildSettings struct { //By default this will only be configurable for users in a role with the server admin permission
	AllowVoice              bool                  `json:"allowVoice,omitempty"`              //Whether voice commands should be usable in this guild
	BotAdminRoles           []string              `json:"adminRoles,omitempty"`              //An array of role IDs that can admin the bot without the guild administrator permission
	BotAdminUsers           []string              `json:"adminUsers,omitempty"`              //An array of user IDs that can admin the bot without a guild administrator role
	BotOptions              BotOptions            `json:"botOptions,omitempty"`              //The bot options to use in this guild (true gets overridden if global bot config is false)
	BotPrefix               string                `json:"botPrefix",omitempty`               //The bot prefix to use in this guild
	CustomResponses         []CustomResponseQuery `json:"customResponses,omitempty"`         //An array of custom responses specific to the guild
	LogSettings             LogSettings           `json:"logSettings,omitempty"`             //Logging settings
	SwearFilter             SwearFilter           `json:"swearFilter,omitempty"`             //The swear filter settings specific to this guild
	TipsChannel             string                `json:"tipsChannel,omitempty"`             //The channel to post tip messages to
	UserJoinMessage         string                `json:"userJoinMessage,omitempty"`         //A message to send when a user joins
	UserJoinMessageChannel  string                `json:"userJoinMessageChannel,omitempty"`  //The channel to send the user join message to
	UserLeaveMessage        string                `json:"userLeaveMessage,omitempty"`        //A message to send when a user leaves
	UserLeaveMessageChannel string                `json:"userLeaveMessageChannel,omitempty"` //The channel to send the user leave message to
	RoleMeList              []*RoleMe             `json:"roleMeList,omitempty"`              //An array of rolemes specific to this guild
	AutoSendNowPlaying      bool                  `json:"disableNowPlaying,omitempty"`       //Whether or not the Now Playing embed should be sent each time a new track is automatically started without user interaction
	APIInviteChannel        string                `json:"apiInviteChannel,omitempty"`        //The channel to use for server-side invite link generation
	APIInviteKey            string                `json:"apiInviteKey,omitempty"`            //The key to use for server-side invite link generation
	Feeds                   []*Feed               `json:"feeds,omitempty"`                   //A list of feeds for the current guild
}

// UserSettings holds settings specific to a user
type UserSettings struct {
	Balance   int       `json:"balance,omitempty"`   //A balance to use as virtual currency for up-and-coming ideas
	DailyNext time.Time `json:"dailyNext,omitempty"` //The next time the user is able to use the daily credits command

	//Basic info
	AboutMe  string `json:"description,omitempty"` //An aboutme set by the user
	Timezone string `json:"timezone,omitempty"`    //A timezone set by the user to use in other functions

	//Socials
	Socials Socials `json:"socials,omitempty"` //Social media, gamertags, etc
}

// Socials holds socials information
type Socials struct {
	SwitchFC string `json:"switchFC,omitempty"` //Nintendo Switch friend code
	NNID     string `json:"nintyID,omitempty"`  //Nintendo Network ID
	PSN      string `json:"psn,omitempty"`      //PlayStation Network
	Xbox     string `json:"xbox,omitempty"`     //Xbox Live
}

// LogSettings holds settings specific to logging
type LogSettings struct {
	LoggingEnabled bool      `json:"loggingEnabled"` //Whether or not logging enabled
	LoggingChannel string    `json:"loggingChannel"` //The channel to log guild events to
	LoggingEvents  LogEvents `json:"loggingEvents"`  //The events to log
}

// LogEvents holds logging events and whether or not they're enabled
type LogEvents struct {
	//Events received from Discord
	ChannelCreate     bool `json:"channelCreate"`
	ChannelDelete     bool `json:"channelDelete"`
	ChannelUpdate     bool `json:"channelUpdate"`
	GuildBanAdd       bool `json:"guildBanAdd"`
	GuildBanRemove    bool `json:"guildBanRemove"`
	GuildEmojisUpdate bool `json:"guildEmojisUpdate"`
	GuildMemberAdd    bool `json:"guildMemberAdd"`
	GuildMemberRemove bool `json:"guildMemberRemove"`
	GuildRoleCreate   bool `json:"guildRoleCreate"`
	GuildRoleDelete   bool `json:"guildRoleDelete"`
	GuildRoleUpdate   bool `json:"guildRoleUpdate"`
	GuildUpdate       bool `json:"guildUpdate"`
	UserUpdate        bool `json:"userUpdate"`
	VoiceStateUpdate  bool `json:"voiceStateUpdate"`

	//Custom events
	SwearDetect bool `json:"swearDetect"` //Triggered if a user uses a blacklisted (swear) word
	UserModlog  bool `json:"userModlog"`  //Triggered if a user's modlog is updated globally
}

func commandSettingsBot(args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
	switch args[0] {
	case "prefix":
		if len(args) > 1 {
			if args[1] == botData.CommandPrefix {
				guildSettings[env.Guild.ID].BotPrefix = ""
			} else {
				guildSettings[env.Guild.ID].BotPrefix = args[1]
			}
			return NewGenericEmbed("Bot Settings - Command Prefix", "Successfully set the command prefix to ``"+strings.Replace(args[1], "`", "\\`", -1)+"``.")
		}
		if guildSettings[env.Guild.ID].BotPrefix != "" {
			return NewGenericEmbed("Bot Settings - Command Prefix", "Current command prefix:\n\n"+guildSettings[env.Guild.ID].BotPrefix)
		}
		return NewGenericEmbed("Bot Settings - Command Prefix", "Current command prefix:\n\n"+botData.CommandPrefix)
	}
	return NewErrorEmbed("Bot Settings Error", "Error finding the setting ``"+args[0]+"``.")
}

func commandSettingsUser(args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
	//We're getting there (⟃ ͜ʖ ⟄)

	switch args[0] {
	case "about", "aboutme", "description", "desc", "info":
		if len(args) <= 1 {
			if userSettings[env.User.ID].AboutMe == "" {
				return NewErrorEmbed("User Settings - About Me Error", "You must specify an aboutme to view it.")
			}
			return aboutMe(env.User.ID)
		}
		if len(args) == 2 && len(env.Message.Mentions) > 0 {
			return aboutMe(env.Message.Mentions[0].ID)
		}
		userSettings[env.User.ID].AboutMe = strings.Join(args[1:], " ")
		return NewGenericEmbed("User Settings - About Me", "Successfully set your about me!")
	case "timezone", "tz":
		if len(args) <= 1 {
			if userSettings[env.User.ID].Timezone == "" {
				return NewErrorEmbed("User Settings - Timezone Error", "You must specify a timezone to view it.")
			}
			location, err := tz.LoadLocation(userSettings[env.User.ID].Timezone)
			if err != nil {
				return NewErrorEmbed("User Settings - Timezone Error", "You have an invalid timezone set, please set a new one first.\n\nEx: ``"+env.BotPrefix+"user timezone America/New_York``")
			}
			return NewGenericEmbed("User Settings - Timezone", "Your current timezone is set to ``"+userSettings[env.User.ID].Timezone+"``.\nYour current time is ``"+time.Now().In(location).String()+"``.")
		}
		location, err := tz.LoadLocation(args[1])
		if err != nil {
			return NewErrorEmbed("User Settings - Timezone Error", "Invalid timezone.")
		}
		userSettings[env.User.ID].Timezone = args[1]
		return NewGenericEmbed("User Settings - Timezone", "Successfully set your timezone to ``"+args[1]+"``.\nYour current time is ``"+time.Now().In(location).String()+"``.")
	case "social", "socials":
		/*
		* cli$user social add switchfc SW-0000-0000-0000
		* cli$user social list
		* cli$user social clear switchfc
		* cli$user social available
		 */

		socialCommand := &Command{
			HelpText: "Manages your socials.",
			RequiredArguments: []string{
				"setting (value(s))",
			},
			Arguments: []CommandArgument{
				{Name: "set {social}", Description: "Sets a social", ArgType: "social code/name"},
				{Name: "list", Description: "Lists your socials", ArgType: "this"},
				{Name: "remove", Description: "Removes a social", ArgType: "social code/name"},
				{Name: "available", Description: "Lists available socials", ArgType: "this"},
			},
		}
		cmdUsage := getCustomCommandUsage(socialCommand, "user "+args[0], "User Settings - Socials Help", env)

		if len(args) < 2 {
			return cmdUsage
		}

		switch args[1] {
		case "set", "add":
			if len(args) < 4 {
				return NewErrorEmbed("User Settings - Socials", "You must specify a social identifier to set it.")
			}
			switch args[2] {
			case "switchfc":
				if !regexpSwitchFC.MatchString(args[3]) {
					return NewErrorEmbed("User Settings - Socials", "Invalid Switch friend code.")
				}
				if userSettings[env.User.ID].Socials.SwitchFC == args[3] {
					return NewErrorEmbed("User Settings - Socials", "You have already set that Switch friend code.")
				}
				userSettings[env.User.ID].Socials.SwitchFC = args[3]
				return NewGenericEmbed("User Settings - Socials", "Successfully set your Switch friend code to ``"+args[3]+"``.")
			case "nintendoid", "nintyid", "nnid":
				if userSettings[env.User.ID].Socials.NNID == args[3] {
					return NewErrorEmbed("User Settings - Socials", "You have already set that NNID.")
				}
				exists, _, err := botData.BotClients.Ninty.DoesUserExist(args[3])
				if err != nil {
					return NewErrorEmbed("User Settings - Social Error", "There was an error checking if that NNID exists.")
				}
				if !exists {
					return NewErrorEmbed("User Settings - Social Error", "That NNID doesn't exist!")
				}
				userSettings[env.User.ID].Socials.NNID = args[3]
				return NewGenericEmbed("User Settings - Socials", "Successfully set your NNID to ``"+args[3]+"``.")
			case "psn":
				if userSettings[env.User.ID].Socials.PSN == args[3] {
					return NewErrorEmbed("User Settings - Socials", "You have already set that PSN.")
				}
				userSettings[env.User.ID].Socials.PSN = args[3]
				return NewGenericEmbed("User Settings - Socials", "Successfully set your PSN to ``"+args[3]+"``.")
			case "xbox", "gamertag":
				if userSettings[env.User.ID].Socials.Xbox == args[3] {
					return NewErrorEmbed("User Settings - Socials", "You have already set that Xbox Live gamertag.")
				}
				userSettings[env.User.ID].Socials.Xbox = args[3]
				return NewGenericEmbed("User Settings - Socials", "Successfully set your Xbox Live gamertag to ``"+args[3]+"``.")
			}
			return NewErrorEmbed("User Settings - Socials Error", "Unknown social "+args[2]+"``.")
		case "list":
			socialsEmbed := NewEmbed().
				SetTitle("Socials").
				SetDescription("Below are all of the socials you have added.").MessageEmbed

			socials := userSettings[env.User.ID].Socials
			socialsFields := make([]*discordgo.MessageEmbedField, 0)

			if socials.SwitchFC != "" {
				socialsFields = append(socialsFields, &discordgo.MessageEmbedField{
					Name:  "Switch Friend Code",
					Value: socials.SwitchFC,
				})
			}
			if socials.NNID != "" {
				socialsFields = append(socialsFields, &discordgo.MessageEmbedField{
					Name:  "Nintendo Network ID",
					Value: socials.NNID,
				})
			}
			if socials.PSN != "" {
				socialsFields = append(socialsFields, &discordgo.MessageEmbedField{
					Name:  "PSN",
					Value: socials.PSN,
				})
			}
			if socials.Xbox != "" {
				socialsFields = append(socialsFields, &discordgo.MessageEmbedField{
					Name:  "Xbox Live Gamertag",
					Value: socials.Xbox,
				})
			}

			if len(socialsFields) == 0 {
				return NewGenericEmbed("User Settings - Socials", "You don't have any socials yet!")
			}

			socialsEmbed.Fields = socialsFields

			return socialsEmbed
		case "clear", "remove":
			if len(args) < 3 {
				return cmdUsage
			}
			switch args[2] {
			case "switchfc":
				if userSettings[env.User.ID].Socials.SwitchFC == "" {
					return NewErrorEmbed("User Settings - Socials", "You don't have a Switch friend code set.")
				}
				userSettings[env.User.ID].Socials.SwitchFC = ""
				return NewGenericEmbed("User Settings - Socials", "Cleared your Switch friend code.")
			case "nintendoid", "nintyid", "nnid":
				if userSettings[env.User.ID].Socials.NNID == "" {
					return NewErrorEmbed("User Settings - Socials", "You don't have an NNID set.")
				}
				userSettings[env.User.ID].Socials.NNID = ""
				return NewGenericEmbed("User Settings - Socials", "Cleared your NNID.")
			case "psn":
				if userSettings[env.User.ID].Socials.PSN == "" {
					return NewErrorEmbed("User Settings - Socials", "You don't have a PSN set.")
				}
				userSettings[env.User.ID].Socials.PSN = ""
				return NewGenericEmbed("User Settings - Socials", "Cleared your PSN.")
			case "xbox":
				if userSettings[env.User.ID].Socials.Xbox == "" {
					return NewErrorEmbed("User Settings - Socials", "You don't have an Xbox Live gamertag set.")
				}
				userSettings[env.User.ID].Socials.Xbox = ""
				return NewGenericEmbed("User Settings - Socials", "Cleared your Xbox Live gamertag.")
			}
			return NewErrorEmbed("User Settings - Socials Error", "Unknown social ``"+args[2]+"``.")
		case "available", "types":
			return NewGenericEmbed("User Settings - Socials - Types", "These are the available socials you can use:\n\n"+
				"``switchfc`` - Nintendo Switch friend code\n"+
				"``nnid`` - Nintendo Network ID\n"+
				"``psn`` - PlayStation Network\n"+
				"``xbox`` - Xbox Live Gamertag",
			)
		}
		return NewErrorEmbed("User Settings - Socials Error", "Unknown socials command ``"+args[1]+"``.")
	}
	return NewErrorEmbed("User Settings Error", "Error finding the setting ``"+args[0]+"``.")
}

func aboutMe(userID string) *discordgo.MessageEmbed {
	settings, found := userSettings[userID]
	if !found {
		return NewErrorEmbed("About Me - Error", "Error finding the aboutme for <@!"+userID+">.")
	}

	user, err := botData.DiscordSession.User(userID)
	if err != nil {
		return NewErrorEmbed("About Me - Error", "Error finding the user <@!"+userID+">.")
	}

	return NewEmbed().
		SetAuthor(user.Username+"#"+user.Discriminator, user.AvatarURL("2048")).
		AddField("About Me", settings.AboutMe).
		SetColor(0x1C1C1C).MessageEmbed
}

func commandSettingsServer(args []string, env *CommandEnvironment) *discordgo.MessageEmbed {
	switch args[0] {
	case "joinmsg":
		guildSettings[env.Guild.ID].UserJoinMessage = strings.Join(args[1:], " ")
		guildSettings[env.Guild.ID].UserJoinMessageChannel = env.Channel.ID
		return NewGenericEmbed("Server Settings - Join Message", "Successfully set the join message to this channel.")
	case "leavemsg":
		guildSettings[env.Guild.ID].UserLeaveMessage = strings.Join(args[1:], " ")
		guildSettings[env.Guild.ID].UserLeaveMessageChannel = env.Channel.ID
		return NewGenericEmbed("Server Settings - Leave Message", "Successfully set the leave message to this channel.")
	case "tips":
		if len(args) <= 1 {
			if guildSettings[env.Guild.ID].TipsChannel != "" {
				return NewGenericEmbed("Server Settings - Tips", "Tips are enabled for this server.")
			}
			return NewGenericEmbed("Server Settings - Tips", "Tips are disabled for this server.")
		}
		switch args[1] {
		case "enable":
			guildSettings[env.Guild.ID].TipsChannel = env.Channel.ID
			return NewGenericEmbed("Server Settings - Tips", "Successfully enabled hourly tips for this channel.")
		case "disable":
			guildSettings[env.Guild.ID].TipsChannel = ""
			return NewGenericEmbed("Server Settings - Tips", "Successfully disabled hourly tips for this channel.")
		}
		return NewErrorEmbed("Server Settings - Tips Error", "Unknown tips command ``"+args[1]+"``.")
	case "autosendnowplaying":
		switch args[1] {
		case "enable":
			guildSettings[env.Guild.ID].AutoSendNowPlaying = true
			return NewGenericEmbed("Server Settings - Auto Send Now Playing", "Successfully enabled sending now playing messages each time a new track is started without user interaction.")
		case "disable":
			guildSettings[env.Guild.ID].AutoSendNowPlaying = false
			return NewGenericEmbed("Server Settings - Auto Send Now Playing", "Successfully disabled sending now playing messages each time a new track is started without user interaction.")
		}
		return NewErrorEmbed("Server Settings - Auto Send Now Playing Error", "Unknown ASNP command ``"+args[1]+"``.")
	case "invitegen":
		if len(args) < 2 {
			invitegenHelpCmd := &Command{
				HelpText: "Manages invite link generation via the API.",
				RequiredArguments: []string{
					"setting (value(s))",
				},
				Arguments: []CommandArgument{
					{Name: "setchannel", Description: "Sets the invite link channel to the current channel", ArgType: "this"},
					{Name: "key", Description: "Displays or sets the key to use for invite link generation", ArgType: "this/string"},
				},
			}
			return getCustomCommandUsage(invitegenHelpCmd, "server invitegen", "Server Settings - API Invite Generation Help", env)
		}

		switch args[1] {
		case "setchannel":
			guildSettings[env.Guild.ID].APIInviteChannel = env.Channel.ID
			return NewGenericEmbed("Server Settings - API Invite Generation", "Successfully set the channel to use for generating invite links to this channel.")
		case "key":
			if len(args) > 2 {
				guildSettings[env.Guild.ID].APIInviteKey = strings.Join(args[2:], " ")
				return NewGenericEmbed("Server Settings - API Invite Generation", "Successfully set the key to use for generating invite links to ``"+guildSettings[env.Guild.ID].APIInviteKey+"``.")
			}
			if guildSettings[env.Guild.ID].APIInviteKey == "" {
				return NewGenericEmbed("Server Settings - API Invite Generation", "No key is currently set for generating invite links!")
			}
			return NewGenericEmbed("Server Settings - API Invite Generation", "The current key for generating invite links is ``"+guildSettings[env.Guild.ID].APIInviteKey+"``.")
		}
		return NewErrorEmbed("Server Settings - API Invite Generation Error", "Unknown invitegen command ``"+args[1]+"``.")
	case "filter":
		if len(args) < 2 {
			filterHelpCmd := &Command{
				HelpText: "Manages the swear filter for this server.",
				RequiredArguments: []string{
					"setting (value(s))",
				},
				Arguments: []CommandArgument{
					{Name: "enable", Description: "Enables the swear filter for this server", ArgType: "this"},
					{Name: "disable", Description: "Disables the swear filter for this server", ArgType: "this"},
					{Name: "timeout", Description: "Displays or sets the timeout for deleting warning messages", ArgType: "this/number"},
					{Name: "words", Description: "Lists filtered words, or adds/removes specified words/clears all words", ArgType: "this/(add word1)/(remove word2)/clear"},
				},
			}
			return getCustomCommandUsage(filterHelpCmd, "server filter", "Server Settings - Swear Filter Help", env)
		}

		switch args[1] {
		case "enable":
			guildSettings[env.Guild.ID].SwearFilter.Enabled = true
			return NewGenericEmbed("Server Settings - Swear Filter", "Successfully enabled the swear filter.")
		case "disable":
			guildSettings[env.Guild.ID].SwearFilter.Enabled = false
			return NewGenericEmbed("Server Settings - Swear Filter", "Successfully disabled the swear filter.")
		case "words":
			if len(args) < 3 {
				words := "No words are in the swear filter!"
				if len(guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords) > 0 {
					words = strings.Join(guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords, ", ")
				}
				wordListEmbed := NewEmbed().
					SetTitle("Server Settings - Swear Filter").
					AddField("Filtered Words", words).
					SetColor(0x1C1C1C).MessageEmbed
				return wordListEmbed
			}
			switch args[2] {
			case "add":
				if len(args) < 4 {
					return NewErrorEmbed("Server Settings - Swear Filter Error", "You must specify one or more words to add to the filter.")
				}
				guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords = append(guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords, args[3:]...)
				return NewGenericEmbed("Server Settings - Swear Filter", "Successfully added the provided words to the filter.")
			case "remove":
				if len(args) < 4 {
					return NewErrorEmbed("Server Settings - Swear Filter Error", "You must specify one or more words to remove from the filter.")
				}
				for _, word := range guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords {
					guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords = remove(guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords, word)
				}
				return NewGenericEmbed("Server Settings - Swear Filter", "Successfully removed the provided words from the filter.")
			case "clear":
				guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords = make([]string, 0)
				return NewGenericEmbed("Server Settings - Swear Filter", "Successfully cleared all words from the filter.")
			}
		case "timeout":
			if len(args) < 3 {
				if guildSettings[env.Guild.ID].SwearFilter.WarningDeleteTimeout == 0 {
					return NewGenericEmbed("Server Settings - Swear Filter", "The timeout for deleting warning messages is disabled.")
				}
				timeout := strconv.Itoa(int(guildSettings[env.Guild.ID].SwearFilter.WarningDeleteTimeout))
				return NewGenericEmbed("Server Settings - Swear Filter", "The current timeout for deleting warning messages is set to "+timeout+" seconds.")
			}
			timeout, err := strconv.Atoi(args[2])
			if err != nil {
				return NewErrorEmbed("Server Settings - Swear Filter Error", "``"+args[2]+"`` is not a valid number.")
			}
			guildSettings[env.Guild.ID].SwearFilter.WarningDeleteTimeout = time.Duration(timeout)
			return NewGenericEmbed("Server Settings - Swear Filter", "Successfully set he timeout for deleting warning messages to "+args[2]+" seconds.")
		}
		return NewErrorEmbed("Server Settings - Swear Filter Error", "Unknown filter command ``"+args[1]+"``.")
	case "log":
		if len(args) < 2 {
			logHelpCmd := &Command{
				HelpText: "Sets the logging capabilities for this server.",
				RequiredArguments: []string{
					"setting (value(s))",
				},
				Arguments: []CommandArgument{
					{Name: "set", Description: "Sets the logging channel to the current channel", ArgType: "this"},
					{Name: "enable", Description: "Enables logging for the server (to this channel if not set), enabling any optionally specified events", ArgType: "this/event(s)"},
					{Name: "disable", Description: "Disables logging for the server, disabling any optionally specified events", ArgType: "this/event(s)"},
					{Name: "unset", Description: "Unsets the current logging channel and disables logging", ArgType: "this"},
					{Name: "events", Description: "Returns a list of available events to enable/disable", ArgType: "this"},
				},
			}
			return getCustomCommandUsage(logHelpCmd, "server log", "Server Settings - Log Help", env)
		}

		LoggingEventsTmp := &guildSettings[env.Guild.ID].LogSettings.LoggingEvents

		switch args[1] {
		case "set":
			guildSettings[env.Guild.ID].LogSettings.LoggingChannel = env.Channel.ID
			return NewGenericEmbed("Server Settings - Log", "Successfully set the logging channel to this channel.")
		case "enable":
			guildSettings[env.Guild.ID].LogSettings.LoggingEnabled = true

			if len(args) == 3 {
				switch args[2] {
				case "all":
					events := structs.New(LoggingEventsTmp)
					fields := events.Fields()

					for _, event := range fields {
						err := event.Set(true)
						if err != nil {
							return NewErrorEmbed("Server Settings - Log", "Unable to enable all logging events.")
						}
					}

					guildSettings[env.Guild.ID].LogSettings.LoggingEvents = *LoggingEventsTmp

					if guildSettings[env.Guild.ID].LogSettings.LoggingChannel == "" {
						guildSettings[env.Guild.ID].LogSettings.LoggingChannel = env.Channel.ID
						return NewGenericEmbed("Server Settings - Log", "Successfully enabled all logging events and set the logging channel to this channel.")
					}

					return NewGenericEmbed("Server Settings - Log", "Successfully enabled all logging events.")
				case "recommended":
					guildSettings[env.Guild.ID].LogSettings.LoggingEvents = LogEventsRecommended

					if guildSettings[env.Guild.ID].LogSettings.LoggingChannel == "" {
						guildSettings[env.Guild.ID].LogSettings.LoggingChannel = env.Channel.ID
						return NewGenericEmbed("Server Settings - Log", "Successfully toggled all logging events to their recommended states and set the logging channel to this channel.")
					}

					return NewGenericEmbed("Server Settings - Log", "Successfully toggled all logging events to their recommended states.")
				}
			}

			eventsToEnable := make([]string, 0)
			if len(args) > 2 {
				eventsToEnable = args[2:]
			}
			eventsEnabled := make([]string, 0)
			eventsFailed := make([]string, 0)

			if len(eventsToEnable) > 0 {
				events := structs.New(LoggingEventsTmp)

				for _, eventName := range eventsToEnable {
					event, ok := events.FieldOk(eventName)
					if ok {
						event.Set(true)
						eventsEnabled = append(eventsEnabled, eventName)
					} else {
						eventsFailed = append(eventsFailed, eventName)
					}
				}
			}

			guildSettings[env.Guild.ID].LogSettings.LoggingEvents = *LoggingEventsTmp

			responseMessage := "Successfully enabled logging"
			if guildSettings[env.Guild.ID].LogSettings.LoggingChannel != "" {
				responseMessage += "."
			} else {
				responseMessage += " and set the logging channel to this channel."
			}
			if len(eventsToEnable) > 0 {
				responseMessage += "\n"
				if len(eventsEnabled) > 0 {
					responseMessage += "\nEnabled the following events: " + strings.Join(eventsEnabled, ", ")
				}
				if len(eventsFailed) > 0 {
					responseMessage += "\nFailed to find the following events: " + strings.Join(eventsFailed, ", ")
				}
			}
			return NewGenericEmbed("Server Settings - Log", responseMessage)
		case "disable":
			if len(args) == 3 && args[2] == "all" {
				guildSettings[env.Guild.ID].LogSettings.LoggingEvents = LogEvents{}
				return NewGenericEmbed("Server Settings - Log", "Successfully disabled all logging events.")
			}

			eventsToDisable := make([]string, 0)
			if len(args) > 2 {
				eventsToDisable = args[2:]
			}
			eventsDisabled := make([]string, 0)
			eventsFailed := make([]string, 0)

			if len(eventsToDisable) > 0 {
				events := structs.New(LoggingEventsTmp)

				for _, eventName := range eventsToDisable {
					event, ok := events.FieldOk(eventName)
					if ok {
						event.Set(false)
						eventsDisabled = append(eventsDisabled, eventName)
					} else {
						eventsFailed = append(eventsFailed, eventName)
					}
				}
			} else {
				guildSettings[env.Guild.ID].LogSettings.LoggingEnabled = false
				return NewGenericEmbed("Server Settings - Log", "Successfully disabled logging.")
			}

			guildSettings[env.Guild.ID].LogSettings.LoggingEvents = *LoggingEventsTmp

			responseMessage := ""
			if len(eventsToDisable) > 0 {
				if len(eventsDisabled) > 0 {
					responseMessage += "\nDisabled the following events: " + strings.Join(eventsDisabled, ", ")
				}
				if len(eventsFailed) > 0 {
					responseMessage += "\nFailed to find the following events: " + strings.Join(eventsFailed, ", ")
				}
			}
			return NewGenericEmbed("Server Settings - Log", responseMessage)
		case "events":
			responseMessage := "__Event states__\n"

			events := structs.New(guildSettings[env.Guild.ID].LogSettings.LoggingEvents)
			eventFields := events.Fields()

			for _, event := range eventFields {
				responseMessage += "\n" + event.Name() + ": **" + strconv.FormatBool(event.Value().(bool)) + "**"
			}

			return NewGenericEmbed("Server Settings - Log", responseMessage)
		}
		return NewErrorEmbed("Server Settings - Log Error", "Unknown log command ``"+args[1]+"``.")
	case "reset":
		if len(args) < 2 {
			return NewErrorEmbed("Server Settings - Reset Error", "You must specify a setting to reset.")
		}
		switch args[1] {
		case "joinmsg":
			guildSettings[env.Guild.ID].UserJoinMessage = ""
			guildSettings[env.Guild.ID].UserJoinMessageChannel = ""
		case "leavemsg":
			guildSettings[env.Guild.ID].UserLeaveMessage = ""
			guildSettings[env.Guild.ID].UserLeaveMessageChannel = ""
		case "log":
			guildSettings[env.Guild.ID].LogSettings.LoggingChannel = ""
			guildSettings[env.Guild.ID].LogSettings.LoggingEnabled = false
			guildSettings[env.Guild.ID].LogSettings.LoggingEvents = LogEvents{}
		case "filter":
			guildSettings[env.Guild.ID].SwearFilter.Enabled = false
			guildSettings[env.Guild.ID].SwearFilter.BlacklistedWords = make([]string, 0)
			guildSettings[env.Guild.ID].SwearFilter.DisableNormalize = false
			guildSettings[env.Guild.ID].SwearFilter.DisableSpacedTab = false
			guildSettings[env.Guild.ID].SwearFilter.DisableMultiWhitespaceStripping = false
			guildSettings[env.Guild.ID].SwearFilter.DisableZeroWidthStripping = false
			guildSettings[env.Guild.ID].SwearFilter.DisableSpacedBypass = false
			guildSettings[env.Guild.ID].SwearFilter.WarningDeleteTimeout = time.Duration(0)
			guildSettings[env.Guild.ID].SwearFilter.AllowAdminBypass = false
			guildSettings[env.Guild.ID].SwearFilter.AllowBotOwnerBypass = false
		case "invitegen":
			guildSettings[env.Guild.ID].APIInviteChannel = ""
			guildSettings[env.Guild.ID].APIInviteKey = ""
		default:
			return NewErrorEmbed("Server Settings - Reset Error", "Error finding the setting ``"+args[1]+"``.")
		}
		return NewGenericEmbed("Server Settings - Reset", "Successfully reset the settings for ``"+args[1]+"``.")
	}
	return NewErrorEmbed("Server Settings Error", "Error finding the setting ``"+args[0]+"``.")
}
