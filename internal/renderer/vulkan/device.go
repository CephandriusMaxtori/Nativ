//go:build ignore
package vulkan

import (
	"fmt"

	vk "github.com/vulkan-go/vulkan"
)

type queueFamilyIndices struct {
	graphics uint32
	present  uint32
	separate bool
}

func pickPhysicalDevice(instance vk.Instance, surface vk.SurfaceKHR) (vk.PhysicalDevice, queueFamilyIndices, error) {
	var count uint32
	vk.EnumeratePhysicalDevices(instance, &count, nil)
	if count == 0 {
		return nil, queueFamilyIndices{}, fmt.Errorf("no vulkan physical devices found")
	}

	devices := make([]vk.PhysicalDevice, count)
	vk.EnumeratePhysicalDevices(instance, &count, devices)

	// Score devices: prefer discrete GPU with graphics+present
	type scored struct {
		device vk.PhysicalDevice
		score  int
		families queueFamilyIndices
	}

	var best scored
	for _, dev := range devices {
		var props vk.PhysicalDeviceProperties
		vk.GetPhysicalDeviceProperties(dev, &props)

		families, ok := findQueueFamilies(dev, surface)
		if !ok {
			continue
		}

		score := 0
		if props.DeviceType == vk.PhysicalDeviceTypeDiscreteGpu {
			score += 1000
		}
		score += int(props.Limits.MaxImageDimension2D)

		if score > best.score {
			best = scored{dev, score, families}
		}
	}

	if best.device == nil {
		return nil, queueFamilyIndices{}, fmt.Errorf("no suitable vulkan device found")
	}
	return best.device, best.families, nil
}

func findQueueFamilies(device vk.PhysicalDevice, surface vk.SurfaceKHR) (queueFamilyIndices, bool) {
	var count uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(device, &count, nil)
	props := make([]vk.QueueFamilyProperties, count)
	vk.GetPhysicalDeviceQueueFamilyProperties(device, &count, props)

	var indices queueFamilyIndices
	foundGraphics := false
	foundPresent := false

	for i, prop := range props {
		if prop.QueueCount == 0 {
			continue
		}

		if prop.QueueFlags&vk.QueueFlags(vk.QueueGraphicsBit) != 0 {
			indices.graphics = uint32(i)
			foundGraphics = true
		}

		var presentSupport uint32
		result := vk.GetPhysicalDeviceSurfaceSupportKHR(device, uint32(i), surface, &presentSupport)
		if result == vk.Success && presentSupport != 0 {
			indices.present = uint32(i)
			foundPresent = true
		}

		if foundGraphics && foundPresent {
			break
		}
	}

	if !foundGraphics || !foundPresent {
		return queueFamilyIndices{}, false
	}

	indices.separate = indices.graphics != indices.present
	return indices, true
}

func createLogicalDevice(physicalDevice vk.PhysicalDevice, families queueFamilyIndices) (vk.Device, vk.Queue, vk.Queue, error) {
	uniqueQueues := []uint32{families.graphics}
	if families.separate {
		uniqueQueues = append(uniqueQueues, families.present)
	}

	queueCreateInfos := make([]vk.DeviceQueueCreateInfo, len(uniqueQueues))
	for i, qf := range uniqueQueues {
		priority := float32(1.0)
		queueCreateInfos[i] = vk.DeviceQueueCreateInfo{
			SType:            vk.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: qf,
			QueueCount:       1,
			PQueuePriorities: []float32{priority},
		}
	}

	extensions := []string{"VK_KHR_swapchain"}

	createInfo := &vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount:    uint32(len(queueCreateInfos)),
		PQueueCreateInfos:       queueCreateInfos,
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
	}

	var device vk.Device
	result := vk.CreateDevice(physicalDevice, createInfo, nil, &device)
	if result != vk.Success {
		return nil, nil, nil, fmt.Errorf("failed to create logical device: %v", vk.Error(result))
	}

	var graphicsQueue vk.Queue
	vk.GetDeviceQueue(device, families.graphics, 0, &graphicsQueue)

	var presentQueue vk.Queue
	if families.separate {
		vk.GetDeviceQueue(device, families.present, 0, &presentQueue)
	} else {
		presentQueue = graphicsQueue
	}

	return device, graphicsQueue, presentQueue, nil
}

func destroyDevice(device vk.Device) {
	if device != nil {
		vk.DestroyDevice(device, nil)
	}
}
