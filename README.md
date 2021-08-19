# gbkToJSON

Ever wanted to turn your gbk plasmid files into JSON? 
Look no further, just use this library.

Installation 

`go get https://github.com/cruzryan/gbkToJSON`

Usage

```go 
import (
"fmt"
"github.com/cruzryan/gbkToJSON"
)

func main(){
	data := " your gbk string data goes here! "

	string_json := gbkToJSON.getAsString(data)
	fmt.Println(string_json)
	
    
	plasmid := gbkToJSON.getAsPlasmidStruct(data)
    fmt.Println(plasmid.DNA)

}
```
