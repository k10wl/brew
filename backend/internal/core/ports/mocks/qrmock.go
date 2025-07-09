package mocks

import (
	"context"

	"brew/internal/core/ports"
)

var _ ports.QRCodeGenerator = (*QRCodeGenerator)(nil)

type QRCodeGenerator struct {
	GenerateQRCodeFunc func(ctx context.Context, jarID string) ([]byte, error)
	ParseQRCodeFunc    func(ctx context.Context, qrData []byte) (string, error)
}

func (m *QRCodeGenerator) GenerateQRCode(ctx context.Context, jarID string) ([]byte, error) {
	if m.GenerateQRCodeFunc != nil {
		return m.GenerateQRCodeFunc(ctx, jarID)
	}
	return nil, nil
}

func (m *QRCodeGenerator) ParseQRCode(ctx context.Context, qrData []byte) (string, error) {
	if m.ParseQRCodeFunc != nil {
		return m.ParseQRCodeFunc(ctx, qrData)
	}
	return "", nil
}
