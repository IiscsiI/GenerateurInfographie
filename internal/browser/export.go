package browser

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"math"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ExportConfig holds export parameters.
type ExportConfig struct {
	URL         string
	Format      string  // "A3", "A4", "A2"
	Orientation string  // "portrait", "landscape"
	DPI         float64 // 72, 150, 300
	BrowserPath string
}

// paperSizes in inches (width, height) for portrait orientation.
var paperSizes = map[string][2]float64{
	"A4": {8.27, 11.69},
	"A3": {11.69, 16.54},
	"A2": {16.54, 23.39},
}

// paperViewport returns CSS pixel dimensions matching the exact paper ratio.
// baseWidth is the reference width; height is computed from the ratio.
func paperViewport(format, orientation string, baseWidth int) (int, int) {
	size, ok := paperSizes[format]
	if !ok {
		size = paperSizes["A3"]
	}
	w, h := size[0], size[1]
	if orientation == "landscape" {
		w, h = h, w
	}
	ratio := h / w // height/width ratio
	vpH := int(math.Round(float64(baseWidth) * ratio))
	return baseWidth, vpH
}

// buildRenderURL appends format and orientation query params.
func buildRenderURL(baseURL, format, orientation string) string {
	sep := "?"
	if strings.Contains(baseURL, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%sformat=%s&orientation=%s", baseURL, sep, format, orientation)
}

// ExportPNG captures a full-page screenshot with layout matching the paper ratio.
func ExportPNG(ctx context.Context, cfg ExportConfig) ([]byte, error) {
	dpi := cfg.DPI
	if dpi == 0 {
		dpi = 150
	}
	scale := dpi / 96.0
	if scale < 1 {
		scale = 1
	}

	orientation := cfg.Orientation
	if orientation == "" {
		orientation = "portrait"
	}

	baseWidth := 1200
	if orientation == "landscape" {
		baseWidth = 1600
	}
	vpW, vpH := paperViewport(cfg.Format, orientation, baseWidth)
	renderURL := buildRenderURL(cfg.URL, cfg.Format, orientation)

	opts := defaultOpts(cfg.BrowserPath, vpW)
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()
	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()
	taskCtx, timeoutCancel := context.WithTimeout(taskCtx, 60*time.Second)
	defer timeoutCancel()

	var buf []byte

	log.Printf("Export PNG: %s (%s %s %gDPI, viewport %dx%d)", renderURL, cfg.Format, orientation, dpi, vpW, vpH)
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(renderURL),
		chromedp.WaitVisible("#infographic-ready", chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		chromedp.EmulateViewport(int64(vpW), int64(vpH), chromedp.EmulateScale(scale)),
		chromedp.Sleep(500*time.Millisecond),
		// Capture exactly the viewport (not beyond) to match paper ratio
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithCaptureBeyondViewport(false).
				WithFromSurface(true).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("export PNG: %w", err)
	}

	log.Printf("Export PNG: OK (%d bytes)", len(buf))
	return buf, nil
}

// ExportPDF captures a screenshot at exact paper ratio then wraps in single-page PDF.
func ExportPDF(ctx context.Context, cfg ExportConfig) ([]byte, error) {
	size, ok := paperSizes[cfg.Format]
	if !ok {
		size = paperSizes["A3"]
	}

	paperWidthIn := size[0]
	paperHeightIn := size[1]
	orientation := cfg.Orientation
	if orientation == "" {
		orientation = "portrait"
	}
	if orientation == "landscape" {
		paperWidthIn, paperHeightIn = paperHeightIn, paperWidthIn
	}

	pageW := paperWidthIn * 72.0 // PDF points
	pageH := paperHeightIn * 72.0

	baseWidth := 1200
	if orientation == "landscape" {
		baseWidth = 1600
	}
	vpW, vpH := paperViewport(cfg.Format, orientation, baseWidth)
	renderURL := buildRenderURL(cfg.URL, cfg.Format, orientation)
	scale := 2.0

	opts := defaultOpts(cfg.BrowserPath, vpW)
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()
	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()
	taskCtx, timeoutCancel := context.WithTimeout(taskCtx, 60*time.Second)
	defer timeoutCancel()

	var jpegBuf []byte

	log.Printf("Export PDF: %s (%s %s, viewport %dx%d)", renderURL, cfg.Format, orientation, vpW, vpH)
	err := chromedp.Run(taskCtx,
		chromedp.Navigate(renderURL),
		chromedp.WaitVisible("#infographic-ready", chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		chromedp.EmulateViewport(int64(vpW), int64(vpH), chromedp.EmulateScale(scale)),
		chromedp.Sleep(500*time.Millisecond),
		// Capture exactly the viewport — matches paper ratio
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			jpegBuf, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatJpeg).
				WithQuality(95).
				WithCaptureBeyondViewport(false).
				WithFromSurface(true).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("export PDF screenshot: %w", err)
	}

	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(jpegBuf))
	if err != nil {
		return nil, fmt.Errorf("export PDF decode: %w", err)
	}

	log.Printf("Export PDF: screenshot %dx%d, page %.0fx%.0f pt", imgCfg.Width, imgCfg.Height, pageW, pageH)

	// Image fills the entire page — no aspect ratio adjustment needed
	// because viewport was set to exact paper ratio
	pdfBytes := buildPDFWithJPEG(jpegBuf, imgCfg.Width, imgCfg.Height, pageW, pageH)

	log.Printf("Export PDF: OK (%d bytes)", len(pdfBytes))
	return pdfBytes, nil
}

