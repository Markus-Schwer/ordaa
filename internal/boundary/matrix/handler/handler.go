package handler

const (
	MatrixCommandPrefix      = ".ordaa"
	MatrixCommandPrefixRegex = "\\.ordaa"
)

type CommandResponse struct {
	Msg    string
	AsHTML bool
}
