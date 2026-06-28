package okxapi

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile fetches url into outDir with a content-addressed filename.
func DownloadFile(ctx context.Context, url, outDir, ext string) (path, sumHex string, size int64, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", 0, err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", 0, fmt.Errorf("GET %s: HTTP %d", url, resp.StatusCode)
	}

	tmp, err := os.CreateTemp(outDir, "okx-dl-*")
	if err != nil {
		return "", "", 0, err
	}
	tmpName := tmp.Name()
	defer func() {
		if err != nil {
			os.Remove(tmpName)
		}
	}()

	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(tmp, h), resp.Body)
	tmp.Close()
	if err != nil {
		return "", "", 0, err
	}
	if ext == "" {
		ext = filepath.Ext(url)
	}
	if ext == "" {
		ext = ".bin"
	}
	sumHex = hex.EncodeToString(h.Sum(nil))
	final := filepath.Join(outDir, sumHex+ext)
	if err := os.Rename(tmpName, final); err != nil {
		return "", "", 0, err
	}
	return final, sumHex, n, nil
}
