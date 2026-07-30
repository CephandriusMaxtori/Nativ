//go:build ignore
package vulkan

import vk "github.com/vulkan-go/vulkan"

func createInstance() (vk.Instance, error) {
	err := vk.Init()
	if err != nil {
		return nil, err
	}

	appInfo := &vk.ApplicationInfo{
		SType:              vk.StructureTypeApplicationInfo,
		PApplicationName:   "Nativ",
		ApplicationVersion: vk.MakeVersion(1, 0, 0),
		PEngineName:        "Nativ",
		EngineVersion:      vk.MakeVersion(1, 0, 0),
		ApiVersion:         vk.MakeVersion(1, 2, 0),
	}

	extensions := []string{
		"VK_KHR_surface",
		"VK_KHR_win32_surface",
	}

	createInfo := &vk.InstanceCreateInfo{
		SType:                   vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo:        appInfo,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
	}

	var instance vk.Instance
	result := vk.CreateInstance(createInfo, nil, &instance)
	if result != vk.Success {
		return nil, vk.Error(result)
	}
	return instance, nil
}

func destroyInstance(instance vk.Instance) {
	if instance != nil {
		vk.DestroyInstance(instance, nil)
	}
}
