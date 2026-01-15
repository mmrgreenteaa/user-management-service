package redis

import (
	"testing"
)

func TestConnect(t *testing.T) {

	_ = Connect(&RCongif{})

}

func TestAddJwt(t *testing.T) {

	db := Connect(&RCongif{})

	tests := []struct {
		name   string
		jwt    string
		userid string
		err    error
	}{
		{
			name:   "succes add",
			jwt:    "testtesttesttest",
			userid: "5",
			err:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddJwt(tt.jwt, tt.userid)
			if err != tt.err {
				if err != nil {
					t.Error(err)
				}
			}
		})
	}

}

func TestGetjwt(t *testing.T) {

	db := Connect(&RCongif{})

	tests := []struct {
		name   string
		jwt    string
		userid string
		err    error
	}{
		{
			name:   "succes to get",
			jwt:    "testtesttesttest",
			userid: "5",
			err:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userid, err := db.GetJwt(tt.jwt)
			if err != nil {
				t.Error(err)
			}
			t.Log(userid)

			if userid != tt.userid {
				t.Fail()
			}

		})
	}

}
