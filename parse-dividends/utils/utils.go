package utils

import "os"

// Checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// Returns slice of strings containing filenames within the specified directory path.
func GetFileNames(dirPath string) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	fis, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	data := make([]string, len(fis))
	for i, fi := range fis {
		data[i] = fi.Name()
	}

	return data, nil
}
