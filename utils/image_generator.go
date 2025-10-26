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
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		log.Printf("Font file not found at %s. Using default font.", fontPath)
		// Fallback to default font if custom font is not found
		if err := dc.LoadFontFace(gg.DefaultFont, 24); err != nil {
			log.Printf("Failed to load default font: %v", err)
		}
	} else {
		if err := dc.LoadFontFace(fontPath, 24); err != nil {
			log.Printf("Failed to load font from %s: %v", fontPath, err)
			// Fallback to default font
			if err := dc.LoadFontFace(gg.DefaultFont, 24); err != nil {
				log.Printf("Failed to load default font: %v", err)
			}
		}
	}

	dc.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}) // White text

	// Title
	dc.SetFontFace(dc.FontContext.LoadFontFace("sans", 36)) // Re-load with a larger size for title
	dc.DrawStringAnchored("Country Data Summary", imageWidth/2, 50, 0.5, 0.5)

	// Total Countries
	dc.SetFontFace(dc.FontContext.LoadFontFace("sans", 24)) // Reset font size
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
