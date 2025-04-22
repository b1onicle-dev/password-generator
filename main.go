package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	appVersion = "0.1.0"                                              // Application version
	devName    = "b1onicle-dev"                                       // Developer name
	websiteURL = "https://github.com/b1onicle-dev/password-generator" // Website URL
)

// Password generation function
func generatePassword(length int, useUpper, useLower, useDigits, useSymbols bool) (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		symbols   = "!@#$%^&*()_+=-`~[]{};':\",./<>?"
	)

	var charset strings.Builder
	if useLower {
		charset.WriteString(lowercase)
	}
	if useUpper {
		charset.WriteString(uppercase)
	}
	if useDigits {
		charset.WriteString(digits)
	}
	if useSymbols {
		charset.WriteString(symbols)
	}

	finalCharset := charset.String()
	if finalCharset == "" {
		return "", errors.New("no character set selected") // Error: No character set selected
	}

	if length <= 0 {
		return "", errors.New("password length must be greater than 0") // Error: Password length must be > 0
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = finalCharset[rng.Intn(len(finalCharset))]
	}

	return string(password), nil
}

func main() {
	// Initialize Fyne application
	a := app.New()

	// Create main window
	w := a.NewWindow("Password Generator") // Window Title: Password Generator

	// Label and slider for password length
	lengthLabel := widget.NewLabel("Password length:")                         // Label: Password length:
	initialLength := 12.0                                                      // Initial value for the slider
	lengthValueLabel := widget.NewLabel(fmt.Sprintf("%d", int(initialLength))) // Label to display the value
	lengthSlider := widget.NewSlider(4, 64)                                    // Slider from 4 to 64
	lengthSlider.Value = initialLength
	lengthSlider.Step = 1 // Step 1
	lengthSlider.OnChanged = func(value float64) {
		lengthValueLabel.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Buttons (define generateButton earlier for use in checkboxChanged)
	var generateButton *widget.Button

	// Checkboxes for options
	// Common handler to remove focus from checkboxes
	checkboxChanged := func(b bool) {
		// Try to remove focus from all elements
		w.Canvas().Focus(nil)
	}
	upperCheck := widget.NewCheck("Uppercase (A-Z)", checkboxChanged)    // Checkbox: Uppercase (A-Z)
	lowerCheck := widget.NewCheck("Lowercase (a-z)", checkboxChanged)    // Checkbox: Lowercase (a-z)
	lowerCheck.SetChecked(true)                                          // Checked by default
	digitsCheck := widget.NewCheck("Digits (0-9)", checkboxChanged)      // Checkbox: Digits (0-9)
	symbolsCheck := widget.NewCheck("Symbols (!@#...)", checkboxChanged) // Checkbox: Symbols (!@#...)

	// Label for displaying the result
	passwordLabel := widget.NewLabel("")
	passwordLabel.Wrapping = fyne.TextWrapWord           // Word wrap for long passwords
	passwordScroll := container.NewScroll(passwordLabel) // Add scroll for very long passwords
	passwordScroll.SetMinSize(fyne.NewSize(380, 50))     // Limit the height of the password area

	// Create data binding for the error text
	errorBinding := binding.NewString()

	// Label for errors (now with data binding)
	errorLabel := widget.NewLabelWithData(errorBinding)
	errorLabel.TextStyle.Bold = true
	// Set red color for errors
	errorLabel.Importance = widget.DangerImportance

	// Buttons
	copyButton := widget.NewButton("Copy", nil) // Button: Copy
	copyButton.Disable()                        // Initially disabled

	// --- "About" Button ---
	aboutButton := widget.NewButton("About", func() { // Button: About
		// Create URL
		parsedURL, _ := url.Parse(websiteURL)

		// Create dialog content
		aboutContent := container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Version: %s", appVersion)), // Label: Version
			widget.NewLabel(fmt.Sprintf("Developer: %s", devName)),  // Label: Developer
			widget.NewHyperlink(websiteURL, parsedURL),
		)

		// Show dialog
		dialog.ShowCustom("About", "Close", aboutContent, w) // Dialog Title: About, Button: Close
	})

	// --- Button Logic ---
	generateButton = widget.NewButton("Generate", func() { // Button: Generate
		errorBinding.Set("") // Clear previous generation errors

		// Get length from the slider
		length := int(lengthSlider.Value)

		if length <= 0 { // Add a check just in case, though slider shouldn't allow 0 or less
			errorBinding.Set("Error: Invalid length") // Error: Invalid length
			passwordLabel.SetText("")
			copyButton.Disable()
			return
		}

		useUpper := upperCheck.Checked
		useLower := lowerCheck.Checked
		useDigits := digitsCheck.Checked
		useSymbols := symbolsCheck.Checked

		password, err := generatePassword(length, useUpper, useLower, useDigits, useSymbols)
		if err != nil {
			errorBinding.Set("Error: " + err.Error()) // Error: <original error>
			passwordLabel.SetText("")
			copyButton.Disable()
			return
		}

		passwordLabel.SetText(password)
		copyButton.Enable()
	})

	copyButton.OnTapped = func() {
		password := passwordLabel.Text
		if password != "" {
			w.Clipboard().SetContent(password)
			// Short success message via binding
			errorBinding.Set("Password copied!") // Message: Password copied!
			go func() {
				time.Sleep(2 * time.Second) // Delay before hiding the message
				// Get the current value and check
				currentMsg, _ := errorBinding.Get()
				if currentMsg == "Password copied!" {
					// Clear the message via binding
					errorBinding.Set("")
				}
			}()
		}
	}

	// --- Widget Layout ---
	w.SetContent(container.NewVBox(
		widget.NewLabel("Generation Settings:"), // Label: Generation Settings:
		// Use Border layout for the slider so it stretches
		container.NewBorder(nil, nil, lengthLabel, lengthValueLabel, lengthSlider),
		upperCheck,
		lowerCheck,
		digitsCheck,
		symbolsCheck,
		generateButton,
		widget.NewLabel("Generated Password:"), // Label: Generated Password:
		passwordScroll,                         // Use scroll container
		copyButton,
		errorLabel, // Label for errors/messages
		aboutButton,
	))

	// Set a minimum window size to fit everything
	w.Resize(fyne.NewSize(400, 450)) // Slightly increase height

	// Start the application loop
	w.ShowAndRun()
}
