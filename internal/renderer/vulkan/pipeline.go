//go:build ignore
package vulkan

import (
	"fmt"
	"unsafe"

	vk "github.com/vulkan-go/vulkan"
)

const maxQuads = 1024

type graphicsPipeline struct {
	pipeline       vk.Pipeline
	pipelineLayout vk.PipelineLayout
	renderPass     vk.RenderPass
	vertexBuffer   vk.Buffer
	vertexMemory   vk.DeviceMemory
	commandPool    vk.CommandPool
	commandBuffers []vk.CommandBuffer
	imageAvailable vk.Semaphore
	renderFinished vk.Semaphore
	inFlightFence  vk.Fence
}

func createRenderPass(device vk.Device, format vk.Format) (vk.RenderPass, error) {
	colorAttachment := vk.AttachmentDescription{
		Format:        format,
		Samples:       vk.SampleCount1Bit,
		LoadOp:        vk.AttachmentLoadOpClear,
		StoreOp:       vk.AttachmentStoreOpStore,
		StencilLoadOp: vk.AttachmentLoadOpDontCare,
		StencilStoreOp: vk.AttachmentStoreOpDontCare,
		InitialLayout: vk.ImageLayoutUndefined,
		FinalLayout:   vk.ImageLayoutPresentSrc,
	}

	colorAttachmentRef := vk.AttachmentReference{
		Attachment: 0,
		Layout:     vk.ImageLayoutColorAttachmentOptimal,
	}

	subpass := vk.SubpassDescription{
		PipelineBindPoint:    vk.PipelineBindPointGraphics,
		ColorAttachmentCount: 1,
		PColorAttachments:    []vk.AttachmentReference{colorAttachmentRef},
	}

	createInfo := &vk.RenderPassCreateInfo{
		SType:           vk.StructureTypeRenderPassCreateInfo,
		AttachmentCount: 1,
		PAttachments:    []vk.AttachmentDescription{colorAttachment},
		SubpassCount:    1,
		PSubpasses:      []vk.SubpassDescription{subpass},
	}

	var renderPass vk.RenderPass
	result := vk.CreateRenderPass(device, createInfo, nil, &renderPass)
	if result != vk.Success {
		return nil, fmt.Errorf("failed to create render pass: %v", vk.Error(result))
	}
	return renderPass, nil
}

func createShaderModule(device vk.Device, code []byte) (vk.ShaderModule, error) {
	createInfo := &vk.ShaderModuleCreateInfo{
		SType:    vk.StructureTypeShaderModuleCreateInfo,
		CodeSize: uintptr(len(code)),
		PCode:    unsafe.Pointer(&code[0]),
	}

	var module vk.ShaderModule
	result := vk.CreateShaderModule(device, createInfo, nil, &module)
	if result != vk.Success {
		return nil, fmt.Errorf("failed to create shader module: %v", vk.Error(result))
	}
	return module, nil
}

