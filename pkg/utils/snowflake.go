package utils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
)

var snowflakeNode *snowflake.Node

func init() {
	var err error
	snowflakeNode, err = snowflake.NewNode(1)
	if err != nil {
		panic(fmt.Errorf("init snowflake failed: %v", err))
	}
}

func SnowflakeId() snowflake.ID {
	next := snowflakeNode.Generate()
	return next
}
