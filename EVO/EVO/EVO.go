package EVO

import (
    "fmt"
    "strconv"
    //"time"
    "math"
    "strings"
    "github.com/tebeka/selenium"
)


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

func FloatRound(f float64, n int) float64 {
    format := "%." + strconv.Itoa(n) + "f"
    res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
    return res
}

func Iterate(wd selenium.WebDriver, A_constant float64,IOL string)(Ref float64){
   Back, err := wd.FindElement(selenium.ByID, "btnBack")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   if err := Back.Click(); err != nil {
        panic(err)
   }

   btn, err := wd.FindElement(selenium.ByID, "txtAConstant")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   err =  btn.Clear()
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   err =  btn.SendKeys(strconv.FormatFloat(A_constant,'f', 3, 64))
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   //time.Sleep(3 * time.Second)

   Calc, err := wd.FindElement(selenium.ByID, "btnCalculate")
   if err != nil {
      //panic(err)
      fmt.Println(err)
   }
   if err := Calc.Click(); err != nil {
        panic(err)
   }
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
   Power_Refs := strings.Split(Results,"\n")

   for i,Power_Ref := range Power_Refs{
       if i<=5 {
            continue
        } else if i<=10 {
            Power_Ref_:=strings.Split(strings.TrimSpace(Power_Ref)," ")
            Power_ := Power_Ref_[0]
            Ref_ := Power_Ref_[1]
            if Power_ == IOL{
                Ref,_ = strconv.ParseFloat(Ref_,64)
            }
        } 
   }
   return Ref 
}

func Ajust(wd selenium.WebDriver, A_constant float64,IOL string,Ref_Post float64, Ref float64)(float64, float64){
   A_constant = FloatRound(A_constant + Ref_Post - Ref,3)
   Ref = Iterate(wd, A_constant, IOL)
   return Ref, A_constant
}

func Micro_Ajust(wd selenium.WebDriver, A_constant float64, IOL string, _Ref float64, Ref float64,Step float64)(float64,float64, float64){
   A_constant =FloatRound(A_constant + Step,3)
   Ref = Iterate(wd, A_constant, IOL)
   //fmt.Println(_Ref,Ref)
   if _Ref > Ref{
    _Step := FloatRound(math.Abs(Step)/2,3)
    Step = _Step
    fmt.Println("Upaward!")
    //fmt.Println(Step)
   }else if _Ref < Ref{
    _Step := FloatRound(math.Abs(Step)/2,3)
    Step = - _Step
    fmt.Println("Downaward!")
    //fmt.Println(Step)
   }else{
    fmt.Println("Keep original direction!")
    //fmt.Println(Step)
   }
   //fmt.Println(Ref, A_constant, Step)
   fmt.Printf("Refraction (SE): %.3f A_constant: %.3f next step: %.3f\n", Ref, A_constant, Step)
   return Ref, A_constant, Step

}

func Micro_Ajust_UpAndDown(wd selenium.WebDriver, A_constant float64,IOL string, Ref float64)(float64){
    fmt.Println("Micro_Ajust_UpAndDown!")
    A_constant_max,A_constant_min :=A_constant,A_constant
    Ref0 :=Ref
    Step := 0.01
    for{
        _Ref := Ref
        Ref, A_constant_max, Step = Micro_Ajust(wd, A_constant_max,IOL, _Ref, Ref, Step) 
        if _Ref >Ref {
            break
        }
        
    }
    fmt.Println("A_constant_max:" +strconv.FormatFloat(A_constant_max, 'f', 3, 64))

    Ref = Ref0
    Step = -0.01
    for{
        _Ref := Ref
        Ref, A_constant_min, Step = Micro_Ajust(wd, A_constant_min,IOL,_Ref, Ref, Step)
        if _Ref < Ref {
            break
        }
    }
    fmt.Println("A_constant_min:" + strconv.FormatFloat(A_constant_min, 'f', 3, 64))
    A_constant = FloatRound ((A_constant_max+A_constant_min)/2,4)
    return A_constant
}

