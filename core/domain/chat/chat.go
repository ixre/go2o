package chat

import "github.com/ixre/go2o/core/domain/interface/chat"

var _ chat.IChatUserAggregateRoot = new(chatUserAggregateRoot)

type chatUserAggregateRoot struct {
}

// GetAggregateRootId implements chat.IChatUserAggregateRoot.
func (c *chatUserAggregateRoot) GetAggregateRootId() int {
	panic("unimplemented")
}
