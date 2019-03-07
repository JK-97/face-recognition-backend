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

	"gitlab.jiangxingai.com/luyor/face-recognition-backend/log"
)

var uploadMap= map[string]int{}
var uploadMutex = &sync.Mutex{}

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


// ListImgs 获取一个路径下(ceph/local) 上的所有图片名字列表
// 并且异步检查{文件缺少、新增} local 和 remote， 使得两边状态一致
// 如果当前路径存在上传操作，不同步
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

    go func() {
        var remotepath = fmt.Sprintf("/data/edgebox/remote/%s", path)
        uploadMutex.Lock()
        _, uploading := uploadMap[path]
        uploadMutex.Unlock()

        if !uploading {
            remotefiles, err := ioutil.ReadDir(remotepath)
            if err != nil {
                return
            }

            var remoteStats map[string]string
            for _, rf := range remotefiles {
                remotefile := fmt.Sprintf("%s/%s", remotepath, rf.Name())

                fileMeta := strings.Split(".", rf.Name())
                if len(fileMeta) != 3 {
                    continue
                }

                if err != nil {
                    log.Warn("Sync remote files Failed: ", err)
                    return
                }
                localfile := fmt.Sprintf("%s/%s", localpath, rf.Name())
                if err != nil && os.IsNotExist(err) {
                    copyFile(remotefile, localfile)
                } else if err != nil {
                    log.Warn("Sync remote files stats Failed: ", err)
                }

                index := fileMeta[0]
                checkSum := fileMeta[2]
                remoteStats[index] = checkSum
            }

            for _, lf := range localfiles {
                fileMeta := strings.Split(".", lf.Name())
                if len(fileMeta) != 3 {
                    continue
                }

                index := fileMeta[0]
                checkSum := fileMeta[2]

                _, exists := remoteStats[index]
                if !exists || checkSum != remoteStats[index] {
                    os.Remove(fmt.Sprintf("%s/%s", localpath, lf.Name()))
                }
            }
        }
    } ()

    return imgs, nil
}

// SaveImg saves a image to local dir
// TODO: save file in database to survive from app container removal
// 上传文件过程中， 需要将local path标记为只读, 不能再次开启从远端同步过程
func SaveImg(img *image.Image, path string, image_seq int) string {

    local := createFileName(path, string(image_seq), "local")
    outLocalFile, _ := os.Create(local)
    jpeg.Encode(outLocalFile, *img, nil)
    outLocalFile.Close()

    uid, _ := uuid.NewUUID()
    checkSuffix := uid.String()
    os.Rename(local, fmt.Sprintf("%s.%s", local, checkSuffix))

    go func() {
        uploadMutex.Lock()
        _, exists := uploadMap[path]
        if !exists {
            uploadMap[path] = 1
        } else {
            uploadMap[path] += 1
        }
        uploadMutex.Unlock()

        remote := createFileName(path, string(image_seq), "remote")
        outRemoteFile, _ := os.Create(fmt.Sprintf("%s.%s", remote, checkSuffix))
        jpeg.Encode(outRemoteFile, *img, nil)
        outRemoteFile.Close()

        uploadMutex.Lock()
        uploadMap[path] -= 1
        if uploadMap[path] == 0 {
            delete (uploadMap, path)
        }
        uploadMutex.Unlock()
    } ()

    return fmt.Sprint("%s.jpeg.%s", image_seq, checkSuffix)
}

// RemoveImg 删除远端图像
func RemoveImg(path string, image_id string) error {
    var remotepath = fmt.Sprintf("/data/edgebox/remote/%s", path)
    var localpath = fmt.Sprintf("/data/edgebox/local/%s", path)

    if image_id == "" {
        remotefiles, err := ioutil.ReadDir(remotepath)
        if err != nil {
            return err
        }
        for _, rf := range remotefiles {
            os.Remove(fmt.Sprintf("%s/%s", remotepath, rf.Name()))
            os.Remove(fmt.Sprintf("%s/%s", localpath, rf.Name()))
        }
    } else {
        os.Remove(fmt.Sprintf("%s/%s", remotepath, image_id))
        os.Remove(fmt.Sprintf("%s/%s", localpath, image_id))
    }
    return nil
}


// GetImg read image from file (local or ceph)
func GetImg(fileName string) (*image.Image, error) {
    localFile := fmt.Sprint("/data/edgebox/local/%s", fileName)
    existingImageFile, err := os.Open(localFile)
    if err == nil {
        defer existingImageFile.Close()
        imageData, _, err := image.Decode(existingImageFile)
        if err != nil {
            return nil, err
        }
        return &imageData, nil
    }

    remoteFile := fmt.Sprint("/data/edgebox/remote/%s", fileName)
    existingImageFile, err = os.Open(remoteFile)
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
