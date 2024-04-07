package database

import (
	"github.com/alicebob/miniredis/v2"
)

func InitializeTestingRedis() *miniredis.Miniredis {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	return s
}
