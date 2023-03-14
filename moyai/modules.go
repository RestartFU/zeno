package moyai

import "github.com/moyai-studio/practice-revamp/moyai/module"

// registerAllModules registers all modules.
func (m *Moyai) registerAllModules() {
	m.registerModules(
		&module.Click{
			Lobby: m.Server.World(),
		},
		&module.ItemInteraction{},
		&module.Welcome{Host: m, Moyai: m, StaffMap: m.Staffs()},
		&module.Protection{
			Lobby: m.DefaultWorld(),
		},
		&module.Game{
			Host: m,
		},
		&module.Data{Remover: m, DB: m.db},
		&module.Chat{Perm: m},
	)

}

// registerModules registers the modules provided, ensuring that they are ordered by priority.
func (m *Moyai) registerModules(modules ...module.Module) {
	modules = append(modules, m.modules...)

	var sortedModules []module.Module

	sort := func(priority module.Priority) {
		for _, mod := range modules {
			if mod.Priority() == priority {
				sortedModules = append(sortedModules, mod)
			}
		}
	}

	for _, priority := range module.Priorities() {
		sort(priority)
	}

	m.modules = sortedModules
}
