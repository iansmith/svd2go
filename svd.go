package svd

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"strings"
	"unicode"
)

const arrayDesignator = "[%s]"
type usedFields [32]bool

func (u *usedFields) Use(n,m int) bool {
	result:=true
	for i:=m; i<=n; i++ {
		if u[i]==true {
			result=false
		} else {
			u[i]=true
		}
	}
	return result
}
func (u *usedFields) Clear() {
	for i:=0; i<32; i++ {
		u[i]=false
	}
}

type constantDef struct {
	Name string
	Value string
}
var constants = []*constantDef{}

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

	//
	// Create templates
	//
	deviceTemplate:=template.New("device")
	deviceTemplate = template.Must(deviceTemplate.Parse(deviceTemplateText))

	bitFieldDeclTemplate:=template.New("bitFieldDecl")
	bitFieldDeclTemplate = template.Must(bitFieldDeclTemplate.Parse(bitFieldDeclTemplateText))

	preambleTemplate:=template.New("preamble")
	preambleTemplate = template.Must(preambleTemplate.Parse(preambleTemplateText))

	constantTemplate:=template.New("constant")
	constantTemplate = template.Must(constantTemplate.Parse(constantTemplateText))

	//
	// Process what we got in the SVD
	//
	addReservedRegisters(&device)
	makeObjectsExported(&device)
	makeBitfieldDecl(&device)

	//
	//assign things we got from user
	//
	device.Package=opts.Pkg
	device.Tags = opts.Tags
	device.SourceFilename=opts.InputFilename
	device.Import = opts.Import

	////////// EXECUTE TEMPLATES //////////////
	if err:=preambleTemplate.Execute(opts.Out,device); err!=nil {
		log.Fatal(err)
	}

	for _, p := range device.Peripheral {
		if !p.AddressBlock.Size.IsSet() {
			log.Fatalf("peripheral %s's address block has no size", p.Name)
		}
		for _, r := range p.Register {
			if !r.AddressOffset.IsSet() {
				log.Fatalf("peripheral %s, register %s has no address offset", p.Name,r.Name)
			}
			if !r.Size.IsSet() {
				log.Fatalf("peripheral %s, register %s has no size", p.Name,r.Name)
			}
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
	if err:=constantTemplate.Execute(opts.Out,constants); err!=nil {
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
				f.RegName = makeExported(r.Name)
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
// is respected on the peripheral.  Checks that all the names do not
// have spaces in them (since that's not valid go).  Check that dimensioned
// arrays have valid size and that the name of array registers is right.
func makeObjectsExported(d *DeviceDef) {
	for _, peripheral := range d.Peripheral {
		peripheral.Name = strings.TrimSpace(peripheral.Name)
		peripheral.TypeName = strings.TrimSpace(peripheral.TypeName)
		if !isValidGoIdentifier(peripheral.Name){
			log.Fatalf("names in an SVD document must be valid identifiers in go:  peripheral '%s'",
				peripheral.Name)
		}
		if peripheral.TypeName==""  {
			peripheral.TypeName=peripheral.Name
		}
		if !isValidGoIdentifier(peripheral.TypeName){
			log.Fatalf("names in an SVD document must be valid identifiers in go:  peripheral typeName '%s'",
				peripheral.TypeName)
		}
		peripheral.TypeName = makeExported(peripheral.Name)+"Def"
		if peripheral.HeaderStructName!="" {
			peripheral.TypeName = makeExported(peripheral.HeaderStructName)
		}
		peripheral.Name = makeExported(peripheral.Name)
		for _, r:=range peripheral.Register {
			r.Name = strings.TrimSpace(r.Name)
			if !isValidGoIdentifier(r.Name){
				log.Fatalf("names in an SVD document must be valid identifiers in go:  register '%s'",
					r.Name)
			}
			if strings.HasPrefix(r.Name,"reserved") || strings.HasSuffix(r.Name,arrayDesignator){
				r.TypeName=""
			} else {
				r.TypeName = makeExported(r.Name)
			}
			//
			// Handle Dims for registers only
			//
			if r.Dim.IsSet() {
				if r.Dim.Get()<1 {
					log.Fatalf("size of dimensioned array must be at least 1 (register %s in peripheral %s)",
						r.Name, peripheral.Name)
				}
				if r.DimIncrement.Get()!=4 {
					failNotImplemented("we can only generate registers that are 32 bits or dimIncrement=4 (register %s in peripheral %s)",
						r.Name, peripheral.Name)
				}
				if !strings.HasSuffix(r.Name,arrayDesignator) {
					log.Fatalf("dimension given as %d but register is not declared as array with '[%%s]' (register %s in peripheral %s)",
						r.Dim.Get(), r.Name, peripheral.Name)
				} else {
					if r.TypeName!="" {
						log.Fatalf("cannot use typeName with a dimensioned array (register %s in peripheral %s)",
							r.Name,peripheral.Name)
					}
					r.Name = strings.TrimSuffix(r.Name,arrayDesignator)
					for n, v:= range strings.Split(r.DimIndex,",") {
						c:=&constantDef{
							Name:r.Name+strings.TrimSpace(v),
							Value: fmt.Sprint(n),
						}
						constants=append(constants,c)
					}
				}
			}
			for _, f:=range r.Field {
				f.Name = strings.TrimSpace(f.Name)
				if !isValidGoIdentifier(r.Name){
					log.Fatalf("names in an SVD document must be valid identifiers in go:  field '%s'",
						f.Name)
				}
				f.Name = makeExported(f.Name)
				for _, e:=range f.EnumeratedValue {
					e.Name = strings.TrimSpace(e.Name)
					if !isValidGoIdentifier(e.Name){
						log.Fatalf("names in an SVD document must be valid identifiers in go:  enumerated value '%s'",
							e.Name)
					}
					e.Name=makeExported(e.Name)
				}
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
						Size: MultiformatInt{v:32, isSet:true},
						AddressOffset: MultiformatInt{v:int64(current), isSet:true},
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
		n:=len(regsOutput)<<2 //4 bytes each
		//fill ones _after_ last declared register
		for i:=int64(n); i<peripheral.AddressBlock.Size.Get(); i+=4{
			reserved := RegisterDef{
				Name: fmt.Sprintf("reserved%03d", reservedCount),
				Size: MultiformatInt{v:32, isSet:true},
				AddressOffset: MultiformatInt{v:int64(i), isSet:true},
			}
			regsOutput = append(regsOutput, &reserved)
			reservedCount++
		}
		peripheral.Register=regsOutput
	}
}

func isValidGoIdentifier(s string) bool {
	if s=="" {
		return false
	}
	if strings.Index(s," ")!=-1 {
		return false
	}
	r:=[]rune(s)
	if unicode.IsDigit(r[0]){
		return false
	}
	switch s {
	case "_", "true","false","iota","nil":
		return false
	case "int", "int8", "int16", "int32", "int64", "uint",
		"uint8", "uint16", "uint32", "uint64", "uintptr",
		"float32", "float64", "complex128", "complex64",
		"bool", "byte", "rune", "string", "error":
		return false
	case "make", "len", "cap", "new", "append", "copy", "close",
			"delete", "complex", "real", "imag", "panic", "recover":
		return false
	}

	return true
}