package domain

import (
	"testing"

	"github.com/ixre/go2o/core/domain/interface/chat"
	"github.com/ixre/go2o/core/inject"
)

func TestChat(t *testing.T) {
	repo := inject.GetChatRepo()
	sender := repo.GetChatUser(848)
	replayer := repo.GetChatUser(854)
	ic, err := sender.BuildConversation(replayer.GetAggregateRootId(), 0, "")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	rc, _ := replayer.BuildConversation(sender.GetAggregateRootId(), 0, "")
	if rc.GetDomainId() != ic.GetDomainId() {
		t.Error("conversation id not equal")
		t.FailNow()
	}
	t.Logf("welcome: #%d", ic.GetDomainId())
	ic.Greet("-> hello 你好! 我们已经是好友了")
	msgId, err := rc.Send(&chat.MsgBody{
		MsgType: 1,
		Content: "<- 你好，AI ROBOT",
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("<- send msg: #%d", msgId)
	msgId, err = ic.Send(&chat.MsgBody{
		MsgType: 1,
		Content: "-> I'am not a rebot!",
	})
	t.Logf("-> send msg: #%d", msgId)
	rc.Send(&chat.MsgBody{
		MsgType: 1,
		Content: "that's so cool",
	})
}
