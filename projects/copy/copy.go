package copy

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Serial copies files from source to destination
func Serial(source, dest string) error {
	cleanedRootPath := path.Clean(source)
	err := filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		cleanedSourcePath := path.Clean(filePath)
		destPath := filepath.Join(dest, strings.TrimPrefix(cleanedSourcePath, cleanedRootPath))

		// 如果是文件夹，跳过
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode().Perm())
		}

		// 打开源文件
		srcFile, err := os.Open(cleanedSourcePath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 创建目标文件
		distFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer distFile.Close()

		// 复制
		_, err = io.Copy(distFile, srcFile)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// Concurrent copies files concurrently
func Concurrent(source, dest string) error {
	cleanedRootPath := path.Clean(source)

	waiter := make(chan error)
	count := 0

	err := filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		cleanedSourcePath := path.Clean(filePath)
		destPath := filepath.Join(dest, strings.TrimPrefix(cleanedSourcePath, cleanedRootPath))
		// 如果是文件夹，跳过
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode().Perm())
		}

		count++
		go func() {

			// 打开源文件
			srcFile, err := os.Open(cleanedSourcePath)
			if err != nil {
				waiter <- err
			}
			defer srcFile.Close()

			// 创建目标文件
			distFile, err := os.Create(destPath)
			if err != nil {
				waiter <- err
			}
			defer distFile.Close()

			// 复制
			_, err = io.Copy(distFile, srcFile)
			if err != nil {
				waiter <- err
			}

			waiter <- nil
		}()

		return nil
	})
	if err != nil {
		return err
	}

	for i := 0; i < count; i++ {
		err := <-waiter
		if err != nil {
			return err
		}
	}
	return nil
}
