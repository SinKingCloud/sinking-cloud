package web_auth

import (
	"archive/zip"
	"bufio"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/service/rpc"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/util/str"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/model"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_app"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-service/service-open/app/service/web_auth/auth_version"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type pkg struct {
}

var Package *pkg

// BuildRelease 构建正式包
func (p *pkg) BuildRelease(id int, userAuth *model.UserAuth) (string, error) {
	if userAuth == nil {
		return "", errors.New("获取许可证信息失败")
	}
	authVersion, err := auth_version.GetIns().FindByIdWithCache(id, false)
	if err != nil || authVersion == nil || authVersion.ReleaseUrl == "" {
		return "", errors.New("获取版本信息失败")
	}
	buildPkg, err := p.buildPkg(0, authVersion, userAuth)
	if err != nil {
		return "", err
	}
	//完整包下载数量+1
	go func(id int) {
		_ = auth_app.GetIns().UpdateById(authVersion.AuthId, map[string]interface{}{
			"release_down_num": gorm.Expr("`release_down_num`+1"),
		})
		_ = auth_version.GetIns().UpdateById(authVersion.AuthId, map[string]interface{}{
			"release_down_num": gorm.Expr("`release_down_num`+1"),
		})
	}(authVersion.Id)
	return buildPkg, nil
}

// BuildUpdate 构建更新包
func (p *pkg) BuildUpdate(id int, userAuth *model.UserAuth) (string, error) {
	if userAuth == nil {
		return "", errors.New("获取许可证信息失败")
	}
	authVersion, err := auth_version.GetIns().FindByIdWithCache(id, false)
	if err != nil || authVersion == nil || authVersion.UpdateUrl == "" {
		return "", errors.New("获取版本信息失败")
	}
	buildPkg, err := p.buildPkg(1, authVersion, userAuth)
	if err != nil {
		return "", err
	}
	//更新包下载数量+1
	go func(id int) {
		_ = auth_app.GetIns().UpdateById(authVersion.AuthId, map[string]interface{}{
			"update_down_num": gorm.Expr("`update_down_num`+1"),
		})
		_ = auth_version.GetIns().UpdateById(authVersion.AuthId, map[string]interface{}{
			"update_down_num": gorm.Expr("`update_down_num`+1"),
		})
	}(authVersion.Id)
	return buildPkg, nil
}

