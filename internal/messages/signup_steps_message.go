package messages

type SignupStepsMessage struct {
	Steps []SignupStep
}

type SignupStep struct {
	Heading string
	Body    string
}

func (s *SignupStepsMessage) AddStep(step SignupStep) {
	s.Steps = append(s.Steps, step)
}
