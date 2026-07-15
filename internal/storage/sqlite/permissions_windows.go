//go:build windows

package sqlite

func secureDatabaseFile(string) error {
	return nil
}
