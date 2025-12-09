package comm

import "os"

func FileCount(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, v := range entries {
		if v.IsDir() {
			continue
		}

		count++
	}

	return count, nil
}
