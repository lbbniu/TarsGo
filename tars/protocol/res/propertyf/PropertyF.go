// Package propertyf comment
// This file was generated by tars2go 1.1.5
// Generated from PropertyF.tars
package propertyf

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars/protocol/codec"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = fmt.Errorf
var _ = codec.FromInt8

// StatPropMsgHead struct implement
type StatPropMsgHead struct {
	ModuleName   string `json:"moduleName"`
	Ip           string `json:"ip"`
	PropertyName string `json:"propertyName"`
	SetName      string `json:"setName"`
	SetArea      string `json:"setArea"`
	SetID        string `json:"setID"`
	SContainer   string `json:"sContainer"`
	IPropertyVer int32  `json:"iPropertyVer"`
}

func (st *StatPropMsgHead) ResetDefault() {
	st.IPropertyVer = 1
}

//ReadFrom reads  from _is and put into struct.
func (st *StatPropMsgHead) ReadFrom(_is *codec.Reader) error {
	var (
		err    error
		length int32
		have   bool
		ty     byte
	)
	st.ResetDefault()

	err = _is.Read_string(&st.ModuleName, 0, true)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.Ip, 1, true)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.PropertyName, 2, true)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.SetName, 3, false)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.SetArea, 4, false)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.SetID, 5, false)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.SContainer, 6, false)
	if err != nil {
		return err
	}

	err = _is.Read_int32(&st.IPropertyVer, 7, false)
	if err != nil {
		return err
	}

	_ = err
	_ = length
	_ = have
	_ = ty
	return nil
}

//ReadBlock reads struct from the given tag , require or optional.
func (st *StatPropMsgHead) ReadBlock(_is *codec.Reader, tag byte, require bool) error {
	var (
		err  error
		have bool
	)
	st.ResetDefault()

	have, err = _is.SkipTo(codec.STRUCT_BEGIN, tag, require)
	if err != nil {
		return err
	}
	if !have {
		if require {
			return fmt.Errorf("require StatPropMsgHead, but not exist. tag %d", tag)
		}
		return nil
	}

	err = st.ReadFrom(_is)
	if err != nil {
		return err
	}

	err = _is.SkipToStructEnd()
	if err != nil {
		return err
	}
	_ = have
	return nil
}

//WriteTo encode struct to buffer
func (st *StatPropMsgHead) WriteTo(_os *codec.Buffer) error {
	var err error

	err = _os.Write_string(st.ModuleName, 0)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.Ip, 1)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.PropertyName, 2)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.SetName, 3)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.SetArea, 4)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.SetID, 5)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.SContainer, 6)
	if err != nil {
		return err
	}

	err = _os.Write_int32(st.IPropertyVer, 7)
	if err != nil {
		return err
	}

	_ = err

	return nil
}

//WriteBlock encode struct
func (st *StatPropMsgHead) WriteBlock(_os *codec.Buffer, tag byte) error {
	var err error
	err = _os.WriteHead(codec.STRUCT_BEGIN, tag)
	if err != nil {
		return err
	}

	err = st.WriteTo(_os)
	if err != nil {
		return err
	}

	err = _os.WriteHead(codec.STRUCT_END, 0)
	if err != nil {
		return err
	}
	return nil
}

// StatPropInfo struct implement
type StatPropInfo struct {
	Policy string `json:"policy"`
	Value  string `json:"value"`
}

func (st *StatPropInfo) ResetDefault() {
}

//ReadFrom reads  from _is and put into struct.
func (st *StatPropInfo) ReadFrom(_is *codec.Reader) error {
	var (
		err    error
		length int32
		have   bool
		ty     byte
	)
	st.ResetDefault()

	err = _is.Read_string(&st.Policy, 0, true)
	if err != nil {
		return err
	}

	err = _is.Read_string(&st.Value, 1, true)
	if err != nil {
		return err
	}

	_ = err
	_ = length
	_ = have
	_ = ty
	return nil
}

