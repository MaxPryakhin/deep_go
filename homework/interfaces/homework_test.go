package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UserService struct {
	// not need to implement
	NotEmptyStruct bool
}
type MessageService struct {
	// not need to implement
	NotEmptyStruct bool
}

type Container struct {
	constructors map[string]any
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]any),
	}
}

func (c *Container) RegisterType(name string, constructor any) {
	if reflect.TypeOf(constructor).Kind() != reflect.Func {
		panic("constructor must be a function")
	}
	c.constructors[name] = constructor
}

func (c *Container) Resolve(name string) (any, error) {
	constructorRaw, found := c.constructors[name]
	if !found {
		return nil, fmt.Errorf("no constructor registered for type: %s", name)
	}

	constructor := reflect.ValueOf(constructorRaw)
	args := make([]reflect.Value, 0)
	results := constructor.Call(args)

	return results[0].Interface(), nil
}

func TestDIContainer(t *testing.T) {
	container := NewContainer()
	container.RegisterType("UserService", func() interface{} {
		return &UserService{}
	})
	container.RegisterType("MessageService", func() interface{} {
		return &MessageService{}
	})

	userService1, err := container.Resolve("UserService")
	assert.NoError(t, err)
	userService2, err := container.Resolve("UserService")
	assert.NoError(t, err)

	u1 := userService1.(*UserService)
	u2 := userService2.(*UserService)
	assert.False(t, u1 == u2)

	messageService, err := container.Resolve("MessageService")
	assert.NoError(t, err)
	assert.NotNil(t, messageService)

	paymentService, err := container.Resolve("PaymentService")
	assert.Error(t, err)
	assert.Nil(t, paymentService)
}
