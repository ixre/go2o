package sensitive

import (
	"io/ioutil"
	"os"
	"strings"
)

// ref：https://github.com/TomatoMr/SensitiveWords
type SensitiveMap struct {
	sensitiveNode map[string]interface{}
	isEnd         bool
}

var s *SensitiveMap

func Singleton() *SensitiveMap {
	if s == nil {
		s = initDictionary("./assets/sensitive_dict.txt")
	}
	return s
}


// 初始化敏感词词典，根据DFA算法构建trie
func initDictionary(dictionaryPath string) *SensitiveMap {
	s := &SensitiveMap{
		sensitiveNode: make(map[string]interface{}),
		isEnd:         false,
	}
	dictionary := readDictionary(dictionaryPath)
	for _, words := range dictionary {
		sMapTmp := s
		w := []rune(words)
		wordsLength := len(w)
		for i := 0; i < wordsLength; i++ {
			t := string(w[i])
			isEnd := false
			//如果是敏感词的最后一个字，则确定状态
			if i == (wordsLength - 1) {
				isEnd = true
			}
			func(tx string) {
				if _, ok := sMapTmp.sensitiveNode[tx]; !ok { //如果该字在该层级索引中找不到，则创建新的层级
					sMapTemp := new(SensitiveMap)
					sMapTemp.sensitiveNode = make(map[string]interface{})
					sMapTemp.isEnd = isEnd
					sMapTmp.sensitiveNode[tx] = sMapTemp
				}
				sMapTmp = sMapTmp.sensitiveNode[tx].(*SensitiveMap) //进入下一层级
				sMapTmp.isEnd = isEnd
			}(t)
		}
	}
	return s
}

// 读取词典文件
func readDictionary(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	str, err := ioutil.ReadAll(file)
	return strings.Fields(string(str))
}

// 返回文本中的所有敏感词
// 返回值：数组，格式为“["敏感词"][敏感词在检测文本中的索引位置，敏感词长度]”
type Target struct {
	Indexes []int
	Len     int
}

// 检查是否含有敏感词，仅返回检查到的第一个敏感词,返回值：敏感词，是否含有敏感词
func (s *SensitiveMap) CheckSensitive(text string) (string, bool) {
	content := []rune(text)
	contentLength := len(content)
	result := false
	ta := ""
	for index := range content {
		sMapTmp := s
		target := ""
		in := index
		for {
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp.sensitiveNode[wo]; ok {
				if sMapTmp.sensitiveNode[wo].(*SensitiveMap).isEnd {
					result = true
					break
				}
				if in == contentLength-1 {
					break
				}
				sMapTmp = sMapTmp.sensitiveNode[wo].(*SensitiveMap) //进入下一层级
				in++
			} else {
				break
			}
		}
		if result {
			ta = target
			break
		}
	}
	return ta, result
}

func (s *SensitiveMap) FindAllSensitive(text string) map[string]*Target {
	content := []rune(text)
	contentLength := len(content)
	result := false

	ta := make(map[string]*Target)
	for index := range content {
		sMapTmp := s
		target := ""
		in := index
		result = false
		for {
			wo := string(content[in])
			target += wo
			if _, ok := sMapTmp.sensitiveNode[wo]; ok {
				if sMapTmp.sensitiveNode[wo].(*SensitiveMap).isEnd {
					result = true
					break
				}
				if in == contentLength-1 {
					break
				}
				sMapTmp = sMapTmp.sensitiveNode[wo].(*SensitiveMap) //进入下一层级
				in++
			} else {
				break
			}
		}
		if result {
			if _, targetInTa := ta[target]; targetInTa {
				ta[target].Indexes = append(ta[target].Indexes, index)
			} else {
				ta[target] = &Target{
					Indexes: []int{index},
					Len:     len([]rune(target)),
				}
			}
		}
	}
	return ta
}
