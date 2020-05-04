package svd

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type SVDDecodeError struct {
	val string
}
func (s *SVDDecodeError) Error() string {
	return fmt.Sprintf("MultiformatInt DecodeError: unable to understand input value: %s",s.val)
}

type MultiformatInt struct {
	v int64
}
func (m* MultiformatInt) Get() int64 {
	return m.v
}

func (m *MultiformatInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if strings.HasPrefix(s,"0x") {
		v, err:=strconv.ParseInt(s[2:],16,64)
		//fmt.Printf("*** 0x %s => %s\n",s,fmt.Sprint(v))
		if err!=nil {
			return &SVDDecodeError{s}
		}
		m.v = v
		return nil
	}
	if strings.HasPrefix(s,"#") {
		v, err:=strconv.ParseInt(s[1:],2,64)
		//fmt.Printf("*** # %s => %s\n",s,fmt.Sprint(v))
		if err!=nil {
			return &SVDDecodeError{s}
		}
		m.v = v
		return nil
	}
	v, err:=strconv.ParseInt(s,10, 64)
	//fmt.Printf("*** %s => %s\n",s,fmt.Sprint(v))
	if err!=nil {
		return &SVDDecodeError{s}
	}
	m.v = v
	return nil
}


type Boolean struct {
	v bool
}
func (b Boolean) Get() bool {
	return b.v
}

func (b *Boolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if strings.ToLower(s)=="true" {
		b.v = true
		return nil
	}
	if strings.ToLower(s)=="false" {
		b.v = false
		return nil
	}
	return &SVDDecodeError{s}
}

type Access struct {
	read bool
	write bool
}
func (a Access) Get() (bool, bool){
	return a.read, a.write
}
func (a Access) String() string {
	result:="read-only"
	if a.read && a.write{
		result = "read-write"
	} else if a.write {
		result = "write-only"
	}
	return result
}

func (b *Access) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if strings.ToLower(s)=="read-write" {
		b.read=true
		b.write=true
		return nil
	}
	if strings.ToLower(s)=="read-only" {
		b.read=true
		return nil
	}
	if strings.ToLower(s)=="write-only" {
		b.write=true
		return nil
	}
	return &SVDDecodeError{s}
}

type BitRange struct {
	lsb int
	msb int
}
func (b *BitRange) String() string{
	return fmt.Sprintf("[%d:%d]",b.msb,b.lsb)
}

func (b *BitRange) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	if len(s)>2 && s[0]=='[' && s[len(s)-1]==']' {
		pieces:=strings.Split(s[1:len(s)-1],":")
		if len(pieces)!=2 {
			return &SVDDecodeError{s}
		}
		first, err:=strconv.ParseInt(pieces[0],10, 64)
		if err!=nil {
			return &SVDDecodeError{s}
		}
		second, err:=strconv.ParseInt(pieces[1],10, 64)
		if err!=nil {
			return &SVDDecodeError{s}
		}
		lsb:=int(first)
		msb:=int(second)
		if first>second {
			lsb=int(second)
			msb=int(first)
		}
		b.lsb=lsb
		b.msb=msb
		return nil
	}
	return &SVDDecodeError{s}
}
