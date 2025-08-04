package infrastructure

type EmailInfrastructure struct {}

func NewEmailInfrastructure() *EmailInfrastructure {
	return &EmailInfrastructure{}
}

// TODO : Implement the Email service
func (ei *EmailInfrastructure) SendEmail(from string, to []string, content string) error {
	return nil
}