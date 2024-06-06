package main

import (
    "crypto/md5"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "path/filepath"
    "strings"
    "time"
)

var (
    sourceDir = flag.String("s", "", "Source Directory")
    targetDir = flag.String("t", "", "Target Directory")
)

func main() {
    flag.Parse()
    if *sourceDir == "" || *targetDir == "" {
        flag.Usage()
        log.Fatalf("[ERRO] %v", "Invalid Flag(s)")
    }
    log.Printf("[INFO] %v --> %v", *sourceDir, *targetDir)
    for {
        backupFiles(*sourceDir, *targetDir)
        time.Sleep(1 * time.Hour)
    }
}

func backupFiles(sourceDir, targetDir string) {
    if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
        log.Fatalf("[ERRO] %v", err)
    }
    filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil || info.IsDir() {
            return err
        }
        sourcePath := path
        targetPath := filepath.Join(targetDir, info.Name())
        return copyFileWithUniqueName(sourcePath, targetPath)
    })
}

func copyFileWithUniqueName(sourcePath, targetPath string) error {
    if _, err := os.Stat(targetPath); err == nil {
        same, err := compareFiles(sourcePath, targetPath)
        if err != nil {
            return err
        }
        if same {
            return nil
        }
        targetPath = getUniqueFileName(targetPath)
    }
    return copyFile(sourcePath, targetPath)
}

func getUniqueFileName(filePath string) string {
    ext := filepath.Ext(filePath)
    base := strings.TrimSuffix(filePath, ext)
    for i := 1; ; i++ {
        newPath := fmt.Sprintf("%s-%d%s", base, i, ext)
        if _, err := os.Stat(newPath); os.IsNotExist(err) {
            return newPath
        }
    }
}

func compareFiles(file1, file2 string) (bool, error) {
    hash1, err := getFileHash(file1)
    if err != nil {
        return false, err
    }
    hash2, err := getFileHash(file2)
    if err != nil {
        return false, err
    }
    return hash1 == hash2, nil
}

func getFileHash(filePath string) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hash := md5.New()
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func copyFile(sourcePath, targetPath string) error {
    sourceFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    defer sourceFile.Close()

    targetFile, err := os.Create(targetPath)
    if err != nil {
        return err
    }
    defer targetFile.Close()

    if _, err := io.Copy(targetFile, sourceFile); err != nil {
        return err
    }
    return targetFile.Sync()
}