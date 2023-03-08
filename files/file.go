package files

import (
	"archive/zip"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 判断文件夹是否存在（支持多个路径）
func IsExists(paths ...string) bool {
	if len(paths) == 0 {
		return false
	}

	for _, path := range paths {
		_, err := os.Stat(path)
		exists := err == nil || os.IsExist(err)

		if !exists {
			return false
		}
	}

	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 创建文件夹
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return err
	}

	return os.Chmod(path, os.ModePerm)
}

// 根据文件名称取文件后缀
func GetExtension(fileName string) string {
	if !strings.Contains(fileName, ".") {
		return ""
	}

	fileNameArr := strings.Split(fileName, ".")

	return fileNameArr[len(fileNameArr)-1]
}

// 获取文件大小
func GetFileSize(filePath string) int64 {
	var result int64
	filepath.Walk(filePath, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

// 删除文件（支持多个文件）
func RemoveFile(filePaths ...string) error {
	for _, filePath := range filePaths {
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func Zip(src, dst string) error {
	// 初始化给定的目录
	baseDir := CreateDirIfNotExists(src)
	CreateDirIfNotExists(dst)
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fw.Close()

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			fmt.Printf("[Zip]关闭文件失败: %v", err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) error {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}

		// 替换文件信息中的文件名(去除baseDir)
		fh.Name = strings.TrimPrefix(path, baseDir)

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return err
		}
		// 输出压缩的内容
		fmt.Printf("[Zip]成功压缩文件: %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}

func UnZip(src, dst string) ([]string, error) {
	// 记录全部被解压的文件
	files := make([]string, 0)
	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(src)
	if err != nil {
		return files, err
	}
	defer zr.Close()

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dst != "" {
		if err := os.MkdirAll(dst, os.ModePerm); err != nil {
			return files, nil
		}
	}

	// 遍历 zr，将文件写入到磁盘
	for _, file := range zr.File {
		var decodeName string
		if file.Flags == 0 {
			// 如果标致位是0, 则是默认的本地编码gbk
			i := bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			// 如果标志为是 1 << 11也就是2048, 则是utf-8编码
			decodeName = file.Name
		}
		path := filepath.Join(dst, decodeName)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return files, nil
			}
			// 因为是目录，跳过当前循环，因为后面都是文件的处理
			continue
		}

		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			return files, nil
		}

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return files, nil
		}

		n, err := io.Copy(fw, fr)
		if err != nil {
			return files, nil
		}

		// 将解压的结果输出
		fmt.Printf("[UnZip]成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)
		// 记录解压文件名
		files = append(files, path)
		// 因为是在循环中，无法使用 defer ，直接放在最后
		// 不过这样也有问题，当出现 err 的时候就不会执行这个了，
		// 可以把它单独放在一个函数中，这里是个实验，就这样了
		fw.Close()
		fr.Close()
	}
	return files, nil
}

// 创建不存在的目录(给定文件名或文件夹名, 判断对应目录是否存在)
func CreateDirIfNotExists(name string) string {
	info, err := os.Stat(name)
	// 存在
	if err == nil {
		// 是目录
		if info.IsDir() {
			return name
		}
		// 是文件
		dir, _ := filepath.Split(name)
		return dir
	} else {
		dir, filename := filepath.Split(name)
		// 在unix系统下以.开头的为隐藏目录
		if strings.HasPrefix(filename, ".") {
			if !strings.Contains(strings.TrimPrefix(filename, "."), ".") {
				dir = name
			}
		} else {
			// 最后一级只要包含"."则认为是文件, 否则是目录
			if !strings.Contains(filename, ".") {
				dir = name
			}
		}
		// 创建不存在的目录
		os.MkdirAll(dir, os.ModePerm)
		return dir
	}
}

// 获取程序运行目录
func GetWorkDir() string {
	pwd, _ := os.Getwd()
	return pwd
}

func DirFileInfo(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}
	return entries
}

func GetDirFileNames(dir string) []string {
	names := make([]string, 0)
	infos := DirFileInfo(dir)
	for _, info := range infos {
		names = append(names, info.Name())
	}
	return names
}

var ostype = runtime.GOOS

func NormalPath(filepath string) string {
	if ostype == "windows" {
		filepath = strings.Replace(filepath, "/", "\\", -1)
	}
	return filepath
}

func NormalPathF(path string, args ...interface{}) string {
	str := fmt.Sprintf(path, args...)
	return NormalPath(str)
}

func GetApplicationDir() string {
	tmpDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		tfile, _ := filepath.Abs(file)
		tmpDir, _ = filepath.Split(tfile)
	}
	return tmpDir
}

func GetAbsolutePath(path string, args ...interface{}) string {
	appDir := GetApplicationDir()
	str := fmt.Sprintf(path, args...)
	str = appDir + str
	return NormalPath(str)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func DeleteFileRecursion(folder string) error {
	return filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if info.IsDir() {
			return os.RemoveAll(path)
		} else {
			return os.Remove(path)
		}
		return nil
	})
}
func Copy(sourceFile, targetFile string) error {
	file1, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	file2, err := os.OpenFile(targetFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file1.Close()
	defer file2.Close()
	_, err = io.Copy(file2, file1)
	if err != nil {
		return err
	}
	return nil
}

func GetPath() string {
	d, err := os.Getwd()
	if err != nil {
		return ""
	}
	return d
}

func WriteStrToFile(f, content string) {
	err := os.WriteFile(f, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
func WriteBytesToFile(f string, b []byte) {
	err := os.WriteFile(f, b, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
