package domain

type ASCIIText struct {
	Text string
}

type ASCIITextRequest struct {
	Text   string `json:"text"`
	Banner string `json:"banner"`
}

type AsciiTextResponse struct {
	AsciiArt string `json:"ascii,omitempty"`
	Error    string `json:"error,omitempty"`
}
