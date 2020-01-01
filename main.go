package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
)

type Node struct {
	Value      string `xml:"value,attr"`
	Name       string `xml:"name,attr"`
	Abbr       string `xml:"abbr,attr"`
	Desc       string `xml:"desc,attr"`
	CustomDesc string `xml:"customdesc,attr"`
}

type Data struct {
	Node
	Param []Node `xml:"PARAM"`
}

type Header struct {
	Data []Data `xml:"DATA"`
}
type LSB struct {
	Node
	Data []Data `xml:"LSB>DATA"`
}
type Structure struct {
	LSB []LSB `xml:"LSB"`
}
type SysX struct {
	Header    Header    `xml:"Header"`
	Structure Structure `xml:"Structure"`
}

type Param struct {
	Address    uint16
	Name       string
	Desc       string
	CustomDesc string
}

func main() {

	data, err := ioutil.ReadFile("/Users/witek/Downloads/Katana editor - FW3 2/midi.xml")
	if err != nil {
		panic(err)
	}

	sysx := SysX{}
	err = xml.Unmarshal(data, &sysx)
	if err != nil {
		panic(err)
	}
	//fmt.Println(sysx)
	//var equ []Param
	idx := 0
	for _, lsb := range sysx.Structure.LSB {
		fmt.Println("----------", lsb.Value, lsb.Name)
		for _, data := range lsb.Data {
			idx++
			/*if data.Name == "EFFECTS" {
				address1, _ := strconv.ParseUint(lsb.Value, 10, 8)
				address2, _ := strconv.ParseUint(data.Value, 10, 8)
				equ = append(equ, Param{
					Address:    uint16(address1<<8 | address2),
					Name:       data.Name,
					Desc:       data.Desc,
					CustomDesc: data.CustomDesc,
				})
			}*/
			params := ""
			if len(data.Param) == 1 && data.Param[0].Value == "range" {
				tokens := strings.Split(data.Param[0].Name, "/")
				params = fmt.Sprintf("\n     %s -> %s (%s -> %s)", tokens[0], tokens[1], tokens[2], tokens[3])
			} else if len(data.Param) > 1 && data.Name != "Text" && !strings.HasPrefix(data.Name, "Name") {
				params = "\n"
				for _, param := range data.Param {
					params += "     " + param.Value + " " + param.Name + "\n"
				}
			}
			fmt.Printf("%s %s %s %s %s %s\n", lsb.Value, data.Value, data.Name, data.Desc, data.CustomDesc, params)

			/*for _, param := range data.Param {
				fmt.Println(param.Value, param.Name, param.Desc, param.CustomDesc)
			}*/
		}
	}
	fmt.Println("Bytes count:", idx)
	/*fmt.Println("Equalizer")
	for _, data := range equ {
		fmt.Printf("%04X %s %s %s\n", data.Address, data.Name, data.Desc, data.CustomDesc)
	}*/
}
