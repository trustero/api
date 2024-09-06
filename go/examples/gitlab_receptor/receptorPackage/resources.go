// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.
package receptorPackage

import (
	"embed"
)

//go:embed resources/*
var embedFS embed.FS

func readEmbedFile(path string) (string, error) {
	data, err := embedFS.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func GetInstructionsImpl() (string, error) {
	return readEmbedFile("resources/trr-gitlab.md")
}

func GetLogoImpl() (string, error) {
	return readEmbedFile("resources/GitLab_Logo.svg")
}
