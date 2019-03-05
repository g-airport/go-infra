package match

import (
	"regexp"
	"sync"
)

type Match struct {
	Data         interface{}
	Key          string
	matchMap     map[string]*regexp.Regexp
	fastMatchMap map[string]string
	rw           sync.RWMutex
}

func (b *Match) GetData(dataKey string) (string, bool) {
	b.rw.RLock()
	v, ok := b.fastMatchMap[dataKey]
	b.rw.RUnlock()
	if ok {
		return v, true
	}
	for fv, reg := range b.matchMap {
		ok = reg.MatchString(dataKey)
		if ok {
			v = fv
			break
		}
	}
	if !ok {
		return v, false
	}

	b.rw.Lock()
	b.fastMatchMap[dataKey] = v
	b.rw.Unlock()
	return v, true
}

type MFactory struct {
	m  map[string]*Match
	rw sync.RWMutex
}

func (m *MFactory) Get(s string) (*Match, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	v, ok := m.m[s]
	return v, ok
}

func (m *MFactory) ForeachMatch(f func(key *Match) error) error {
	m.rw.RLock()
	defer m.rw.RUnlock()
	for _, key := range m.m {
		err := f(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MFactory) Reset(dataList []*Match) {
	m.rw.Lock()
	defer m.rw.Unlock()
	for _, d := range dataList {
		m.m[d.Key] = d
	}
}
