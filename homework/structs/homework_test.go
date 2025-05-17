package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

const (
	ManaMask    = 0x3FF
	ManaShift   = 0
	HealthMask  = 0x3FF
	HealthShift = 10
	HouseFlag   = 1 << 20
	GunFlag     = 1 << 21
	FamilyFlag  = 1 << 22
	TypeMask    = 0x03
	TypeShift   = 23

	RespectMask     = 0x0F
	RespectShift    = 0
	StrengthMask    = 0x0F
	StrengthShift   = 4
	ExperienceMask  = 0x0F
	ExperienceShift = 8
	LevelMask       = 0x0F
	LevelShift      = 12
)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copyLen := len(name)
		if copyLen > len(person.name) {
			copyLen = len(person.name)
		}

		for i := range person.name {
			person.name[i] = 0
		}

		copy(person.name[:], name[:copyLen])
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
		if gold < 0 {
			gold = 0
		}
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana < 0 {
			mana = 0
		} else if mana > 1023 {
			mana = 1023
		}
		manaValue := uint32(mana) & ManaMask
		person.attributesMask = (person.attributesMask &^ ManaMask) | manaValue
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health < 0 {
			health = 0
		} else if health > 1023 {
			health = 1023
		}
		healthValue := uint32(health) & HealthMask
		person.attributesMask = (person.attributesMask &^ (HealthMask << HealthShift)) | (healthValue << HealthShift)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect < 0 {
			respect = 0
		} else if respect > 15 {
			respect = 15
		}
		respectValue := uint16(respect) & RespectMask
		person.statsMask = (person.statsMask &^ RespectMask) | respectValue
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength < 0 {
			strength = 0
		} else if strength > 15 {
			strength = 15
		}
		strengthValue := uint16(strength) & StrengthMask
		person.statsMask = (person.statsMask &^ (StrengthMask << StrengthShift)) | (strengthValue << StrengthShift)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience < 0 {
			experience = 0
		} else if experience > 15 {
			experience = 15
		}
		expValue := uint16(experience) & ExperienceMask
		person.statsMask = (person.statsMask &^ (ExperienceMask << ExperienceShift)) | (expValue << ExperienceShift)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level < 0 {
			level = 0
		} else if level > 15 {
			level = 15
		}
		levelValue := uint16(level) & LevelMask
		person.statsMask = (person.statsMask &^ (LevelMask << LevelShift)) | (levelValue << LevelShift)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= HouseFlag
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= GunFlag
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributesMask |= FamilyFlag
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType < 0 || personType > 3 {
			personType = 0
		}
		person.attributesMask = (person.attributesMask &^ (TypeMask << TypeShift)) | (uint32(personType&TypeMask) << TypeShift)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z        int32
	gold           uint32
	attributesMask uint32
	statsMask      uint16
	name           [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	var length int
	for i, b := range p.name {
		if b == 0 {
			length = i
			break
		}
	}
	if length == 0 && p.name[0] != 0 {
		length = len(p.name)
	}
	return string(p.name[:length])
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
	return int(p.attributesMask & ManaMask)
}

func (p *GamePerson) Health() int {
	return int((p.attributesMask >> HealthShift) & HealthMask)
}

func (p *GamePerson) Respect() int {
	return int(p.statsMask & RespectMask)
}

func (p *GamePerson) Strength() int {
	return int((p.statsMask >> StrengthShift) & StrengthMask)
}

func (p *GamePerson) Experience() int {
	return int((p.statsMask >> ExperienceShift) & ExperienceMask)
}

func (p *GamePerson) Level() int {
	return int((p.statsMask >> LevelShift) & LevelMask)
}

func (p *GamePerson) HasHouse() bool {
	return (p.attributesMask & HouseFlag) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.attributesMask & GunFlag) != 0
}

func (p *GamePerson) HasFamilty() bool {
	return (p.attributesMask & FamilyFlag) != 0
}

func (p *GamePerson) Type() int {
	return int((p.attributesMask >> TypeShift) & TypeMask)
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
