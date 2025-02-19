package Fs

import "os"

func DirAva(path string) bool {
	st, err := os.Stat(path)
	if err != nil || !st.IsDir() {
		return false
	}
	return true
}
