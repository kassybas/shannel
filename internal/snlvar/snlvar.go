package snlvar

import (
	"fmt"
	"strings"
	"sync"
)

type VarTable struct {
	sync.RWMutex
	internal map[string]string
}

func NewVarTable() *VarTable {
	vt := VarTable{}
	vt.internal = make(map[string]string)
	return &vt
}

func (vt *VarTable) Get(key string) (value string, ok bool) {
	vt.RLock()
	defer vt.RUnlock()
	value, ok = vt.internal[key]
	return
}

func (vt *VarTable) Set(key, value string) {
	vt.Lock()
	defer vt.Unlock()
	vt.internal[key] = value
}

func (vt *VarTable) Eval(s string) (string, bool) {
	if strings.HasPrefix(s, "$") {
		return vt.Get(strings.TrimPrefix(s, "$"))
	}
	return s, true
}

func (vt *VarTable) DumpShellFormat() []string {
	envVars := []string{}
	vt.RLock()
	defer vt.RUnlock()
	for k, v := range vt.internal {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}
	return envVars
}
