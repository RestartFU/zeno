package main

import (
	"github.com/RestartFU/gophig"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/jmoiron/sqlx"
	"github.com/moyai-studio/practice-revamp/moyai"
	"github.com/moyai-studio/practice-revamp/moyai/command"
	"github.com/sirupsen/logrus"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

func init() {
	chat.Global.Subscribe(chat.StdoutSubscriber{})
}

func main() {
	db, err := sqlx.Open("sqlite", "./data/players.db")
	if err != nil {
		panic(err)
	}
	db.Exec("CREATE TABLE IF NOT EXISTS playerdata(xuid VARCHAR(200) PRIMARY KEY,username text,role text,cps boolean,scoreboard boolean,kills int,deaths int);")

	//db.Exec("CREATE TABLE IF NOT EXISTS bans(xuid VARCHAR(200) PRIMARY KEY,reason text,expiration timestamp)")
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	logger.Level = logrus.InfoLevel
	m := moyai.New(readConfig(), logger, db)
	registerAllCommands(m)
	loadWorld("./data/arenas/nodebuff", "NoDebuff", m)
	setWorldSettings(m)

	bans := m.Bans()
	m.Allow(&Allower{list: bans})
	registerCommands(command.Ban(bans, m), command.Unban(bans, m))

	err = m.Start()
	if err != nil {
		panic(err)
	}
}
func registerAllCommands(m *moyai.Moyai) {
	registerCommands(
		command.Rank(m),
		command.Freeze(m),
		command.Spawn(m, m.DefaultWorld()),
		command.Stats(),
		command.Rekit(m),
		command.SetWorldSpawn(m),
		command.Ping(),
		command.Kick(m),
		command.Settings(),
		command.GameMode(m),
		command.Whisper(m),
		command.Reply(),
		command.Logs(m),
		command.Teleport(m),
		command.Disguise(m),
		command.Nick(m),
		command.PotionCount(m),
	)
}
func registerCommands(cmds ...cmd.Command) {
	for _, c := range cmds {
		cmd.Register(c)
	}
}

func readConfig() *server.Config {
	var config = server.DefaultConfig()

	var goph = gophig.NewGophig("./config", "toml", 0777)

	if err := goph.GetConf(&config); os.IsNotExist(err) {
		err := goph.SetConf(config)
		if err != nil {
			logrus.Error(err)
			return &config
		}
		return &config
	}
	return &config
}

func loadWorld(path, name string, m *moyai.Moyai) {
	settings := &world.Settings{
		Name: name,
	}
	if err := m.LoadWorld(path, settings, world.Overworld); err != nil {
		log.Fatalln(err)
	}
}

func setWorldSettings(m *moyai.Moyai) {
	for _, w := range m.Worlds() {
		w.Handle(NewWorldHandler(w.World))
		w.SetTime(0)
		w.StopTime()
		w.StopRaining()
		w.StopWeatherCycle()
		w.StopThundering()
		w.SetDifficulty(world.DifficultyHard)
		w.SetRandomTickSpeed(0)
	}
}

type WorldHandler struct {
	world.NopHandler
	world *world.World
}

func NewWorldHandler(w *world.World) *WorldHandler { return &WorldHandler{world: w} }

func (h *WorldHandler) HandleSound(ctx *event.Context, s world.Sound, _ mgl64.Vec3) {
	switch s.(type) {
	case sound.Attack:
		ctx.Cancel()
	}
}
