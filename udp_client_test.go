package common

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func aTestMain(m *testing.M) {
	InitUdpClient("9067aadfbee458e0d2e1f9876e898882aab5f5876f617634e7e84ec9be5bded0")
	os.Exit(m.Run())
}

func aTest_Dingding(t *testing.T) {
	for i := 0; i < 1; i += 1 {
		SendDingDing(fmt.Sprintf("log:%d", i), fmt.Sprintf("message_seq:%d", i))
		time.Sleep(time.Second)
	}
}
