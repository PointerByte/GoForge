package models

import "github.com/PointerByte/GoForge/encrypt/common"

type KeyData struct {
	PublicKey string
	KeyID     string
	KeyRef    string
	Provider  string
}

type RotateKeyRequest struct {
	UID   string
	KeyID string
}

type GetKeyRequest struct {
	UID   string
	KeyID string
}

type DeactivateKeyRequest struct {
	UID   string
	KeyID string
}

type GenerateSymmetricKeyRequest struct {
	UID  string
	Size common.SizeSymetrycKey
}

type EncryptAESRequest struct {
	UID        string
	SecretKey  string
	Value      string
	Additional *string
}

type DecryptAESRequest struct {
	UID         string
	SecretKey   string
	CipherValue string
	Additional  *string
}

type GenerateRSAKeyRequest struct {
	UID  string
	Size common.SizeAsymetrycKey
}

type GenerateECDHCurveKeyRequest struct {
	UID   string
	Curve common.CurveAsymmetricKey
}

type RSAOAEPEncodeRequest struct {
	UID       string
	PublicKey string
	Text      string
}

type RSAOAEPDecodeRequest struct {
	UID        string
	PrivateKey string
	CipherText string
}

type ECDHEncodeRequest struct {
	UID       string
	PublicKey string
	Text      string
}

type ECDHDecodeRequest struct {
	UID        string
	PrivateKey string
	CipherText string
}
