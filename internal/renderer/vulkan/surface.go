//go:build ignore
package vulkan

import (
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

func createSurface(instance vk.Instance, hwnd, hinstance uintptr) (vk.SurfaceKHR, error) {
	surfaceInfo := &vk.Win32SurfaceCreateInfoKHR{
		SType:     vk.StructureTypeWin32SurfaceCreateInfoKHR,
		Hinstance: unsafe.Pointer(hinstance),
		Hwnd:      unsafe.Pointer(hwnd),
	}

	var surface vk.SurfaceKHR
	result := vk.CreateWin32SurfaceKHR(instance, surfaceInfo, nil, &surface)
	if result != vk.Success {
		return nil, vk.Error(result)
	}
	return surface, nil
}

func destroySurface(instance vk.Instance, surface vk.SurfaceKHR) {
	if surface != nil {
		vk.DestroySurfaceKHR(instance, surface, nil)
	}
}
