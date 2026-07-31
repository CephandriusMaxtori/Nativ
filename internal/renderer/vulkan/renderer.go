//go:build ignore
package vulkan

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"

	"github.com/CephandriusMaxtori/Nativ/internal/platform"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
	vk "github.com/vulkan-go/vulkan"
)

func init() {
	renderer.RegisterBackend("vulkan", func() renderer.Renderer {
		return &vulkanRenderer{}
	})
}

type vulkanRenderer struct {
	instance       vk.Instance
	physicalDevice vk.PhysicalDevice
	device         vk.Device
	graphicsQueue  vk.Queue
	presentQueue   vk.Queue
	surface        vk.SurfaceKHR
	swapchainInfo  *swapchainInfo
	pipeline       *graphicsPipeline
	families       queueFamilyIndices
	width, height int
	initialized    bool
}

func (r *vulkanRenderer) Init(plat platform.Platform, width, height int) error {
	var err error

	r.instance, err = createInstance()
	if err != nil {
		return fmt.Errorf("vulkan: %w", err)
	}

	r.surface, err = createSurface(r.instance, plat.Hwnd(), plat.Hdc())
	if err != nil {
		r.destroy()
		return fmt.Errorf("vulkan: %w", err)
	}

	r.physicalDevice, r.families, err = pickPhysicalDevice(r.instance, r.surface)
	if err != nil {
		r.destroy()
		return fmt.Errorf("vulkan: %w", err)
	}

	r.device, r.graphicsQueue, r.presentQueue, err = createLogicalDevice(r.physicalDevice, r.families)
	if err != nil {
		r.destroy()
		return fmt.Errorf("vulkan: %w", err)
	}

	// Default size; will be updated on first Draw
	r.width = 800
	r.height = 600

	r.initialized = true
	return nil
}

func (r *vulkanRenderer) ensureSwapchain(width, height int) error {
	changed := r.swapchainInfo == nil || r.width != width || r.height != height
	if !changed {
		return nil
	}
	if r.swapchainInfo != nil {
		if r.pipeline != nil {
			destroyPipelineObjects(r.device, r.pipeline)
			r.pipeline = nil
		}
		destroySwapchain(r.device, r.swapchainInfo)
		r.swapchainInfo = nil
	}

	r.width = width
	r.height = height

	info, err := createSwapchain(r.device, r.physicalDevice, r.surface, width, height)
	if err != nil {
		return err
	}
	r.swapchainInfo = info

	renderPass, err := createRenderPass(r.device, info.imageFormat)
	if err != nil {
		return err
	}

	gp, err := createGraphicsPipeline(r.device, renderPass, info.extent)
	if err != nil {
		vk.DestroyRenderPass(r.device, renderPass, nil)
		return err
	}
	gp.renderPass = renderPass

	framebuffers, err := createFramebuffers(r.device, info.imageViews, renderPass, info.extent)
	if err != nil {
		return err
	}
	info.framebuffers = framebuffers

	gp.vertexBuffer, gp.vertexMemory, err = createVertexBuffer(r.device, r.physicalDevice)
	if err != nil {
		return err
	}

	// Use graphics queue family for command pool
	gp.commandPool, err = createCommandPool(r.device, r.families.graphics)
	if err != nil {
		return err
	}

	gp.commandBuffers, err = createCommandBuffers(r.device, gp.commandPool, info.imageCount)
	if err != nil {
		return err
	}

	gp.imageAvailable, gp.renderFinished, gp.inFlightFence, err = createSyncObjects(r.device)
	if err != nil {
		return err
	}

	r.pipeline = gp
	return nil
}