// buildPDFWithJPEG creates a minimal single-page PDF with the JPEG stretched to fill the page.
func buildPDFWithJPEG(jpegData []byte, imgW, imgH int, pageW, pageH float64) []byte {
	var buf bytes.Buffer
	offsets := make([]int, 7)

	// Image fills entire page (viewport already matches paper ratio)
	contentStream := fmt.Sprintf("q\n%.4f 0 0 %.4f 0 0 cm\n/Img1 Do\nQ\n",
		pageW, pageH)

	buf.WriteString("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")

	offsets[1] = buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	offsets[2] = buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	offsets[3] = buf.Len()
	buf.WriteString(fmt.Sprintf("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.4f %.4f] /Contents 4 0 R /Resources << /XObject << /Img1 5 0 R >> >> >>\nendobj\n",
		pageW, pageH))

	offsets[4] = buf.Len()
	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n", len(contentStream)))
	buf.WriteString(contentStream)
	buf.WriteString("endstream\nendobj\n")

	offsets[5] = buf.Len()
	buf.WriteString(fmt.Sprintf("5 0 obj\n<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode /Length %d >>\nstream\n",
		imgW, imgH, len(jpegData)))
	buf.Write(jpegData)
	buf.WriteString("\nendstream\nendobj\n")

	offsets[6] = buf.Len()
	buf.WriteString("6 0 obj\n<< /Producer (Generateur Infographie v2.0) /Creator (Generateur Infographie) >>\nendobj\n")

	xrefOffset := buf.Len()
	buf.WriteString("xref\n")
	buf.WriteString(fmt.Sprintf("0 %d\n", len(offsets)))
	buf.WriteString("0000000000 65535 f \n")
	for i := 1; i < len(offsets); i++ {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	buf.WriteString(fmt.Sprintf("trailer\n<< /Size %d /Root 1 0 R /Info 6 0 R >>\n", len(offsets)))
	buf.WriteString(fmt.Sprintf("startxref\n%d\n%%%%EOF\n", xrefOffset))

	return buf.Bytes()
}

func defaultOpts(browserPath string, vpWidth int) []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.WindowSize(vpWidth, 900),
	)
	if browserPath != "" {
		opts = append(opts, chromedp.ExecPath(browserPath))
	}
	return opts
}
