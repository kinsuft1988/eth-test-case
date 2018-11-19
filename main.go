package main

import (
	"bitbucket.org/oudmondev/ethereum-test/report"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	baseUrl        = "http://192.168.31.246:3001"
	addTxUrl       = baseUrl + "/addOneTx"
	unlockUrl      = baseUrl + "/unlock"
	getTxCountsUrl = baseUrl + "/getTxCounts"
	testNetUrl     = "http://192.168.31.246:8080/job/eth-test-net/build?token=eth-test-net?branch="
	branches       = []string{"bloc-test-hard-10"}
)

func main() {

	for _, branch := range branches {

		//Env config start
		startTestNet(branch)

		time.Sleep(time.Minute * 10)
		//Env config end

		unlock()

		time.Sleep(time.Second * 3)

		initCount, _ := getTxCounts()

		fmt.Printf("initCount : %d", initCount)

		startTime := time.Now().Unix()

		for initTps := 10; initTps < 100; initTps += 5 {

			for i := 0; i < 4000; i++ {

				time.Sleep(time.Millisecond * time.Duration(1000/initTps))

				go addTx()

			}

			time.Sleep(time.Second * 10)

			endTime := time.Now().Unix()

			fmt.Printf("div time : %d", endTime-startTime)

			resultCount, _ := getTxCounts()
			totalSendCount := initCount + 4000

			if totalSendCount-resultCount < 400 {
				fmt.Printf("tps: %d 通过了验证", initTps)
				time.Sleep(time.Second * 20)
			} else {
				fmt.Printf("tps: %d 没有通过验证，还有%d的交易没有处理", initTps, totalSendCount-resultCount)

				if initTps != 10 {
					fmt.Printf("最终的tps结果为%d", initTps-5)

					resultMsg := fmt.Sprintf("最终的tps结果为%d", initTps-5)
					report := blocReport.Report{}
					report.SendMail(resultMsg)

				} else {
					fmt.Printf("最终的tps结果低于10")

					resultMsg := fmt.Sprintf("最终的tps结果低于10")
					report := blocReport.Report{}
					report.SendMail(resultMsg)
				}

			}
		}

	}

}

func test() {
	report := blocReport.Report{}
	report.SendMail("test")
}

func addTx() (uint64, error) {

	resp, err := http.Get(addTxUrl)

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}
	defer resp.Body.Close()

	return 0, nil
}

func unlock() (uint64, error) {

	resp, err := http.Get(unlockUrl)

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}
	defer resp.Body.Close()

	return 0, nil
}

func getTxCounts() (int64, error) {

	clientHttp := &http.Client{}
	reqest, _ := http.NewRequest("GET", getTxCountsUrl, nil)

	resp, err := clientHttp.Do(reqest)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	str := string(body[:])
	result, _ := strconv.ParseInt(str, 10, 64)

	return result, nil
}

func startTestNet(branch string) (uint64, error) {

	clientHttp := &http.Client{}
	testNetUrl += branch
	reqest, _ := http.NewRequest("GET", testNetUrl, nil)
	reqest.SetBasicAuth("blockcloud", "blockcloud2018")
	var err error
	_, err = clientHttp.Do(reqest)

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}

	return 0, nil
}
