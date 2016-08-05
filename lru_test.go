package lru

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "testing"
)

var _ = Describe("LRU", func() {
  var lru *Cache
  const limit int = 1
  key, value := "foo", "bar"
  key2, value2 := "foo2", "bar2"

  Context("Basic", func() {
    BeforeEach(func() {
      lru = New(limit)
    })

    It("Can Set and Get value to LRU", func() {
      Expect(lru.Set(key, value)).Should(Equal(true))
      Expect(lru.Get(key)).Should(Equal(value))
    })

    It("Can RemoveElement from list by key", func() {
      lru.Set(key, value)
      Expect(lru.remove(key)).Should(BeTrue())
      Expect(lru.Get(key)).Should(BeNil())
    })

    It("Can check OverLimit", func() {
      Expect(lru.overLimit()).Should(BeFalse())
    })

    It("Should remove eldest element if over limit", func() {
      lru.Set(key, value)
      lru.Set(key2, value2)

      Expect(lru.overLimit()).Should(BeFalse())
      Expect(lru.storage.Len()).Should(Equal(limit))
      Expect(lru.Get(key2)).Should(Equal(value2))
    })
  })
  
  Context("Ordering", func() {
    BeforeEach(func() {
      lru = New(2)

      lru.Set(key, value)
      lru.Set(key2, value2)
    })
    
    It("Should add element to top of the list", func() {
      Expect(lru.storage.Front().Value.(*Entry).value).Should(Equal(value2))
    })

    It("Should move up element to top of the list when used", func() {
      lru.Get(key)
      Expect(lru.storage.Front().Value.(*Entry).value).Should(Equal(value))
    })
  })
})

func TestLRU(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "LRU")
}
