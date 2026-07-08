package shortener

type Link struct {
	OriginalURL string
	Code        string
}

func NewLink(originalURL string, codeLength int) (Link, error) {
	code, err := GenerateCode(codeLength)

	if err != nil {
		return Link{}, err
	}

	result := Link{
		OriginalURL: originalURL,
		Code:        code,
	}

	return result, nil
}
