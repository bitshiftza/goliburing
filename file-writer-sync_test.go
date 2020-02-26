package goliburing

import (
	"os"
	"testing"
)

func TestFileWriterSync(t *testing.T) {
	ring, err := NewRing(128, nil)
	if err != nil {
		t.Fatal(err)
	}

	f, err := os.OpenFile("/tmp/file-writer-sync-test", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	writer, err := NewFileWriterSync(ring, f)
	if err != nil {
		t.Fatal(err)
	}

	data := make([]byte, 256)
	bytesWritten, err := writer.Write(data)
	if err != nil {
		t.Fatal(err)
	}

	if want, have := len(data), bytesWritten; want != have {
		t.Fatalf("Write: want %d, have %d", want, have)
	}
}

func BenchmarkFileWriterSync(b *testing.B) {
	ring, err := NewRing(128, nil)
	if err != nil {
		b.Fatal(err)
	}

	f, err := os.OpenFile("tmp-file-writer-sync-test", os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	writer, err := NewFileWriterSync(ring, f)
	if err != nil {
		b.Fatal(err)
	}

	data := make([]byte, 256)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := writer.Write(data)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	os.Remove(f.Name())
}
