package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	data := []int{0, 1, 2, 3, 4, 52} // вводим числа, хэш которых будем считать

	results := make(chan string, len(data))

	for _, v := range data {
		wg.Add(1)
		go proccess(v, &wg, results) // Запускаем горутины с вычислениями
	}

	wg.Wait()
	close(results)

	res := ""
	for result := range results {
		res += result
	}
	res = res[:len(res)-1] // Форматируем вывод

	fmt.Println(res)
	fmt.Println("Duration:", time.Since(start)) // Считаем время работы функции
}


func proccess(data int, wg *sync.WaitGroup, results chan string) {
	defer wg.Done()
	dataStr := strconv.Itoa(data)

	crcData := DataSignerCrc323(dataStr)
	md5Hash := DataSignerMd51(dataStr)
	crcMd5Hash := DataSignerCrc323(md5Hash)

	combined := crcData + "~" + crcMd5Hash

	var innerWg sync.WaitGroup
	res := make([]string, 6)
	for i := 0; i < 6; i++ {
		innerWg.Add(1)
		go func(i int) {
			defer innerWg.Done()
			crc := DataSignerCrc323(strconv.Itoa(i) + combined)
			res[i] = crc
		}(i)
	}
	innerWg.Wait()


	combinedRes := strings.Join(res, "")
	results <- combinedRes + "_"
}

func DataSignerCrc323(data string) string {
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second) // Симуляция долгих вычислений (1 сек)
	return dataHash
}

func DataSignerMd51(data string) string {
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond) // Симуляция вычислений (10 мс)
	return dataHash
}
