package goconsist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDistribute(t *testing.T) {
	t.Run("should generate ring with 4 sections and step 1", func(t *testing.T) {
		singleRange := uint32(1)
		rangesCount := uint32(4)

		expected := []section{
			{min: 0, max: 1},
			{min: 2, max: 3},
			{min: 4, max: 5},
			{min: 6, max: 0},
		}

		res := distribute(singleRange, rangesCount)
		assert.Equal(t, expected, res)
	})

	t.Run("should generate ring with 6 sections and step 15", func(t *testing.T) {
		singleRange := uint32(15)
		rangesCount := uint32(6)

		expected := []section{
			{min: 0, max: 15},
			{min: 16, max: 31},
			{min: 32, max: 47},
			{min: 48, max: 63},
			{min: 64, max: 79},
			{min: 80, max: 0},
		}

		res := distribute(singleRange, rangesCount)
		assert.Equal(t, expected, res)
	})
}

func TestVshard_More(t *testing.T) {
	shard := section{max: 5}

	t.Run("should return true if specified x is more then section max", func(t *testing.T) {
		res := shard.more(6)
		assert.True(t, res)
	})

	t.Run("should return false if specified x is less then section max", func(t *testing.T) {
		res := shard.more(4)
		assert.False(t, res)
	})

	t.Run("should return false if specified x is equals to section max", func(t *testing.T) {
		res := shard.more(5)
		assert.False(t, res)
	})
}

func TestVshard_Less(t *testing.T) {
	shard := section{min: 5}

	t.Run("should return true if specified x is less then section min", func(t *testing.T) {
		res := shard.less(4)
		assert.True(t, res)
	})

	t.Run("should return false if specified x is more then section min", func(t *testing.T) {
		res := shard.less(6)
		assert.False(t, res)
	})

	t.Run("should return false if specified x is equals to section min", func(t *testing.T) {
		res := shard.less(5)
		assert.False(t, res)
	})
}

func TestVshard_Between(t *testing.T) {
	shard := section{min: 4, max: 6}

	t.Run("should return true if specified x is more then min factor and current section is last", func(t *testing.T) {
		shard := section{min: 4, max: 0}

		res := shard.between(5)
		assert.True(t, res)
	})

	t.Run("should return true if specified x is between section min and max", func(t *testing.T) {
		res := shard.between(5)
		assert.True(t, res)
	})

	t.Run("should return true if specified x is equals to section min", func(t *testing.T) {
		res := shard.between(4)
		assert.True(t, res)
	})

	t.Run("should return true if specified x is equals to section max", func(t *testing.T) {
		res := shard.between(6)
		assert.True(t, res)
	})

	t.Run("should return false if specified x is outside section min and max range", func(t *testing.T) {
		res := shard.between(3)
		assert.False(t, res)
	})
}

func TestSearch(t *testing.T) {
	shards := []section{
		{min: 0, max: 15},
		{min: 16, max: 31},
		{min: 32, max: 47},
		{min: 48, max: 63},
		{min: 64, max: 79},
		{min: 80, max: 0},
	}

	t.Run("should return section and true if target found", func(t *testing.T) {
		expected := section{min: 32, max: 47}

		res, ok := search(shards, 34)
		assert.True(t, ok)
		assert.Equal(t, expected, res)
	})

	t.Run("should return true if target more then min of last section", func(t *testing.T) {
		expected := section{min: 80, max: 0}

		res, ok := search(shards, 100)
		assert.True(t, ok)
		assert.Equal(t, expected, res)
	})

	t.Run("should return true if target equals to first section", func(t *testing.T) {
		expected := section{min: 0, max: 15}

		res, ok := search(shards, 0)
		assert.True(t, ok)
		assert.Equal(t, expected, res)
	})

	t.Run("should return true if target equals to last section", func(t *testing.T) {
		expected := section{min: 80, max: 0}

		res, ok := search(shards, 80)
		assert.True(t, ok)
		assert.Equal(t, expected, res)
	})
}
