package main

import (
    "fmt"			//fmt包是输入输出格式化包
    "io/ioutil"		//ioutil包实现了一些I/O实用程序功能。
    "regexp"		//regexp包实现了正则表达式搜索。
)

var (
    NWORDS map[string]int
)

const (
	//字母表
    alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func words(text string) []string {
	regex, _ := regexp.Compile("[a-z]+")
    //编译将解析正则表达式，如果成功，则返回可用于与文本匹配的Regexp对象。
	return regex.FindAllString(text, -1)
	// regexp.Regexp.FindAllString()
	//它返回表达式的所有连续匹配的一部分，如包注释中的“全部”描述所定义。
}

func train(features []string) map[string]int {
	//用来建立一个"字典"结构。
	//文本库的每一个词，都是这个"字典"的键；
	//它们所对应的值，就是这个词在文本库的出现频率。
	result := make(map[string]int)
	//每一个词的默认出现频率为1。这是针对那些没有出现在文本库的词。
	//如果一个词没有在文本库出现，我们并不能认定它就是一个不存在的词，
	//因此将每个词出现的默认频率设为1。以后每出现一次，频率就增加1。
    for i := range features {
        _, isexist := result[features[i]]
        if !isexist {
			//如果词不存在，将频率置为1
            result[features[i]] = 1
        } else {
			//出现一次，频率 += 1
            result[features[i]] += 1
        }
    }
    fmt.Println(result);//打印出【单词-频率】表
    return result
}

func edit1(word string) []string {
	//定义edits1()函数，用来生成所有与输入参数word的"编辑距离"为1的词。
    type tuple struct{ a, b string }
	//splits：将word依次按照每一位分割成前后两半。
	//比如，'abc'会被分割成 [('', 'abc'), ('a', 'bc'), ('ab', 'c'), ('abc', '')] 。
	var splits []tuple
    for i := 0; i < len(word)+1; i++ {
        splits = append(splits, tuple{word[:i], word[i:]})
	}
	//deletes：依次删除word的每一位后、所形成的所有新词。
	//比如，'abc'对应的deletes就是 ['bc', 'ac', 'ab'] 。
    var deletes []string
    for _, t := range splits {
        if len(t.b) > 0 {
            deletes = append(deletes, t.a+t.b[1:])
        }
	}
	//transposes：依次交换word的邻近两位，所形成的所有新词。
	//比如，'abc'对应的transposes就是 ['bac', 'acb'] 。
    var transposes []string
    for _, t := range splits {
        if len(t.b) > 1 {
            transposes = append(transposes, t.a+string(t.b[1])+string(t.b[0])+t.b[2:])
        }
	}
	//replaces：将word的每一位依次替换成其他25个字母，所形成的所有新词。
	//比如，'abc'对应的replaces就是 ['abc', 'bbc', 'cbc', ... , 'abx', ' aby', 'abz' ]
	//一共包含78个词（26 × 3）。
    var replaces []string
    for _, c := range alphabet {
        for _, t := range splits {
            if len(t.b) > 0 {
                replaces = append(replaces, t.a+string(c)+t.b[1:])
            }
        }
	}
	//inserts：在word的邻近两位之间依次插入一个字母，所形成的所有新词。
	//比如，'abc' 对应的inserts就是['aabc', 'babc', 'cabc', ..., 'abcx', 'abcy', 'abcz']
	//一共包含104个词（26 × 4）。
    var inserts []string
    for _, c := range alphabet {
        for _, t := range splits {
            inserts = append(inserts, t.a+string(c)+t.b)
        }
    }
    //concat this slice 
    deletes = append(deletes, transposes...)
    deletes = append(deletes, replaces...)
    deletes = append(deletes, inserts...)
    fmt.Println(deletes);
	return set(deletes)
	//最后，edit1()返回deletes、transposes、replaces、inserts的合集，
	//这就是与word"编辑距离"等于1的所有词。对于一个n位的词，会返回54n+25个词。
}

func known_edits2(word string) []string {
	//定义edit2()函数，用来生成所有与word的"编辑距离"为2的词语。
	//但是这样的话，会返回一个 (54n+25) * (54n+25) 的数组，实在是太大了。
	//因此，我将edit2()改为known_edits2()函数，将返回的词限定为在文本库中出现过的词。

	//_,ok是go语言中的Comma-ok断言格式
	//element必须是接口类型的变量，T是普通类型。
	//如果断言失败，ok为false，否则ok为true并且value为变量的值
    var arr []string
    for _, e1 := range edit1(word) {
        for _, e2 := range edit1(e1) {
            if _, ok := NWORDS[e2]; ok {
                arr = append(arr, e2)
            }
        }
    }
    return set(arr)
}

//定义known()，correct()函数，用来从所有备选的词中，选出用户最可能想要拼写的词。
func known(words []string) []string {
    var knows []string
    for _, value := range words {
        if _, ok := NWORDS[value]; ok {
            knows = append(knows, value)
        }
    }
    return knows
}

func appendIfMissing(slice []string, i string) []string {
	//如果单词找不到，添加到slice切片内
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

func set(arr []string) []string {
    var result []string
    for _, ele := range arr {
        result = appendIfMissing(result, ele)
    }
    return result
}

func correct(word string) string {
	//如果word是文本库现有的词，说明该词拼写正确，直接返回这个词
	//如果word不是现有的词，则返回"编辑距离"为1的词之中，在文本库出现频率最高的那个词
	//如果"编辑距离"为1的词，都不是文本库现有的词，则返回"编辑距离"为2的词中，出现频率最高的那个词
	//如果上述三条规则，都无法得到结果，则直接返回word。
    candidates := known([]string{word})
    if len(candidates) <= 0 {
        candidates = known(edit1(word))
        if len(candidates) <= 0 {
            candidates = known(known_edits2(word))
        }
    }
    return max(candidates, NWORDS)
}

func max(arr []string, dict map[string]int) string {
	//排序找到最大概率
    flag := 0
    index := 0
    for ix, value := range arr {
        if v, ok := dict[value]; ok {
            if v > flag {
                flag = v
                index = ix
            }
        }
    }
    fmt.Println("最接近的词,频率为：")
    fmt.Println(arr[index],flag)

    return arr[index]
}

func main() {
	var word string                         //定义一个word用来供用户输入单词                 
	buf, _ := ioutil.ReadFile("big.txt")	//打开大词典文件
	NWORDS = train(words(string(buf)))		//使用words()和train()函数，生成"词频字典"，放入变量NWORDS。
	//fmt.Println(NWORDS);					//函数主体流程开始
	for word != "#" {
        fmt.Println("please input words: (press # to exit)")
        /*if word == "#" {
			break
		}*/
		fmt.Scanf("%s",&word)
        fmt.Println("input:", word, "except word:", correct(word))
        
	}
}