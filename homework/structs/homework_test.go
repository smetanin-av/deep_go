package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	maskMana       uint16 = 0b0000_0011_1111_1111
	maskNameLen    uint16 = 0b1111_1100_0000_0000
	shiftNameLen          = 10
	maskHealth     uint16 = 0b0000_0011_1111_1111
	maskHouse      uint16 = 0b0000_0100_0000_0000
	maskGun        uint16 = 0b0000_1000_0000_0000
	maskFamily     uint16 = 0b0001_0000_0000_0000
	maskType       uint16 = 0b0110_0000_0000_0000
	shiftType             = 13
	maskRespect    byte   = 0b0000_1111
	maskStrength   byte   = 0b1111_0000
	shiftStrength         = 4
	maskExperience byte   = 0b0000_1111
	maskLevel      byte   = 0b1111_0000
	shiftLevel            = 4
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
		size := min(len(person.name), len(name))
		person.f1 = (person.f1 & ^maskNameLen) | uint16(size<<shiftNameLen)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 || mana > 1000 {
			return
		}
		person.f1 = (person.f1 & ^maskMana) | uint16(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 || health > 1000 {
			return
		}
		person.f2 = (person.f2 & ^maskHealth) | uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect < 0 || respect > 10 {
			return
		}
		person.f3 = (person.f3 & ^maskRespect) | byte(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength < 0 || strength > 10 {
			return
		}
		person.f3 = (person.f3 & ^maskStrength) | byte(strength<<shiftStrength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience < 0 || experience > 10 {
			return
		}
		person.f4 = (person.f4 & ^maskExperience) | byte(experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level < 0 || level > 10 {
			return
		}
		person.f4 = (person.f4 & ^maskLevel) | byte(level<<shiftLevel)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.f2 = person.f2 | maskHouse
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.f2 = person.f2 | maskGun
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.f2 = person.f2 | maskFamily
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType < BuilderGamePersonType || personType > WarriorGamePersonType {
			return
		}
		person.f2 = (person.f2 & ^maskType) | uint16(personType<<shiftType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	name    [42]byte
	f1, f2  uint16
	f3, f4  byte
	x, y, z int32
	gold    uint32
}

func NewGamePerson(options ...Option) GamePerson {
	var res GamePerson
	for _, opt := range options {
		opt(&res)
	}
	return res
}

func (p *GamePerson) Name() string {
	nameLen := (p.f1 & maskNameLen) >> shiftNameLen
	return unsafe.String(unsafe.SliceData(p.name[:]), nameLen)
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.f1 & maskMana)
}

func (p *GamePerson) Health() int {
	return int(p.f2 & maskHealth)
}

func (p *GamePerson) Respect() int {
	return int(p.f3 & maskRespect)
}

func (p *GamePerson) Strength() int {
	return int(p.f3&maskStrength) >> shiftStrength
}

func (p *GamePerson) Experience() int {
	return int(p.f4 & maskExperience)
}

func (p *GamePerson) Level() int {
	return int(p.f3&maskLevel) >> shiftLevel
}

func (p *GamePerson) HasHouse() bool {
	return p.f2&maskHouse != 0
}

func (p *GamePerson) HasGun() bool {
	return p.f2&maskGun != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.f2&maskFamily != 0
}

func (p *GamePerson) Type() int {
	return int(p.f2&maskType) >> shiftType
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
