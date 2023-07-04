package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/reconquest/karma-go"
)

// useCache executes do() function and caches its result in path.
// if cache is older than ttl, do() function is executed again.
// if do() function returns error, cache is not updated.
func useCache[T any](
	do func() (T, error),
	path string,
	ttl time.Duration,
) (T, error) {
	var result T

	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return result, karma.Format(
			err,
			"unable to get user cache dir",
		)
	}

	cacheDir := filepath.Join(userCacheDir, "tubekit")

	cachePath := filepath.Join(cacheDir, path)

	cacheFile, err := os.OpenFile(cachePath, os.O_RDWR, 0644)
	switch {
	case err == nil:
		defer cacheFile.Close()

		stat, err := cacheFile.Stat()
		if err != nil {
			return result, karma.Format(
				err,
				"unable to stat cache file: %s", cachePath,
			)
		}

		if time.Since(stat.ModTime()) < ttl {
			var result T
			err := json.NewDecoder(cacheFile).Decode(&result)
			if err != nil {
				return result, karma.Format(
					err,
					"unable to decode cache file: %s", cachePath,
				)
			}

			return result, nil
		}

	case os.IsNotExist(err):
		err = os.MkdirAll(filepath.Dir(cachePath), 0755)
		if err != nil {
			return result, karma.Format(
				err,
				"unable to create cache dir for: %s", cachePath,
			)
		}

		cacheFile, err = os.Create(cachePath)
		if err != nil {
			return result, karma.Format(
				err,
				"unable to create cache file: %s", cachePath,
			)
		}

		defer cacheFile.Close()

	default:
		return result, karma.Format(
			err,
			"unable to open cache file: %s", cachePath,
		)
	}

	result, err = do()
	if err != nil {
		return result, err
	}

	encoder := json.NewEncoder(cacheFile)

	err = encoder.Encode(result)
	if err != nil {
		return result, karma.Format(
			err,
			"unable to encode cache file: %s", path,
		)
	}

	return result, nil
}
