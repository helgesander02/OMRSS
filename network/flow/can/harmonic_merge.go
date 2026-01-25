package can

// ============================================================================
// Harmonic Merge 相關輔助函數
// ============================================================================

// gcd 計算兩個整數的最大公約數 (Greatest Common Divisor)
// 使用歐幾里得算法
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// gcdMultiple 計算多個整數的最大公約數
func gcdMultiple(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	result := nums[0]
	for i := 1; i < len(nums); i++ {
		result = gcd(result, nums[i])
	}
	return result
}

// isHarmonicPeriod 檢查兩個週期是否為調和關係
// 調和週期：一個週期是另一個週期的整數倍
func isHarmonicPeriod(period1, period2 int) bool {
	if period1 == 0 || period2 == 0 {
		return false
	}
	larger := period1
	smaller := period2
	if period2 > period1 {
		larger = period2
		smaller = period1
	}
	return larger%smaller == 0
}

// canMergeFlows 檢查兩個Flow是否可以合併
// 條件：相同源域、相同目的域、調和週期、合併後不超過MTU
func canMergeFlows(flow1, flow2 *Flow, MTU float64) bool {
	// 檢查源域和目的域
	if flow1.Source != flow2.Source || flow1.Destination != flow2.Destination {
		return false
	}

	// 檢查調和週期
	if !isHarmonicPeriod(flow1.Period, flow2.Period) {
		return false
	}

	// 檢查合併後的大小
	if flow1.DataSize+flow2.DataSize > MTU {
		return false
	}

	return true
}

// FlowGroup 用於表示可以合併的Flow組
type FlowGroup struct {
	Flows       []*Flow
	TotalSize   float64
	MinDeadline int
	GCDPeriod   int
}

// harmonicMerge 實現論文第二階段的Harmonic Merge算法
// 1. 識別效率差的TSN消息 (DataSize < MTU/2)
// 2. 尋找具有相同源域、目的域和調和週期的消息
// 3. 重新聚合並更新Flow的屬性
func (method *Method) harmonicMerge(MTU float64) {
	threshold := MTU / 2.0 // 750 bytes

	// 步驟1: 識別需要優化的Flow（DataSize < MTU/2）
	inefficientFlows := make([]*Flow, 0)
	efficientFlows := make([]*Flow, 0)

	for _, flow := range method.CAN2TTFlows {
		if flow.DataSize < threshold {
			inefficientFlows = append(inefficientFlows, flow)
		} else {
			efficientFlows = append(efficientFlows, flow)
		}
	}

	// 如果沒有效率差的Flow，直接返回
	if len(inefficientFlows) == 0 {
		return
	}

	// 步驟2: 尋找可合併的Flow組
	mergeGroups := method.findMergeableGroups(inefficientFlows, MTU)

	// 步驟3: 執行合併
	newFlows := make([]*Flow, 0)
	merged := make(map[*Flow]bool) // 記錄已經合併的Flow

	for _, group := range mergeGroups {
		if len(group.Flows) > 1 {
			// 合併這組Flow
			mergedFlow := method.mergeFlowGroup(group)
			newFlows = append(newFlows, mergedFlow)

			// 標記已合併的Flow
			for _, flow := range group.Flows {
				merged[flow] = true
			}
		}
	}

	// 步驟4: 重建CAN2TTFlows列表
	// 保留效率高的Flow和未被合併的Flow
	finalFlows := make([]*Flow, 0)
	finalFlows = append(finalFlows, efficientFlows...)

	for _, flow := range inefficientFlows {
		if !merged[flow] {
			finalFlows = append(finalFlows, flow)
		}
	}

	finalFlows = append(finalFlows, newFlows...)
	method.CAN2TTFlows = finalFlows
}

// findMergeableGroups 尋找所有可合併的Flow組
func (method *Method) findMergeableGroups(flows []*Flow, MTU float64) []*FlowGroup {
	groups := make([]*FlowGroup, 0)
	processed := make(map[*Flow]bool)

	for i := 0; i < len(flows); i++ {
		if processed[flows[i]] {
			continue
		}

		// 創建新組
		group := &FlowGroup{
			Flows:       []*Flow{flows[i]},
			TotalSize:   flows[i].DataSize,
			MinDeadline: flows[i].Deadline,
			GCDPeriod:   flows[i].Period,
		}
		processed[flows[i]] = true

		// 嘗試添加其他可合併的Flow
		for j := i + 1; j < len(flows); j++ {
			if processed[flows[j]] {
				continue
			}

			// 檢查是否可以加入這個組
			if method.canAddToGroup(group, flows[j], MTU) {
				group.Flows = append(group.Flows, flows[j])
				group.TotalSize += flows[j].DataSize
				if flows[j].Deadline < group.MinDeadline {
					group.MinDeadline = flows[j].Deadline
				}
				// 更新GCD
				group.GCDPeriod = gcd(group.GCDPeriod, flows[j].Period)
				processed[flows[j]] = true
			}
		}

		groups = append(groups, group)
	}

	return groups
}

// canAddToGroup 檢查Flow是否可以加入組
func (method *Method) canAddToGroup(group *FlowGroup, flow *Flow, MTU float64) bool {
	if len(group.Flows) == 0 {
		return false
	}

	// 檢查源域和目的域
	firstFlow := group.Flows[0]
	if flow.Source != firstFlow.Source || flow.Destination != firstFlow.Destination {
		return false
	}

	// 檢查調和週期（與組中所有Flow的週期）
	for _, groupFlow := range group.Flows {
		if !isHarmonicPeriod(flow.Period, groupFlow.Period) {
			return false
		}
	}

	// 檢查大小限制
	if group.TotalSize+flow.DataSize > MTU {
		return false
	}

	return true
}

// mergeFlowGroup 合併一組Flow成為一個新的Flow
func (method *Method) mergeFlowGroup(group *FlowGroup) *Flow {
	if len(group.Flows) == 0 {
		return nil
	}

	// 計算合併後的屬性
	firstFlow := group.Flows[0]
	periods := make([]int, len(group.Flows))
	for i, flow := range group.Flows {
		periods[i] = flow.Period
	}

	mergedFlow := &Flow{
		Source:      firstFlow.Source,
		Destination: firstFlow.Destination,
		Period:      gcdMultiple(periods),  // 新週期 = GCD(所有週期)
		Deadline:    group.MinDeadline,     // 最小截止時間
		DataSize:    group.TotalSize,       // 總數據大小
		HyperPeriod: firstFlow.HyperPeriod, // 保持相同的HyperPeriod
		Streams:     make([]*Stream, 0),
	}

	// 合併所有Stream（根據新的週期重新生成）
	mergedFlow.Streams = method.regenerateStreams(group, mergedFlow.Period, mergedFlow.HyperPeriod)

	return mergedFlow
}

// regenerateStreams 根據新的週期重新生成Stream
func (method *Method) regenerateStreams(group *FlowGroup, newPeriod int, hyperPeriod int) []*Stream {
	streams := make([]*Stream, 0)

	// 按照新週期生成Stream實例
	for currentTime := 0; currentTime < hyperPeriod; currentTime += newPeriod {
		stream := &Stream{
			ArrivalTime: currentTime,
			Deadline:    group.MinDeadline,
			DataSize:    group.TotalSize,
			FinishTime:  currentTime + group.MinDeadline,
		}
		streams = append(streams, stream)
	}

	return streams
}
