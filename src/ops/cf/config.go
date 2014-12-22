/**
 * config file end with enter line
 */

package cf

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
)

const lineEnd byte = '\n'

var (
	//regex = regexp.MustCompile("^(?:!#)\\s*(.+)\\s*=\\s*(.+?)\\s*$")
	regex = regexp.MustCompile("^\\s*([^#\\s]+)\\s*=\\s*([^#\\s]*)\\s*$")
)

//配置
type Config struct {
	configDict map[string]interface{}
}

// 从文件中加载配置
func NewConfig(file string) (cfg *Config, err error) {
	s := &Config{}
	_err := s.load(file)
	return s, _err
}

//从配置中读取数据
func (this *Config) Get(key string) string {
	k, e := this.configDict[key]
	if e {
		v, _ := k.(string)
		return v
	}
	return ""
}

func (this *Config) Set(key string, v interface{}) {
	if _, ok := this.configDict[key]; ok {
		panic("Key '" + key + "' is exist in config")
	}
	this.configDict[key] = v
}

func (this *Config) GetInt(key string) (int, error) {
	k, e := this.configDict[key]
	if e {
		v, _ := k.(string)
		return strconv.Atoi(v)
	}
	return 0, errors.New("no such key")
}

func (this *Config) GetFloat(key string) (float64, error) {
	k, e := this.configDict[key]
	if e {
		v, _ := k.(string)
		return strconv.ParseFloat(v, 64)
	}
	return 0, errors.New("no such key")
}

//从文件中加载配置
func (this *Config) load(file string) (err error) {
	this.configDict = make(map[string]interface{})
	//var allContent string = ""
	f, _err := os.Open(file)
	if _err != nil {
		return _err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, _err := reader.ReadString(lineEnd)
		if _err == io.EOF {
			break
		}

		if regex.Match([]byte(line)) {
			mathches := regex.FindStringSubmatch(line)
			//this.configDict[mathches[1]] = mathches[2]
			this.configDict[mathches[1]] = mathches[2]
		}
		//allContent = allContent + line + "\n"
	}
	return nil
}
