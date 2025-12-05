package cmd

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const (
	githubRepo   = "omisai-tech/sshy"
	githubAPIURL = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update sshy to the latest version",
	Long:  `Check for updates and update sshy to the latest release from GitHub.`,
	RunE:  runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	currentVersion := version
	fmt.Printf("Current version: %s\n", currentVersion)

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}
	fmt.Printf("Latest version:  %s\n", latestVersion)

	if isLatestVersion(currentVersion, latestVersion) {
		fmt.Println("\nYou are already running the latest version.")
		return nil
	}

	fmt.Printf("\nA new version is available: %s -> %s\n", currentVersion, latestVersion)
	fmt.Print("Are you sure you want to update? [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Update cancelled.")
		return nil
	}

	if err := performUpdate(latestVersion); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	fmt.Printf("\nSuccessfully updated to %s!\n", latestVersion)
	return nil
}

func fetchLatestVersion() (string, error) {
	resp, err := http.Get(githubAPIURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func isLatestVersion(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")
	return current == latest
}

func performUpdate(version string) error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	goos := runtime.GOOS
	goarch := runtime.GOARCH

	versionNum := strings.TrimPrefix(version, "v")
	var archiveExt string
	if goos == "windows" {
		archiveExt = "zip"
	} else {
		archiveExt = "tar.gz"
	}

	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/%s/sshy_%s_%s_%s.%s",
		githubRepo, version, versionNum, goos, goarch, archiveExt,
	)
	fmt.Printf("Downloading from: %s\n", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	tmpDir, err := os.MkdirTemp("", "sshy-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "sshy."+archiveExt)
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("failed to create archive file: %w", err)
	}

	if _, err := io.Copy(archiveFile, resp.Body); err != nil {
		archiveFile.Close()
		return fmt.Errorf("failed to write archive: %w", err)
	}
	archiveFile.Close()

	var binaryPath string
	if goos == "windows" {
		binaryPath, err = extractFromZip(archivePath, tmpDir)
	} else {
		binaryPath, err = extractFromTarGz(archivePath, tmpDir)
	}
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}

	if err := os.Chmod(binaryPath, 0755); err != nil {
		return fmt.Errorf("failed to set permissions: %w", err)
	}

	if err := os.Rename(binaryPath, execPath); err != nil {
		oldPath := execPath + ".old"
		if renameErr := os.Rename(execPath, oldPath); renameErr != nil {
			return fmt.Errorf("failed to backup old binary: %w", renameErr)
		}
		if copyErr := copyFile(binaryPath, execPath); copyErr != nil {
			os.Rename(oldPath, execPath)
			return fmt.Errorf("failed to install update: %w", copyErr)
		}
		os.Remove(oldPath)
	}

	return nil
}

func extractFromTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var binaryPath string

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if header.Typeflag == tar.TypeReg && (header.Name == "sshy" || filepath.Base(header.Name) == "sshy") {
			binaryPath = filepath.Join(destDir, "sshy")
			outFile, err := os.Create(binaryPath)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return "", err
			}
			outFile.Close()
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("sshy binary not found in archive")
	}
	return binaryPath, nil
}

func extractFromZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var binaryPath string
	for _, f := range r.File {
		if f.Name == "sshy.exe" || filepath.Base(f.Name) == "sshy.exe" {
			binaryPath = filepath.Join(destDir, "sshy.exe")
			rc, err := f.Open()
			if err != nil {
				return "", err
			}
			outFile, err := os.Create(binaryPath)
			if err != nil {
				rc.Close()
				return "", err
			}
			if _, err := io.Copy(outFile, rc); err != nil {
				outFile.Close()
				rc.Close()
				return "", err
			}
			outFile.Close()
			rc.Close()
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("sshy.exe binary not found in archive")
	}
	return binaryPath, nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return os.Chmod(dst, 0755)
}
