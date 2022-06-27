// Package notifyf comment
// This file was generated by tars2go 1.1.7
// Generated from NotifyF.tars
package notifyf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	m "github.com/TarsCloud/TarsGo/tars/model"
	"github.com/TarsCloud/TarsGo/tars/protocol/codec"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/basef"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/requestf"
	"github.com/TarsCloud/TarsGo/tars/protocol/tup"
	"github.com/TarsCloud/TarsGo/tars/util/current"
	"github.com/TarsCloud/TarsGo/tars/util/tools"
	"unsafe"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = fmt.Errorf
	_ = codec.FromInt8
	_ = unsafe.Pointer(nil)
	_ = bytes.ErrTooLarge
)

// Notify struct
type Notify struct {
	servant m.Servant
}

// ReportServer is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) ReportServer(sServerName string, sThreadId string, sMessage string, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = buf.WriteString(sServerName, 1)
	if err != nil {
		return err
	}

	err = buf.WriteString(sThreadId, 2)
	if err != nil {
		return err
	}

	err = buf.WriteString(sMessage, 3)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}
	tarsResp := new(requestf.ResponsePacket)
	tarsCtx := context.Background()

	err = obj.servant.TarsInvoke(tarsCtx, 0, "reportServer", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// ReportServerWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) ReportServerWithContext(tarsCtx context.Context, sServerName string, sThreadId string, sMessage string, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = buf.WriteString(sServerName, 1)
	if err != nil {
		return err
	}

	err = buf.WriteString(sThreadId, 2)
	if err != nil {
		return err
	}

	err = buf.WriteString(sMessage, 3)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 0, "reportServer", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// ReportServerOneWayWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) ReportServerOneWayWithContext(tarsCtx context.Context, sServerName string, sThreadId string, sMessage string, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = buf.WriteString(sServerName, 1)
	if err != nil {
		return err
	}

	err = buf.WriteString(sThreadId, 2)
	if err != nil {
		return err
	}

	err = buf.WriteString(sMessage, 3)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 1, "reportServer", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// NotifyServer is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) NotifyServer(sServerName string, level NOTIFYLEVEL, sMessage string, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = buf.WriteString(sServerName, 1)
	if err != nil {
		return err
	}

	err = buf.WriteInt32(int32(level), 2)
	if err != nil {
		return err
	}

	err = buf.WriteString(sMessage, 3)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}
	tarsResp := new(requestf.ResponsePacket)
	tarsCtx := context.Background()

	err = obj.servant.TarsInvoke(tarsCtx, 0, "notifyServer", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// NotifyServerWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) NotifyServerWithContext(tarsCtx context.Context, sServerName string, level NOTIFYLEVEL, sMessage string, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = buf.WriteString(sServerName, 1)
	if err != nil {
		return err
	}

	err = buf.WriteInt32(int32(level), 2)
	if err != nil {
		return err
	}

	err = buf.WriteString(sMessage, 3)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 0, "notifyServer", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// NotifyServerOneWayWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) NotifyServerOneWayWithContext(tarsCtx context.Context, sServerName string, level NOTIFYLEVEL, sMessage string, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = buf.WriteString(sServerName, 1)
	if err != nil {
		return err
	}

	err = buf.WriteInt32(int32(level), 2)
	if err != nil {
		return err
	}

	err = buf.WriteString(sMessage, 3)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 1, "notifyServer", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// GetNotifyInfo is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) GetNotifyInfo(stKey *NotifyKey, stInfo *NotifyInfo, opts ...map[string]string) (ret int32, err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = stKey.WriteBlock(buf, 1)
	if err != nil {
		return ret, err
	}

	err = (*stInfo).WriteBlock(buf, 2)
	if err != nil {
		return ret, err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}
	tarsResp := new(requestf.ResponsePacket)
	tarsCtx := context.Background()

	err = obj.servant.TarsInvoke(tarsCtx, 0, "getNotifyInfo", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return ret, err
	}

	readBuf := codec.NewReader(tools.Int8ToByte(tarsResp.SBuffer))
	err = readBuf.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}

	err = (*stInfo).ReadBlock(readBuf, 2, true)
	if err != nil {
		return ret, err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return ret, nil
}

// GetNotifyInfoWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) GetNotifyInfoWithContext(tarsCtx context.Context, stKey *NotifyKey, stInfo *NotifyInfo, opts ...map[string]string) (ret int32, err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = stKey.WriteBlock(buf, 1)
	if err != nil {
		return ret, err
	}

	err = (*stInfo).WriteBlock(buf, 2)
	if err != nil {
		return ret, err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 0, "getNotifyInfo", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return ret, err
	}

	readBuf := codec.NewReader(tools.Int8ToByte(tarsResp.SBuffer))
	err = readBuf.ReadInt32(&ret, 0, true)
	if err != nil {
		return ret, err
	}

	err = (*stInfo).ReadBlock(readBuf, 2, true)
	if err != nil {
		return ret, err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return ret, nil
}

// GetNotifyInfoOneWayWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) GetNotifyInfoOneWayWithContext(tarsCtx context.Context, stKey *NotifyKey, stInfo *NotifyInfo, opts ...map[string]string) (ret int32, err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = stKey.WriteBlock(buf, 1)
	if err != nil {
		return ret, err
	}

	err = (*stInfo).WriteBlock(buf, 2)
	if err != nil {
		return ret, err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 1, "getNotifyInfo", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return ret, err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return ret, nil
}

// ReportNotifyInfo is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) ReportNotifyInfo(info *ReportInfo, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = info.WriteBlock(buf, 1)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}
	tarsResp := new(requestf.ResponsePacket)
	tarsCtx := context.Background()

	err = obj.servant.TarsInvoke(tarsCtx, 0, "reportNotifyInfo", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// ReportNotifyInfoWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) ReportNotifyInfoWithContext(tarsCtx context.Context, info *ReportInfo, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = info.WriteBlock(buf, 1)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 0, "reportNotifyInfo", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// ReportNotifyInfoOneWayWithContext is the proxy function for the method defined in the tars file, with the context
func (obj *Notify) ReportNotifyInfoOneWayWithContext(tarsCtx context.Context, info *ReportInfo, opts ...map[string]string) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	buf := codec.NewBuffer()
	err = info.WriteBlock(buf, 1)
	if err != nil {
		return err
	}

	var statusMap map[string]string
	var contextMap map[string]string
	if len(opts) == 1 {
		contextMap = opts[0]
	} else if len(opts) == 2 {
		contextMap = opts[0]
		statusMap = opts[1]
	}

	tarsResp := new(requestf.ResponsePacket)
	err = obj.servant.TarsInvoke(tarsCtx, 1, "reportNotifyInfo", buf.ToBytes(), statusMap, contextMap, tarsResp)
	if err != nil {
		return err
	}

	if len(opts) == 1 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
	} else if len(opts) == 2 {
		for k := range contextMap {
			delete(contextMap, k)
		}
		for k, v := range tarsResp.Context {
			contextMap[k] = v
		}
		for k := range statusMap {
			delete(statusMap, k)
		}
		for k, v := range tarsResp.Status {
			statusMap[k] = v
		}
	}
	_ = length
	_ = have
	_ = ty
	return nil
}

// SetServant sets servant for the service.
func (obj *Notify) SetServant(servant m.Servant) {
	obj.servant = servant
}

// TarsSetTimeout sets the timeout for the servant which is in ms.
func (obj *Notify) TarsSetTimeout(timeout int) {
	obj.servant.TarsSetTimeout(timeout)
}

// TarsSetProtocol sets the protocol for the servant.
func (obj *Notify) TarsSetProtocol(p m.Protocol) {
	obj.servant.TarsSetProtocol(p)
}

