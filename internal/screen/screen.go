package screen

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	_ "image/png"
	"log/slog"
)

type Screen struct {
	app       fyne.App
	window    fyne.Window
	container *fyne.Container
	progress  *canvas.Text
}

func New() *Screen {
	r := Screen{
		app: app.New(),
	}
	r.window = r.app.NewWindow("Booting")
	r.progress = canvas.NewText("Starting", color.Black)
	r.window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			r.window.Close()
		} else {
			slog.Debug("Key Pressed",
				slog.String("Key", string(k.Name)),
				slog.Int("Code", k.Physical.ScanCode),
			)
		}
	})
	r.window.Canvas().SetOnTypedRune(func(r rune) {
		slog.Debug("Rune Pressed", slog.String("Rune", string(r)))
	})
	return &r
}

// GetWindow returns the window
func (r *Screen) GetWindow() fyne.Window {
	return r.window
}

// SetImage sets the image to be displayed
func (r *Screen) SetImage(imagePath string) *Screen {
	image := canvas.NewImageFromFile(imagePath)
	image.FillMode = canvas.ImageFillOriginal
	r.container = container.NewBorder(nil, r.progress, nil, nil, image)
	r.window.SetContent(r.container)
	return r
}

// SetProgress sets the progress text
func (r *Screen) SetProgress(text string) *Screen {
	r.progress.Text = text
	r.progress.Refresh()
	return r
}

// FullScreen - Sets up the window to be full screen
func (r *Screen) FullScreen() *Screen {
	r.window.SetFullScreen(true)
	r.window.RequestFocus()
	return r
}

// Close closes the window
func (r *Screen) Close() *Screen {
	r.window.Close()
	return r
}
