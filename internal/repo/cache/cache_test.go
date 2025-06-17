package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndGetAll(t *testing.T) {
	c := New()
	item := "Test item"
	id := c.Add(item)
	items := c.GetAll()
	assert.Greater(t, id, -1)
	assert.NotEmpty(t, items)
	assert.Contains(t, items, item)
}

func TestAddMultipleItems(t *testing.T) {
	c := New()
	item1 := "Test item 1"
	item2 := "Test item 2"
	id1 := c.Add(item1)
	id2 := c.Add(item2)
	items := c.GetAll()
	assert.Contains(t, items, item1)
	assert.Contains(t, items, item2)
	assert.NotEqual(t, id1, id2)
}