func Get_A_constant(wd selenium.WebDriver, DataMap map[string]string)(A_constant float64) {
   IOL := DataMap["IOL"]
   Ref_PostOP := DataMap["Ref_PostOP"]
   A_constant, _ = strconv.ParseFloat(DataMap["txtAConstant"],64)
   for k,v := range DataMap{
    switch k {
    case "IOL","Ref_PostOP":
        continue
    default:
        Find_Send(wd,k,v)
    }
   }
   //time.Sleep(3* time.Second)
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
    return 0
   }else{
        fmt.Printf("IOL in Powers:%s\n",IOL)
        fmt.Printf("Refraction post operation:%s\n",DataMap["Ref_PostOP"])
        Ref, _:= strconv.ParseFloat(Power_Refs_Map[IOL],64)
        Ref_Post, _ := strconv.ParseFloat(Ref_PostOP,64)
        //fmt.Println(Ref,Ref_Post)
        for {
                A_constant_ := FloatRound(A_constant + Ref_Post - Ref,3)
                //fmt.Println(A_constant_)
                if !(A_constant_ >=110 && A_constant_<=125){
                    fmt.Println("A_constant out of boundary!")
                    Back, _ := wd.FindElement(selenium.ByID, "btnBack")
                    if err := Back.Click(); err != nil {
                        panic(err)
                    }
                    return 0
                    break
                }else{
                    _Ref_D :=FloatRound(Ref_Post - Ref,3)
                    Ref, A_constant = Ajust(wd, A_constant,IOL,Ref_Post, Ref)
                    Ref_D := FloatRound(Ref_Post - Ref,3)
                    //fmt.Println(_Ref_D,Ref_D)
                    if math.Abs(Ref_D)<=0.02{
                        if Ref_D ==0 {
                            A_constant = Micro_Ajust_UpAndDown(wd, A_constant, IOL, Ref)
                            break
                        }else if _Ref_D*Ref_D<0 {

                            if Ref_D >0 {
                                fmt.Println("Micro_Ajust_Up")
                                Step := 0.01
                                for{
                                    _Ref := Ref
                                    _A_constant := A_constant
                                    Ref, A_constant, Step = Micro_Ajust(wd, A_constant,IOL, _Ref, Ref, Step)
                                    Ref_D = FloatRound(Ref_Post - Ref,2)
                                    if Ref_D ==0 {
                                        A_constant = Micro_Ajust_UpAndDown(wd, A_constant, IOL, Ref)
                                        break
                                    }
                                    if _Ref > Ref {
                                        fmt.Println(_A_constant,A_constant)
                                        A_constant = FloatRound ((_A_constant +A_constant)/2,4)
                                        break
                                    }
                                }
                                break
                            }
                            if Ref_D <0{
                                fmt.Println("Micro_Ajust_Down")
                                Step :=-0.01
                                for{
                                    _Ref := Ref
                                    _A_constant := A_constant
                                    Ref, A_constant, Step = Micro_Ajust(wd, A_constant,IOL, _Ref, Ref, Step)
                                    Ref_D = FloatRound(Ref_Post - Ref,2)
                                    if Ref_D ==0 {
                                        A_constant = Micro_Ajust_UpAndDown(wd, A_constant, IOL, Ref)
                                        break
                                    }
                                    if _Ref < Ref {
                                        fmt.Println(_A_constant,A_constant)
                                        A_constant = FloatRound ((_A_constant +A_constant)/2,4)
                                        break
                                    }
                                }
                                break
                            }
                        }
                    }
                }
            //time.Sleep(10* time.Second)
        }
        
    }
   //time.Sleep(15 * time.Second)
   //wd.Quit()
   Back, _ := wd.FindElement(selenium.ByID, "btnBack")
   if err := Back.Click(); err != nil {
        panic(err)
   }
   return A_constant
}




