/*
@Time       : 2020/4/2 0:21
@Author     : stevinpan
@File       : Service
@Software   : GoLand
@Description: <>
*/
package service

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/reusee/mmh3"
	"strings"
)

// 计算长链接的哈希码
func getHashCode(longLink string) uint32 {
	h := mmh3.New32()
	h.Write([]byte(longLink))

	return h.Sum32()
}

// 将10进制的哈希码转换为62进制，进一步缩短字符长度
func hasCodeTransform(hashCode uint32) string {
	var chars = []string{
		"0","1","2","3","4","5","6","7","8","9",
		"a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z",
		"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z",
	}

	var arr = make([]string, 0)
	// 哈希码有可能就是0
	if hashCode == 0 {
		arr = append(arr, chars[0])
	} else {
		for hashCode > 0 {
			var remainder = hashCode % 62	// 余数
			hashCode = hashCode / 62		// 整数
			arr = append(arr, chars[remainder])
		}
	}

	// 进制转换后得到的数组元素是倒序的，可以选择反转过来；也可以不转

	// 输出字符串
	return strings.Join(arr, "")
}

// 查询数据库，检查是否存在冲突
func isExist(shortLink, longLink string) (bool, error) {
	var result bool
	var err error

	// 连接redis
	client, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Printf("Connect to redis, err=%+v\n", err)
		return result, err
	}
	defer client.Close()

	// 查询Redis中是否存在短连接
	keyExist, err := redis.Bool(client.Do("EXISTS", shortLink))
	if err == nil {
		if keyExist {	// 短链接冲突
			rlongLink, err := redis.String(client.Do("GET", shortLink))
			if err == nil && rlongLink != longLink {	// 不是同一个长链接
				result = true
			} else {	// 是同一个长链接
				result = false
			}
		} else {	// 短链接不冲突
			result = false
		}
	}

	return result, err
}

// 写数据库：将<shortLink, longLink对应关系写入数据库>
func setToDB(shortLink, longLink string) error {
	// 连接数据库
	client, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	defer client.Close()

	// 将短连接写入数据库
	_, err = client.Do("SET", shortLink, longLink)
	return err
}

func Process(longLink, duplication string) (string, error) {
	// 对长链接进行Hash
	var hashCode = getHashCode(longLink)
	var shortLink = hasCodeTransform(hashCode)
	var err error
	fmt.Printf("shortLink=%s, hashCode=%d, longLink=%s\n", shortLink, hashCode, longLink)

	// 检查是否有hash冲突
	result, err := isExist(shortLink, longLink)
	if err == nil {
		if result {
			duplication += "aaa"
			shortLink, err = Process(longLink, duplication)
		} else {
			err = setToDB(shortLink, longLink)
		}
	}

	return shortLink, err
}