func createGraphicsPipeline(device vk.Device, renderPass vk.RenderPass, extent vk.Extent2D) (*graphicsPipeline, error) {
	vertCode := vertexShaderSPIRV()
	fragCode := fragmentShaderSPIRV()

	vertModule, err := createShaderModule(device, vertCode)
	if err != nil {
		return nil, err
	}
	defer vk.DestroyShaderModule(device, vertModule, nil)

	fragModule, err := createShaderModule(device, fragCode)
	if err != nil {
		return nil, err
	}
	defer vk.DestroyShaderModule(device, fragModule, nil)

	mainFn := "main"
	vertStage := vk.PipelineShaderStageCreateInfo{
		SType:  vk.StructureTypePipelineShaderStageCreateInfo,
		Stage:  vk.ShaderStageFlags(vk.ShaderStageVertexBit),
		Module: vertModule,
		PName:  mainFn,
	}

	fragStage := vk.PipelineShaderStageCreateInfo{
		SType:  vk.StructureTypePipelineShaderStageCreateInfo,
		Stage:  vk.ShaderStageFlags(vk.ShaderStageFragmentBit),
		Module: fragModule,
		PName:  mainFn,
	}

	vertexBinding := vk.VertexInputBindingDescription{
		Binding:   0,
		Stride:    24, // 6 floats * 4 bytes
		InputRate: vk.VertexInputRateVertex,
	}

	positionAttrib := vk.VertexInputAttributeDescription{
		Binding:  0,
		Location: 0,
		Format:   vk.FormatR32g32Sfloat,
		Offset:   0,
	}

	colorAttrib := vk.VertexInputAttributeDescription{
		Binding:  0,
		Location: 1,
		Format:   vk.FormatR32g32b32a32Sfloat,
		Offset:   8, // after vec2 (8 bytes)
	}

	vertexInput := vk.PipelineVertexInputStateCreateInfo{
		SType:                           vk.StructureTypePipelineVertexInputStateCreateInfo,
		VertexBindingDescriptionCount:   1,
		PVertexBindingDescriptions:      []vk.VertexInputBindingDescription{vertexBinding},
		VertexAttributeDescriptionCount: 2,
		PVertexAttributeDescriptions:    []vk.VertexInputAttributeDescription{positionAttrib, colorAttrib},
	}

	inputAssembly := vk.PipelineInputAssemblyStateCreateInfo{
		SType:                  vk.StructureTypePipelineInputAssemblyStateCreateInfo,
		Topology:               vk.PrimitiveTopologyTriangleList,
		PrimitiveRestartEnable: vk.False,
	}

	viewport := vk.Viewport{
		X:        0,
		Y:        0,
		Width:    float32(extent.Width),
		Height:   float32(extent.Height),
		MinDepth: 0,
		MaxDepth: 1,
	}

	scissor := vk.Rect2D{
		Offset: vk.Offset2D{X: 0, Y: 0},
		Extent: extent,
	}

	viewportState := vk.PipelineViewportStateCreateInfo{
		SType:         vk.StructureTypePipelineViewportStateCreateInfo,
		ViewportCount: 1,
		PViewports:    []vk.Viewport{viewport},
		ScissorCount:  1,
		PScissors:     []vk.Rect2D{scissor},
	}

	rasterizer := vk.PipelineRasterizationStateCreateInfo{
		SType:                   vk.StructureTypePipelineRasterizationStateCreateInfo,
		DepthClampEnable:        vk.False,
		RasterizerDiscardEnable: vk.False,
		PolygonMode:             vk.PolygonModeFill,
		CullMode:                vk.CullModeFlags(vk.CullModeNone),
		FrontFace:               vk.FrontFaceClockwise,
		DepthBiasEnable:         vk.False,
		LineWidth:               1,
	}

	multisampling := vk.PipelineMultisampleStateCreateInfo{
		SType:                vk.StructureTypePipelineMultisampleStateCreateInfo,
		SampleShadingEnable:  vk.False,
		RasterizationSamples: vk.SampleCount1Bit,
	}

	depthStencil := vk.PipelineDepthStencilStateCreateInfo{
		SType:                 vk.StructureTypePipelineDepthStencilStateCreateInfo,
		DepthTestEnable:       vk.False,
		DepthWriteEnable:      vk.False,
		DepthCompareOp:        vk.CompareOpLess,
		DepthBoundsTestEnable: vk.False,
		StencilTestEnable:     vk.False,
	}

	colorBlendAttachment := vk.PipelineColorBlendAttachmentState{
		BlendEnable: vk.True,
		SrcColorBlendFactor: vk.BlendFactorSrcAlpha,
		DstColorBlendFactor: vk.BlendFactorOneMinusSrcAlpha,
		ColorBlendOp:        vk.BlendOpAdd,
		SrcAlphaBlendFactor: vk.BlendFactorOne,
		DstAlphaBlendFactor: vk.BlendFactorZero,
		AlphaBlendOp:        vk.BlendOpAdd,
		ColorWriteMask:      vk.ColorComponentFlags(vk.ColorComponentRBit | vk.ColorComponentGBit | vk.ColorComponentBBit | vk.ColorComponentABit),
	}

	colorBlend := vk.PipelineColorBlendStateCreateInfo{
		SType:           vk.StructureTypePipelineColorBlendStateCreateInfo,
		LogicOpEnable:   vk.False,
		AttachmentCount: 1,
		PAttachments:    []vk.PipelineColorBlendAttachmentState{colorBlendAttachment},
		BlendConstants:  [4]float32{0, 0, 0, 0},
	}

	pipelineLayoutInfo := &vk.PipelineLayoutCreateInfo{
		SType: vk.StructureTypePipelineLayoutCreateInfo,
	}

	var pipelineLayout vk.PipelineLayout
	result2 := vk.CreatePipelineLayout(device, pipelineLayoutInfo, nil, &pipelineLayout)
	if result2 != vk.Success {
		return nil, fmt.Errorf("failed to create pipeline layout: %v", vk.Error(result2))
	}

	pipelineInfo := &vk.GraphicsPipelineCreateInfo{
		SType:               vk.StructureTypeGraphicsPipelineCreateInfo,
		StageCount:          2,
		PStages:             []vk.PipelineShaderStageCreateInfo{vertStage, fragStage},
		PVertexInputState:   &vertexInput,
		PInputAssemblyState: &inputAssembly,
		PViewportState:      &viewportState,
		PRasterizationState: &rasterizer,
		PMultisampleState:   &multisampling,
		PDepthStencilState:  &depthStencil,
		PColorBlendState:    &colorBlend,
		Layout:              pipelineLayout,
		RenderPass:          renderPass,
		Subpass:             0,
	}

	var pipeline vk.Pipeline
	result3 := vk.CreateGraphicsPipelines(device, nil, 1, []*vk.GraphicsPipelineCreateInfo{pipelineInfo}, &pipeline)
	if result3 != vk.Success {
		vk.DestroyPipelineLayout(device, pipelineLayout, nil)
		return nil, fmt.Errorf("failed to create graphics pipeline: %v", vk.Error(result3))
	}

	return &graphicsPipeline{
		pipeline:       pipeline,
		pipelineLayout: pipelineLayout,
		renderPass:     renderPass,
	}, nil
}

