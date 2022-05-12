package main

import (
    "os"
    "fmt"
    "time"
    "sync"
    "strconv"
    "strings"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/tebeka/selenium"
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
    chanDatamap  = make(chan map[string]string)
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
                f.SetCellValue("Sheet1", "R1", "Ref with A_constant Optimized ") 
                continue
            }else if i >=n{
                /*
                fmt.Println("Processing patient:")
                fmt.Println(i)
                */
                DataMap=map[string]string{
                            "num":strconv.Itoa(i+1),
                            "MainContent_PatientName":row[0],
                            //"MainContent_PatientNo":row[1],
                            "MainContent_Aconstant":"118.56",
                            "MainContent_Axlength":row[1],
                            "MainContent_MeasuredK1":row[3],
                            "MainContent_MeasuredK2":row[5],
                            "MainContent_OpticalACD":row[2],
                            "MainContent_Refraction":"0",
                            "MainContent_LensThickness":row[10],
                            "MainContent_WTW":row[9],
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
    err = wd.Get("https://calc.apacrs.org/barrett_universal2105/")
    if err != nil {
      fmt.Println("get page faild", err.Error())
      //return
    }
    time.Sleep(5*time.Second)
    for DataMap := range chanDatamap{
        Ref := Get_Ref(wd, DataMap)
        lock.Lock()
        f, err := excelize.OpenFile("Data.xlsx")
        if err != nil {
            fmt.Println(err)
            return
        }
        //fmt.Printf("Processing Patient:%d\n",i-1)
        f.SetCellValue("Sheet1", "R"+DataMap["num"], Ref) 
        fmt.Println("Ref with A_constant Optimized :" + Ref)
        err = f.Save()
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        lock.Unlock()
    }
    waitGroup.Done()
}




func Find_Send(wd selenium.WebDriver, ID string, key string){
   btn, err := wd.FindElement(selenium.ByID, ID)
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   err =  btn.Clear()
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   err =  btn.SendKeys(key)
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
}

func ConvertStrSlice2Map(sl []string) map[string]struct{} {
    set := make(map[string]struct{}, len(sl))
    for _, v := range sl {
        set[v] = struct{}{}
    }
    return set
}


func InMap(m map[string]struct{}, s string) bool {
    _, ok := m[s]
    return ok
}


func Get_Ref(wd selenium.WebDriver, DataMap map[string]string)(Ref string) {
    IOL := DataMap["IOL"]
    for k,v := range DataMap{
        switch k {
        case "IOL","Ref_PostOP","num":
            continue
        default:
            Find_Send(wd,k,v)
        }
    }
   //time.Sleep(15*time.Second)
   Calc, err := wd.FindElement(selenium.ByID, "MainContent_Button1")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   if err := Calc.Click(); err != nil {
        panic(err)
    }
   Formula, err := wd.FindElement(selenium.ByLinkText, "Universal Formula")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   if err := Formula.Click(); err != nil {
        panic(err)
    }

   t,err :=wd.FindElement(selenium.ByTagName, "tbody")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   Results,err :=t.Text()
   if err != nil {
     panic(err)
   }
   //fmt.Println(Results)
   Power_Refs := strings.Split(Results,"\n")
   //fmt.Println(Power_Refs)
   Power_Refs_Map := map[string](string){}
   Powers := []string{}
   for i,Power_Ref := range Power_Refs{
       if i==0 {
            continue
        } else {
            _Power_Ref:=strings.Split(Power_Ref," ")
            Power := _Power_Ref[0]
            Powers = append(Powers,Power)
            Ref:=_Power_Ref[2]
            Power_Refs_Map[Power] = Ref
        } 
   }
   /*
   for Power,Ref := range Power_Refs_Map{
    fmt.Println(Power,Ref)
   }
   */
   //fmt.Println(Powers)
   set := ConvertStrSlice2Map(Powers)
   if !InMap(set,IOL){

    fmt.Println("IOL not in Powers")
    Patient_Data, _:= wd.FindElement(selenium.ByLinkText, "Patient Data")
    if err := Patient_Data.Click(); err != nil {
        panic(err)
    }
    return "---"
   }else{
        fmt.Printf("IOL in Powers:%s\n",IOL)
        Ref = Power_Refs_Map[IOL]
        Patient_Data, _:= wd.FindElement(selenium.ByLinkText, "Patient Data")
        if err := Patient_Data.Click(); err != nil {
            panic(err)
        }
        return Ref
    }
}