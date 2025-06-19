package ports

type SecureTokenGenerator interface {
	Generate() (string, error)
}
