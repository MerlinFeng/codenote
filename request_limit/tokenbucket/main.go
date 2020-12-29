// 令牌桶
// @create_time : 2020/12/29 7:55 下午
// @author : fengqiang

package main

import (
	"fmt"
	"time"
)

type TokenBucket struct {
	rate      int // 令牌放入速度
	tokenSize int // 令牌桶的容量
	curNum    int // 当前桶中token
}

func NewTokenBucket(rate, tokenSize int) *TokenBucket {
	return &TokenBucket{
		rate:      rate,
		tokenSize: tokenSize,
	}
}

func (t *TokenBucket) pushToken() chan struct{} {
	tb := make(chan struct{}, t.tokenSize)
	ticker := time.NewTicker(time.Duration(1000) * time.Microsecond)
	//初始化token
	for i := 0; i < t.tokenSize; i++ {
		tb <- struct{}{}
	}
	t.curNum = t.tokenSize

	// 指定速率放入token
	go func() {
		for {
			for i := 0; i < t.rate; i++ {
				tb <- struct{}{}
			}
			t.curNum += t.rate
			if t.curNum > t.tokenSize {
				t.curNum = t.tokenSize
			}
			<-ticker.C
		}
	}()
	return tb
}

// 取出token
func (t *TokenBucket) popToken(bucket chan struct{}, n int) {
	for i := 0; i < n; i++ {
		_, ok := <-bucket
		if ok {
			t.curNum -= 1
			fmt.Println("get  token  success")
		} else {
			fmt.Println("get  token  fail")
		}
	}
}

func main() {
	tokenBucket := NewTokenBucket(10, 20)
	tb := tokenBucket.pushToken()
	reqLen := 100
	for i := 0; i <= reqLen; i += tokenBucket.rate {
		fmt.Println(i)
		tokenBucket.popToken(tb, tokenBucket.rate)

	}
}
