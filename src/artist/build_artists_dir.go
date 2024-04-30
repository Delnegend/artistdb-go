package artist

import (
	"artistdb-go/src/socials"
	"fmt"
	"os"
	"path/filepath"
)

// BuildArtistsDir writes the artists to the output directory, with everything formatted as

// displayName,avatarUrl
// socialUrl1,description1
// socialUrl2,description2
// ...
func BuildArtistsDir(supportedSocials *socials.Supported, artists *[]Artist, outputDir string) error {
	outputDirInfo, err := os.Stat(outputDir)
	if err != nil {
		if os.IsExist(err) {
			if err := os.Remove(outputDir); err != nil {
				return err
			}
		}
		if err := os.Mkdir(outputDir, 0755); err != nil {
			return err
		}
	} else {
		if !outputDirInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", outputDir)
		}
	}

	// filename -> content
	result := make(map[string]string, 0)
	for _, artist := range *artists {
		artistContent, err := artist.Marshal(supportedSocials)
		if err != nil {
			return err
		}
		for fileName, content := range artistContent {
			result[fileName] = content
		}
	}

	for fileName, content := range result {
		filePath := filepath.Join(outputDir, fileName)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return nil
}
