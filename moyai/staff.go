package moyai

import (
	"github.com/moyai-studio/practice-revamp/moyai/user"
	"sync"
)

// StaffMap is the custom type with custom methods for the staff session map.
type StaffMap struct {
	staffs   map[string]*user.User
	staffsMu sync.Mutex
}

func NewStaffMap() *StaffMap { return &StaffMap{staffs: make(map[string]*user.User)} }

func (m *StaffMap) AddStaff(u *user.User) {
	m.staffsMu.Lock()
	defer m.staffsMu.Unlock()
	m.staffs[u.Name()] = u
}
func (m *StaffMap) RemoveStaff(u *user.User) {
	m.staffsMu.Lock()
	defer m.staffsMu.Unlock()
	delete(m.staffs, u.Name())
}
func (m *StaffMap) Staffs() map[string]*user.User {
	m.staffsMu.Lock()
	defer m.staffsMu.Unlock()
	return m.staffs
}

func (m *StaffMap) Staff(name string) bool {
	m.staffsMu.Lock()
	defer m.staffsMu.Unlock()
	_, ok := m.staffs[name]
	return ok
}

// Message will send a message to all online staff.
func (m *StaffMap) Message(a ...interface{}) {
	for _, v := range m.staffs {
		v.Player.Message(a...)
	}
}

// Messagef will send a formatted message to all online staff.
func (m *StaffMap) Messagef(f string, a ...interface{}) {
	for _, v := range m.staffs {
		v.Player.Messagef(f, a...)
	}
}

// Whisper will send a gray italic message to all online staff.
func (m *StaffMap) Whisper(f string, a ...interface{}) {
	for _, v := range m.staffs {
		v.Player.Messagef("§7§o["+f+"]", a...)
	}
}
