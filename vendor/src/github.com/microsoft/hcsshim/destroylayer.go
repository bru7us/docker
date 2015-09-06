package hcsshim

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/Sirupsen/logrus"
)

// DestroyLayer will remove the on-disk files representing the layer with the given
// id, including that layer's containing folder, if any.
func DestroyLayer(info DriverInfo, id string) error {
	title := "hcsshim::DestroyLayer "
	logrus.Debugf(title+"Flavour %s ID %s", info.Flavour, id)

	// Load the DLL and get a handle to the procedure we need
	dll, proc, err := loadAndFind(procDestroyLayer)
	if dll != nil {
		defer dll.Release()
	}
	if err != nil {
		return err
	}

	// Convert id to uint16 pointer for calling the procedure
	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		err = fmt.Errorf(title+" - Failed conversion of id %s to pointer %s", id, err)
		logrus.Error(err)
		return err
	}

	// Convert info to API calling convention
	infop, err := convertDriverInfo(info)
	if err != nil {
		err = fmt.Errorf(title+" - Failed conversion info struct %s", err)
		logrus.Error(err)
		return err
	}

	// Call the procedure itself.
	r1, _, _ := proc.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)))

	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))

	if r1 != 0 {
		err = fmt.Errorf(title+" - Win32 API call returned error r1=%d err=%s id=%s flavour=%d",
			r1, syscall.Errno(r1), id, info.Flavour)
		logrus.Error(err)
		return err
	}

	logrus.Debugf(title+" - succeeded id=%s flavour=%d", id, info.Flavour)
	return nil
}