func (p *pkg) buildPkg(types int, authVersion *model.AuthVersion, userAuth *model.UserAuth) (string, error) {
	var u string
	var c string
	if types == 0 {
		u = authVersion.ReleaseUrl
		c = authVersion.ReleaseHash
	} else {
		u = authVersion.UpdateUrl
		c = authVersion.UpdateHash
	}
	filePath, err := p.download(u)
	if err != nil {
		return "", errors.New("获取源文件错误")
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(filePath)
	//验证文件安全性
	if p.Md5(filePath) != c {
		return "", errors.New("文件hash验证失败")
	}
	folder := "./temp/" + userAuth.Code + "/"
	//解压
	if p.unZip(filePath, folder) != nil {
		return "", errors.New("读取源文件错误")
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			return
		}
	}(folder)
	//构建文件
	if authVersion.BuildScript != "" {
		var buildScript map[string]map[string]string
		err = json.Unmarshal([]byte(authVersion.BuildScript), &buildScript)
		if err == nil {
			for k, v := range buildScript {
				if _, e := os.Stat(folder + k); os.IsNotExist(e) {
					continue
				}
				bytes, e := os.ReadFile(folder + k)
				if e != nil {
					continue
				}
				text := string(bytes)
				for k1, v1 := range v {
					var s string
					switch v1 {
					case "{$id}":
						s = strconv.Itoa(userAuth.Id)
						break
					case "{$auth_id}":
						s = strconv.Itoa(userAuth.AuthId)
						break
					case "{$user_id}":
						s = strconv.Itoa(userAuth.UserId)
						break
					case "{$domain}":
						s = userAuth.Domain
						break
					case "{$code}":
						s = userAuth.Code
						break
					case "{$ip}":
						s = userAuth.Ip
						break
					case "{$version}":
						s = authVersion.Version
						break
					case "{$version_number}":
						s = strconv.Itoa(authVersion.VersionNumber)
						break
					default:
						s = ""
					}
					if s != "" && k1 != "" {
						text = strings.Replace(text, k1, s, -1)
					}
				}
				_ = os.WriteFile(folder+k, []byte(text), 0777)
			}
		}
	}
	//打包文件
	pkgFile := "./temp/" + userAuth.Code + "_" + strconv.FormatInt(str.NewStrTool().RandInt64(100000, 999999), 10) + ".zip"
	if p.zip(folder, pkgFile) != nil {
		return "", errors.New("构建源文件失败")
	}
	//重命名文件
	name := p.Md5(pkgFile) + "_" + strconv.Itoa(authVersion.VersionNumber) + ".zip"
	pkgFile2 := "./temp/" + name
	if os.Rename(pkgFile, pkgFile2) != nil {
		return "", errors.New("构建源文件失败")
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(pkgFile2)
	//上传至oss
	bytes, e := os.ReadFile(pkgFile2)
	if e != nil {
		return "", errors.New("生成源文件失败")
	}
	fileUrl := rpc.Upload.Temp(name, bytes)
	if fileUrl == "" {
		return "", errors.New("上传程序包失败")
	}
	return fileUrl, nil
}

// zip 压缩文件
func (*pkg) zip(srcFile string, destZip string) error {
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	err = filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == srcFile {
			if info.IsDir() {
				return nil
			}
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		path2 := strings.Replace(path, "\\", "/", -1)
		path2 = strings.Replace("./"+path2, srcFile, "", 1)
		header.Name = strings.TrimPrefix(path2, string(os.PathSeparator))
		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
	return err
}

// unZip 解压文件
func (*pkg) unZip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	for _, f := range zipReader.File {
		fPath := ""
		name := f.Name
		if find := strings.Contains(name, ":"); find {
			name = name[strings.Index(name, ":")+2:]
			fPath = filepath.Join(destDir, name)
		} else {
			fPath = filepath.Join(destDir, name)
		}
		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fPath, os.ModePerm)
			if err != nil {
				continue
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()
			outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (*pkg) download(u string) (string, error) {
	uri, err := url.ParseRequestURI(u)
	if err != nil {
		return "", errors.New("网址格式错误")
	}
	client := &http.Client{
		Timeout: 180 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
					return nil
				},
			},
		},
	}
	req, _ := http.NewRequest("GET", u, strings.NewReader(""))
	response, err := client.Do(req)
	if err != nil {
		return "", errors.New("下载文件错误")
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	filePath := "./temp/"
	filename := strconv.FormatInt(str.NewStrTool().RandInt64(100000, 999999), 10) + "_" + path.Base(uri.Path)
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0777)
		if err != nil {
			return "", errors.New("创建文件夹错误")
		}
		err = os.Chmod(filePath, 0777)
		if err != nil {
			return "", errors.New("修改文件夹权限错误")
		}
	}
	file, err := os.OpenFile(filePath+filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return "", errors.New("创建文件错误")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)
	writer := bufio.NewWriter(file)
	reader := bufio.NewReaderSize(response.Body, 52*1024)
	_, err = io.Copy(writer, reader)
	if err != nil {
		return "", errors.New("获取源文件错误")
	}
	return filePath + filename, nil
}

func (*pkg) Md5(path string) string {
	var fileChunk uint64 = 10485760
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	// calculate the file size
	info, _ := file.Stat()
	fileSize := info.Size()
	blocks := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
	h := md5.New()
	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(float64(fileChunk), float64(fileSize-int64(i*fileChunk))))
		buf := make([]byte, blockSize)
		_, err := file.Read(buf)
		if err != nil {
			return ""
		}
		_, err = io.WriteString(h, string(buf)) // append into the hash
		if err != nil {
			return ""
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}
