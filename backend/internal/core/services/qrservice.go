package services

import (
	"context"
	"log/slog"

	"brew/internal/core/ports"
	_ "brew/internal/util"
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
	slog.Debug("Generating QR code", "brew_id", brewID)
	qrData, err := s.qrGenerator.GenerateQRCode(ctx, brewID)
	if err != nil {
		return nil, err
	}
	slog.Debug("QR code generated successfully", "brew_id", brewID, "data_size", len(qrData))
	return qrData, nil
}

func (s *QRService) ParseQRCode(
	ctx context.Context,
	qrData []byte,
) (string, error) {
	slog.Debug("Parsing QR code", "data_size", len(qrData))
	result, err := s.qrGenerator.ParseQRCode(ctx, qrData)
	if err != nil {
		return "", err
	}
	slog.Debug("QR code parsed successfully", "result", result, "data_size", len(qrData))
	return result, nil
}
