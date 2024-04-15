//go:build windows

package edge

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type _ICoreWebView2NewWindowRequestedEventArgsVtbl struct {
	_IUnknownVtbl
	GetUri                     ComProc
	PutNewWindow               ComProc
	GetNewWindow               ComProc
	PutHandled                 ComProc
	GetHandled                 ComProc
	GetIsUserInitiated         ComProc
	GetWindowFeatures          ComProc // Might not be ordered correctly
	GetOriginalSourceFrameInfo ComProc // Might not be ordered correctly
	GetName                    ComProc
	GetDeferral                ComProc // Might not be ordered correctly, but it's the only actual method so I assume it's last :grimace:
}

type ICoreWebView2NewWindowRequestedEventArgs struct {
	vtbl *_ICoreWebView2NewWindowRequestedEventArgsVtbl
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) AddRef() uintptr {
	return i.AddRef()
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetHandled() (bool, error) {
	// var result bool
	var _result int32
	_, _, err := i.vtbl.GetHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_result)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}

	result := _result != 0
	return result, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) PutHandled(handled bool) error {
	var err error
	_, _, err = i.vtbl.PutHandled.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(boolToInt(handled)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetNewWindow() (*ICoreWebView2, error) {
	var err error
	var window *ICoreWebView2
	_, _, err = i.vtbl.GetNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&window)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	return window, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) PutNewWindow(window *ICoreWebView2) error {
	var err error
	_, _, err = i.vtbl.PutNewWindow.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(window)),
	)
	if err != windows.ERROR_SUCCESS {
		return err
	}
	return nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetIsUserInitiated() (bool, error) {
	var err error
	var _result int32
	_, _, err = i.vtbl.GetIsUserInitiated.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_result)),
	)
	if err != windows.ERROR_SUCCESS {
		return false, err
	}
	result := _result != 0
	return result, nil
}

// func (i *ICoreWebView2NewWindowRequestedEventArgs) GetDeferral() (*ICoreWebView2Deferral, error) {
// 	var deferral *ICoreWebView2Deferral
//
// 	hr, _, err := i.vtbl.GetDeferral.Call(
// 		uintptr(unsafe.Pointer(i)),
// 		uintptr(unsafe.Pointer(&deferral)),
// 	)
// 	if windows.Handle(hr) != windows.S_OK {
// 		return nil, syscall.Errno(hr)
// 	}
//
// 	if deferral == nil {
// 		if err == nil {
// 			err = fmt.Errorf("unknown error")
// 		}
// 		return nil, err
// 	}
// 	return deferral, nil
// }

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetName() (string, error) {
	var _name *uint16
	res, _, err := i.vtbl.GetName.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_name)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	}
	if windows.Handle(res) != windows.S_OK {
		return "", syscall.Errno(res)
	}

	name := UTF16PtrToString(_name)
	CoTaskMemFree(unsafe.Pointer(_name))
	return name, nil
}

func (i *ICoreWebView2NewWindowRequestedEventArgs) GetUri() (string, error) {
	var err error
	// Create *uint16 to hold result
	var _uri *uint16
	res, _, err := i.vtbl.GetUri.Call(
		uintptr(unsafe.Pointer(i)),
		uintptr(unsafe.Pointer(&_uri)),
	)
	if err != windows.ERROR_SUCCESS {
		return "", err
	} // Get result and cleanup
	if windows.Handle(res) != windows.S_OK {
		return "", syscall.Errno(res)
	}
	uri := windows.UTF16PtrToString(_uri)
	windows.CoTaskMemFree(unsafe.Pointer(_uri))
	return uri, nil
}
