package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateVerificationHTML(ctx context.Context, token string, endpoint string, message string) string {

	verificationURL := fmt.Sprintf("ordindia.foundation/%s?token=%s", endpoint, token)
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Verify your email</title>
		</head>
		<body>
			<h2>%s</h2>
			<p>Please click the link below to verify your email address:</p>
			<p><a href="%s" style="color: #3498db; text-decoration: none;">Verify Email</a></p>
			<p>This link is valid for 15 minutes.</p>
			<p>If you did not request this, you can ignore this email.</p>
		</body>
		</html>
`, message, verificationURL)

	return htmlBody
}

// Generate verification token
func GenerateVerificationtoken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	// Encode the byte slice to a URL-safe base64 string
	return base64.RawURLEncoding.EncodeToString(b), nil
}
