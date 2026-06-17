package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
)

const tokenStoreService = "\u0042\u0050\u0042-Wizard"

type tokenStore struct{}

func newTokenStore() tokenStore {
	return tokenStore{}
}

func (s tokenStore) Load() (*oauth2.Token, error) {
	data, err := os.ReadFile(tokenFilePath())
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (s tokenStore) Save(token *oauth2.Token) error {
	data, err := json.Marshal(token)
	if err != nil {
		return err
	}

	path := tokenFilePath()

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}

	return os.Chmod(path, 0600)
}

func (s tokenStore) Delete() error {
	err := os.Remove(tokenFilePath())
	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func tokenFilePath() string {
	dir, err := os.UserConfigDir()
	if err != nil || dir == "" {
		if home, homeErr := os.UserHomeDir(); homeErr == nil {
			dir = filepath.Join(home, ".config")
		} else {
			dir = "."
		}
	}

	return filepath.Join(dir, tokenStoreService, "oauth-token.json")
}
