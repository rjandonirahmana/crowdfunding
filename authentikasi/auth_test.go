package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateValidateToken(t *testing.T) {
	// Testcases list
	testCases := []struct {
		testName string
		id       uint
		exp      int64
	}{
		{
			testName: "success",
			id:       10001786343,
			exp:      time.Now().Add(time.Hour * 10).Unix(),
		},
		{
			testName: "success",
			id:       866464744,
			exp:      time.Now().Add(time.Hour * 10).Unix(),
		}, {
			testName: "success",
			id:       7334161711,
			exp:      time.Now().Add(time.Hour * 8).Unix(),
		},
	}

	for _, testCase := range testCases {
		newToken, err := NewAuthentication("coba").GenerateToken(testCase.id)
		assert.Nil(t, err)
		fmt.Println(newToken)

		user_id, er := NewAuthentication("coba").ValidateToken(newToken)
		assert.Nil(t, er)

		assert.Equal(t, testCase.id, user_id)

	}
}
