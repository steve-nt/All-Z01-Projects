package handlers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"strings"
	"sync"
	"time"

	"sgougoupractice/assets"
)

//var contentFS embed.FS

var BuildStamp string

func debugDumpFS() {
	fmt.Println("---- embedded paths ----")
	fs.WalkDir(assets.ContentFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("ERR:", err)
			return nil
		}
		fmt.Println(p)
		return nil
	})
	fmt.Println("------------------------")
}

func buildTime() time.Time {
	t, _ := time.Parse(time.RFC3339, BuildStamp)
	if t.IsZero() {
		t = time.Now()
	}
	return t
}

type asset struct {
	data         []byte
	etag         string
	mod          time.Time
	hashedPublic string
}

var (
	once       sync.Once
	assetTable map[string]asset
)

func loadAssets() {
	assetTable = make(map[string]asset)
	const root = "static"

	err := fs.WalkDir(assets.ContentFS, root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // propagate → WalkDir stops and returns this error
		}
		if d.IsDir() {
			return nil
		}
		b, _ := assets.ContentFS.ReadFile(p)

		sum := sha256.Sum256(b)
		tag := hex.EncodeToString(sum[:4])
		ext := filepath.Ext(p)
		base := strings.TrimSuffix(strings.TrimPrefix(p, "static/"), ext)
		logical := strings.TrimPrefix(p, root+"/")
		hashedFile := fmt.Sprintf("%s.%s%s", base, tag, ext)
		publicURL := "/static/" + hashedFile

		a := asset{
			data:         b,
			etag:         `W/"` + tag + `"`,
			mod:          buildTime(),
			hashedPublic: publicURL,
		}

		assetTable[logical] = a
		assetTable[hashedFile] = a

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func StaticHandler() Apphandler {
	debugDumpFS()
	once.Do(loadAssets)

	return func(w http.ResponseWriter, r *http.Request) error {
		key := strings.TrimPrefix(r.URL.Path, "/")
		if a, ok := assetTable[key]; ok {
			if r.Header.Get("If-None-Match") == a.etag {
				w.WriteHeader(http.StatusNotModified)
				return nil
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Header().Set("ETag", a.etag)
			http.ServeContent(w, r, key, a.mod, bytes.NewReader(a.data))
			return nil
		}

		if idx := strings.LastIndex(key, "."); idx != -1 {
			ext := key[idx:]
			front := key[:idx]
			if j := strings.LastIndex(front, "."); j != -1 {
				baseKey := front[:j] + ext
				if _, ok := assetTable[baseKey]; ok {
					return &HTTPError{
						Status:  http.StatusInternalServerError,
						Message: "asset hash mismatch",
					}
				}

			}

		}
		return &HTTPError{
			Status:  http.StatusNotFound,
			Message: "asset not found",
		}
	}

}
