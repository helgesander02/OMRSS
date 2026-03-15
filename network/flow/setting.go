package flow

var (
	bgTSNFlowsEnd int
	bgAVBFlowsEnd int
)

func fillFlowBreakpoint(bgTSN int, bgAVB int) {
	bgTSNFlowsEnd = bgTSN
	bgAVBFlowsEnd = bgAVB
}
