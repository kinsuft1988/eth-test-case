package main

import (
	"fmt"
	"github.com/bndr/gojenkins"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"bitbucket.org/oudmondev/ethereum-test/report"
)

var (
	baseUrl                   = "http://192.168.31.246:3001"
	addTxUrl                  = baseUrl + "/addOneTx"
	unlockUrl                 = baseUrl + "/unlock"
	getTxCountsUrl            = baseUrl + "/getTxCounts"
	getAvgBlockUrl            = baseUrl + "/getArverageBlockTime"
	getBlockNumberUrl         = baseUrl + "/getBlockNumber"
	jenkinsUrl                = "http://192.168.31.246:8080/"
	jenkinsUser               = "blockcloud"
	jenkinsPassword           = "blockcloud2018"
	jenkinsTestNetJob         = "eth-test-net"
	branches                  = []string{"bloc-test-hard-5", "bloc-test-hard-10"}
	nodeNumbers               = []int{1, 1}
	testTxTotalNumber         = int64(40000)
	testTxToleranceNumber     = int64(1000)
	oneCaseWaitTimeForStatics = 10
	oneCaseWaitTime           = 20
	resultMsgArray            = []string{}
)

func main() {

	for index, branch := range branches {

		//Env config start
		println(index)
		println(branch)
		startTestNet(branch, nodeNumbers[index])

		time.Sleep(time.Minute * 2)

		for {

			number, _ := getBlockNumber()

			if number > 0 {
				break
			}

			time.Sleep(time.Second * 5)

		}
		//Env config end

		unlock()

		time.Sleep(time.Second * 3)

		initCount, _ := getTxCounts()

		fmt.Printf("initCount : %d", initCount)

		startTime := time.Now().Unix()


		for i := int64(0); i < testTxTotalNumber; i++ {

			time.Sleep(time.Millisecond * time.Duration(10))

			go addTx()

		}

		for {
			time.Sleep(time.Second * 3)
			endCount, _ := getTxCounts()

			fmt.Printf("endCount : %d", endCount)

			if initCount + testTxTotalNumber - endCount  < testTxToleranceNumber{
				endTime := time.Now().Unix()
				tps := (endCount - initCount)/(endTime-startTime)

				avgTime, _ := getAvgBlockTime()

				resultMsg := fmt.Sprintf("最终的tps结果为%d ,环境 branch:%s nodeNumber:%d,avgTime:%f\n", tps, branch, nodeNumbers[index], avgTime)
				resultMsg += fmt.Sprintf("testTxTotalNumber: %d \n", testTxTotalNumber)
				resultMsg += fmt.Sprintf("testTxToleranceNumber: %d \n", testTxToleranceNumber)
				resultMsg += fmt.Sprintf("oneCaseWaitTimeForStatics: %d \n", oneCaseWaitTimeForStatics)
				resultMsg += fmt.Sprintf("oneCaseWaitTime: %d \n", oneCaseWaitTime)
				report := blocReport.Report{}
				report.SendMail(resultMsg)

				resultMsgArray = append(resultMsgArray,resultMsg)

				for _,msg := range resultMsgArray {
					fmt.Printf(msg)
				}

				return
			}

		}

	}

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

func getBlockNumber() (float64, error) {

	clientHttp := &http.Client{}
	reqest, _ := http.NewRequest("GET", getBlockNumberUrl, nil)

	resp, err := clientHttp.Do(reqest)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	str := string(body[:])
	fmt.Printf("str:%s",str)
	result, _ := strconv.ParseFloat(str, 64)

	return result, nil
}

func getAvgBlockTime() (float64, error) {

	clientHttp := &http.Client{}
	reqest, _ := http.NewRequest("GET", getAvgBlockUrl, nil)

	resp, err := clientHttp.Do(reqest)
	defer resp.Body.Close()

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	str := string(body[:])
	result, _ := strconv.ParseFloat(str, 64)

	return result, nil
}

func startTestNet(branch string, nodeNumber int) (uint64, error) {

	jenkins := gojenkins.CreateJenkins(nil, jenkinsUrl, jenkinsUser, jenkinsPassword)
	_, err := jenkins.Init()

	if err != nil {
		fmt.Printf("error : %s", err)
		return 0, err
	}

	params := make(map[string]string)
	params["branch"] = branch
	//params["nodeNumber"] = nodeNumber

	jenkins.BuildJob(jenkinsTestNetJob, params)

	return 0, nil
}
