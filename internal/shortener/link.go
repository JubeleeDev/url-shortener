package shortener

type Link struct {
	OriginalURL string
	Code        string
}

func NewLink(originalURL string, codeLenght int) (Link, error) {
	code, err := GenerateCode(codeLenght)

	if err != nil {
		return Link{}, err
	}

	result := Link{
		OriginalURL: originalURL,
		Code:        code,
	}

	return result, nil
}
