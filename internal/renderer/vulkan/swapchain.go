//go:build ignore
package vulkan

import (
	"fmt"
	"math"

	vk "github.com/vulkan-go/vulkan"
)

type swapchainSupport struct {
	capabilities vk.SurfaceCapabilitiesKHR
	formats      []vk.SurfaceFormatKHR
	presentModes []vk.PresentModeKHR
}

func querySwapchainSupport(device vk.PhysicalDevice, surface vk.SurfaceKHR) swapchainSupport {
	var support swapchainSupport
	vk.GetPhysicalDeviceSurfaceCapabilitiesKHR(device, surface, &support.capabilities)

	var formatCount uint32
	vk.GetPhysicalDeviceSurfaceFormatsKHR(device, surface, &formatCount, nil)
	if formatCount > 0 {
		support.formats = make([]vk.SurfaceFormatKHR, formatCount)
		vk.GetPhysicalDeviceSurfaceFormatsKHR(device, surface, &formatCount, support.formats)
	}

	var presentCount uint32
	vk.GetPhysicalDeviceSurfacePresentModesKHR(device, surface, &presentCount, nil)
	if presentCount > 0 {
		support.presentModes = make([]vk.PresentModeKHR, presentCount)
		vk.GetPhysicalDeviceSurfacePresentModesKHR(device, surface, &presentCount, support.presentModes)
	}

	return support
}

type swapchainInfo struct {
	swapchain       vk.SwapchainKHR
	images          []vk.Image
	imageViews      []vk.ImageView
	framebuffers    []vk.Framebuffer
	extent          vk.Extent2D
	imageFormat     vk.Format
	imageCount      uint32
}

func chooseSwapSurfaceFormat(formats []vk.SurfaceFormatKHR) vk.SurfaceFormatKHR {
	for _, f := range formats {
		if f.Format == vk.FormatB8g8r8a8Srgb && f.ColorSpace == vk.ColorSpaceSrgbNonlinear {
			return f
		}
	}
	return formats[0]
}

func chooseSwapPresentMode(modes []vk.PresentModeKHR) vk.PresentModeKHR {
	for _, m := range modes {
		if m == vk.PresentModeMailbox {
			return m
		}
	}
	return vk.PresentModeFifo
}

func chooseSwapExtent(caps vk.SurfaceCapabilitiesKHR, width, height int) vk.Extent2D {
	if caps.CurrentExtent.Width != math.MaxUint32 {
		return caps.CurrentExtent
	}
	w := clamp(uint32(width), caps.MinImageExtent.Width, caps.MaxImageExtent.Width)
	h := clamp(uint32(height), caps.MinImageExtent.Height, caps.MaxImageExtent.Height)
	return vk.Extent2D{Width: w, Height: h}
}

func clamp(v, min, max uint32) uint32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func createSwapchain(device vk.Device, physicalDevice vk.PhysicalDevice, surface vk.SurfaceKHR, width, height int) (*swapchainInfo, error) {
	support := querySwapchainSupport(physicalDevice, surface)

	format := chooseSwapSurfaceFormat(support.formats)
	presentMode := chooseSwapPresentMode(support.presentModes)
	extent := chooseSwapExtent(support.capabilities, width, height)

	imageCount := support.capabilities.MinImageCount + 1
	if support.capabilities.MaxImageCount > 0 && imageCount > support.capabilities.MaxImageCount {
		imageCount = support.capabilities.MaxImageCount
	}

	createInfo := &vk.SwapchainCreateInfoKHR{
		SType:           vk.StructureTypeSwapchainCreateInfoKHR,
		Surface:         surface,
		MinImageCount:   imageCount,
		ImageFormat:     format.Format,
		ImageColorSpace: format.ColorSpace,
		ImageExtent:     extent,
		ImageArrayLayers: 1,
		ImageUsage:      vk.ImageUsageFlags(vk.ImageUsageColorAttachmentBit),
		PreTransform:    support.capabilities.CurrentTransform,
		CompositeAlpha:  vk.CompositeAlphaOpaqueBit,
		PresentMode:     presentMode,
		Clipped:         vk.True,
		OldSwapchain:    nil,
	}

	var swapchain vk.SwapchainKHR
	result := vk.CreateSwapchainKHR(device, createInfo, nil, &swapchain)
	if result != vk.Success {
		return nil, fmt.Errorf("failed to create swapchain: %v", vk.Error(result))
	}

	vk.GetSwapchainImagesKHR(device, swapchain, &imageCount, nil)
	images := make([]vk.Image, imageCount)
	vk.GetSwapchainImagesKHR(device, swapchain, &imageCount, images)

	imageViews := make([]vk.ImageView, len(images))
	for i, img := range images {
		viewInfo := &vk.ImageViewCreateInfo{
			SType:    vk.StructureTypeImageViewCreateInfo,
			Image:    img,
			ViewType: vk.ImageViewType2d,
			Format:   format.Format,
			Components: vk.ComponentMapping{
				R: vk.ComponentSwizzleIdentity,
				G: vk.ComponentSwizzleIdentity,
				B: vk.ComponentSwizzleIdentity,
				A: vk.ComponentSwizzleIdentity,
			},
			SubresourceRange: vk.ImageSubresourceRange{
				AspectMask:     vk.ImageAspectFlags(vk.ImageAspectColorBit),
				BaseMipLevel:   0,
				LevelCount:     1,
				BaseArrayLayer: 0,
				LayerCount:     1,
			},
		}
		vk.CreateImageView(device, viewInfo, nil, &imageViews[i])
	}

	return &swapchainInfo{
		swapchain:   swapchain,
		images:      images,
		imageViews:  imageViews,
		extent:      extent,
		imageFormat: format.Format,
		imageCount:  imageCount,
	}, nil
}

func createFramebuffers(device vk.Device, imageViews []vk.ImageView, renderPass vk.RenderPass, extent vk.Extent2D) ([]vk.Framebuffer, error) {
	framebuffers := make([]vk.Framebuffer, len(imageViews))
	for i, view := range imageViews {
		info := &vk.FramebufferCreateInfo{
			SType:           vk.StructureTypeFramebufferCreateInfo,
			RenderPass:      renderPass,
			AttachmentCount: 1,
			PAttachments:    []vk.ImageView{view},
			Width:           extent.Width,
			Height:          extent.Height,
			Layers:          1,
		}
		result := vk.CreateFramebuffer(device, info, nil, &framebuffers[i])
		if result != vk.Success {
			return nil, fmt.Errorf("failed to create framebuffer: %v", vk.Error(result))
		}
	}
	return framebuffers, nil
}

func destroySwapchain(device vk.Device, info *swapchainInfo) {
	if info == nil {
		return
	}
	for _, fb := range info.framebuffers {
		vk.DestroyFramebuffer(device, fb, nil)
	}
	for _, view := range info.imageViews {
		vk.DestroyImageView(device, view, nil)
	}
	if info.swapchain != nil {
		vk.DestroySwapchainKHR(device, info.swapchain, nil)
	}
}
