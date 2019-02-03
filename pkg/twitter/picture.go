package twitter

import (
	"github.com/processone/dpk/pkg/semweb"
)

// Picture is a dedicated Semantic Web client, specialized in interacting with Twitter Twitter pictures.
type Picture struct {
	semweb.Client
}

func NewPicture() *Picture {
	var pic Picture
	pic.Client = semweb.NewClient()
	return &pic
}

// TODO(mr): Naming is akward, fixme.
func (p *Picture) GetImageURL(imageUrl string) (string, error) {
	// Discover image URL
	body, _, err := p.Get(imageUrl)
	if err != nil {
		return imageUrl, err
	}
	defer body.Close()

	page, err := semweb.ReadPage(body)
	if err != nil {
		return imageUrl, err
	}

	img := page.Properties["og:image"]
	return img, nil
}
