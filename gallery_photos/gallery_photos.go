package gallery_photos

import (
	"errors"

	"git.dsrt-int.net/actionmc/actionmc-site-go/authdatabase"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
)

type GalleryPhoto struct {
	Name             string
	Description      string
	ImageURL         string
	HasCredit        bool
	CreditedAuthor   string
	PublishTimestamp int64
}

type GalleryPhotoHandler struct {
	photoMap   map[int]GalleryPhoto
	ImageCount int
	logger     *logging.Logger
}

func (GPH *GalleryPhotoHandler) Get(id int) (GalleryPhoto, error) {
	data, ok := GPH.photoMap[id]

	if !ok {
		return *new(GalleryPhoto), errors.New("Unknown image.")
	}

	return data, nil
}

func NewGalleryPhotos(db *authdatabase.MCAuthDB_sqlite3) *GalleryPhotoHandler {

	handler := GalleryPhotoHandler{
		photoMap: map[int]GalleryPhoto{},
		logger:   logging.New(),
	}

	images := db.GetPhotos()
	defer images.Close()

	for images.Next() {
		var imageID int
		var name string
		var description string
		var imageURL string
		var hasCredit int
		var creditedAuthor interface{}
		var publishedTimestamp int64

		var finalAuthor string

		err := images.Scan(&imageID, &name, &description, &imageURL, &hasCredit, &creditedAuthor, &publishedTimestamp)

		switch creditedAuthor.(type) {
		case string:
			finalAuthor = creditedAuthor.(string)

		case nil:
			finalAuthor = ""

		default:
			finalAuthor = ""
		}
		if err != nil {
			// this is not good. throw a fatal error.
			handler.logger.Fatal.Fatalln(err)
		}

		handler.photoMap[imageID] = GalleryPhoto{
			Name:             name,
			Description:      description,
			ImageURL:         imageURL,
			HasCredit:        (hasCredit == 1),
			CreditedAuthor:   finalAuthor,
			PublishTimestamp: publishedTimestamp,
		}
	}

	handler.ImageCount = len(handler.photoMap)
	return &handler
}
