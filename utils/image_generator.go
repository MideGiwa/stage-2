package utils

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	"stage-2/models"

	"github.com/fogleman/gg"
)

const (
	imageWidth  = 800
	imageHeight = 600
	cacheDir    = "cache"
	imagePath   = cacheDir + "/summary.png"
)

// GenerateSummaryImage generates a summary image with total countries, top GDP countries, and refresh timestamp.
func GenerateSummaryImage(totalCountries int, topCountries []models.Country, lastRefreshedAt time.Time) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	dc := gg.NewContext(imageWidth, imageHeight)

	// Set background color
	dc.SetColor(color.RGBA{R: 30, G: 30, B: 30, A: 255}) // Dark background
	dc.Clear()

	// Load font
	fontPath := "assets/Roboto-Bold.ttf" // Assuming a font is available or added
	// Load font
	// Default to a system font or a bundled font if Roboto-Bold.ttf is not found.
	// gg.LoadFontFace expects a valid font path. If `fontPath` fails, it falls back to a generic sans-serif.
	// For better control, one might embed the font or ensure its presence.
	currentFontSize := 24.0
	err := dc.LoadFontFace(fontPath, currentFontSize)
	if err != nil {
		log.Printf("Warning: Could not load font from %s: %v. Using a default system font.", fontPath, err)
		// Try to load a generic system font if the specific one fails
		err = dc.LoadFontFace("sans", currentFontSize)
		if err != nil {
			log.Printf("Error: Could not load any font, text rendering might be affected: %v", err)

		}
	}

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White text

	// Title
	titleFontSize := 36.0
	// Load and set font for the title
	titleLoadErr := dc.LoadFontFace(fontPath, titleFontSize)
	if titleLoadErr != nil {
		log.Printf("Warning: Could not load preferred font from %s for title: %v. Falling back to 'sans' font.", fontPath, titleLoadErr)
		// Fallback to a generic sans-serif font for the title
		if fallbackErr := dc.LoadFontFace("sans", titleFontSize); fallbackErr != nil {
			log.Printf("Error: Could not load 'sans' font for title either: %v. Title rendering might be affected.", fallbackErr)
		}
	}
	dc.DrawStringAnchored("Country Data Summary", imageWidth/2, 50, 0.5, 0.5)

	// Reset to the base font for the rest of the text
	// We need to reload it explicitly as the title font changed the current font face
	baseFontSize := currentFontSize // Assuming baseFontSize should be the same as the initial currentFontSize
	resetBaseFontErr := dc.LoadFontFace(fontPath, baseFontSize)
	if resetBaseFontErr != nil {
		if fallbackErr := dc.LoadFontFace("sans", baseFontSize); fallbackErr != nil {
			log.Printf("Error: Could not load 'sans' font for resetting base font: %v. Text rendering might be affected.", fallbackErr)
		}
	}

	dc.DrawString(fmt.Sprintf("Total Countries: %d", totalCountries), 50, 120)

	// Last Refreshed At
	dc.DrawString(fmt.Sprintf("Last Refreshed: %s", lastRefreshedAt.Format("2006-01-02 15:04:05 MST")), 50, 160)

	// Top 5 Countries by Estimated GDP
	dc.DrawString("Top 5 Countries by Estimated GDP:", 50, 220)

	y := 260.0
	for i, country := range topCountries {
		if i >= 5 {
			break
		}
		gdp := "N/A"
		if country.EstimatedGDP != nil {
			gdp = fmt.Sprintf("%.2f", *country.EstimatedGDP)
		}
		dc.DrawString(fmt.Sprintf("%d. %s (GDP: %s)", i+1, country.Name, gdp), 70, y)
		y += 30
	}

	// Save the image
	outputPath := filepath.Join(cacheDir, "summary.png")
	if err := dc.SavePNG(outputPath); err != nil {
		return fmt.Errorf("failed to save summary image: %w", err)
	}

	log.Printf("Summary image generated successfully at %s", outputPath)
	return nil
}

// GetSummaryImagePath returns the path to the summary image
func GetSummaryImagePath() string {
	return imagePath
}

// EnsureFontExists checks if the required font exists and downloads it if not.
// For simplicity, this example assumes a font exists or uses a default.
// In a real-world scenario, you might want to embed the font or provide a download mechanism.
func EnsureFontExists() {
	fontPath := "assets/Roboto-Bold.ttf"
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		log.Printf("Font file not found at %s. Please ensure 'assets/Roboto-Bold.ttf' exists for optimal image generation.", fontPath)
		log.Printf("You can download it from: https://fonts.google.com/specimen/Roboto and place it in an 'assets' directory.")
	}
}
