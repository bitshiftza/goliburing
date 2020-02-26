package goliburing

import (
	"os"
	"testing"
)

func TestNewRing(t *testing.T) {
	var queueDepth uint32 = 128
	ring, err := NewRing(queueDepth, nil)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := queueDepth, ring.QueueDepth(); want != have {
		t.Fatalf("queueDepth: want %d, have %d", want, have)
	}
}

func TestPrepWriteV(t *testing.T) {
	ring, err := NewRing(128, nil)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile("tmp", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	data := make([]byte, 256)
	sqe, err := ring.GetEmptySQE()
	if err != nil {
		t.Fatal(err)
	}
	sqe.PrepWriteV(int(f.Fd()), data, 0)
	ring.Submit()
	cqe, err := ring.WaitCQE()
	if err != nil {
		t.Fatal(err)
	}
	cqe.Seen()
	os.Remove(f.Name())
}

func BenchmarkPrepWriteV(b *testing.B) {
	ring, err := NewRing(128, nil)
	if err != nil {
		b.Error(err)
	}
	defer ring.Destroy()

	f, err := os.OpenFile("tmp-benchmark-prep-write-v", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		b.Error(err)
	}
	defer f.Close()

	data := make([]byte, 256)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sqe, err := ring.GetEmptySQE()
		if err != nil {
			b.Error(err)
			return
		}
		sqe.PrepWriteV(int(f.Fd()), data, 0)
		ring.Submit()
		cqe, err := ring.WaitCQE()
		if err != nil {
			b.Error(err)
			return
		}
		cqe.Seen()
	}
	b.StopTimer()
	os.Remove(f.Name())
}

func BenchmarkNormalWrite(b *testing.B) {
	f, err := os.OpenFile("/tmp/benchmark-normal-write", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		b.Error(err)
	}
	defer f.Close()

	data := make([]byte, 256)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = f.Write(data)
		if err != nil {
			b.Error(err)
		}
		// f.Sync()
	}
	b.StopTimer()
	os.Remove(f.Name())
}
