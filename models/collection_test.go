package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatsCollection_calculateNewKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		collection *StatsCollection
		key        string
		want       string
	}{
		{
			name:       "first_time_key",
			collection: NewCollection(),
			key:        "key",
			want:       "key (1)",
		},
		{
			name:       "second_time_key",
			collection: NewCollection(),
			key:        "key (1)",
			want:       "key (2)",
		},
		{
			name:       "third_time_key",
			collection: NewCollection(),
			key:        "key (2)",
			want:       "key (3)",
		},
		{
			name:       "key_no_numeric_suffix",
			collection: NewCollection(),
			key:        "key (n)",
			want:       "key (n) (1)",
		},
		{
			name:       "key_value_suffix",
			collection: NewCollection(),
			key:        "key ()",
			want:       "key () (1)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.collection.calculateNewKey(test.key)
			assert.Equal(t, test.want, got)
		})
	}
}
