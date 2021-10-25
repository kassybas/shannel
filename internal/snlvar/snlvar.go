package snlvar

import (
	"fmt"
	"strings"
	"sync"
)

// VarTable is a concurrent safe map to store shannel variables
type VarTable struct {
	m        sync.RWMutex
	internal map[string]string
}

// NewVarTable initializes a VarTable instance
func NewVarTable() *VarTable {
	vt := VarTable{}
	vt.internal = make(map[string]string)
	return &vt
}

// Get returns the value stored in the VarTable for a key.
// The ok result indicates whether value was found in the map.
func (vt *VarTable) Get(key string) (value string, ok bool) {
	vt.m.RLock()
	defer vt.m.RUnlock()
	value, ok = vt.internal[key]
	return
}

// Set sets the value for a key
func (vt *VarTable) Set(key, value string) {
	vt.m.Lock()
	defer vt.m.Unlock()
	vt.internal[key] = value
}

// Eval evaluates a string. If the string is in variable format, the variable is evaluated, otherwise the
// string is returned. The last bool return value indicates wheter the evaluation was successful
func (vt *VarTable) Eval(s string) (string, bool) {
	if strings.HasPrefix(s, "$") {
		return vt.Get(strings.TrimPrefix(s, "$"))
	}
	return s, true
}

// DumpShellFormat returns the contents of the vartable in a shell format of `key=value` as a slice of strings
func (vt *VarTable) DumpShellFormat() []string {
	envVars := []string{}
	vt.m.RLock()
	defer vt.m.RUnlock()
	for k, v := range vt.internal {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}
