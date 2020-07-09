package awskms

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/vaulty/vaulty"
	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/encryption"
)

type Params struct {
	KeyID  string
	Region string
}

func Factory(config *vaulty.Config) *encryption.Encrypter {
	return nil
}

type payload struct {
	Key  []byte
	Data []byte
}

type encrypter struct {
	client *kms.KMS
	keyID  string
}

func NewEncrypter(params *Params) (*encrypter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: &params.Region,
	})
	if err != nil {
		return nil, err
	}

	enc := &encrypter{
		client: kms.New(sess),
		keyID:  params.KeyID,
	}

	return enc, nil
}

func (e *encrypter) Encrypt(plaintext []byte) ([]byte, error) {
	output, err := e.client.GenerateDataKey(&kms.GenerateDataKeyInput{
		KeyId:   &e.keyID,
		KeySpec: aws.String("AES_256"),
	})

	enc, err := encrypt.NewAesGcm(output.Plaintext)
	if err != nil {
		return nil, err
	}

	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}

	// encode encrypted key together with ciphertext
	buf := &bytes.Buffer{}
	err = gob.NewEncoder(buf).Encode(&payload{
		Key:  output.CiphertextBlob,
		Data: ciphertext,
	})
	if err != nil {
		return nil, err
	}

	// convert it to hex to use in http
	hexPayload := make([]byte, hex.EncodedLen(buf.Len()))
	hex.Encode(hexPayload, buf.Bytes())

	return hexPayload, nil
}

func (e *encrypter) Decrypt(message []byte) ([]byte, error) {
	p := &payload{}

	plainPayload := make([]byte, hex.DecodedLen(len(message)))
	hex.Decode(plainPayload, message)

	if err := gob.NewDecoder(bytes.NewReader(plainPayload)).Decode(p); err != nil {
		return nil, err
	}

	output, err := e.client.Decrypt(&kms.DecryptInput{
		CiphertextBlob: p.Key,
	})

	enc, err := encrypt.NewAesGcm(output.Plaintext)
	if err != nil {
		return nil, err
	}

	plain, err := enc.Decrypt(p.Data)
	if err != nil {
		return nil, err
	}

	return plain, nil
}
