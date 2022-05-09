package bot

type FullResponse struct {
	Blocks []string
}

func (fr FullResponse) ToString() string {
	var fullResponse string
	for _, block := range fr.Blocks {
		fullResponse += block + "\n"
	}

	return fullResponse
}
