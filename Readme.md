# 注意事项
##  1. chromedriver
+ 请将chromedriver.exe更换为自己电脑chrome版本对应的驱动
+ 下载地址
[https://chromedriver.storage.googleapis.com/index.html](https://chromedriver.storage.googleapis.com/index.html)

## 2. Data.xlsx
+ 更换为自己的数据
+ 请勿调整列顺序，增删列
+ 程序运行时，请不要打开Data.xlsx

## 4. GetRef.go 
+ 该程序为得到优化A常数后，用优化后的A常数的平均值代入在线计算器以计算每个晶体的预留度数，请自行修改对应的A常数值

## 5. 其他
+ Barret Universal II 和 EVO多线程版为10线程并行
+ Kane计算网站对固定IP多次计算可能有封锁机制，请自行变动IP；反应较慢，多线程没戏
+ 多线程版为了保证运行的稳定，通道常开，因此运行结束并不会自行退出，请手动结束