//ReadBlock reads struct from the given tag , require or optional.
func (st *StatPropInfo) ReadBlock(_is *codec.Reader, tag byte, require bool) error {
	var (
		err  error
		have bool
	)
	st.ResetDefault()

	have, err = _is.SkipTo(codec.STRUCT_BEGIN, tag, require)
	if err != nil {
		return err
	}
	if !have {
		if require {
			return fmt.Errorf("require StatPropInfo, but not exist. tag %d", tag)
		}
		return nil
	}

	err = st.ReadFrom(_is)
	if err != nil {
		return err
	}

	err = _is.SkipToStructEnd()
	if err != nil {
		return err
	}
	_ = have
	return nil
}

//WriteTo encode struct to buffer
func (st *StatPropInfo) WriteTo(_os *codec.Buffer) error {
	var err error

	err = _os.Write_string(st.Policy, 0)
	if err != nil {
		return err
	}

	err = _os.Write_string(st.Value, 1)
	if err != nil {
		return err
	}

	_ = err

	return nil
}

//WriteBlock encode struct
func (st *StatPropInfo) WriteBlock(_os *codec.Buffer, tag byte) error {
	var err error
	err = _os.WriteHead(codec.STRUCT_BEGIN, tag)
	if err != nil {
		return err
	}

	err = st.WriteTo(_os)
	if err != nil {
		return err
	}

	err = _os.WriteHead(codec.STRUCT_END, 0)
	if err != nil {
		return err
	}
	return nil
}

// StatPropMsgBody struct implement
type StatPropMsgBody struct {
	VInfo []StatPropInfo `json:"vInfo"`
}

func (st *StatPropMsgBody) ResetDefault() {
}

//ReadFrom reads  from _is and put into struct.
func (st *StatPropMsgBody) ReadFrom(_is *codec.Reader) error {
	var (
		err    error
		length int32
		have   bool
		ty     byte
	)
	st.ResetDefault()

	_, ty, err = _is.SkipToNoCheck(0, true)
	if err != nil {
		return err
	}

	if ty == codec.LIST {
		err = _is.Read_int32(&length, 0, true)
		if err != nil {
			return err
		}

		st.VInfo = make([]StatPropInfo, length)
		for i0, e0 := int32(0), length; i0 < e0; i0++ {

			err = st.VInfo[i0].ReadBlock(_is, 0, false)
			if err != nil {
				return err
			}

		}
	} else if ty == codec.SIMPLE_LIST {
		err = fmt.Errorf("not support simple_list type")
		if err != nil {
			return err
		}

	} else {
		err = fmt.Errorf("require vector, but not")
		if err != nil {
			return err
		}

	}

	_ = err
	_ = length
	_ = have
	_ = ty
	return nil
}

//ReadBlock reads struct from the given tag , require or optional.
func (st *StatPropMsgBody) ReadBlock(_is *codec.Reader, tag byte, require bool) error {
	var (
		err  error
		have bool
	)
	st.ResetDefault()

	have, err = _is.SkipTo(codec.STRUCT_BEGIN, tag, require)
	if err != nil {
		return err
	}
	if !have {
		if require {
			return fmt.Errorf("require StatPropMsgBody, but not exist. tag %d", tag)
		}
		return nil
	}

	err = st.ReadFrom(_is)
	if err != nil {
		return err
	}

	err = _is.SkipToStructEnd()
	if err != nil {
		return err
	}
	_ = have
	return nil
}

//WriteTo encode struct to buffer
func (st *StatPropMsgBody) WriteTo(_os *codec.Buffer) error {
	var err error

	err = _os.WriteHead(codec.LIST, 0)
	if err != nil {
		return err
	}

	err = _os.Write_int32(int32(len(st.VInfo)), 0)
	if err != nil {
		return err
	}

	for _, v := range st.VInfo {

		err = v.WriteBlock(_os, 0)
		if err != nil {
			return err
		}

	}

	_ = err

	return nil
}

//WriteBlock encode struct
func (st *StatPropMsgBody) WriteBlock(_os *codec.Buffer, tag byte) error {
	var err error
	err = _os.WriteHead(codec.STRUCT_BEGIN, tag)
	if err != nil {
		return err
	}

	err = st.WriteTo(_os)
	if err != nil {
		return err
	}

	err = _os.WriteHead(codec.STRUCT_END, 0)
	if err != nil {
		return err
	}
	return nil
}
