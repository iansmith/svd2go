package svd

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"strings"
)

func ProcessSVD(reader io.Reader, write io.Writer) {
	var device DeviceDef
	decoder:=xml.NewDecoder(reader)
	if err:=decoder.Decode(&device); err!=nil {
		log.Fatalf("decoding error: %v",err)
	}
	for _, peripheral:=range device.Peripheral {
		groupName:=""
		if peripheral.GroupName!="" {
			groupName="["+peripheral.GroupName+"]"
		}
		fmt.Printf(">>> PERIPHERAL %s %s  %s <<<\n",strings.ToUpper(peripheral.Name), groupName,peripheral.Description)
		fmt.Printf("%20s 0x%x\n","base address",peripheral.BaseAddress.Get())
		fmt.Printf("%20s %d\n","size",peripheral.Size.Get())
		fmt.Printf("%20s %s\n","access",peripheral.Access)
		fmt.Printf("%20s: 0x%x\n","addr block => base addr",peripheral.AddressBlock.BaseAddress.Get())
		fmt.Printf("%20s: %d\n","addr block => size",peripheral.AddressBlock.Size.Get())
		fmt.Printf("%20s: %s\n","addr block => usage",peripheral.AddressBlock.Usage)
		fmt.Printf("\tINTERRUPT %s  %s\n",peripheral.Interrupt.Name,peripheral.Interrupt.Description)
		fmt.Printf("\t%20s: %d\n","interrupt value",peripheral.Interrupt.Value.Get())
		for _, reg:=range peripheral.Register {
			fmt.Printf("\t\tREGISTER %s  %s\n",reg.Name,reg.Description)
			fmt.Printf("\t\t%20s : 0x%x\n","  offset", reg.AddressOffset.Get())
			fmt.Printf("\t\t%20s : %d\n","  size", reg.Size.Get())
			fmt.Printf("\t\t%20s : %s\n","  access", reg.Access.String())
			fmt.Printf("\t\t%20s : 0x%x\n","  reset value", reg.ResetValue.Get())
			fmt.Printf("\t\t%20s : 0x%x\n","  reset mask", reg.ResetMask.Get())
			for _, field:=range reg.Field {
				fmt.Printf("\t\t\tFIELD %s   %s\n",field.Name,field.Description)
				fmt.Printf("\t\t\t%20s : %s\n","bit range",field.BitRange.String())
				fmt.Printf("\t\t\t%20s : %s\n","access", field.Access.String())
				for _, enumval:=range field.EnumeratedValue {
					fmt.Printf("\t\t\t\t%d: %s  %s\n",enumval.Value.Get(),enumval.Name,enumval.Description)
				}
			}
		}
		fmt.Printf("\n")
	}
}
