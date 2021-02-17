package service

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// WaveJSON ...
type WaveJSON struct {
	SampleRate int   `json:"sample_rate"`
	Data       []int `json:"data"`
}

// MP3ToJSON ...
func MP3ToJSON(filepath string) (duration uint32, waveformURL string, err error) {
	ossHost := os.Getenv("OSS_ENDPOINT_INTERNAL")
	ossID := os.Getenv("OSS_ID")
	ossSecret := os.Getenv("OSS_SECRET")
	bucketName := os.Getenv("OSS_BUCKET")
	waveformURL = strings.Replace(filepath, "mp3", "json", 1)
	filename := filepath[7:]
	trackID := filename[:len(filename)-4]
	dir := "temp/"
	os.MkdirAll(dir, os.ModePerm)
	waveformURL = "jsons/" + trackID + ".json"
	path, _ := os.Getwd()
	path = path + "/"
	tempMP3 := path + dir + filename
	tempJSON := path + dir + trackID + ".json"

	ossClient, err := oss.New(ossHost, ossID, ossSecret)
	if err != nil {
		return
	}

	// 获取存储空间。
	bucket, err := ossClient.Bucket(bucketName)
	if err != nil {
		return
	}

	// 下载文件到本地。
	err = bucket.GetObjectToFile(filepath, dir+filename)
	if err != nil {
		return
	}

	cmd := exec.Command("audiowaveform", "-i", tempMP3, "-o", tempJSON, "-b", "8", "--pixels-per-second", "10")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return
	}
	outStr := string(out)
	// fmt.Printf("combined out:\n%s\n", outStr)

	index := strings.Index(outStr, "Frames decoded")
	if index < 0 {
		return
	}
	outStr = outStr[index:]
	index = strings.Index(outStr, "(")
	endIndex := strings.Index(outStr, ")")
	if index < 0 || endIndex < 0 {
		return
	}
	outStr = outStr[index+1 : endIndex]

	s, err := strToSecond(outStr)
	if err != nil {
		return
	}
	duration = uint32(s)
	// fmt.Println(duration)

	data, err := ioutil.ReadFile(tempJSON)
	if err != nil {
		return
	}

	waveJSON := &WaveJSON{}
	err = json.Unmarshal(data, waveJSON)
	if err != nil {
		return
	}

	arr := []int{}
	for k, v := range waveJSON.Data {
		if (k+1)%2 == 0 {
			arr = append(arr, v)
		}
	}

	fp, err := os.Create(tempJSON)
	if err != nil {
		return
	}
	defer fp.Close()
	newJSON := &WaveJSON{}
	newJSON.SampleRate = waveJSON.SampleRate
	newJSON.Data = arr
	jsonData, err := json.Marshal(newJSON)
	if err != nil {
		return
	}
	_, err = fp.Write(jsonData)
	if err != nil {
		return
	}

	err = bucket.PutObjectFromFile(waveformURL, tempJSON)
	os.Remove(tempMP3)
	os.Remove(tempJSON)
	return
}

func strToSecond(str string) (s int, err error) {
	ss := strings.Split(str, ".")
	ts := strings.Split(ss[0], ":")
	j := 0
	for i := len(ts) - 1; i >= 0; i-- {
		var ti int64
		ti, err = strconv.ParseInt(ts[i], 10, 32)
		if err != nil {
			return
		}
		if j == 0 {
			s = s + int(ti)
		}
		if j == 1 {
			s = s + int(ti)*60
		}
		if j == 2 {
			s = s + int(ti)*60*60
		}
		j = j + 1
	}
	if len(ss) > 1 {
		var ms int64
		ms, err = strconv.ParseInt(ss[1][0:1], 10, 32)
		if err != nil {
			return
		}
		if ms >= 5 {
			s = s + 1
		}
	}
	return
}
