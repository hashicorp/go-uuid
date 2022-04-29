package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// GenerateRandomBytes is used to generate random bytes of given size.
func GenerateRandomBytes(size int) ([]byte, error) {
	return GenerateRandomBytesWithReader(size, rand.Reader)
}

// GenerateRandomBytesWithReader is used to generate random bytes of given size read from a given reader.
func GenerateRandomBytesWithReader(size int, reader io.Reader) ([]byte, error) {
	if reader == nil {
		return nil, fmt.Errorf("provided reader is nil")
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}
	return buf, nil
}


const uuidLen = 16

// GenerateUUID is used to generate a random UUID
func GenerateUUID() (string, error) {
	return GenerateUUIDWithReader(rand.Reader)
}

// GenerateSecureUUID panics in the event reading from the secure crypto/rand
// reader fails. This is because entropy errors usually aren't userspace
// recoverable; we can potentially wait before re-reading from /dev/urandom, but
// once initialized, it should not cause an error.
//
// See: https://github.com/golang/go/blob/master/src/crypto/rand/rand_unix.go#L46-L47
//
// Most call sites of GenerateUUID in Vault do one of two things:
//
//  1. Bubble the error up to their parent to handle.
//  2. Ignore the error and hard-code a resulting UUID.
//
// In both cases, rather than making the caller think about how to handle
// entropy request errors (the only possible case as crypto/rand.Reader
// should never be nil), panic()'ing is a cleaner and safer approach as
// the entropy source is safer and most callers should not continue to
// operate until this has been addressed.
//
// It also has the side-effect of returning only the requested information,
// simplifying the call signature.
func GenerateSecureUUID() string {
	ret, err := GenerateUUID()
	if err != nil {
		panic("crypto/rand returned fatal error: " + err.Error())
	}

	return ret
}

// GenerateUUIDWithReader is used to generate a random UUID with a given Reader
func GenerateUUIDWithReader(reader io.Reader) (string, error) {
	if reader == nil {
		return "", fmt.Errorf("provided reader is nil")
	}
	buf, err := GenerateRandomBytesWithReader(uuidLen, reader)
	if err != nil {
		return "", err
	}
	return FormatUUID(buf)
}

func FormatUUID(buf []byte) (string, error) {
	if buflen := len(buf); buflen != uuidLen {
		return "", fmt.Errorf("wrong length byte slice (%d)", buflen)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16]), nil
}

func ParseUUID(uuid string) ([]byte, error) {
	if len(uuid) != 2 * uuidLen + 4 {
		return nil, fmt.Errorf("uuid string is wrong length")
	}

	if uuid[8] != '-' ||
		uuid[13] != '-' ||
		uuid[18] != '-' ||
		uuid[23] != '-' {
		return nil, fmt.Errorf("uuid is improperly formatted")
	}

	hexStr := uuid[0:8] + uuid[9:13] + uuid[14:18] + uuid[19:23] + uuid[24:36]

	ret, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	if len(ret) != uuidLen {
		return nil, fmt.Errorf("decoded hex is the wrong length")
	}

	return ret, nil
}
