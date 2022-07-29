package stress

import (
	"context"
)

func RunStressTest(ctx context.Context, test string, iGoroutine GoroutineIndex) uint32 {
	stressFun, ok := tests[test]
	if !ok {
		flushAndPanic("illegal test name")
	}

	var i _RunIndex = 0
	for {
		select {
		case <-ctx.Done():
			return uint32(i)
		default:
			stressFun(iGoroutine, i)
			i++
		}
	}
}

func GetStressTestNames() []string {
	names := make([]string, 0, len(tests))
	for k, _ := range tests {
		names = append(names, k)
	}
	return names
}
