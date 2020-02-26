package goliburing

import (
	"os"
	"testing"
)

func TestFileWriterAasync(t *testing.T) {
	ring, err := NewRing(128, nil)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile("/tmp/file-writer-async-test", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	writer, err := NewFileWriterAsync(ring, f)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 256)
	err = writer.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	bytesWritten, err := writer.WaitForCompletion()
	if err != nil {
		t.Fatal(err)
	}

	if want, have := len(data), bytesWritten; want != have {
		t.Fatalf("Write: want %d, have %d", want, have)
	}
}

func BenchmarkFileWriterAasync(b *testing.B) {
	ring, err := NewRing(128, nil)
	if err != nil {
		b.Fatal(err)
	}

	f, err := os.OpenFile("/tmp/file-writer-async-test", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	writer, err := NewFileWriterAsync(ring, f)
	if err != nil {
		b.Fatal(err)
	}

	data := make([]byte, 256)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := writer.Write(data)
		if err != nil {
			b.Fatal(err)
		}
		_, err = writer.WaitForCompletion()
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	os.Remove(f.Name())
}
