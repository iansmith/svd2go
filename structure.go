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
	Peripheral []PeripheralDef `xml:"peripherals>peripheral"`
}

type CpuDef struct {
	Name string `xml:"name"`
	Revision string `xml:"revision"`
	Endian string `xml:"endian"`
	MpuPresent Boolean `xml:"mpuPresent"`
	FpuPresent Boolean `xml:"fpuPresent"`
	NvicPrioBits string `xml:"NvicPrioBits"`
	VendorSystickConfig string `xml:"VendorSystickConfig"`
}

type PeripheralDef struct {
	Name string `xml:"name"`
	Version string `xml:"version"`
	Description string `xml:"description"`
	GroupName string `xml:"groupName"`
	AddressBlock AddressBlockDef `xml:"addressBlock"`
	Interrupt InterruptDef `xml:"interrupt"`
	BaseAddress MultiformatInt `xml:"baseAddress"`
	Size MultiformatInt `xml:"size"`
	Access Access `xml:"access"`
	Register []RegisterDef `xml:"registers>register"`
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
	Field []FieldDef `xml:"fields>field"`
}


type FieldDef struct {
	Name string `xml:"name"`
	Description string `xml:"description"`
	BitRange BitRange `xml:"bitRange"`
	Access Access `xml:"access"`
	EnumeratedValue []EnumeratedValueDef `xml:"enumeratedValues>enumeratedValue"`
}

type EnumeratedValueDef struct {
	Name string `xml:"name"name"`
	Description string `xml:"description"`
	Value MultiformatInt `xml:"value"`
}

