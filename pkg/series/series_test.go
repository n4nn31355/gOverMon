package series

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEntry(t *testing.T) {
	type args struct {
		size int
	}
	tests := []struct {
		name  string
		args  args
		panic bool
	}{
		{
			name:  "size: 0",
			args:  args{size: 0},
			panic: true,
		},
		{
			name: "size: 1",
			args: args{size: 1},
		},
		{
			name: "size: 10",
			args: args{size: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				assert.Panics(t, func() { NewEntry(tt.args.size) })
				return
			}
			entry := NewEntry(tt.args.size)
			assert.Equal(t, tt.args.size, entry.size)
			assert.NotNil(t, entry.subscribers)
			assert.False(t, entry.IsActive())
		})
	}
}

func TestEntry_Subscriber_SamePtr(t *testing.T) {
	sub1 := &Subscriber{}
	sub2 := &Subscriber{}

	// NOTE: Do not delete. Printf is a part of check
	fmt.Printf("%p, %p\n", &sub1, &sub2)

	require.NotSame(t, sub1, sub2)
	require.NotSame(t, &sub1, &sub2)
}

func TestEntry_Subscribe_Unsubscribe(t *testing.T) {
	entry := &Entry{
		size:        1,
		subscribers: make(SubscribersMap),
	}
	assert.Empty(t, entry.subscribers)
	assert.Equal(t, SubscribersMap{}, entry.subscribers)
	assert.False(t, entry.IsActive())

	sub1 := &Subscriber{}
	sub2 := &Subscriber{}

	entry.Subscribe(sub1)
	assert.Len(t, entry.subscribers, 1)
	assert.Equal(t, SubscribersMap{sub1: struct{}{}}, entry.subscribers)
	assert.True(t, entry.IsActive())

	assert.Panics(t, func() { entry.Subscribe(sub1) })
	assert.Len(t, entry.subscribers, 1)
	assert.Equal(t, SubscribersMap{sub1: struct{}{}}, entry.subscribers)
	assert.True(t, entry.IsActive())

	entry.Subscribe(sub2)
	assert.Len(t, entry.subscribers, 2)
	subMap := SubscribersMap{sub1: struct{}{}, sub2: struct{}{}}
	assert.Equal(t, subMap, entry.subscribers)
	assert.True(t, entry.IsActive())

	assert.Panics(t, func() { entry.Unsubscribe(&Subscriber{}) })
	assert.Equal(t, subMap, entry.subscribers)
	assert.True(t, entry.IsActive())

	entry.Unsubscribe(sub1)
	assert.Equal(t, SubscribersMap{sub2: struct{}{}}, entry.subscribers)
	assert.True(t, entry.IsActive())

	entry.Unsubscribe(sub2)
	assert.Equal(t, SubscribersMap{}, entry.subscribers)
	assert.False(t, entry.IsActive())
}
