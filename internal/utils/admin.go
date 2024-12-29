package utils

import (
	"context"
	"fmt"
)

func GenerateWelcomeAdminHTML(ctx context.Context, endpoint string, email string) string {

	welcomeURL := fmt.Sprintf("internal.ordindia.foundation%s?email=%s", endpoint, email)
	htmlBody := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Onboard to ORDI</title>
		</head>
		<body>
			<h2>Welcome to ORDI team</h2>
			<p>Hi, Hope this finds you in good health. You've been added to the ORDI team.</p>
			<p>The email includes the details that have been used to onboard you on the team.</p>
			<p>Please click on the link below to set up your credentials and onboard:</p>
			<p><a href="%s">Set up your credentials</a></p>
			<p>If you are not aware of this, you can ignore this email.</p>
		</body>
		</html>
	`, welcomeURL)

	return htmlBody
}
