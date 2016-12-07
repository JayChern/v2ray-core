package encoding

import (
	"hash/fnv"

	"crypto/md5"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/serial"
)

func Authenticate(b []byte) uint32 {
	fnv1hash := fnv.New32a()
	fnv1hash.Write(b)
	return fnv1hash.Sum32()
}

type FnvAuthenticator struct {
}

func (v *FnvAuthenticator) NonceSize() int {
	return 0
}

func (v *FnvAuthenticator) Overhead() int {
	return 4
}

func (v *FnvAuthenticator) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	dst = serial.Uint32ToBytes(Authenticate(plaintext), dst[:0])
	return append(dst, plaintext...)
}

func (v *FnvAuthenticator) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if serial.BytesToUint32(ciphertext[:4]) != Authenticate(ciphertext[4:]) {
		return dst, crypto.ErrAuthenticationFailed
	}
	return append(dst[:0], ciphertext[4:]...), nil
}

func GenerateChacha20Poly1305Key(b []byte) []byte {
	key := make([]byte, 32)
	t := md5.Sum(b)
	copy(key, t[:])
	t = md5.Sum(key[:16])
	copy(key[16:], t[:])
	return key
}
