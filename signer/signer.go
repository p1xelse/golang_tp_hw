package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	in := make(chan interface{})
	var out chan interface{}

	for _, job := range jobs {
		out = make(chan interface{})
		wg.Add(1)
		go runJob(job, in, out, wg)
		in = out
	}

	wg.Wait()
}

func runJob(j job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	j(in, out)
	close(out)
}

func asyncCrc32(data string, result *string, wg *sync.WaitGroup) {
	defer wg.Done()
	*result = DataSignerCrc32(data)
}

func SingleHash(input, output chan interface{}) {
	wgOuter := &sync.WaitGroup{}
	for val := range input {
		data := strconv.Itoa(val.(int))
		md5res := DataSignerMd5(data)

		wgOuter.Add(1)
		go func(md5res string, data string, out chan interface{}) {
			defer wgOuter.Done()

			var crc32data, crc32md5data string

			wgInner := &sync.WaitGroup{}
			wgInner.Add(2)
			go asyncCrc32(data, &crc32data, wgInner)
			go asyncCrc32(md5res, &crc32md5data, wgInner)
			wgInner.Wait()

			out <- crc32data + "~" + crc32md5data
		}(md5res, data, output)

	}
	wgOuter.Wait()
}

func MultiHash(input, output chan interface{}) {
	wgOuter := &sync.WaitGroup{}

	for val := range input {
		strdata := val.(string)

		wgOuter.Add(1)
		go func(data string) {
			defer wgOuter.Done()

			result := make([]string, 6)
			wgInner := &sync.WaitGroup{}
			mutex := &sync.Mutex{}

			for th := 0; th < 6; th++ {
				wgInner.Add(1)
				go func(th int) {
					defer wgInner.Done()
					hash := DataSignerCrc32(strconv.Itoa(th) + strdata)

					mutex.Lock()
					result[th] = hash
					mutex.Unlock()
				}(th)
			}

			wgInner.Wait()
			output <- strings.Join(result, "")
		}(strdata)
	}
	wgOuter.Wait()
}

func CombineResults(input, output chan interface{}) {
	var res []string

	for val := range input {
		resElem := val.(string)
		res = append(res, resElem)
	}

	sort.Strings(res)
	output <- strings.Join(res, "_")
}
