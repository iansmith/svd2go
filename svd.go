package svd

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
)

func ProcessSVD(reader io.Reader, opts *UserOptions) {
	var device DeviceDef
	decoder:=xml.NewDecoder(reader)
	if err:=decoder.Decode(&device); err!=nil {
		log.Fatalf("decoding error: %v",err)
	}
	if opts.Dump {
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

	deviceTemplate:=template.New("device")
	deviceTemplate = template.Must(deviceTemplate.Parse(deviceTemplateText))

	bitFieldDeclTemplate:=template.New("bitFieldDecl")
	bitFieldDeclTemplate = template.Must(bitFieldDeclTemplate.Parse(bitFieldDeclTemplateText))

	preambleTemplate:=template.New("preamble")
	preambleTemplate = template.Must(preambleTemplate.Parse(preambleTemplateText))

	addReservedRegisters(&device)
	makeObjectsExported(&device)
	makeBitfieldDecl(&device)
	device.Package=opts.Pkg
	device.Tags = opts.Tags
	device.SourceFilename=opts.InputFilename

	////////// EXECUTE TEMPLATES //////////////
	if err:=preambleTemplate.Execute(opts.Out,device); err!=nil {
		log.Fatal(err)
	}

	for _, p := range device.Peripheral {
		for _, r := range p.Register {
			for _, f := range r.Field {
				if err:=bitFieldDeclTemplate.Execute(opts.Out,f); err!=nil {
					log.Fatal(err)
				}

			}
		}
	}
	if err:=deviceTemplate.Execute(opts.Out,device); err!=nil {
		log.Fatal(err)
	}
}

func failNotImplemented(format string, params ...interface{}) {
	s:=fmt.Sprintf(format,params...)
	panic("not implemented:"+s)
}

// makeExported capitalizes the first letter only.
func makeExported(s string) (string) {
	if len(s)==0 {
		panic("attempt to process empty string with makeExported")
	}
	if strings.HasPrefix(s,"reserved") {
		return s
	}
	if len(s)==1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[0:1])+s[1:]
}

// Create parent reg for a bitfield, setup read/write attributes
func makeBitfieldDecl(d *DeviceDef) {
	for _, p := range d.Peripheral {
		for _, r := range p.Register {
			if strings.HasPrefix(r.Name,"reserved") {
				continue
			}
			var regRead, regWrite bool
			if !r.Access.IsSet() {
				//assume read-write
				regRead=true
				regWrite=true
			} else {
				regRead, regWrite = r.Access.Get()
			}
			for _, f := range r.Field {
				f.RegName = r.Name
				if !f.Access.IsSet() {
					f.CanRead=regRead
					f.CanWrite=regWrite
				} else {
					f.CanRead, f.CanWrite = f.Access.Get()
				}
				for _,ev:=range f.EnumeratedValue {
					ev.Field = f
				}
			}
		}
	}
}

// Make sure the objects named in the file are exported and that the HeaderStructName
// is respected on the peripheral.
func makeObjectsExported(d *DeviceDef) {
	for _, peripheral := range d.Peripheral {
		peripheral.TypeName = makeExported(peripheral.Name)+"Def"
		if peripheral.HeaderStructName!="" {
			peripheral.TypeName = makeExported(peripheral.HeaderStructName)
		}
		peripheral.Name = makeExported(peripheral.Name)
		for _, r:=range peripheral.Register {
			if strings.HasPrefix(r.Name,"reserved") {
				r.TypeName=""
			} else {
				r.TypeName = makeExported(r.Name)
			}
		}
	}
}

// addReservedRegisters does two checks, and generates all the intermediate registers
// that were not in the svd file. One check is the peripheral address block is for
// registers and the other is that the distance between registers is 32bits.  This function
// assumes registers that were not in the svd file are reserved/unused and are not
// exported out of the package.
func addReservedRegisters(d *DeviceDef) {
	//add in
	reservedCount:=0
	for _, peripheral := range d.Peripheral {
		current:=0x0
		regsOutput:=[]*RegisterDef{}
		for _, register:= range peripheral.Register {
			if peripheral.AddressBlock.Usage!="registers" {
				failNotImplemented("only able to generate code for register address blocks, %s in %s is not registers",
					peripheral.Name, d.Name)
			}
			for current <  int(register.AddressOffset.Get()) {
				for current != int(register.AddressOffset.Get()) {
					//too far
					if int(register.AddressOffset.Get()) < current {
						failNotImplemented("unable to generate registers for anything except 32 bit units, %s in %s is not aligned (expected 0x%x but got 0x%x)",
							register.Name, peripheral.Name, current, register.AddressOffset.Get())
					}
					//add in a register to make up the numbers
					reserved := RegisterDef{
						Name: fmt.Sprintf("reserved%03d", reservedCount),
						Size: MultiformatInt{32},
						AddressOffset: MultiformatInt{int64(current)},
					}
					regsOutput = append(regsOutput, &reserved)
					reservedCount++
					current+=4
				}
			}
			//we reached the next register
			regsOutput = append(regsOutput, register)
			current+=4 //only 32 bit regs
		}
		peripheral.Register=regsOutput
	}
}