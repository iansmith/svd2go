package svd


type DeviceDef struct {
	Vendor string `xml:"vendor"`
	VendorID string `xml:"vendorID"`
	Name string `xml:"name"`
	Series string `xml:"series"`
	Version string `xml:"version"`
	Description string `xml:"description"`
	LicenseText string `xml:"licenseText"`
	Cpu CpuDef `xml:"cpu"`
	AddressUnitBits int `xml:"addressUnitBits"`
	Width MultiformatInt `xml:"width"`
	Size MultiformatInt `xml:"size"`
	Access Access `xml:"access"`
	ResetValue MultiformatInt `xml:"resetValue"`
	ResetMask MultiformatInt `xml:"resetMask"`
	Peripheral []*PeripheralDef `xml:"peripherals>peripheral"`
	Package string // this comes from the user opts
	SourceFilename string // this is the filename used to create all this
	Tags string //this comes from the command line option
	Import string //this comes from command line option
}

type CpuDef struct {
	Name                string         `xml:"name"`
	Revision            string         `xml:"revision"`
	Endian              string         `xml:"endian"`
	MpuPresent          Boolean        `xml:"mpuPresent"`
	FpuPresent          Boolean        `xml:"fpuPresent"`
	NvicPrioBits        string         `xml:"NvicPrioBits"` //nested vector interrupt controller
	VendorSystickConfig string         `xml:"VendorSystickConfig"`
	FPUDP               Boolean        `xml:"fpuDP"`
	DspPresent          Boolean        `xml:"dspPresent"`
	ICachePresent       Boolean        `xml:"icachePresent"`
	DCachePresent       Boolean        `xml:"dcachePresent"`
	ItcmPresent         Boolean        `xml:"itcmPresent"` //TCM = tightly coupled memory
	DtcmPresent         Boolean        `xml:"dtcmPresent"`
	VtorPresent         Boolean        `xml:"vtorPresent"` // vector table offset register
	DeviceNumInterrupts MultiformatInt `xml:"deviceNumInterrupts"`

}

type PeripheralDef struct {
	Name string `xml:"name"`
	Version string `xml:"version"`
	Description string `xml:"description"`
	PrependToName string `xml:"prependToName"`
	AppendToName string `xml:"appendToName"`
	HeaderStructName string `xml:"headerStructName"`
	GroupName string `xml:"groupName"`
	AddressBlock *AddressBlockDef `xml:"addressBlock"`
	Interrupt InterruptDef `xml:"interrupt"`
	BaseAddress MultiformatInt `xml:"baseAddress"`
	Size MultiformatInt `xml:"size"`
	Access Access `xml:"access"`
	Register []*RegisterDef `xml:"registers>register"`
	TypeName string //generated by our processing
	SourceFilename string //from user options or the composite file
}

type AddressBlockDef struct {
	BaseAddress MultiformatInt `xml:"baseAddress"`
	Size MultiformatInt `xml:"size"`
	Usage string `xml:"usage"`
}

type InterruptDef struct {
	Name string `xml:"name"`
	Description string `xml:"descripton"`
	Value MultiformatInt `xml:"value"`
}

type RegisterDef struct {
	Name string `xml:"name"`
	Description string `xml:"description"`
	AddressOffset MultiformatInt `xml:"addressOffset"`
	Size  MultiformatInt `xml:"size"`
	Access Access `xml:"access"`
	ResetValue MultiformatInt `xml:"resetValue"`
	ResetMask MultiformatInt `xml:"resetMask"`
	Field []*FieldDef `xml:"fields>field"`
	TypeName string //created when processing, for reserved regs this is ""
	Dim MultiformatInt `xml:"dim"`
	DimIncrement MultiformatInt `xml:"dimIncrement"`
	DimIndex string `xml:"dimIndex"`
}


type FieldDef struct {
	Name string `xml:"name"`
	Description string `xml:"description"`
	BitRange BitRange `xml:"bitRange"`
	Access Access `xml:"access"`
	EnumeratedValue []*EnumeratedValueDef `xml:"enumeratedValues>enumeratedValue"`
	RegName string //created during processing
	CanRead, CanWrite bool //created during processing
}

type EnumeratedValueDef struct {
	Name string `xml:"name"name"`
	Description string `xml:"description"`
	Value MultiformatInt `xml:"value"`
	Field *FieldDef //created during processing
}

