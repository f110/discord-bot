package amesh

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"time"
)

type Generator struct {
	mapImage  image.Image
	maskImage image.Image
}

func NewGenerator() (*Generator, error) {
	mapImage, err := getImage("https://tokyo-ame.jwa.or.jp/map/map100.jpg")
	if err != nil {
		return nil, err
	}
	maskImage, err := getImage("https://tokyo-ame.jwa.or.jp/map/msk100.png")
	if err != nil {
		return nil, err
	}

	return &Generator{
		mapImage:  mapImage,
		maskImage: maskImage,
	}, nil
}

func (g *Generator) LatestTime() string {
	return time.Now().Add(time.Duration(-1) * time.Minute).Truncate(5 * time.Minute).Format("200601021504")
}

func (g *Generator) Latest() (string, *bytes.Buffer, error) {
	date := g.LatestTime()
	img, err := g.Generate(date)
	return date, img, err
}

func (g *Generator) Generate(date string) (*bytes.Buffer, error) {
	rainImage, err := getImage("http://tokyo-ame.jwa.or.jp/mesh/100/" + date + ".gif")
	if err != nil {
		return nil, err
	}

	imageSize := g.mapImage.Bounds()
	sp := image.Pt(0, 0)
	base := image.NewRGBA(imageSize)
	draw.Draw(base, imageSize, g.mapImage, sp, draw.Over)
	draw.Draw(base, imageSize, rainImage, sp, draw.Over)
	draw.Draw(base, imageSize, g.maskImage, sp, draw.Over)

	buf := new(bytes.Buffer)
	crop := image.Rect(900, 650, 2200, 1250)
	if err := png.Encode(buf, base.SubImage(crop)); err != nil {
		return nil, err
	}

	return buf, nil
}

func getImage(url string) (image.Image, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var img image.Image
	switch res.Header.Get("Content-Type") {
	case "image/gif":
		img, err = gif.Decode(res.Body)
	case "image/jpeg":
		img, err = jpeg.Decode(res.Body)
	case "image/png":
		img, err = png.Decode(res.Body)
	default:
		return nil, fmt.Errorf("unknown image type: %s", res.Header.Get("Content-Type"))
	}
	if err != nil {
		return nil, err
	}

	return img, nil
}
