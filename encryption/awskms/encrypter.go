package awskms

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/vaulty/vaulty/config"
	"github.com/vaulty/vaulty/encrypt"
	"github.com/vaulty/vaulty/encryption"
)

var _ encryption.Encrypter = (*AwsKms)(nil)

func Factory(conf *config.Config) (encryption.Encrypter, error) {
	keyID := conf.Encryption.AWSKMSKeyID
	if keyID == "" {
		keyID = "alias/" + conf.Encryption.AWSKMSKeyAlias
	}

	params := &Params{
		Region: conf.Encryption.AWSKMSRegion,
		KeyID:  keyID,
	}

	return NewEncrypter(params)
}

type Params struct {
	KeyID  string
	Region string
}

type AwsKms struct {
	client *kms.KMS
	keyID  string
}

func NewEncrypter(params *Params) (*AwsKms, error) {
	if params.Region == "" {
		return nil, errors.New("AWS KMS Region is not confured")
	}

	if params.KeyID == "" {
		return nil, errors.New("AWS KMS Key ID is not confured")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: &params.Region,
	})
	if err != nil {
		return nil, err
	}

	enc := &AwsKms{
		client: kms.New(sess),
		keyID:  params.KeyID,
	}

	return enc, nil
}

type payload struct {
	EncryptedKey []byte
	Data         []byte
}

func (e *AwsKms) Encrypt(plaintext []byte) ([]byte, error) {
	// API call to AWS KMS
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

	buf := &bytes.Buffer{}
	// link together encrypted encryption key and ciphertext
	err = gob.NewEncoder(buf).Encode(&payload{
		EncryptedKey: output.CiphertextBlob,
		Data:         ciphertext,
	})
	if err != nil {
		return nil, err
	}

	// convert it to hex to use in http
	hexPayload := make([]byte, hex.EncodedLen(buf.Len()))
	hex.Encode(hexPayload, buf.Bytes())

	return hexPayload, nil
}

func (e *AwsKms) Decrypt(message []byte) ([]byte, error) {
	p := &payload{}

	plainPayload := make([]byte, hex.DecodedLen(len(message)))
	hex.Decode(plainPayload, message)

	if err := gob.NewDecoder(bytes.NewReader(plainPayload)).Decode(p); err != nil {
		return nil, err
	}

	// decrypt encrypted key (API call to AWS KMS)
	output, err := e.client.Decrypt(&kms.DecryptInput{
		CiphertextBlob: p.EncryptedKey,
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
