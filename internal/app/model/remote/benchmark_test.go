package remote

import (
	"testing"
)

func BenchmarkDetect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestDetect()
	}
}

func BenchmarkRecord(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TestRecord()
	}
}

func BenchmarkCapture(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Capture()
	}
}
