package handlers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCepInfo(t *testing.T) {
	cepInfo, err := getCepInfo("01311200")

	assert.NotNil(t, cepInfo)
	assert.Nil(t, err)
	fmt.Println(cepInfo)
}
