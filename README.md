# Password Generator (Go + Fyne)

A simple, cross-platform password generator application built with Go and the Fyne GUI toolkit.

## Features

*   **Password Generation:**
    *   Adjustable password length using a slider (4-64 characters).
    *   Options to include uppercase letters, lowercase letters, digits, and symbols.
    *   Generates cryptographically secure random passwords.
    *   Password strength indicator (Weak, Medium, Strong, Very Strong) with color coding.
    *   "Copy to Clipboard" button with icon.
    *   "Generate" button with icon.
*   **Password History:**
    *   Automatically saves generated passwords locally.
    *   History view with individual password deletion.
    *   "Clear History" button.
*   **Settings:**
    *   Switch between Light and Dark themes.
    *   Set maximum history size (20, 50, 100, Unlimited).
    *   Option to require confirmation before clearing history.
    *   Define a custom set of symbols for generation.
*   **Interface:**
    *   Modern tabbed interface (Generate, History, Settings, About).
    *   Structured layout using cards and forms.
    *   Cross-platform native look and feel.

## Screenshot

*(Optional: Add a screenshot of the application here)*

```
[Screenshot placeholder]
```

## Installation

### Prerequisites

*   **Go:** Version 1.18 or later. ([Download Go](https://go.dev/dl/))
*   **C Compiler:** Fyne requires a C compiler for building.
    *   **Linux:** GCC is usually pre-installed or available via package manager (`build-essential` on Debian/Ubuntu, `"Development Tools"` group on Fedora).
    *   **Windows:** You need to install a GCC compiler, for example, via [MSYS2/MinGW-w64](https://www.msys2.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/). Ensure `gcc` is in your system's PATH.

### Linux Instructions

1.  **Install Go:** Follow the official Go installation guide or use your distribution's package manager (e.g., `sudo dnf install golang` on Fedora, `sudo apt install golang` on Debian/Ubuntu).

2.  **Install Fyne Dependencies:** You need development packages for graphics libraries.
    *   **Fedora:**
        ```bash
        sudo dnf install mesa-libGL-devel libXcursor-devel libXi-devel libXrandr-devel libXinerama-devel libXxf86vm-devel
        ```
    *   **Debian/Ubuntu:**
        ```bash
        sudo apt-get update
        sudo apt-get install libgl1-mesa-dev xorg-dev
        ```
    *   *(Note: Package names might differ slightly on other distributions.)*

3.  **Clone the Repository (Optional):**
    ```bash
    git clone https://github.com/b1onicle-dev/password-generator # Replace with your actual repo URL
    cd password-generator
    ```

4.  **Build the Application:**
    ```bash
    go build
    ```

### Windows Instructions

1.  **Install Go:** Download and run the installer from the [official Go website](https://go.dev/dl/).

2.  **Install C Compiler (MinGW-w64):**
    *   Download and run the installer from [MSYS2](https://www.msys2.org/).
    *   After installation, open the MSYS2 MinGW 64-bit terminal (mingw64.exe).
    *   Update the package database and core packages:
        ```bash
        pacman -Syu
        ```
        (You might need to close and reopen the terminal if prompted).
    *   Install the MinGW-w64 GCC toolchain:
        ```bash
        pacman -S mingw-w64-x86_64-gcc
        ```
    *   Add the MinGW-w64 bin directory to your Windows PATH environment variable (e.g., `C:\msys64\mingw64\bin`).

3.  **Clone the Repository (Optional):**
    ```bash
    git clone https://github.com/b1onicle/password-generator # Replace with your actual repo URL
    cd password-generator
    ```

4.  **Build the Application:** Open a standard Command Prompt or PowerShell (not the MSYS2 terminal) and run:
    ```bash
    go build
    ```

## Usage

After building, run the executable file created:

*   **Linux:** `./passwordgen` (or whatever the executable is named)
*   **Windows:** `passwordgen.exe`

Navigate through the tabs:

*   **Generate:** Use the slider and checkboxes to configure password options. Click "Generate". Use the "Copy" button to copy the password. The strength indicator will update automatically.
*   **History:** View previously generated passwords. Click the delete icon next to a password to remove it individually. Use the "Clear History" button at the bottom to remove all entries (respects the confirmation setting).
*   **Settings:** Adjust application theme, history limit, confirmation preference, and custom symbols.
*   **About:** View application version, developer info, and a link to the GitHub repository.

## Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 