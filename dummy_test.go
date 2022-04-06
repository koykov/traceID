package traceID

import (
	"math"
	"testing"
	"time"
)

type testDEQInterface interface {
	DeepEqual(interface{}, interface{}) bool
}

type testDEQ struct {
}

func (d testDEQ) DeepEqual(_, _ interface{}) bool { return true }

func BenchmarkDummy(b *testing.B) {
	b.Run("alloc-free", func(b *testing.B) {
		optDEQ := Option("deq")
		c := DummyClock{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var opt testDEQInterface
			opt = &testDEQ{}
			dc := DummyCtx{}
			dc.SetID("foobar").
				SetService("server").
				SetFlag(FlagOverwrite, true).
				Watch(LevelAll).
				SetBroadcastTimeout(time.Second).
				SetClock(c).
				Debug("trace0").Var("v0", "foo").With(OptionIndent, true)
			dc.Info("trace1").
				Var("v1", math.MaxInt32).With(optDEQ, opt).
				Var("v2", 256.56)
			dc.Warn("trace2").
				Slug("t2").
				Var("v3", math.MaxFloat64).
				Comment("something went strange").
				Var("v4", -15)
		}
	})
}
