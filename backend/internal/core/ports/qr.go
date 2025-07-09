package ports

import "context"

type QRCodeGenerator interface {
	GenerateQRCode(
		ctx context.Context,
		jarID string,
	) ([]byte, error)
	ParseQRCode(
		ctx context.Context,
		qrData []byte,
	) (string, error)
}
