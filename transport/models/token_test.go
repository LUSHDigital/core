package models

import (
	"fmt"
	"testing"
)

func TestToken_PrepareForHttp(t *testing.T) {
	tt := []struct {
		name                  string
		token                 Token
		expectedPreparedToken string
	}{
		{
			name: "Normal data JWT",
			token: Token{
				Type:  "JWT",
				Value: "sdfsdfsdfdsfsdf.sdfsdfdsfsdfsdf.sdfdsfsdfsf",
			},
			expectedPreparedToken: "Bearer sdfsdfsdfdsfsdf.sdfsdfdsfsdfsdf.sdfdsfsdfsf",
		},
		{
			name: "Extreme data JWT",
			token: Token{
				Type:  "JWT",
				Value: "*(*)*(D*FDF*)*DF.@H£@£HH@KH@£@H£.()_)(DDFJDFDJKHSF",
			},
			expectedPreparedToken: "Bearer *(*)*(D*FDF*)*DF.@H£@£HH@KH@£@H£.()_)(DDFJDFDJKHSF",
		},
		{
			name: "Normal data random",
			token: Token{
				Type:  "random",
				Value: "jjsjsjsjsjs.kskskskksksks.pqpqpqppqpq",
			},
			expectedPreparedToken: "Bearer jjsjsjsjsjs.kskskskksksks.pqpqpqppqpq",
		},
		{
			name: "Extreme data random",
			token: Token{
				Type:  "random",
				Value: "()**SD()*SS*D*(SD.*&&XC(&X&*X&(X&X&C(XC&.S£S$£$£S%£S%£",
			},
			expectedPreparedToken: "Bearer ()**SD()*SS*D*(SD.*&&XC(&X&*X&(X&X&C(XC&.S£S$£$£S%£S%£",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			preparedToken := tc.token.PrepareForHTTP()
			if preparedToken != tc.expectedPreparedToken {
				t.Errorf("TestToken_PrepareForHttp: %s: expected %v got %v", tc.name, tc.expectedPreparedToken, preparedToken)
			}
		})
	}
}

func ExampleToken_PrepareForHttp() {
	t := Token{
		Type:  "JWT",
		Value: "xxxxxx.xxxxxx.xxxxxx",
	}

	fmt.Println(t.PrepareForHTTP())

	// Output: Bearer xxxxxx.xxxxxx.xxxxxx
}
