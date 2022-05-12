package main

import (
    "os"
    "fmt"
    "strconv"
    "time"
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
    err = wd.Get("https://www.iolformula.com/")
    if err != nil {
      fmt.Println("get page faild", err.Error())
      //return
   }
   time.Sleep(5*time.Second)
   Agree, err := wd.FindElement(selenium.ByCSSSelector, `div[class="btn btn-primary btn_agreement"]`)
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   if err := Agree.Click(); err != nil {
        panic(err)
    }
   time.Sleep(5*time.Second)

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
                            "Patient":row[0],
                            "Sex":row[14],
                            "A-Constant1":"118.80",
                            "al-right":row[1],
                            "k1-right":row[3],
                            "k2-right":row[5],
                            "acd-right":row[2],
                            "right-target":"0",
                            "lt-right":row[10],
                            //"cct-right":row[],
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
   if len(IOL)<=2{IOL = IOL +".0"}
   //fmt.Println(IOL)
   for k,v := range DataMap{
    switch k {
    case "Sex","IOL","Ref_PostOP":
        continue
    default:
        Find_Send(wd,k,v)
    }
   }
   time.Sleep(time.Second)
   //time.Sleep(3* time.Second)
   Sex, _ := wd.FindElement(selenium.ByCSSSelector, `div[class="btn-group radio-group h-gender"]`)
   male, _ := Sex.FindElement(selenium.ByCSSSelector, "label:nth-child(1)")
   female, _ := Sex.FindElement(selenium.ByCSSSelector, "label:nth-child(2)")
   if DataMap["Sex"] =="1"{
       if err := male.Click(); err != nil {
        panic(err)
       }
   }else {
       if err := female.Click(); err != nil {
        panic(err)
       }
   }
   Submit, _ := wd.FindElement(selenium.ByCSSSelector, `div[class ="button_submit_block form-group submit row jq_class_1"]`)
   Calc, _ := Submit.FindElement(selenium.ByCSSSelector, "div:nth-child(1)")
   if err := Calc.Click(); err != nil {
        panic(err)
    }
   time.Sleep(5*time.Second)
   t,err :=wd.FindElement(selenium.ByCSSSelector, `div[class="res_nontoric"]`)
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
       if i==0 {
            continue
        } else {
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
    Menu, _ := wd.FindElement(selenium.ByCSSSelector, `div[class="form-group submit row"]`)
    NewPatient, _ := Menu.FindElement(selenium.ByCSSSelector, "div:nth-child(3)")
    if err := NewPatient.Click(); err != nil {
        panic(err)
    }
    return "---"
   }else{
        fmt.Printf("IOL in Powers:%s\n",IOL)
        Ref = Power_Refs_Map[IOL]
        Menu, _ := wd.FindElement(selenium.ByCSSSelector, `div[class="form-group submit row"]`)
        NewPatient, _ := Menu.FindElement(selenium.ByCSSSelector, "div:nth-child(3)")
        if err := NewPatient.Click(); err != nil {
            panic(err)
        }
        return Ref
    }
}