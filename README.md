# ip-place
input IP ，return province and city

go get github.com/gitxiaolin/ip_place

# example 
package main

import {

  "github.com/gitxiaolin/ip_place"
  
  "fmt"
  
}

func main(){

  ip := "119.39.23.134"
  
  province,city,err := ip_place.GetPlaceNameByIP(ip)
  
  if err!= nil{
  
    fmt.Println(err)
    
    return 
    
  }
  
  fmt.Println(province,city)
  
}
