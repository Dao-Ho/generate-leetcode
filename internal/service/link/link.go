package link

type LinkService struct {
}

func NewLinkService() *LinkService {
	return &LinkService{}
}

func (s *LinkService) GetLink(messageText string) (string, error) {
	// hehehehe
	return "Here's your link: https://www.youtube.com/watch?v=dQw4w9WgXcQ&list=RDdQw4w9WgXcQ&start_radio=1", nil
}
