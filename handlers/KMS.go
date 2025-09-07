package handlers

import (
	"context"
	"encoding/base64"
	"fmt"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

// 暗号化
func encryptWithKMS(ctx context.Context, keyName, plaintext string) (string, error) {
	fmt.Println("plaintext:", plaintext)
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		fmt.Println("KMS client error:", err)
		return "", err
	}
	defer client.Close()

	req := &kmspb.EncryptRequest{
		Name:      keyName,
		Plaintext: []byte(plaintext),
	}
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		fmt.Println("KMS Encrypt error:", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(resp.Ciphertext), nil
}

// 復号
func decryptWithKMS(ctx context.Context, keyName, ciphertext string) (string, error) {
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	decoded, _ := base64.StdEncoding.DecodeString(ciphertext)

	req := &kmspb.DecryptRequest{
		Name:       keyName,
		Ciphertext: decoded,
	}
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return "", err
	}
	return string(resp.Plaintext), nil
}
