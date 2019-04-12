package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test_FromRau(t *testing.T) {
	str, err := FromRau("oliver", "abc")
	assert.NoError(t, err)
	assert.Equal(t, nil, str)
}

func Test_ToRau(t *testing.T) {
	str, err := ToRau("oliver", "qiuchang")
	assert.NoError(t, err)
	assert.Equal(t, nil, str)
}