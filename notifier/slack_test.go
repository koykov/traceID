package notifier

import "testing"

func TestSlack(t *testing.T) {
	assertStr := func(t testing.TB, n, a, b string) {
		if a != b {
			t.Errorf("%s mismatch: need '%s', got '%s'", n, a, b)
		}
	}
	t.Run("parse addr", func(t *testing.T) {
		n := Slack{}
		n.SetAddr("channel=XXX; username=YYY; https://hooks.slack.com/services/QWE/RTY/fake-token")
		channel, username, url, err := n.parseAddr()
		if err != nil {
			t.Log(err)
		}
		assertStr(t, "channel", "XXX", channel)
		assertStr(t, "username", "YYY", username)
		assertStr(t, "url", "https://hooks.slack.com/services/QWE/RTY/fake-token", url)
	})
}