func createVertexBuffer(device vk.Device, physicalDevice vk.PhysicalDevice) (vk.Buffer, vk.DeviceMemory, error) {
	bufferSize := uintptr(maxQuads * 6 * 24) // maxQuads * 6 verts per quad * 24 bytes per vert

	bufferInfo := &vk.BufferCreateInfo{
		SType:       vk.StructureTypeBufferCreateInfo,
		Size:        bufferSize,
		Usage:       vk.BufferUsageFlags(vk.BufferUsageVertexBufferBit),
		SharingMode: vk.SharingModeExclusive,
	}

	var buffer vk.Buffer
	result := vk.CreateBuffer(device, bufferInfo, nil, &buffer)
	if result != vk.Success {
		return nil, nil, fmt.Errorf("failed to create vertex buffer: %v", vk.Error(result))
	}

	var memReqs vk.MemoryRequirements
	vk.GetBufferMemoryRequirements(device, buffer, &memReqs)

	memTypeIndex := findMemoryType(physicalDevice, memReqs.MemoryTypeBits, vk.MemoryPropertyHostVisibleBit|vk.MemoryPropertyHostCoherentBit)

	allocInfo := &vk.MemoryAllocateInfo{
		SType:          vk.StructureTypeMemoryAllocateInfo,
		AllocationSize: memReqs.Size,
		MemoryTypeIndex: memTypeIndex,
	}

	var bufferMemory vk.DeviceMemory
	result2 := vk.AllocateMemory(device, allocInfo, nil, &bufferMemory)
	if result2 != vk.Success {
		vk.DestroyBuffer(device, buffer, nil)
		return nil, nil, fmt.Errorf("failed to allocate vertex buffer memory: %v", vk.Error(result2))
	}

	vk.BindBufferMemory(device, buffer, bufferMemory, 0)
	return buffer, bufferMemory, nil
}

func findMemoryType(physicalDevice vk.PhysicalDevice, typeFilter uint32, properties vk.MemoryPropertyFlags) uint32 {
	var memProps vk.PhysicalDeviceMemoryProperties
	vk.GetPhysicalDeviceMemoryProperties(physicalDevice, &memProps)

	for i := uint32(0); i < memProps.MemoryTypeCount; i++ {
		if typeFilter&(1<<i) != 0 && memProps.MemoryTypes[i].PropertyFlags&properties == properties {
			return i
		}
	}
	return 0
}

