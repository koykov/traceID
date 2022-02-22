package notifier

import "testing"

func TestTelegram(t *testing.T) {
	t.Run("parse addr", func(t *testing.T) {
		n := &Telegram{}
		n.SetConfig("chatID=-XXX; token=000:YYY; template=New ID [#{TID}](https://trace.com/{TID}).")
		if err := n.parseAddr(); err != nil {
			t.Error(err)
		}
		assertStr(t, "chatID", "-XXX", n.chatID)
		assertStr(t, "token", "000:YYY", n.token)
		assertStr(t, "template", "New ID [#{TID}](https://trace.com/{TID}).", n.template)
	})
}
