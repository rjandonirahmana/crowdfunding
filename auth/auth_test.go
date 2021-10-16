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
			id:       3,
			exp:      time.Now().Add(time.Hour * 10).Unix(),
		},
		{
			testName: "success",
			id:       4,
			exp:      time.Now().Add(time.Hour * 10).Unix(),
		}, {
			testName: "success",
			id:       1,
			exp:      time.Now().Add(time.Hour * 8).Unix(),
		},
	}

	for _, testCase := range testCases {
		newToken, err := NewAuthentication("coba", "apaaja").GenerateToken(testCase.id)
		assert.Nil(t, err)
		fmt.Println(newToken)

		newTokenAdmin, err := NewAuthentication("coba", "apaaja").GenerateTokenAdmin(testCase.id)
		assert.Nil(t, err)
		fmt.Println(newTokenAdmin)

		user_id, er := NewAuthentication("coba", "apaaja").ValidateToken(newToken)
		assert.Nil(t, er)

		admin_id, er := NewAuthentication("coba", "apaaja").ValidateTokenAdmin(newTokenAdmin)
		assert.Nil(t, er)

		assert.Equal(t, testCase.id, user_id)
		assert.Equal(t, testCase.id, admin_id)

	}
}
