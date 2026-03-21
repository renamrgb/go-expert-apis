package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testEmail = "j@j.com"
const testName = "john Doe"

func TestNewUser(t *testing.T) {
	user, err := NewUser(testName, testEmail, "123456")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.NotEmpty(t, user.Password)
	assert.Equal(t, testName, user.Name)
	assert.Equal(t, testEmail, user.Email)
}

func TestUserValidatePassword(t *testing.T) {
	user, err := NewUser(testName, testEmail, "123456")
	assert.Nil(t, err)
	assert.NotEqual(t, "123456", user.Password)
	assert.True(t, user.ValidatePassword("123456"))
	assert.False(t, user.ValidatePassword("1234567"))
}
