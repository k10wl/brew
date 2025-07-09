package ports

import "context"

type QRCodeGenerator interface {
	GenerateQRCode(
		ctx context.Context,
		brewID string,
	) ([]byte, error)
	ParseQRCode(
		ctx context.Context,
		qrData []byte,
	) (string, error)
}
