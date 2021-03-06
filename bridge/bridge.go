package bridge

import (
	"strings"

	"github.com/42wim/matterbridge/bridge/config"
	"github.com/sirupsen/logrus"
)

type Bridger interface {
	Send(msg config.Message) (string, error)
	Connect() error
	JoinChannel(channel config.ChannelInfo) error
	Disconnect() error
}

type Bridge struct {
	Bridger
	Name     string
	Account  string
	Protocol string
	Channels map[string]config.ChannelInfo
	Joined   map[string]bool
	Log      *logrus.Entry
	Config   config.Config
	General  *config.Protocol
}

type Config struct {
	//	General *config.Protocol
	Remote chan config.Message
	Log    *logrus.Entry
	*Bridge
}

// Factory is the factory function to create a bridge
type Factory func(*Config) Bridger

func New(bridge *config.Bridge) *Bridge {
	b := new(Bridge)
	b.Channels = make(map[string]config.ChannelInfo)
	accInfo := strings.Split(bridge.Account, ".")
	protocol := accInfo[0]
	name := accInfo[1]
	b.Name = name
	b.Protocol = protocol
	b.Account = bridge.Account
	b.Joined = make(map[string]bool)
	return b
}

func (b *Bridge) JoinChannels() error {
	err := b.joinChannels(b.Channels, b.Joined)
	return err
}

func (b *Bridge) joinChannels(channels map[string]config.ChannelInfo, exists map[string]bool) error {
	for ID, channel := range channels {
		if !exists[ID] {
			b.Log.Infof("%s: joining %s (ID: %s)", b.Account, channel.Name, ID)
			err := b.JoinChannel(channel)
			if err != nil {
				return err
			}
			exists[ID] = true
		}
	}
	return nil
}

func (b *Bridge) GetBool(key string) bool {
	val, ok := b.Config.GetBool(b.Account + "." + key)
	if !ok {
		val, _ = b.Config.GetBool("general." + key)
	}
	return val
}

func (b *Bridge) GetInt(key string) int {
	val, ok := b.Config.GetInt(b.Account + "." + key)
	if !ok {
		val, _ = b.Config.GetInt("general." + key)
	}
	return val
}

func (b *Bridge) GetString(key string) string {
	val, ok := b.Config.GetString(b.Account + "." + key)
	if !ok {
		val, _ = b.Config.GetString("general." + key)
	}
	return val
}

func (b *Bridge) GetStringSlice(key string) []string {
	val, ok := b.Config.GetStringSlice(b.Account + "." + key)
	if !ok {
		val, _ = b.Config.GetStringSlice("general." + key)
	}
	return val
}

func (b *Bridge) GetStringSlice2D(key string) [][]string {
	val, ok := b.Config.GetStringSlice2D(b.Account + "." + key)
	if !ok {
		val, _ = b.Config.GetStringSlice2D("general." + key)
	}
	return val
}
