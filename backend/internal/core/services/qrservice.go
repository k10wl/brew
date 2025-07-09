package services

import (
	"context"

	"brew/internal/core/ports"
)

type QRService struct {
	qrGenerator ports.QRCodeGenerator
}

func NewQRService(qrGenerator ports.QRCodeGenerator) *QRService {
	return &QRService{
		qrGenerator: qrGenerator,
	}
}

func (s *QRService) GenerateQRCode(
	ctx context.Context,
	brewID string,
) ([]byte, error) {
	return s.qrGenerator.GenerateQRCode(ctx, brewID)
}

func (s *QRService) ParseQRCode(
	ctx context.Context,
	qrData []byte,
) (string, error) {
	return s.qrGenerator.ParseQRCode(ctx, qrData)
}