func (r *vulkanRenderer) Draw(cmds []renderer.DrawCommand) error {
	if !r.initialized {
		return nil
	}

	// Get current window size from platform if available
	width, height := r.width, r.height
	if err := r.ensureSwapchain(width, height); err != nil {
		return err
	}

	// Wait for previous frame
	_, _ = vk.WaitForFences(r.device, 1, []vk.Fence{r.pipeline.inFlightFence}, vk.True, math.MaxUint64)
	vk.ResetFences(r.device, 1, []vk.Fence{r.pipeline.inFlightFence})

	// Acquire next image
	var imageIndex uint32
	result := vk.AcquireNextImageKHR(r.device, r.swapchainInfo.swapchain, math.MaxUint64, r.pipeline.imageAvailable, nil, &imageIndex)
	if result == vk.ErrorOutOfDateKHR {
		r.swapchainInfo = nil
		return nil
	}
	if result != vk.Success && result != vk.SuboptimalKHR {
		return fmt.Errorf("failed to acquire swapchain image: %v", vk.Error(result))
	}

	// Update vertex buffer
	if err := r.updateVertexBuffer(cmds); err != nil {
		return err
	}

	// Record command buffer
	cb := r.pipeline.commandBuffers[imageIndex]
	vk.ResetCommandBuffer(cb, 0)

	beginInfo := &vk.CommandBufferBeginInfo{
		SType: vk.StructureTypeCommandBufferBeginInfo,
		Flags: vk.CommandBufferUsageFlags(vk.CommandBufferUsageSimultaneousUseBit),
	}
	vk.BeginCommandBuffer(cb, beginInfo)

	clearColor := vk.ClearValue{
		Color: vk.ClearColorValue{
			Float32: [4]float32{0.1, 0.1, 0.12, 1.0},
		},
	}

	renderPassBeginInfo := &vk.RenderPassBeginInfo{
		SType:        vk.StructureTypeRenderPassBeginInfo,
		RenderPass:   r.pipeline.renderPass,
		Framebuffer:  r.swapchainInfo.framebuffers[imageIndex],
		RenderArea: vk.Rect2D{
			Offset: vk.Offset2D{X: 0, Y: 0},
			Extent: r.swapchainInfo.extent,
		},
		ClearValueCount: 1,
		PClearValues:    []vk.ClearValue{clearColor},
	}

	vk.CmdBeginRenderPass(cb, renderPassBeginInfo, vk.SubpassContentsInline)
	vk.CmdBindPipeline(cb, vk.PipelineBindPointGraphics, r.pipeline.pipeline)

	buffers := []vk.Buffer{r.pipeline.vertexBuffer}
	offsets := []uintptr{0}
	vk.CmdBindVertexBuffers(cb, 0, 1, buffers, offsets)

	quadCount := len(cmds)
	if quadCount > maxQuads {
		quadCount = maxQuads
	}
	if quadCount > 0 {
		vk.CmdDraw(cb, uint32(quadCount*6), 1, 0, 0)
	}

	vk.CmdEndRenderPass(cb)
	vk.EndCommandBuffer(cb)

	// Submit
	submitInfo := &vk.SubmitInfo{
		SType:                vk.StructureTypeSubmitInfo,
		WaitSemaphoreCount:   1,
		PWaitSemaphores:      []vk.Semaphore{r.pipeline.imageAvailable},
		PWaitDstStageMask:    []vk.PipelineStageFlags{vk.PipelineStageFlags(vk.PipelineStageColorAttachmentOutputBit)},
		CommandBufferCount:   1,
		PCommandBuffers:      []vk.CommandBuffer{cb},
		SignalSemaphoreCount: 1,
		PSignalSemaphores:    []vk.Semaphore{r.pipeline.renderFinished},
	}

	vk.QueueSubmit(r.graphicsQueue, 1, []*vk.SubmitInfo{submitInfo}, r.pipeline.inFlightFence)

	// Present
	presentInfo := &vk.PresentInfoKHR{
		SType:              vk.StructureTypePresentInfoKHR,
		WaitSemaphoreCount: 1,
		PWaitSemaphores:    []vk.Semaphore{r.pipeline.renderFinished},
		SwapchainCount:     1,
		PSwapchains:        []vk.SwapchainKHR{r.swapchainInfo.swapchain},
		PImageIndices:      []uint32{imageIndex},
	}

	result2 := vk.QueuePresentKHR(r.presentQueue, presentInfo)
	if result2 == vk.ErrorOutOfDateKHR || result2 == vk.SuboptimalKHR {
		r.swapchainInfo = nil
	}

	return nil
}

func (r *vulkanRenderer) updateVertexBuffer(cmds []renderer.DrawCommand) error {
	if len(cmds) == 0 {
		return nil
	}
	if len(cmds) > maxQuads {
		cmds = cmds[:maxQuads]
	}

	vertexSize := 6 * 4 // 6 floats * 4 bytes
	bufferSize := len(cmds) * 6 * vertexSize

	var data unsafe.Pointer
	result := vk.MapMemory(r.device, r.pipeline.vertexMemory, 0, vk.WholeSize, 0, &data)
	if result != vk.Success {
		return fmt.Errorf("failed to map vertex memory: %v", vk.Error(result))
	}

	w := float32(r.width)
	h := float32(r.height)
	halfW := w * 0.5
	halfH := h * 0.5

	buf := unsafe.Slice((*byte)(data), bufferSize)
	offset := 0

	for _, cmd := range cmds {
		x1 := cmd.X
		y1 := cmd.Y
		x2 := cmd.X + cmd.Width
		y2 := cmd.Y + cmd.Height

		// Convert to NDC
		nx1 := (x1 - halfW) / halfW
		ny1 := -(y1 - halfH) / halfH
		nx2 := (x2 - halfW) / halfW
		ny2 := -(y2 - halfH) / halfH

		verts := [6]struct {
			px, py float32
			r, g, b, a float32
		}{
			{nx1, ny1, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx2, ny1, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx2, ny2, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx1, ny1, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx2, ny2, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx1, ny2, cmd.R, cmd.G, cmd.B, cmd.A},
		}

		for _, v := range verts {
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v.px)); offset += 4
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v.py)); offset += 4
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v.r)); offset += 4
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v.g)); offset += 4
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v.b)); offset += 4
			binary.LittleEndian.PutUint32(buf[offset:], math.Float32bits(v.a)); offset += 4
		}
	}

	vk.UnmapMemory(r.device, r.pipeline.vertexMemory)
	return nil
}

func (r *vulkanRenderer) Shutdown() error {
	if r.device != nil {
		vk.DeviceWaitIdle(r.device)
	}
	r.destroy()
	return nil
}

func (r *vulkanRenderer) destroy() {
	if r.pipeline != nil {
		destroyPipelineObjects(r.device, r.pipeline)
	}
	if r.swapchainInfo != nil {
		destroySwapchain(r.device, r.swapchainInfo)
	}
	destroyDevice(r.device)
	if r.surface != nil {
		destroySurface(r.instance, r.surface)
	}
	destroyInstance(r.instance)
	r.initialized = false
}
