package models

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	uniqueKeyRegex = regexp.MustCompile(`(?m)(.+)\((\d+)\)$`)
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

func (c *StatsCollection) Keys() []time.Time {
	c.m.Lock()
	defer c.m.Unlock()

	return c.keys
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

	c.data[key] = stats
	c.keys = append(c.keys, key)
}

func (c *StatsCollection) Get(key time.Time) []*Stats {
	c.m.Lock()
	defer c.m.Unlock()

	return c.data[key]
}

func (c *StatsCollection) calculateNewKey(key string) string {
	index := 1
	prefix := key

	if uniqueKeyRegex.MatchString(key) {
		prefix = strings.TrimSpace(uniqueKeyRegex.FindStringSubmatch(key)[1])
		indexRaw := uniqueKeyRegex.FindStringSubmatch(key)[2]

		currentIndex, err := strconv.Atoi(indexRaw)
		if err != nil {
			log.Println("failed to parse index from key:", key)
		}

		index = currentIndex + 1
	}

	return fmt.Sprintf("%s (%d)", prefix, index)
}

func (c *StatsCollection) containsKey(key time.Time) bool {
	_, ok := c.data[key]

	return ok
}
