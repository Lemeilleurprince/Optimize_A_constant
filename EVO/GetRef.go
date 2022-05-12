package main

import (
    "os"
    "fmt"
    "strconv"
    //"time"
    "strings"
    "github.com/360EntSecGroup-Skylar/excelize"
    "github.com/tebeka/selenium"
)

const (
    chromeDriverPath = "./chromedriver"
    port             = 8080
)

func main() {
    var n int
    fmt.Println("Please input from where to start,default input is 0:")
    fmt.Scanln(&n)

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

    f, err := excelize.OpenFile("Data.xlsx")
    if err != nil {
        fmt.Println(err)
        return
    }
    // Get all the rows in the Sheet1.
    DataMap := make(map[string]string, 0)
    rows, err := f.GetRows("Sheet1")
    for i, row := range rows {
        if i>=n{
            if i == 0 {
                f.SetCellValue("Sheet1", "R1", "Ref with A_constant Optimized ") 
                continue
            }else{

                fmt.Println("Processing patient:")
                fmt.Println(i)
                DataMap=map[string]string{
                            "TextBoxName":row[0],
                            //"TextBoxID":row[1],
                            "txtAConstant":"118.63",
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
                Ref := Get_Ref(wd, DataMap)
                
                f.SetCellValue("Sheet1", "R"+strconv.Itoa(i+1), Ref) 
                fmt.Println("Ref with A_constant Optimized :" + Ref)
                //time.Sleep(5*time.Second)
            }
            err = f.Save()
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
    }
    wd.Quit()
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
        case "IOL","Ref_PostOP":
            continue
        default:
            Find_Send(wd,k,v)
        }
    }
   //time.Sleep(15*time.Second)
   Calc, err := wd.FindElement(selenium.ByID, "btnCalculate")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   if err := Calc.Click(); err != nil {
        panic(err)
    }

   //time.Sleep(3* time.Second)
   div,err :=wd.FindElement(selenium.ByCSSSelector, "#PnPred")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   t,err :=div.FindElement(selenium.ByCSSSelector, "tbody")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   Results,err :=t.Text()
   if err != nil {
     panic(err)
   }
   //fmt.Println(Results)
   //time.Sleep(5* time.Second)
   Power_Refs := strings.Split(Results,"\n")
   //fmt.Println(len(Power_Refs))
   //fmt.Println(Power_Refs)
   Power_Refs_Map := map[string](string){}
   Powers := []string{}
   for i,Power_Ref := range Power_Refs{
       if i<=5 {
            continue
        } else if i<=10 {
            _Power_Ref:=strings.Split(strings.TrimSpace(Power_Ref)," ")
            //fmt.Println(_Power_Ref)
            Power := _Power_Ref[0]
            //fmt.Println(Power)
            Powers = append(Powers,Power)
            Ref:=_Power_Ref[1]
            Power_Refs_Map[Power] = Ref
            //fmt.Println(Ref)
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
    Back, _ := wd.FindElement(selenium.ByID, "btnBack")
    if err := Back.Click(); err != nil {
        panic(err)
    }
    return "---"
   }else{
        fmt.Printf("IOL in Powers:%s\n",IOL)
        Ref = Power_Refs_Map[IOL]
        Back, _ := wd.FindElement(selenium.ByID, "btnBack")
        if err := Back.Click(); err != nil {
            panic(err)
        }
        return Ref
    }
}