type NotifyServant interface {
	ReportServer(sServerName string, sThreadId string, sMessage string) (err error)
	NotifyServer(sServerName string, level NOTIFYLEVEL, sMessage string) (err error)
	GetNotifyInfo(stKey *NotifyKey, stInfo *NotifyInfo) (ret int32, err error)
	ReportNotifyInfo(info *ReportInfo) (err error)
}
type NotifyServantWithContext interface {
	ReportServer(tarsCtx context.Context, sServerName string, sThreadId string, sMessage string) (err error)
	NotifyServer(tarsCtx context.Context, sServerName string, level NOTIFYLEVEL, sMessage string) (err error)
	GetNotifyInfo(tarsCtx context.Context, stKey *NotifyKey, stInfo *NotifyInfo) (ret int32, err error)
	ReportNotifyInfo(tarsCtx context.Context, info *ReportInfo) (err error)
}

// Dispatch is used to call the server side implement for the method defined in the tars file. withContext shows using context or not.
func (obj *Notify) Dispatch(tarsCtx context.Context, val interface{}, tarsReq *requestf.RequestPacket, tarsResp *requestf.ResponsePacket, withContext bool) (err error) {
	var (
		length int32
		have   bool
		ty     byte
	)
	readBuf := codec.NewReader(tools.Int8ToByte(tarsReq.SBuffer))
	buf := codec.NewBuffer()
	switch tarsReq.SFuncName {
	case "reportServer":
		var sServerName string
		var sThreadId string
		var sMessage string

		if tarsReq.IVersion == basef.TARSVERSION {

			err = readBuf.ReadString(&sServerName, 1, true)
			if err != nil {
				return err
			}

			err = readBuf.ReadString(&sThreadId, 2, true)
			if err != nil {
				return err
			}

			err = readBuf.ReadString(&sMessage, 3, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			reqTup := tup.NewUniAttribute()
			reqTup.Decode(readBuf)

			var tupBuffer []byte

			reqTup.GetBuffer("sServerName", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = readBuf.ReadString(&sServerName, 0, true)
			if err != nil {
				return err
			}

			reqTup.GetBuffer("sThreadId", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = readBuf.ReadString(&sThreadId, 0, true)
			if err != nil {
				return err
			}

			reqTup.GetBuffer("sMessage", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = readBuf.ReadString(&sMessage, 0, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.JSONVERSION {
			var jsonData map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(readBuf.ToBytes()))
			decoder.UseNumber()
			err = decoder.Decode(&jsonData)
			if err != nil {
				return fmt.Errorf("decode reqpacket failed, error: %+v", err)
			}
			{
				jsonStr, _ := json.Marshal(jsonData["sServerName"])
				if err = json.Unmarshal(jsonStr, &sServerName); err != nil {
					return err
				}
			}
			{
				jsonStr, _ := json.Marshal(jsonData["sThreadId"])
				if err = json.Unmarshal(jsonStr, &sThreadId); err != nil {
					return err
				}
			}
			{
				jsonStr, _ := json.Marshal(jsonData["sMessage"])
				if err = json.Unmarshal(jsonStr, &sMessage); err != nil {
					return err
				}
			}

		} else {
			err = fmt.Errorf("decode reqpacket fail, error version: %d", tarsReq.IVersion)
			return err
		}

		if !withContext {
			imp := val.(NotifyServant)
			err = imp.ReportServer(sServerName, sThreadId, sMessage)
		} else {
			imp := val.(NotifyServantWithContext)
			err = imp.ReportServer(tarsCtx, sServerName, sThreadId, sMessage)
		}

		if err != nil {
			return err
		}

		if tarsReq.IVersion == basef.TARSVERSION {
			buf.Reset()

		} else if tarsReq.IVersion == basef.TUPVERSION {
			rspTup := tup.NewUniAttribute()

			buf.Reset()
			err = rspTup.Encode(buf)
			if err != nil {
				return err
			}
		} else if tarsReq.IVersion == basef.JSONVERSION {
			rspJson := map[string]interface{}{}

			var rspByte []byte
			if rspByte, err = json.Marshal(rspJson); err != nil {
				return err
			}

			buf.Reset()
			err = buf.WriteSliceUint8(rspByte)
			if err != nil {
				return err
			}
		}
	case "notifyServer":
		var sServerName string
		var level NOTIFYLEVEL
		var sMessage string

		if tarsReq.IVersion == basef.TARSVERSION {

			err = readBuf.ReadString(&sServerName, 1, true)
			if err != nil {
				return err
			}

			err = readBuf.ReadInt32((*int32)(&level), 2, true)
			if err != nil {
				return err
			}

			err = readBuf.ReadString(&sMessage, 3, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			reqTup := tup.NewUniAttribute()
			reqTup.Decode(readBuf)

			var tupBuffer []byte

			reqTup.GetBuffer("sServerName", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = readBuf.ReadString(&sServerName, 0, true)
			if err != nil {
				return err
			}

			reqTup.GetBuffer("level", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = readBuf.ReadInt32((*int32)(&level), 0, true)
			if err != nil {
				return err
			}

			reqTup.GetBuffer("sMessage", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = readBuf.ReadString(&sMessage, 0, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.JSONVERSION {
			var jsonData map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(readBuf.ToBytes()))
			decoder.UseNumber()
			err = decoder.Decode(&jsonData)
			if err != nil {
				return fmt.Errorf("decode reqpacket failed, error: %+v", err)
			}
			{
				jsonStr, _ := json.Marshal(jsonData["sServerName"])
				if err = json.Unmarshal(jsonStr, &sServerName); err != nil {
					return err
				}
			}
			{
				jsonStr, _ := json.Marshal(jsonData["level"])
				if err = json.Unmarshal(jsonStr, &level); err != nil {
					return err
				}
			}
			{
				jsonStr, _ := json.Marshal(jsonData["sMessage"])
				if err = json.Unmarshal(jsonStr, &sMessage); err != nil {
					return err
				}
			}

		} else {
			err = fmt.Errorf("decode reqpacket fail, error version: %d", tarsReq.IVersion)
			return err
		}

		if !withContext {
			imp := val.(NotifyServant)
			err = imp.NotifyServer(sServerName, level, sMessage)
		} else {
			imp := val.(NotifyServantWithContext)
			err = imp.NotifyServer(tarsCtx, sServerName, level, sMessage)
		}

		if err != nil {
			return err
		}

		if tarsReq.IVersion == basef.TARSVERSION {
			buf.Reset()

		} else if tarsReq.IVersion == basef.TUPVERSION {
			rspTup := tup.NewUniAttribute()

			buf.Reset()
			err = rspTup.Encode(buf)
			if err != nil {
				return err
			}
		} else if tarsReq.IVersion == basef.JSONVERSION {
			rspJson := map[string]interface{}{}

			var rspByte []byte
			if rspByte, err = json.Marshal(rspJson); err != nil {
				return err
			}

			buf.Reset()
			err = buf.WriteSliceUint8(rspByte)
			if err != nil {
				return err
			}
		}
	case "getNotifyInfo":
		var stKey NotifyKey
		var stInfo NotifyInfo

		if tarsReq.IVersion == basef.TARSVERSION {

			err = stKey.ReadBlock(readBuf, 1, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			reqTup := tup.NewUniAttribute()
			reqTup.Decode(readBuf)

			var tupBuffer []byte

			reqTup.GetBuffer("stKey", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = stKey.ReadBlock(readBuf, 0, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.JSONVERSION {
			var jsonData map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(readBuf.ToBytes()))
			decoder.UseNumber()
			err = decoder.Decode(&jsonData)
			if err != nil {
				return fmt.Errorf("decode reqpacket failed, error: %+v", err)
			}
			{
				jsonStr, _ := json.Marshal(jsonData["stKey"])
				stKey.ResetDefault()
				if err = json.Unmarshal(jsonStr, &stKey); err != nil {
					return err
				}
			}

		} else {
			err = fmt.Errorf("decode reqpacket fail, error version: %d", tarsReq.IVersion)
			return err
		}

		var funRet int32
		if !withContext {
			imp := val.(NotifyServant)
			funRet, err = imp.GetNotifyInfo(&stKey, &stInfo)
		} else {
			imp := val.(NotifyServantWithContext)
			funRet, err = imp.GetNotifyInfo(tarsCtx, &stKey, &stInfo)
		}

		if err != nil {
			return err
		}

		if tarsReq.IVersion == basef.TARSVERSION {
			buf.Reset()

			err = buf.WriteInt32(funRet, 0)
			if err != nil {
				return err
			}

			err = stInfo.WriteBlock(buf, 2)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			rspTup := tup.NewUniAttribute()

			err = buf.WriteInt32(funRet, 0)
			if err != nil {
				return err
			}

			rspTup.PutBuffer("", buf.ToBytes())
			rspTup.PutBuffer("tars_ret", buf.ToBytes())

			buf.Reset()
			err = stInfo.WriteBlock(buf, 0)
			if err != nil {
				return err
			}

			rspTup.PutBuffer("stInfo", buf.ToBytes())

			buf.Reset()
			err = rspTup.Encode(buf)
			if err != nil {
				return err
			}
		} else if tarsReq.IVersion == basef.JSONVERSION {
			rspJson := map[string]interface{}{}
			rspJson["tars_ret"] = funRet
			rspJson["stInfo"] = stInfo

			var rspByte []byte
			if rspByte, err = json.Marshal(rspJson); err != nil {
				return err
			}

			buf.Reset()
			err = buf.WriteSliceUint8(rspByte)
			if err != nil {
				return err
			}
		}
	case "reportNotifyInfo":
		var info ReportInfo

		if tarsReq.IVersion == basef.TARSVERSION {

			err = info.ReadBlock(readBuf, 1, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.TUPVERSION {
			reqTup := tup.NewUniAttribute()
			reqTup.Decode(readBuf)

			var tupBuffer []byte

			reqTup.GetBuffer("info", &tupBuffer)
			readBuf.Reset(tupBuffer)
			err = info.ReadBlock(readBuf, 0, true)
			if err != nil {
				return err
			}

		} else if tarsReq.IVersion == basef.JSONVERSION {
			var jsonData map[string]interface{}
			decoder := json.NewDecoder(bytes.NewReader(readBuf.ToBytes()))
			decoder.UseNumber()
			err = decoder.Decode(&jsonData)
			if err != nil {
				return fmt.Errorf("decode reqpacket failed, error: %+v", err)
			}
			{
				jsonStr, _ := json.Marshal(jsonData["info"])
				info.ResetDefault()
				if err = json.Unmarshal(jsonStr, &info); err != nil {
					return err
				}
			}

		} else {
			err = fmt.Errorf("decode reqpacket fail, error version: %d", tarsReq.IVersion)
			return err
		}

		if !withContext {
			imp := val.(NotifyServant)
			err = imp.ReportNotifyInfo(&info)
		} else {
			imp := val.(NotifyServantWithContext)
			err = imp.ReportNotifyInfo(tarsCtx, &info)
		}

		if err != nil {
			return err
		}

		if tarsReq.IVersion == basef.TARSVERSION {
			buf.Reset()

		} else if tarsReq.IVersion == basef.TUPVERSION {
			rspTup := tup.NewUniAttribute()

			buf.Reset()
			err = rspTup.Encode(buf)
			if err != nil {
				return err
			}
		} else if tarsReq.IVersion == basef.JSONVERSION {
			rspJson := map[string]interface{}{}

			var rspByte []byte
			if rspByte, err = json.Marshal(rspJson); err != nil {
				return err
			}

			buf.Reset()
			err = buf.WriteSliceUint8(rspByte)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("func mismatch")
	}
	var statusMap map[string]string
	if status, ok := current.GetResponseStatus(tarsCtx); ok && status != nil {
		statusMap = status
	}
	var contextMap map[string]string
	if ctx, ok := current.GetResponseContext(tarsCtx); ok && ctx != nil {
		contextMap = ctx
	}
	*tarsResp = requestf.ResponsePacket{
		IVersion:     tarsReq.IVersion,
		CPacketType:  0,
		IRequestId:   tarsReq.IRequestId,
		IMessageType: 0,
		IRet:         0,
		SBuffer:      tools.ByteToInt8(buf.ToBytes()),
		Status:       statusMap,
		SResultDesc:  "",
		Context:      contextMap,
	}

	_ = readBuf
	_ = buf
	_ = length
	_ = have
	_ = ty
	return nil
}
