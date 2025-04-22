package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url" // Import strconv for history limit
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	// Import for strength color (optional later)
)

const (
	appVersion            = "0.2.0"                                              // Application version
	devName               = "b1onicle-dev"                                       // Developer name
	websiteURL            = "https://github.com/b1onicle-dev/password-generator" // Website URL
	prefsHistoryKey       = "passwordHistory"                                    // Key for preferences
	prefsHistoryLimitKey  = "historyLimit"                                       // Key for history limit
	prefsConfirmClearKey  = "confirmClearHistory"                                // Key for confirm clear
	prefsCustomSymbolsKey = "customSymbols"                                      // Key for custom symbols
	defaultSymbols        = "!@#$%^&*()_+=-`~[]{};':\",./<>?"                    // Default symbols
)

// Password generation function (accepts preferences for custom symbols)
func generatePassword(length int, useUpper, useLower, useDigits, useSymbols bool, prefs fyne.Preferences) (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		// defaultSymbols is now a global const
	)

	// Determine which symbols to use
	symbolsToUse := prefs.String(prefsCustomSymbolsKey)
	if symbolsToUse == "" {
		symbolsToUse = defaultSymbols // Use default if custom is empty
	}

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
		charset.WriteString(symbolsToUse) // Use the determined symbols
	}

	finalCharset := charset.String()
	if finalCharset == "" {
		return "", errors.New("no character set selected")
	}

	if length <= 0 {
		return "", errors.New("password length must be greater than 0")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = finalCharset[rng.Intn(len(finalCharset))]
	}

	return string(password), nil
}

// Function to calculate password strength, now returns text and theme color name
func calculateStrength(length int, useUpper, useLower, useDigits, useSymbols bool) (string, fyne.ThemeColorName) {
	score := 0
	sets := 0

	if useUpper {
		sets++
	}
	if useLower {
		sets++
	}
	if useDigits {
		sets++
	}
	if useSymbols {
		sets++
	}

	// Base score on length
	if length >= 16 {
		score += 3
	} else if length >= 12 {
		score += 2
	} else if length >= 8 {
		score += 1
	}

	// Add score based on character sets used
	if sets == 4 {
		score += 3
	} else if sets == 3 {
		score += 2
	} else if sets >= 2 {
		score += 1
	}

	// Determine strength text and color name based on score
	var strengthText string
	var strengthColorName fyne.ThemeColorName

	if score >= 5 { // Very Strong
		strengthText = "Very Strong"
		strengthColorName = theme.ColorNameSuccess // Use theme color for success (usually green)
	} else if score >= 4 { // Strong
		strengthText = "Strong"
		strengthColorName = theme.ColorNameSuccess // Also green, maybe less intense?
	} else if score >= 2 { // Medium
		strengthText = "Medium"
		strengthColorName = theme.ColorNameWarning // Use theme color for warning (usually yellow/orange)
	} else { // Weak
		strengthText = "Weak"
		strengthColorName = theme.ColorNameError // Use theme color for error (usually red)
	}

	if score == 0 {
		strengthText = "Very Weak"
		strengthColorName = theme.ColorNameError // Still red
	}

	return strengthText, strengthColorName
}

// Function to refresh the history view
func refreshHistoryView(historyBox *fyne.Container, prefs fyne.Preferences) {
	historyList := prefs.StringList(prefsHistoryKey)
	historyBox.RemoveAll()

	for i, password := range historyList {
		// Capture index and password for the button handler
		idx := i
		pwd := password

		passLabel := widget.NewLabel(pwd)
		deleteBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
			currentList := prefs.StringList(prefsHistoryKey)
			if idx < len(currentList) { // Bounds check
				// Create a new list excluding the item at idx
				newList := append(currentList[:idx], currentList[idx+1:]...)
				prefs.SetStringList(prefsHistoryKey, newList)
				refreshHistoryView(historyBox, prefs) // Refresh the view after deletion
			}
		})
		// Make delete button less prominent
		deleteBtn.Importance = widget.LowImportance

		historyCard := container.NewBorder(nil, nil, nil, deleteBtn, passLabel) // Password stretches, delete on right
		historyBox.Add(historyCard)
	}
	historyBox.Refresh()
}

