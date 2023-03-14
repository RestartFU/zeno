package moyai

import (
	"github.com/RestartFU/gophig"
	"github.com/RestartFU/list"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-plus/worldmanager"
	"github.com/jmoiron/sqlx"
	"github.com/moyai-studio/practice-revamp/moyai/game"
	"github.com/moyai-studio/practice-revamp/moyai/module"
	"github.com/moyai-studio/practice-revamp/moyai/permission"
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

//Moyai ...
type Moyai struct {
	*server.Server
	logger *logrus.Logger

	*worldmanager.WorldManager

	modules       []module.Module
	duelProviders map[string]*game.DuelsProvider
	ffaProviders  map[string]*game.FFAProvider

	db *sqlx.DB

	users   map[string]*user.User
	usersMu sync.Mutex

	staffs *StaffMap

	operators *list.List
	bans      *list.List
}

// New returns a new *Moyai.
func New(config *server.Config, logger *logrus.Logger, db *sqlx.DB) *Moyai {
	s := server.New(config, logger)
	op, err := list.New(&list.Settings{Gophig: gophig.NewGophig("operators", "toml", 0777), CacheOnly: false})
	if err != nil {
		panic(err)
	}
	bans, err := list.New(&list.Settings{Gophig: gophig.NewGophig("bans", "json", 0777), CacheOnly: false})
	if err != nil {
		panic(err)
	}
	return &Moyai{
		Server:       s,
		WorldManager: worldmanager.New(s, logger),

		logger: logger,
		db:     db,

		duelProviders: make(map[string]*game.DuelsProvider),
		ffaProviders:  make(map[string]*game.FFAProvider),

		staffs: NewStaffMap(),

		users:     make(map[string]*user.User),
		operators: op,
		bans:      bans,
	}
}

// Start starts the server.
func (m *Moyai) Start() error {
	if err := m.Server.Start(); err != nil {
		return err
	}
	m.registerAllModules()
	m.registerAllItems()
	m.registerAll()
	m.registerAllRanks()
	m.CloseOnProgramEnd()
	for {
		p, err := m.Server.Accept()
		if err != nil {
			break
		}

		go m.handleJoin(p)
	}
	return nil
}

func (m *Moyai) CloseOnProgramEnd() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		if err := m.WorldManager.Close(); err != nil {
			m.logger.Errorf("error shutting down server: %v", err)
		}
		if err := m.Server.Close(); err != nil {
			m.logger.Errorf("error shutting down server: %v", err)
		}
	}()
}

// handleJoin handles when a player joins.
func (m *Moyai) handleJoin(p *player.Player) {
	u := user.New(p, m.db, m, m.Staffs())
	m.usersMu.Lock()
	m.users[p.Name()] = u
	m.usersMu.Unlock()
	p.Handle(&handler{user: u, server: m})
	for _, m := range m.modules {
		m.HandleJoin(u)
	}
}

func (m *Moyai) RemoveUser(u *user.User) {
	m.usersMu.Lock()
	delete(m.users, u.Name())
	m.usersMu.Unlock()
}

func (m *Moyai) Staffs() *StaffMap { return m.staffs }

func (m *Moyai) DisguisedUser(name string) (*user.User, bool) {
	m.usersMu.Lock()
	defer m.usersMu.Unlock()
	for _, u := range m.users {
		if strings.EqualFold(name, u.DisguisedName()) && u.Disguised() {
			return u, true
		}
	}
	return nil, false
}

func (m *Moyai) User(name string) (*user.User, bool) {
	m.usersMu.Lock()
	u, ok := m.users[name]
	m.usersMu.Unlock()
	return u, ok
}

// RequestFFAProvider ...
func (m *Moyai) RequestFFAProvider(game game.Game) (*game.FFAProvider, bool) {
	provider, ok := m.ffaProviders[game.Name()]
	return provider, ok
}

// RequestDuelProvider ...
func (m *Moyai) RequestDuelProvider(game game.Game) (*game.DuelsProvider, bool) {
	provider, ok := m.duelProviders[game.Name()]
	return provider, ok
}

// providers returns all providers.
func (m *Moyai) providers() []game.Provider {
	providers := make([]game.Provider, 0, len(m.ffaProviders))
	for _, p := range m.ffaProviders {
		providers = append(providers, p)
	}
	return providers
}

func (m *Moyai) Playing() (n int) {
	for _, p := range m.providers() {
		n += len(p.Users())
	}
	return
}

func (m *Moyai) HasPermission(src cmd.Source, flag uint64) bool {
	if m.Operator(src.Name()) {
		return true
	}
	if u, ok := m.User(src.Name()); ok {
		if r, ok := u.Rank(); ok {
			return r.Flags()&permission.FlagAdministrator > 0 || r.Flags()&flag > 0
		}
	}
	return false
}

// SearchUser ...
func (m *Moyai) SearchUser(user *user.User) (game.Provider, bool) {
	for _, prov := range m.providers() {
		if prov.HasUser(user) {
			return prov, true
		}
	}
	return nil, false
}

// Operator ...
func (m *Moyai) Operator(name string) bool { return m.operators.Listed(name) }

/*type bannedData struct {
	XUID       string    `db:"xuid"`
	Reason     string    `db:"reason"`
	Expiration time.Time `db:"expiration"`
}

// Banned ...
func (m *Moyai) Banned(xuid string) (bool, *bannedData) {
	data := &bannedData{}
	rows, err := m.db.Queryx("SELECT * FROM bans WHERE xuid = ?", xuid)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return false, data
	}
	if data.XUID == "" {
		return false, data
	}
	for rows.Next() {
		_ = rows.StructScan(&data)
	}
	if data.Expiration.After(time.Now()) {
		return true, data
	}
	return false, data
}

// Ban ...
func (m *Moyai) Ban(xuid string, reason string, d time.Duration) {
	if xuid == "" {
		return
	}
	_, err := m.db.Exec("REPLACE INTO bans VALUES ($1, $2, $3)",
		xuid,
		reason,
		time.Now().Add(d),
	)
	if err != nil {
		fmt.Println(err)
	}
}

// UnBan ...
func (m *Moyai) UnBan(xuid string) {
	if xuid == "" {
		return
	}
	_, err := m.db.Exec("DELETE FROM bans WHERE xuid = ?", xuid)
	if err != nil {
		fmt.Println(err)
	}
}*/
func (m *Moyai) Bans() *list.List { return m.bans }
