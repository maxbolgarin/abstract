package abstract_test

import (
	"testing"

	"github.com/maxbolgarin/abstract"
)

func TestLinkedList(t *testing.T) {
	t.Run("PopBack", func(t *testing.T) {
		l := abstract.NewLinkedList[int]()

		for i := 1; i <= 10; i++ {
			l.InsertBack(i)
		}
		if l.Len() != 10 {
			t.Fatalf("expected 10, got %d", l.Len())
		}
		front, _ := l.Front()
		if front != 1 {
			t.Fatalf("expected 1, got %d", front)
		}
		back, _ := l.Back()
		if back != 10 {
			t.Fatalf("expected 10, got %d", back)
		}

		for i := 0; i < 10; i++ {
			out, ok := l.PopBack()
			if !ok {
				t.Fatal("expected to be ok")
			}
			if out != 10-i {
				t.Fatalf("expected %d, got %d", 10-i, out)
			}
		}

		if l.Len() != 0 {
			t.Fatalf("expected 0, got %d", l.Len())
		}

		_, ok := l.PopFront()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.PopBack()
		if ok {
			t.Fatalf("expected false, got true")
		}

		_, ok = l.Front()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.Back()
		if ok {
			t.Fatalf("expected false, got true")
		}

	})

	t.Run("PopFront", func(t *testing.T) {
		l := abstract.NewLinkedList[int]()

		for i := 1; i <= 10; i++ {
			l.InsertFront(i)
		}
		if l.Len() != 10 {
			t.Fatalf("expected 10, got %d", l.Len())
		}

		front, _ := l.Front()
		if front != 10 {
			t.Fatalf("expected 10, got %d", front)
		}
		back, _ := l.Back()
		if back != 1 {
			t.Fatalf("expected 1, got %d", back)
		}

		for i := 10; i > 0; i-- {
			out, ok := l.PopFront()
			if !ok {
				t.Fatal("expected to be ok")
			}
			if out != i {
				t.Fatalf("expected %d, got %d", i, out)
			}
		}

		if l.Len() != 0 {
			t.Fatalf("expected 0, got %d", l.Len())
		}

		_, ok := l.PopFront()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.PopBack()
		if ok {
			t.Fatalf("expected false, got true")
		}

		_, ok = l.Front()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.Back()
		if ok {
			t.Fatalf("expected false, got true")
		}
	})
}

func TestSafeLinkedList(t *testing.T) {
	t.Run("PopBack", func(t *testing.T) {
		l := abstract.NewSafeLinkedList[int]()

		for i := 1; i <= 10; i++ {
			l.InsertBack(i)
		}
		if l.Len() != 10 {
			t.Fatalf("expected 10, got %d", l.Len())
		}
		front, _ := l.Front()
		if front != 1 {
			t.Fatalf("expected 1, got %d", front)
		}
		back, _ := l.Back()
		if back != 10 {
			t.Fatalf("expected 10, got %d", back)
		}

		for i := 0; i < 10; i++ {
			out, ok := l.PopBack()
			if !ok {
				t.Fatal("expected to be ok")
			}
			if out != 10-i {
				t.Fatalf("expected %d, got %d", 10-i, out)
			}
		}

		if l.Len() != 0 {
			t.Fatalf("expected 0, got %d", l.Len())
		}

		_, ok := l.PopFront()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.PopBack()
		if ok {
			t.Fatalf("expected false, got true")
		}

		_, ok = l.Front()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.Back()
		if ok {
			t.Fatalf("expected false, got true")
		}

	})

	t.Run("PopFront", func(t *testing.T) {
		l := abstract.NewSafeLinkedList[int]()

		for i := 1; i <= 10; i++ {
			l.InsertFront(i)
		}
		if l.Len() != 10 {
			t.Fatalf("expected 10, got %d", l.Len())
		}

		front, _ := l.Front()
		if front != 10 {
			t.Fatalf("expected 10, got %d", front)
		}
		back, _ := l.Back()
		if back != 1 {
			t.Fatalf("expected 1, got %d", back)
		}

		for i := 10; i > 0; i-- {
			out, ok := l.PopFront()
			if !ok {
				t.Fatal("expected to be ok")
			}
			if out != i {
				t.Fatalf("expected %d, got %d", i, out)
			}
		}

		if l.Len() != 0 {
			t.Fatalf("expected 0, got %d", l.Len())
		}

		_, ok := l.PopFront()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.PopBack()
		if ok {
			t.Fatalf("expected false, got true")
		}

		_, ok = l.Front()
		if ok {
			t.Fatalf("expected false, got true")
		}
		_, ok = l.Back()
		if ok {
			t.Fatalf("expected false, got true")
		}
	})
}