func main() {
	// Initialize Fyne application with a unique ID for Preferences
	appID := "com.b1onicle-dev.passwordgenerator" // Use reverse domain name notation or similar unique ID
	a := app.NewWithID(appID)

	// Get preferences handle early
	prefs := a.Preferences()

	// Create main window
	w := a.NewWindow("Password Generator") // Window Title: Password Generator

	// Define historyBox early so it can be used in generateButton handler
	var historyBox *fyne.Container

	// -- Remove unused lengthLabel --
	// lengthLabel := widget.NewLabel("Password length:")
	initialLength := 12.0                                                      // Initial value for the slider
	lengthValueLabel := widget.NewLabel(fmt.Sprintf("%d", int(initialLength))) // Label to display the value
	lengthSlider := widget.NewSlider(4, 64)                                    // Slider from 4 to 64
	lengthSlider.Value = initialLength
	lengthSlider.Step = 1 // Step 1
	lengthSlider.OnChanged = func(value float64) {
		lengthValueLabel.SetText(fmt.Sprintf("%d", int(value)))
	}

	// Define generateButton *before* checkboxChanged so it can be used as focus target
	var generateButton *widget.Button

	// Checkboxes for options
	// Common handler to remove focus from checkboxes
	checkboxChanged := func(b bool) {
		// Try to remove focus from all elements
		// Note: generateButton might be nil if called before button is fully assigned,
		// but Focus(nil) is safe.
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

	// --- Strength Indicator (RichText) ---
	strengthRichText := widget.NewRichText()
	// Устанавливаем перенос слов, если текст будет слишком длинным
	strengthRichText.Wrapping = fyne.TextWrapWord

	// Buttons (with Icons)
	copyButton := widget.NewButtonWithIcon("Copy", theme.ContentCopyIcon(), func() {
		// Existing OnTapped logic for copyButton
		password := passwordLabel.Text
		if password != "" {
			w.Clipboard().SetContent(password)
			errorBinding.Set("Password copied!")
			go func() {
				time.Sleep(2 * time.Second)
				currentMsg, _ := errorBinding.Get()
				if currentMsg == "Password copied!" {
					errorBinding.Set("")
				}
			}()
		}
	})
	copyButton.Disable() // Initially disabled

	// --- "About" Button ---
	// Удаляем старую кнопку About
	// aboutButton := widget.NewButton("About", func() { ... })

	// --- Settings Form ---
	lengthSliderContainer := container.NewBorder(nil, nil, nil, lengthValueLabel, lengthSlider) // Length label is now inside the FormItem
	settingsForm := widget.NewForm(
		widget.NewFormItem("Length:", lengthSliderContainer), // Label text here
		widget.NewFormItem("Uppercase:", upperCheck),
		widget.NewFormItem("Lowercase:", lowerCheck),
		widget.NewFormItem("Digits:", digitsCheck),
		widget.NewFormItem("Symbols:", symbolsCheck),
	)

	// --- Button Logic (Generate) ---
	generateButton = widget.NewButtonWithIcon("Generate", theme.ConfirmIcon(), func() { // Use ConfirmIcon
		errorBinding.Set("")
		passwordLabel.SetText("")
		strengthRichText.Segments = []widget.RichTextSegment{}
		strengthRichText.Refresh()
		copyButton.Disable()

		// Get length from the slider
		length := int(lengthSlider.Value)

		// Add check just in case (needed for checkboxChanged focus target)
		if length <= 0 {
			errorBinding.Set("Error: Invalid length")
			return
		}

		useUpper := upperCheck.Checked
		useLower := lowerCheck.Checked
		useDigits := digitsCheck.Checked
		useSymbols := symbolsCheck.Checked

		// --- Call generatePassword with prefs ---
		password, err := generatePassword(length, useUpper, useLower, useDigits, useSymbols, prefs) // Pass prefs
		if err != nil {
			errorBinding.Set("Error: " + err.Error())
			return
		}

		passwordLabel.SetText(password)
		strengthText, strengthColorName := calculateStrength(length, useUpper, useLower, useDigits, useSymbols)
		strengthSegment := &widget.TextSegment{
			Text: "Strength: " + strengthText,
			Style: widget.RichTextStyle{
				ColorName: strengthColorName,
				TextStyle: fyne.TextStyle{Bold: true},
			},
		}
		strengthRichText.Segments = []widget.RichTextSegment{strengthSegment}
		strengthRichText.Refresh()

		// --- Save to History and Refresh History View ---
		currentHistory := prefs.StringList(prefsHistoryKey)
		newHistory := append([]string{password}, currentHistory...)

		// --- Apply History Limit ---
		limitStr := prefs.StringWithFallback(prefsHistoryLimitKey, "50") // Default to 50
		limit, errConv := strconv.Atoi(limitStr)
		// Use limit only if conversion is successful and limit is positive (0 means no limit here)
		if errConv == nil && limit > 0 && len(newHistory) > limit {
			newHistory = newHistory[:limit]
		}
		// --- End Apply History Limit ---

		prefs.SetStringList(prefsHistoryKey, newHistory)
		refreshHistoryView(historyBox, prefs) // Update the history tab UI

		copyButton.Enable()
	})

	// --- Create Content for "Generate" Tab ---
	// Wrap settings in a Card
	settingsCard := widget.NewCard("Settings", "", settingsForm)

	// Wrap output elements in a Card
	outputBox := container.NewVBox(passwordScroll, strengthRichText, copyButton)
	outputCard := widget.NewCard("Generated Password", "", outputBox)

	generateTabContent := container.NewVBox(
		settingsCard, // Card with settings
		// widget.NewSeparator(), // Optional separator
		generateButton, // Generate button between cards
		// widget.NewSeparator(), // Optional separator
		outputCard,         // Card with output
		layout.NewSpacer(), // Pushes error label to bottom if VBox expands
		errorLabel,         // Error label at the bottom
	)

	// --- Create Content for "About" Tab ---
	parsedURL, _ := url.Parse(websiteURL) // Need URL for hyperlink
	githubButton := widget.NewButtonWithIcon("GitHub", theme.ComputerIcon(), func() {
		a.OpenURL(parsedURL) // Use app instance 'a' to open URL
	})
	aboutTabContent := container.NewCenter( // Center the content
		container.NewVBox(
			widget.NewLabelWithStyle("Password Generator", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			widget.NewLabel(fmt.Sprintf("Version: %s", appVersion)),
			widget.NewLabel(fmt.Sprintf("Developer: %s", devName)),
			githubButton, // Add the GitHub button here
		),
	)

	// --- Create Content for "Settings" Tab ---
	// Theme selector
	themeRadio := widget.NewRadioGroup([]string{"Light", "Dark"}, func(selected string) {
		switch selected {
		case "Light":
			a.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			a.Settings().SetTheme(theme.DarkTheme())
		}
	})
	// Set initial value based on current setting
	if fyne.CurrentApp().Settings().Theme() == theme.DarkTheme() { // Compare with actual dark theme instance
		themeRadio.SetSelected("Dark")
	} else {
		themeRadio.SetSelected("Light") // Default to Light if not Dark
	}

	// History Limit Selector
	historyLimitSelect := widget.NewSelect([]string{"20", "50", "100", "Unlimited"}, func(selected string) {
		limit := selected
		if selected == "Unlimited" {
			limit = "0" // Use "0" to represent unlimited in preferences
		}
		prefs.SetString(prefsHistoryLimitKey, limit)
	})
	currentLimit := prefs.StringWithFallback(prefsHistoryLimitKey, "50") // Default limit 50
	if currentLimit == "0" {
		historyLimitSelect.SetSelected("Unlimited")
	} else {
		historyLimitSelect.SetSelected(currentLimit)
	}

	// Confirm Clear History Checkbox
	confirmClearCheck := widget.NewCheck("Confirm before clearing history", func(checked bool) {
		prefs.SetBool(prefsConfirmClearKey, checked)
	})
	confirmClearCheck.SetChecked(prefs.BoolWithFallback(prefsConfirmClearKey, true)) // Default to true

	// Custom Symbols Entry
	customSymbolsEntry := widget.NewEntry()
	customSymbolsEntry.SetPlaceHolder("Default: " + defaultSymbols)
	customSymbolsEntry.SetText(prefs.String(prefsCustomSymbolsKey)) // Load saved value
	customSymbolsEntry.OnChanged = func(text string) {
		prefs.SetString(prefsCustomSymbolsKey, text) // Save on change
	}

	settingsTabContent := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Theme:", themeRadio),
			widget.NewFormItem("Max History:", historyLimitSelect),
			widget.NewFormItem("", confirmClearCheck), // Checkbox without label in form item
			widget.NewFormItem("Custom Symbols:", customSymbolsEntry),
		),
		layout.NewSpacer(), // Push to top
	)

	// --- Create Content for "History" Tab ---
	historyBox = container.NewVBox() // Assign to the declared variable (Fix linter error)
	historyScroll := container.NewScroll(historyBox)
	clearHistoryButton := widget.NewButtonWithIcon("Clear History", theme.ContentClearIcon(), func() {
		// --- Implement Confirm Clear ---
		if prefs.BoolWithFallback(prefsConfirmClearKey, true) { // Check the preference
			dialog.ShowConfirm("Confirm Clear", "Are you sure you want to clear all password history?", func(confirm bool) {
				if confirm {
					prefs.SetStringList(prefsHistoryKey, []string{}) // Clear the list in preferences
					refreshHistoryView(historyBox, prefs)            // Refresh the view
				}
			}, w)
		} else {
			// Clear directly if confirmation is off
			prefs.SetStringList(prefsHistoryKey, []string{}) // Clear the list in preferences
			refreshHistoryView(historyBox, prefs)            // Refresh the view
		}
		// --- End Implement Confirm Clear ---
	})
	clearHistoryButton.Importance = widget.DangerImportance

	// Wrap the scroll area in a card
	historyCard := widget.NewCard("History", "", historyScroll)

	historyTabContent := container.NewBorder(nil, clearHistoryButton, nil, nil, historyCard) // Card with history in center

	// Initial population of history view
	refreshHistoryView(historyBox, prefs)

	// --- Create Tabs ---
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Generate", theme.HomeIcon(), generateTabContent),
		container.NewTabItemWithIcon("History", theme.HistoryIcon(), historyTabContent), // Add History Tab
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingsTabContent),
		container.NewTabItemWithIcon("About", theme.InfoIcon(), aboutTabContent),
	)

	// --- Widget Layout (Now uses Tabs) ---
	w.SetContent(tabs) // Set tabs as the main content

	// Set a minimum window size to fit everything
	w.Resize(fyne.NewSize(450, 500)) // Size might need adjustment for tabs

	// Start the application loop
	w.ShowAndRun()
}
