package qr

import (
	"context"

	"brew/internal/core/ports"
)

type Service struct {
	qrGenerator ports.QRCodeGenerator
}

func NewService(qrGenerator ports.QRCodeGenerator) *Service {
	return &Service{
		qrGenerator: qrGenerator,
	}
}

func (s *Service) GenerateQRCode(
	ctx context.Context,
	jarID string,
) ([]byte, error) {
	return s.qrGenerator.GenerateQRCode(ctx, jarID)
}

func (s *Service) ParseQRCode(
	ctx context.Context,
	qrData []byte,
) (string, error) {
	return s.qrGenerator.ParseQRCode(ctx, qrData)
}
