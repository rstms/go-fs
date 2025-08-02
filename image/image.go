package image

import (
	"fmt"
	fs "github.com/rstms/go-fs"
	"github.com/rstms/go-fs/fat"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

func walk(dir string, entries []fs.DirectoryEntry, target *string) ([]string, error) {
	//log.Printf("walkFS: %+v\n", entries)
	files := []string{}
	for _, entry := range entries {
		//log.Printf("entry Name=%s isDir=%v %+v\n", entry.Name(), entry.IsDir(), entry)
		name := path.Join(dir, entry.Name())
		if entry.IsDir() {
			if entry.Name() != "." && entry.Name() != ".." {
				if target != nil {
					targetDir := filepath.Join(*target, name)
					log.Printf("MKDIR %s\n", targetDir)
					err := os.Mkdir(targetDir, 0700)
					if err != nil {
						return []string{}, err
					}
				}
				files = append(files, name+"/")
				entryDir, err := entry.Dir()
				if err != nil {
					return []string{}, err
				}
				dirFiles, err := walk(name, entryDir.Entries(), target)
				if err != nil {
					return []string{}, err
				}
				for _, dirFile := range dirFiles {
					if dirFile != name {
						files = append(files, dirFile)
					}
				}
			}
		} else {
			files = append(files, name)
			if target != nil {
				targetName := filepath.Join(*target, name)
				log.Printf("COPY %v -> %s\n", name, targetName)
				entryFile, err := entry.File()
				if err != nil {
					return []string{}, err
				}
				err = copyFile(targetName, entryFile, entry.Size())
				if err != nil {
					return []string{}, err
				}
			}
		}
	}
	return files, nil
}

func copyFile(dstName string, src fs.File, size int64) error {
	log.Printf("copyFile(%s, %+v, %d)\n", dstName, src, size)
	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer dst.Close()
	count, err := io.Copy(dst, src)
	if err != nil {
		return err
	}
	if count != size {
		return fmt.Errorf("write count mismatch size=%d written=%d\n", size, count)
	}
	return nil
}

func scanImage(imageFilename string, target *string) ([]string, error) {

	imageFile, err := os.Open(imageFilename)
	if err != nil {
		return []string{}, err
	}
	defer imageFile.Close()

	// BlockDevice backed by the file for our filesystem
	device, err := fs.NewFileDisk(imageFile)
	if err != nil {
		return []string{}, err
	}

	// The actual FAT filesystem
	ffs, err := fat.New(device)
	if err != nil {
		return []string{}, err
	}

	// Get the root directory to the filesystem
	rootDir, err := ffs.RootDir()
	if err != nil {
		return []string{}, err
	}

	entries := rootDir.Entries()
	if len(entries) < 2 {
		return []string{}, nil
	}
	files, err := walk("/", entries[1:], target)
	if err != nil {
		return []string{}, err
	}

	return files, nil
}

func ListFiles(imageFilename string) ([]string, error) {
	return scanImage(imageFilename, nil)
}

func ExtractFiles(imageFilename, outputDirectory string) error {
	err := os.Mkdir(outputDirectory, 0700)
	if err != nil {
		return err
	}
	_, err = scanImage(imageFilename, &outputDirectory)
	if err != nil {
		return err
	}
	return nil
}
