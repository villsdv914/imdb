package imdblog


import (
"fmt"
"os"
"strconv"
"strings"
"time"
)



func createFile(filename string) (file *os.File, err error){
	pwd, _ := os.Getwd()
	if _, err := os.Stat(pwd+"/log"); os.IsNotExist(err){
		_ = os.MkdirAll(pwd+"/log", os.ModePerm)
	}
	currTime := time.Now()
	currentTime := currTime.Format("02_01_2006")
	if _, errFolder := os.Stat(pwd+"/log/"+ currentTime); os.IsNotExist(errFolder){
		_ = os.MkdirAll(pwd+"/log/"+currentTime, os.ModePerm)
	}
	baseFilename := pwd+"/log/"+currentTime+"/"+filename+".log"
	if _, errFile := os.Stat(baseFilename); os.IsNotExist(errFile){

		_, _ = os.Create(baseFilename)
		return os.OpenFile(baseFilename,os.O_WRONLY|os.O_APPEND,0600)
	}else {
		return os.OpenFile(baseFilename,os.O_WRONLY|os.O_APPEND,0600)
	}

}

func WriteFile(filename , content, process string){
	file,err :=createFile(filename)
	if err != nil{
		fmt.Print(err)
	}
	var dataStr []string
	pid := strconv.Itoa(os.Getpid())
	dataStr = append(dataStr,time.Now().Format("02-01-2006 15:04:05"))
	dataStr = append(dataStr, process, pid, content,"\n")
	result := strings.Join(dataStr,", ")

	_, err =file.WriteString(result)
	if err != nil{
		fmt.Print(err)
	}
	err = file.Sync()
	if err != nil{
		fmt.Print(err)
	}
}

