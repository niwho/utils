package utils

import (
	"github.com/bwmarrin/snowflake"
)

var Default = &UIDFactory{}
type UIDFactory struct {

    IDFactory *snowflake.Node

}

var IDFactory *snowflake.Node

func InitIDFactory(node int64) (err error) {
	IDFactory, err = snowflake.NewNode(node)
	Default = &UIDFactory {
	    IDFactory: IDFactory,
	}
	return
}

func NewIDFactory(node int64) (*UIDFactory, error) {

	factory, err := snowflake.NewNode(node)
	return &UIDFactory {
	    IDFactory: factory,
	}, err
}


func GetNewID() int64 {
	return int64(Default.IDFactory.Generate())
}
