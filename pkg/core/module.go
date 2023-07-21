package core

import ()

const (
	tryReadKeyTimes = 5
)

// Module interface
type Module interface {
	Name() string
	Group() string
	Enable() bool
	Start()
}

var (
	// Modules map
	modules         map[string]Module
	disabledModules map[string]Module
)

func init() {
	modules = make(map[string]Module)
	disabledModules = make(map[string]Module)
}

// Register register module
func Register(m Module) {
	if m.Enable() {
		modules[m.Name()] = m
	} else {
		disabledModules[m.Name()] = m
	}
}

// GetModules gets modules map
func GetModules() map[string]Module {
	return modules
}
