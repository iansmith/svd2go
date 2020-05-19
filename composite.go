package svd

import "os"

type CompositeSVD struct {
	Device *DeviceDef
	PeripheralFiles []string
	Perpipheral []*PeripheralDef
}

type CompositeDef struct {
	Device DeviceDef
	Includes []string
}

func ProcessCSVD(fp *os.File, opts *UserOptions) *CompositeSVD{
	//ok, we have the base elements (device and peripherals, likely just a device)
	result:=ProcessSVD(fp,opts,opts.InputFilename)

	return result
}

func (c *CompositeSVD) Add(p *PeripheralDef) {
	c.Perpipheral=append(c.Perpipheral, p)
}

func (c *CompositeSVD) ProcessPeripheral(fp *os.File, opts *UserOptions, filename string) *PeripheralDef{

}
