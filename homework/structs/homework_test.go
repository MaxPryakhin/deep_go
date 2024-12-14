package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		nameBytes := []byte(name)
		byteArray := make([]byte, 43)
		copy(byteArray, nameBytes)
		person.name = [42]byte(byteArray)
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
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manahealth[0] = byte(mana >> manaShiftTopBits)
		person.manahealth[1] |= byte((mana << manaShiftLowBits) & manaMaskLowBits)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manahealth[1] |= byte(health >> healthShiftTopBits)
		person.manahealth[2] = byte(health & 0xFF)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes = (person.attributes &^ RespectMask) | uint16(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes = (person.attributes &^ StrengthMask) | uint16(strength)<<StrengthShift
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes = (person.attributes &^ ExperienceMask) | uint16(experience)<<ExperienceShift
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attributes = (person.attributes &^ LevelMask) | uint16(level)<<LevelShift
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= HasHouseFlag
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= HasWeaponFlag
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.flags |= HasFamilyFlag
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		if personType == BuilderGamePersonType {
			person.flags |= HasBuilderTypeFlag
		}
		if personType == WarriorGamePersonType {
			person.flags |= HasWarriorTypeFlag
		}
		if personType == BlacksmithGamePersonType {
			person.flags |= HasBlacksmithTypeFlag
		}
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

const (
	HasHouseFlag uint8 = 1 << iota
	HasWeaponFlag
	HasFamilyFlag
	HasBuilderTypeFlag
	HasBlacksmithTypeFlag
	HasWarriorTypeFlag
)

type GamePerson struct {
	x          int32
	y          int32
	z          int32
	gold       int32
	attributes uint16
	flags      uint8
	manahealth [3]byte
	name       [42]byte
}

const (
	RespectMask    = 0xF    // 0000 0000 0000 1111
	StrengthMask   = 0xF0   // 0000 0000 1111 0000
	ExperienceMask = 0xF00  // 0000 1111 0000 0000
	LevelMask      = 0xF000 // 1111 0000 0000 0000

	StrengthShift   = 4
	ExperienceShift = 8
	LevelShift      = 12

	manaShiftTopBits = 2
	manaShiftLowBits = 6
	manaMaskLowBits  = 0b11000000

	healthShiftTopBits = 8
	healthMaskTopBits  = 0b00000011
	healthMaskLowBits  = 0b11111111
)

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	if len(p.name) == 0 {
		return ""
	}

	end := len(p.name)
	for index, v := range p.name {
		if v == 0 {
			end = index
			break
		}
	}

	return string(p.name[:end])
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
	return int(p.manahealth[0])<<manaShiftTopBits + int(p.manahealth[1]&manaMaskLowBits>>manaShiftLowBits)
}

func (p *GamePerson) Health() int {
	return int(p.manahealth[1])&healthMaskTopBits<<healthShiftTopBits + int(p.manahealth[2])
}

func (p *GamePerson) Respect() int {
	return int(p.attributes & RespectMask)
}

func (p *GamePerson) Strength() int {
	return int(p.attributes&StrengthMask) >> StrengthShift
}

func (p *GamePerson) Experience() int {
	return int(p.attributes&ExperienceMask) >> ExperienceShift
}

func (p *GamePerson) Level() int {
	return int(p.attributes&LevelMask) >> LevelShift
}

func (p *GamePerson) HasHouse() bool {
	return p.flags&HasHouseFlag != 0
}

func (p *GamePerson) HasGun() bool {
	return p.flags&HasWeaponFlag != 0
}

func (p *GamePerson) HasFamilty() bool {
	return p.flags&HasFamilyFlag != 0
}

func (p *GamePerson) Type() int {
	if p.flags&HasBlacksmithTypeFlag != 0 {
		return BlacksmithGamePersonType
	}
	if p.flags&HasBuilderTypeFlag != 0 {
		return BuilderGamePersonType
	}
	if p.flags&HasWarriorTypeFlag != 0 {
		return WarriorGamePersonType
	}

	return 0
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
