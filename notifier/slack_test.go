package notifier

import "testing"

func TestSlack(t *testing.T) {
	t.Run("parse addr", func(t *testing.T) {
		n := Slack{}
		n.SetAddr("channel=XXX; username=YYY; template=New traceID <https://trace.com/{TID}|#{TID}>.; https://hooks.slack.com/services/QWE/RTY/fake-token")
		channel, username, template, url, err := n.parseAddr()
		if err != nil {
			t.Log(err)
		}
		assertStr(t, "channel", "XXX", channel)
		assertStr(t, "username", "YYY", username)
		assertStr(t, "template", "New traceID <https://trace.com/{TID}|#{TID}>.", template)
		assertStr(t, "url", "https://hooks.slack.com/services/QWE/RTY/fake-token", url)
	})
}
