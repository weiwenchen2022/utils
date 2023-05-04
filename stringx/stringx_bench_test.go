package stringx_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"
	"unsafe"

	. "github.com/weiwenchen2022/utils/stringx"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytes(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	return string(b)
}

func RandStringBytesMask(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}

	return string(b)
}

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)

	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func RandStringBytesMaskImprSrcSB(n int) string {
	var b strings.Builder
	b.Grow(n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b.WriteByte(letterBytes[idx])
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return b.String()
}

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)

	// A src.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}

		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func Benchmark(b *testing.B) {
	const n = 16

	b.Run("Runes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringRunes(n)
		}
	})

	b.Run("Bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytes(n)
		}
	})

	b.Run("Rmndr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytesRmndr(n)
		}
	})

	b.Run("BytesMask", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytesMask(n)
		}
	})

	b.Run("BytesMaskImpr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytesMaskImpr(n)
		}
	})

	b.Run("BytesMaskImprSrc", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytesMaskImprSrc(n)
		}
	})

	b.Run("BytesMaskImprSrcSB", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytesMaskImprSrcSB(n)
		}
	})

	b.Run("BytesMaskImprSrcUnsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandStringBytesMaskImprSrcUnsafe(n)
		}
	})

	b.Run("BytesMaskImprUnsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = RandString(n)
		}
	})
}
