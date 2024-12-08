package utils

import (
	"ORDI/internal/messages"
)

func CreatePatientSignupStepsMessage() messages.SignupStepsMessage {
	signupMessage := messages.SignupStepsMessage{}
	signupMessage.AddStep(messages.SignupStep{
		Heading: "Build Your Account",
		Body:    "Begin by setting up your account on our platform. Share essential details like your medical history and contact information, so we can better understand and address your needs.",
	})
	signupMessage.AddStep(messages.SignupStep{
		Heading: "Verify Your Account",
		Body:    "After entering your details, verify your account by confirming your email, phone number, and any necessary documents. This ensures your information is secure and accurate.",
	})
	signupMessage.AddStep(messages.SignupStep{
		Heading: "Connect with the Right Doctors and Receive Ongoing Support",
		Body:    "Once verified, we’ll connect you with the right doctors specializing in your condition. We’ll be with you throughout your journey, providing ongoing support and guidance as needed.",
	})
	return signupMessage
}

func CreateDoctorSignupStepsMessage() messages.SignupStepsMessage {
	signupMessage := messages.SignupStepsMessage{}
	signupMessage.AddStep(messages.SignupStep{
		Heading: "Build Your Account",
		Body:    "Begin by setting up your account on our platform. Share essential details like your medical history and contact information, so we can better understand and address your needs.",
	})
	signupMessage.AddStep(messages.SignupStep{
		Heading: "Verify Your Account",
		Body:    "After entering your details, verify your account by confirming your email, phone number, and any necessary documents. This ensures your information is secure and accurate.",
	})
	signupMessage.AddStep(messages.SignupStep{
		Heading: "Help us fight rare Diseases in India",
		Body:    "Once verified, we’ll connect you with the right doctors specializing in your condition. We’ll be with you throughout your journey, providing ongoing support and guidance as needed.",
	})
	return signupMessage
}
