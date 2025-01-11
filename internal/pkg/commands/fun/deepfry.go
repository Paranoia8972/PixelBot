package fun

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"math/rand"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

func applyDeepfry(img image.Image) image.Image {
	bounds := img.Bounds()
	dc := gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawImage(img, 0, 0)

	for x := 0; x < bounds.Dx(); x++ {
		for y := 0; y < bounds.Dy(); y++ {
			r, g, b, a := img.At(x, y).RGBA()
			r = (r * 2) & 0xffff
			g = (g * 2) & 0xffff
			b = (b * 2) & 0xffff
			dc.SetRGBA255(int(r>>8), int(g>>8), int(b>>8), int(a>>8))
			dc.SetPixel(x, y)
		}
	}

	for i := 0; i < (bounds.Dx()*bounds.Dy())/10; i++ {
		x := rand.Intn(bounds.Dx())
		y := rand.Intn(bounds.Dy())
		dc.SetRGBA255(rand.Intn(255), rand.Intn(255), rand.Intn(255), 255)
		dc.SetPixel(x, y)
	}

	img = dc.Image()
	bounds = img.Bounds()
	adjusted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r = r >> 8
			g = g >> 8
			b = b >> 8

			contrast := 2.0
			r = uint32(math.Min(255, math.Max(0, float64(r)*contrast)))
			g = uint32(math.Min(255, math.Max(0, float64(g)*contrast)))
			b = uint32(math.Min(255, math.Max(0, float64(b)*contrast)))

			adjusted.Set(x, y, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a >> 8)})
		}
	}

	dc = gg.NewContext(bounds.Dx(), bounds.Dy())
	dc.DrawImage(adjusted, 0, 0)

	for i := 0; i < 5; i++ {
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, dc.Image(), &jpeg.Options{Quality: 1})
		tmpImg, _ := jpeg.Decode(buf)
		dc = gg.NewContext(bounds.Dx(), bounds.Dy())
		dc.DrawImage(tmpImg, 0, 0)
	}

	return dc.Image()
}

func DeepfryCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return
	}

	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		RespondWithMessage(s, i, "Please provide an image to deepfry")
		return
	}

	attachment := i.ApplicationCommandData().Resolved.Attachments[options[0].Value.(string)]
	if attachment == nil {
		content := "Please provide a valid image attachment"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}

	resp, err := http.Get(attachment.URL)
	if err != nil {
		RespondWithMessage(s, i, "Failed to download image")
		return
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		RespondWithMessage(s, i, "Failed to decode image")
		return
	}

	deepfried := applyDeepfry(img)

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, deepfried, &jpeg.Options{Quality: 1})

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Files: []*discordgo.File{
			{
				Name:   "deepfried.jpg",
				Reader: buf,
			},
		},
	})
	if err != nil {
		content := "Failed to send deepfried image"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}
}

func DeepFryImage(inputPath, outputPath string) error {
	src, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	adjusted := imaging.AdjustBrightness(src, 30)
	adjusted = imaging.AdjustContrast(adjusted, 50)

	adjusted = imaging.AdjustSaturation(adjusted, 80)

	err = imaging.Save(adjusted, outputPath)
	if err != nil {
		return err
	}

	return nil
}
