package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/repository"
)

type AssetService struct {
	repo *repository.AssetRepo
	gcm  cipher.AEAD
}

func NewAssetService(repo *repository.AssetRepo, encryptionKey string) *AssetService {
	key := []byte(encryptionKey)
	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		key = padded
	}
	block, err := aes.NewCipher(key[:32])
	if err != nil {
		panic(fmt.Sprintf("failed to create AES cipher: %v", err))
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("failed to create GCM: %v", err))
	}
	return &AssetService{repo: repo, gcm: gcm}
}

func (s *AssetService) encrypt(plaintext string) (string, error) {
	nonce := make([]byte, s.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := s.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *AssetService) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	nonceSize := s.gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := s.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func (s *AssetService) List(search string) ([]model.Asset, error) {
	return s.repo.FindAll(search)
}

func (s *AssetService) GetByID(id int64) (*model.Asset, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// Return asset without exposing credential
	a.Credential = ""
	return a, nil
}

func (s *AssetService) GetByIDWithCredential(id int64) (*model.Asset, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	cred, err := s.Decrypt(a.Credential)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt credential: %w", err)
	}
	a.Credential = cred
	return a, nil
}

func (s *AssetService) Create(name, host string, port int, username, authType, credential, status string) (*model.Asset, error) {
	encCred, err := s.encrypt(credential)
	if err != nil {
		return nil, err
	}
	a := &model.Asset{
		Name:       name,
		Host:       host,
		Port:       port,
		Username:   username,
		AuthType:   authType,
		Credential: encCred,
		Status:     status,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	a.Credential = ""
	return a, nil
}

func (s *AssetService) Update(id int64, name, host string, port int, username, authType, credential, status string) (*model.Asset, error) {
	a, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	a.Name = name
	a.Host = host
	a.Port = port
	a.Username = username
	a.AuthType = authType
	if credential != "" {
		encCred, err := s.encrypt(credential)
		if err != nil {
			return nil, err
		}
		a.Credential = encCred
	}
	a.Status = status
	a.UpdatedAt = time.Now()
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	a.Credential = ""
	return a, nil
}

func (s *AssetService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *AssetService) TestConnect(id int64) error {
	a, err := s.GetByIDWithCredential(id)
	if err != nil {
		return fmt.Errorf("asset not found: %w", err)
	}
	client, err := s.dial(a)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	client.Close()
	return nil
}

func (s *AssetService) dial(a *model.Asset) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User:            a.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}
	switch a.AuthType {
	case model.AuthTypePassword:
		config.Auth = []ssh.AuthMethod{ssh.Password(a.Credential)}
	case model.AuthTypeKey:
		signer, err := ssh.ParsePrivateKey([]byte(a.Credential))
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %w", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", a.AuthType)
	}
	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", a.Host, a.Port), config)
}
