package main

import (
	"errors"
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
type PaymentService struct {
	// not need to implement
	NotEmptyStruct bool
}

var (
	errInvalidConstructor = errors.New("invalid constructor")
	errUnregisteredType   = errors.New("unregistered type")
)

type Container struct {
	constructors map[string]interface{}
	singletons   map[string]struct{}
	instances    map[string]interface{}
}

func NewContainer() *Container {
	return &Container{
		constructors: make(map[string]interface{}),
		singletons:   make(map[string]struct{}),
		instances:    make(map[string]interface{}),
	}
}

func (c *Container) RegisterType(name string, constructor interface{}) {
	c.constructors[name] = constructor
}

func (c *Container) RegisterSingletonType(name string, constructor interface{}) {
	c.RegisterType(name, constructor)
	c.singletons[name] = struct{}{}
}

func (c *Container) build(name string) (interface{}, error) {
	constructor, ok := c.constructors[name]
	if !ok {
		return nil, errUnregisteredType
	}

	fn, ok := constructor.(func() interface{})
	if !ok {
		return nil, errInvalidConstructor
	}

	return fn(), nil
}

func (c *Container) Resolve(name string) (interface{}, error) {
	if _, ok := c.singletons[name]; !ok {
		return c.build(name)
	}

	instance, ok := c.instances[name]
	if ok {
		return instance, nil
	}

	instance, err := c.build(name)
	if err != nil {
		return nil, err
	}

	c.instances[name] = instance
	return instance, nil
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
	assert.ErrorIs(t, err, errUnregisteredType)
	assert.Nil(t, paymentService)

	container.RegisterSingletonType("PaymentService", func() interface{} {
		return &PaymentService{}
	})

	paymentService1, err := container.Resolve("PaymentService")
	assert.NoError(t, err)
	paymentService2, err := container.Resolve("PaymentService")
	assert.NoError(t, err)

	p1 := paymentService1.(*PaymentService)
	p2 := paymentService2.(*PaymentService)
	assert.True(t, p1 == p2)
}