func createCommandPool(device vk.Device, queueFamilyIndex uint32) (vk.CommandPool, error) {
	poolInfo := &vk.CommandPoolCreateInfo{
		SType:            vk.StructureTypeCommandPoolCreateInfo,
		Flags:            vk.CommandPoolCreateFlags(vk.CommandPoolCreateResetCommandBufferBit),
		QueueFamilyIndex: queueFamilyIndex,
	}

	var pool vk.CommandPool
	result := vk.CreateCommandPool(device, poolInfo, nil, &pool)
	if result != vk.Success {
		return nil, fmt.Errorf("failed to create command pool: %v", vk.Error(result))
	}
	return pool, nil
}

func createCommandBuffers(device vk.Device, pool vk.CommandPool, count uint32) ([]vk.CommandBuffer, error) {
	allocInfo := &vk.CommandBufferAllocateInfo{
		SType:              vk.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        pool,
		Level:              vk.CommandBufferLevelPrimary,
		CommandBufferCount: count,
	}

	buffers := make([]vk.CommandBuffer, count)
	result := vk.AllocateCommandBuffers(device, allocInfo, buffers)
	if result != vk.Success {
		return nil, fmt.Errorf("failed to allocate command buffers: %v", vk.Error(result))
	}
	return buffers, nil
}

func createSyncObjects(device vk.Device) (vk.Semaphore, vk.Semaphore, vk.Fence, error) {
	semaInfo := &vk.SemaphoreCreateInfo{
		SType: vk.StructureTypeSemaphoreCreateInfo,
	}

	var imageAvailable, renderFinished vk.Semaphore
	result1 := vk.CreateSemaphore(device, semaInfo, nil, &imageAvailable)
	if result1 != vk.Success {
		return nil, nil, nil, fmt.Errorf("failed to create semaphore: %v", vk.Error(result1))
	}
	result2 := vk.CreateSemaphore(device, semaInfo, nil, &renderFinished)
	if result2 != vk.Success {
		vk.DestroySemaphore(device, imageAvailable, nil)
		return nil, nil, nil, fmt.Errorf("failed to create semaphore: %v", vk.Error(result2))
	}

	fenceInfo := &vk.FenceCreateInfo{
		SType: vk.StructureTypeFenceCreateInfo,
		Flags: vk.FenceCreateFlags(vk.FenceCreateSignaledBit),
	}
	var fence vk.Fence
	result3 := vk.CreateFence(device, fenceInfo, nil, &fence)
	if result3 != vk.Success {
		vk.DestroySemaphore(device, imageAvailable, nil)
		vk.DestroySemaphore(device, renderFinished, nil)
		return nil, nil, nil, fmt.Errorf("failed to create fence: %v", vk.Error(result3))
	}

	return imageAvailable, renderFinished, fence, nil
}

func destroyPipelineObjects(device vk.Device, gp *graphicsPipeline) {
	if gp == nil {
		return
	}
	if gp.inFlightFence != nil {
		vk.DestroyFence(device, gp.inFlightFence, nil)
	}
	if gp.renderFinished != nil {
		vk.DestroySemaphore(device, gp.renderFinished, nil)
	}
	if gp.imageAvailable != nil {
		vk.DestroySemaphore(device, gp.imageAvailable, nil)
	}
	if gp.commandPool != nil {
		vk.DestroyCommandPool(device, gp.commandPool, nil)
	}
	if gp.vertexBuffer != nil {
		vk.DestroyBuffer(device, gp.vertexBuffer, nil)
	}
	if gp.vertexMemory != nil {
		vk.FreeMemory(device, gp.vertexMemory, nil)
	}
	if gp.pipeline != nil {
		vk.DestroyPipeline(device, gp.pipeline, nil)
	}
	if gp.pipelineLayout != nil {
		vk.DestroyPipelineLayout(device, gp.pipelineLayout, nil)
	}
	if gp.renderPass != nil {
		vk.DestroyRenderPass(device, gp.renderPass, nil)
	}
}
