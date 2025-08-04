package repository

import (
	"fmt"

	"github.com/sony/sonyflake/v2"
)

var (
	sf *sonyflake.Sonyflake
)

func SetupIDGenerator() (err error) {
	sf, err = sonyflake.New(sonyflake.Settings{})
	if err != nil {
		return fmt.Errorf("id generator: %w", err)
	}
	
	return nil
}
