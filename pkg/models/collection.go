package models

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type StatsCollection struct {
	data map[time.Time][]*Stats
	m    sync.Mutex
	keys []time.Time
}

func NewCollection() *StatsCollection {
	return &StatsCollection{
		data: make(map[time.Time][]*Stats),
		keys: make([]time.Time, 0),
	}
}

func (c *StatsCollection) UnmarshalJSON(data []byte) error {
	c.m.Lock()
	defer c.m.Unlock()

	err := json.Unmarshal(data, &c.data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal stats collection data: %w", err)
	}

	c.keys = make([]time.Time, 0, len(c.data))
	for key, _ := range c.data {
		c.keys = append(c.keys, key)
	}

	return nil
}

func (c *StatsCollection) MarshalJSON() ([]byte, error) {
	c.m.Lock()
	defer c.m.Unlock()

	data, err := json.Marshal(c.data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats collection data: %w", err)
	}

	return data, nil
}

func (c *StatsCollection) KeysCount() int {
	c.m.Lock()
	defer c.m.Unlock()

	return len(c.keys)
}

func (c *StatsCollection) KeysSorted() []time.Time {
	c.m.Lock()
	defer c.m.Unlock()

	if len(c.keys) == 0 {
		return nil
	}

	sort.Slice(c.keys, func(i, j int) bool {
		return c.keys[i].Before(c.keys[j])
	})

	return c.keys
}

func (c *StatsCollection) FirstKey() time.Time {
	if c.Len() == 0 {
		return time.Time{}
	}

	return c.KeysSorted()[0]
}

func (c *StatsCollection) LastKey() time.Time {
	if c.KeysCount() == 0 {
		return time.Time{}
	}

	return c.KeysSorted()[c.KeysCount()-1]
}

func (c *StatsCollection) Len() int {
	c.m.Lock()
	defer c.m.Unlock()

	return len(c.data)
}

func (c *StatsCollection) Merge(other *StatsCollection) {
	if other == nil {
		return
	}

	for key, stats := range other.data {
		c.Add(key, stats)
	}
}

func (c *StatsCollection) Add(key time.Time, stats []*Stats) {
	c.m.Lock()
	defer c.m.Unlock()

	key = time.Date(key.Year(), key.Month(), key.Day(), 0, 0, 0, 0, key.Location())

	c.data[key] = stats
	c.keys = append(c.keys, key)
}

func (c *StatsCollection) Get(key time.Time) []*Stats {
	c.m.Lock()
	defer c.m.Unlock()

	return c.data[key]
}

func (c *StatsCollection) DiffPrevious(currentKey time.Time, language string, statType StatType) float64 {
	if c.Len() == 0 || c.Len() == 1 {
		return 0
	}

	previousKey := c.previousKey(currentKey)
	if previousKey.IsZero() {
		return 0
	}

	currentStats := float64(0)
	for _, s := range c.Get(currentKey) {
		if !strings.EqualFold(s.Name, language) {
			continue
		}

		currentStats += s.ValueOf(statType)
	}

	previousStats := float64(0)
	for _, s := range c.Get(previousKey) {
		if !strings.EqualFold(s.Name, language) {
			continue
		}

		previousStats += s.ValueOf(statType)
	}

	return currentStats - previousStats
}

func (c *StatsCollection) previousKey(currentKey time.Time) time.Time {
	if c.Len() == 0 || c.Len() == 1 {
		return time.Time{}
	}

	for i, key := range c.KeysSorted() {
		if key.Equal(currentKey) {
			if i == 0 {
				return time.Time{}
			}

			return c.KeysSorted()[i-1]
		}
	}

	return time.Time{}
}
