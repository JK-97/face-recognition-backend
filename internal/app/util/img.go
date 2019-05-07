package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io/ioutil"
	"strings"
    "sync"
    "fmt"
    "os"
    "io"
    "bufio"
    "github.com/google/uuid"
    "strconv"

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

func createFileName(path string, uid string, mode string) string {
    cpath := fmt.Sprintf("/data/edgebox/%s/%s", mode, path)
    if _, err := os.Stat(cpath); os.IsNotExist(err) {
        os.MkdirAll(cpath, 0755)
    }
    return fmt.Sprintf("%s/%s.jpeg", cpath, uid)
}

func copyFile(src, dst string) error {
    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destination.Close()
    _, err = io.Copy(destination, source)
    return err
}


// ListImgs 获取一个路径下上的所有图片名字列表
func ListImgs(path string) ([]string, error) {
    var imgs []string
    var localpath = fmt.Sprintf("/data/edgebox/local/%s", path)

    localfiles, err := ioutil.ReadDir(localpath)
    if err != nil {
        return nil, err
    }
    for _, f := range localfiles {
        imgs = append(imgs, f.Name())
    }
    return imgs, nil
}

// SaveImg saves a image to local dir
func SaveImg(img *image.Image, path string, image_seq int) string {
    uid, _ := uuid.NewUUID()
    checkSuffix := uid.String()
    filename := fmt.Sprintf("%s.jpeg.%s", strconv.Itoa(image_seq), checkSuffix)
    local := createFileName(path, filename, "local")
    outLocalFile, _ := os.Create(local)
    jpeg.Encode(outLocalFile, *img, nil)
    outLocalFile.Close()
    return filename
}

// RemoveImg 删除远端图像
func RemoveImg(path string, image_id string) error {
    var localpath = fmt.Sprintf("/data/edgebox/local/%s", path)
    if image_id == "" {
        remotefiles, err := ioutil.ReadDir(localpath)
        if err != nil {
            return err
        }
        for _, rf := range remotefiles {
            os.Remove(fmt.Sprintf("%s/%s", localpath, rf.Name()))
        }
    } else {
        os.Remove(fmt.Sprintf("%s/%s", localpath, image_id))
    }
    return nil
}

// GetImg read image from file (local or ceph)
func GetImg(fileName string) (*image.Image, error) {
    localFile := fmt.Sprintf("/data/edgebox/local/%s", fileName)
    existingImageFile, err := os.Open(localFile)
    if err == nil {
        defer existingImageFile.Close()
        imageData, _, err := image.Decode(existingImageFile)
        if err != nil {
            return nil, err
        }
        return &imageData, nil
    }
    return nil, err
}

// Base64Str2Img converts base64 image to image.Image
func Base64Str2Img(str string) (*image.Image, error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(str))
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	return &img, nil
}

// Img2Base64Str converts images to base64 string
func Img2Base64Str(img *image.Image)(string, error) {
    var b bytes.Buffer
    w := bufio.NewWriter(&b)
    jpeg.Encode(w, *img, nil)
    encoded := base64.StdEncoding.EncodeToString(b.Bytes())
    return encoded, nil
}
