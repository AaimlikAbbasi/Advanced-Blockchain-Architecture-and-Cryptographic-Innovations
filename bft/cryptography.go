package bft

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Generate RSA key pair for cryptographic verification
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	pub := &priv.PublicKey
	return priv, pub, nil
}

// Sign a message using RSA
func SignMessage(priv *rsa.PrivateKey, message []byte) ([]byte, error) {
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Verify RSA signature
func VerifySignature(pub *rsa.PublicKey, message, signature []byte) error {
	hashed := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], signature)
}
