package services

import (
	"context"

	"brew/internal/core/ports"
	"brew/internal/util"
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
	util.Debug("Generating QR code", "brew_id", brewID)
	qrData, err := s.qrGenerator.GenerateQRCode(ctx, brewID)
	if err != nil {
		return nil, err
	}
	util.Debug("QR code generated successfully", "brew_id", brewID, "data_size", len(qrData))
	return qrData, nil
}

func (s *QRService) ParseQRCode(
	ctx context.Context,
	qrData []byte,
) (string, error) {
	util.Debug("Parsing QR code", "data_size", len(qrData))
	result, err := s.qrGenerator.ParseQRCode(ctx, qrData)
	if err != nil {
		return "", err
	}
	util.Debug("QR code parsed successfully", "result", result, "data_size", len(qrData))
	return result, nil
}
