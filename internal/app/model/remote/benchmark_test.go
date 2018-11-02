package remote

import (
	"testing"
)

func BenchmarkDetect(b *testing.B) {
	TestDetect()
}

func BenchmarkRecord(b *testing.B) {
	TestRecord()
}

func BenchmarkCapture(b *testing.B) {
	Capture()
}
