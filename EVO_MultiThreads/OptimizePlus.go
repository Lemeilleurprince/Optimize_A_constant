package main

import (
    "os"
    "fmt"
    //"time"
    "sync"
    "strconv"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/tebeka/selenium"
    "./EVO"
)

const (
    chromeDriverPath = "./chromedriver"
    port             = 8080
)
var(
    chanDatamap         chan map[string]string
    waitGroup           sync.WaitGroup
    lock                sync.Mutex
)


func main(){
    var n int
    fmt.Println("Please input from where to start,default input is 0:")
    fmt.Scanln(&n)
    chanDatamap  = make(chan map[string]string, 200)
    go GetDataMap(n)
    for i := 0; i < 10; i++ {
        waitGroup.Add(1)
        go run()
    }
    waitGroup.Wait()
}


func GetDataMap(n int){
   // Get all the rows in the Sheet1.
   DataMap := make(map[string]string, 0)
   f, err := excelize.OpenFile("Data.xlsx")
   if err != nil {
        fmt.Println(err)
        return
    }
   rows, err := f.GetRows("Sheet1")
   for i, row := range rows {
            if i == 0 {
                f.SetCellValue("Sheet1", "N1", "Optimized A_constant") 
                continue
            }else if i>=n{
                /*
                fmt.Println("Processing patient:")
                fmt.Println(i)
                */
                DataMap=map[string]string{
                            "num":strconv.Itoa(i+1),
                            "TextBoxName":row[0],
                            //"TextBoxID":row[1],
                            "txtAConstant":"118.80",
                            "txtAL":row[1],
                            "txtK1":row[3],
                            "txtK2":row[5],
                            "txtACD":row[2],
                            "txtRefraction":"0",
                            "txtLT":row[10],
                            //"txtCCT":row[9],
                            "IOL":row[11],
                            "Ref_PostOP":row[12],
                        }

                /*
                for k,v := range DataMap{
                    fmt.Println(k,v)
                }
                */
                chanDatamap <- DataMap
            }
    }
    //close(chanDatamap)
}

func run() {
    // Start a WebDriver server instance
    opts := []selenium.ServiceOption{
        //selenium.Output(os.Stderr),            // Output debug information to STDERR.
    }
    //selenium.SetDebug(true)
    selenium.SetDebug(false)
    service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
    if err != nil {
        panic(err) // panic is used only as an example and is not otherwise recommended.
    }
    defer service.Stop()

    // Connect to the WebDriver instance running locally.
    caps := selenium.Capabilities{"browserName": "chrome"}
    wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
    if err != nil {
        panic(err)
    }
    defer wd.Quit()

    // Navigate to the simple playground interface.
    // Navigate to the simple playground interface.
    err = wd.Get("https://www.evoiolcalculator.com/calculator.aspx")
    if err != nil {
      fmt.Println("get page faild", err.Error())
      //return
    }
   for DataMap := range chanDatamap{
        i, _ := strconv.Atoi(DataMap["num"])
        fmt.Printf("Processing Patient:%d\n",i-1)
        A_constant := EVO.Get_A_constant(wd, DataMap)
        A_constant_ :=strconv.FormatFloat(A_constant, 'f', 4, 64)
        
        lock.Lock()
        f, err := excelize.OpenFile("Data.xlsx")
        if err != nil {
            fmt.Println(err)
            return
        }
        f.SetCellValue("Sheet1", "N"+DataMap["num"], A_constant_) 
        fmt.Println("Optimized A_constant :" + A_constant_)
        err = f.Save()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        lock.Unlock()
        if i == len(chanDatamap){
            close(chanDatamap)
        }
    }
    wd.Quit()
    waitGroup.Done()
}
