package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/wolf88844/godemo/m3u8/tool"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	CryptMethodAES  CryptMethod = "AES-128"
	CryptMethodNone CryptMethod = "NONE"
)

var lineParameterPattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)

type CryptMethod string

type M3u8 struct {
	Segments           []*Segment
	MasterPlaylistURIS []string
}

type Segment struct {
	URI string
	Key *Key
}

type Key struct {
	URI    string
	IV     string
	key    string
	Method CryptMethod
}
type Result struct {
	URL  *url.URL
	M3u8 *M3u8
	Keys map[*Key]string
}

//解析参数
func parseParameters(line string) map[string]string {
	r := lineParameterPattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}

//解析地址
func parseLines(lines []string) (*M3u8, error) {
	var (
		i       = 0
		lineLne = len(lines)
		m3u8    = &M3u8{}
		key     *Key
		seg     *Segment
	)
	for ; i < lineLne; i++ {
		line := strings.TrimSpace(lines[i])
		if i == 0 {
			if "#EXTM3U" != line {
				return nil, fmt.Errorf("invalid m3u8,missing #EXTM3U in line 1")
			}
			continue
		}
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "#EXT-X-STREAM-INF"):
			i++
			m3u8.MasterPlaylistURIS = append(m3u8.MasterPlaylistURIS, lines[i])
			continue
		case !strings.HasPrefix(line, "#"):
			seg = new(Segment)
			seg.URI = line
			m3u8.Segments = append(m3u8.Segments, seg)
			seg.Key = key
			continue
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			params := parseParameters(line)
			if len(params) == 0 {
				return nil, fmt.Errorf("invalid EXT-X-KEY: %s,line: %d", line, i+1)
			}
			key = new(Key)
			method := CryptMethod(params["METHOD"])
			if method != "" && method != CryptMethodAES && method != CryptMethodNone {
				return nil, fmt.Errorf("invalid EXT-X-KEY method: %s,line:%d", method, i+1)
			}
			key.Method = method
			key.URI = params["URI"]
			key.IV = params["IV"]
		default:
			continue
		}

	}
	return m3u8, nil
}

func fromURL(link string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()
	body, err := tool.Get(link)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	defer body.Close()
	s := bufio.NewScanner(body)
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	m3u8, err := parseLines(lines)
	if err != nil {
		return nil, err
	}
	if m3u8.MasterPlaylistURIS != nil {
		return fromURL(tool.ResolveURL(u, m3u8.MasterPlaylistURIS[0]))
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any segment")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[*Key]string),
	}

	for _, seg := range m3u8.Segments {
		switch {
		case seg.Key == nil || seg.Key.Method == "" || seg.Key.Method == CryptMethodNone:
			continue
		case seg.Key.Method == CryptMethodAES:
			if _, ok := result.Keys[seg.Key]; ok {
				continue
			}
			keyUrl := seg.Key.URI
			keyUrl = tool.ResolveURL(u, keyUrl)
			resp, err := tool.Get(keyUrl)
			if err != nil {
				return nil, fmt.Errorf("extract key failed: %s", err.Error())
			}
			keyByte, err := ioutil.ReadAll(resp)
			_ = resp.Close()
			if err != nil {
				return nil, err
			}
			fmt.Println("decryption key: ", string(keyByte))
			result.Keys[seg.Key] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", seg.Key.Method)
		}
	}
	return result, nil

}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
		}
	}()
	//地址
	m3u8URL := "https://dz3vqpme6j4gm.cloudfront.net/20200327/nhdtb-348/index.m3u8"
	u, err := url.Parse(m3u8URL)
	split := strings.Split(u.Path, "/")
	name := split[2]
	fmt.Println(name)
	result, err := fromURL(m3u8URL)
	if err != nil {
		panic(err)
	}
	storeFolder := "D://ts"
	if err := os.MkdirAll(storeFolder, 0777); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	limitChan := make(chan byte, 64)
	num := len(result.M3u8.Segments)
	for idx, seg := range result.M3u8.Segments {
		wg.Add(1)
		go func(i int, s *Segment) {
			defer func() {
				wg.Done()
				<-limitChan
			}()

			tsFile := filepath.Join(storeFolder, name+"_"+strconv.Itoa(i)+".ts")
			_, err := os.Stat(tsFile)
			if err == nil {
				fmt.Println("%s已存在", tsFile)
				num = num - 1
				return
			}
			tsFileTmpPath := tsFile + "_tmp"
			tsFileTmp, err := os.Create(tsFileTmpPath)
			if err != nil {
				fmt.Printf("Create TS file failed: %s\n", err.Error())
				return
			}
			defer tsFileTmp.Close()

			fullURL := tool.ResolveURL(result.URL, s.URI)
			body, err := tool.Get(fullURL)
			if err != nil {
				fmt.Printf("Download failed [%s] %s\n", err.Error(), fullURL)
			}
			defer body.Close()

			bytes, err := ioutil.ReadAll(body)
			if err != nil {
				fmt.Printf("Read TS file failed: %s\n", err.Error())
				return
			}
			if s.Key != nil {
				key := result.Keys[s.Key]
				if key != "" {
					bytes, err = tool.AES128Decrypt(bytes, []byte(key), []byte(s.Key.IV))
					if err != nil {
						fmt.Printf("decryt TS failed: %s\n", err.Error())
					}
				}
			}
			syncByte := uint8(71)
			bLen := len(bytes)
			for j := 0; j < bLen; j++ {
				if bytes[j] == syncByte {
					bytes = bytes[j:]
					break
				}
			}
			if _, err := tsFileTmp.Write(bytes); err != nil {
				fmt.Printf("Save TS file failed:%s\n", err.Error())
				return
			}
			_ = tsFileTmp.Close()
			if err = os.Rename(tsFileTmpPath, tsFile); err != nil {
				fmt.Printf("Rename TS file failed: %s\n", err.Error())
				return
			}
			num = num - 1
			fmt.Printf("下载成功：%s,还剩%d个\n", fullURL, num)

		}(idx, seg)
		limitChan <- 1
	}
	wg.Wait()
	mainFile, err := os.Create(filepath.Join(storeFolder, name+".mp4"))
	if err != nil {
		panic(err)
	}
	defer mainFile.Close()

	for i := 0; i < len(result.M3u8.Segments); i++ {
		path := filepath.Join(storeFolder, name+"_"+strconv.Itoa(i)+".ts")
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if _, err := mainFile.Write(bytes); err != nil {
			fmt.Println(err.Error())
			continue
		}
		err = os.Remove(path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}
	_ = mainFile.Sync()
	fmt.Println("下载完成")
}
