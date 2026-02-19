package www

import (
	"net/http"
)

func (h *Handlers) handleManualOrder(w http.ResponseWriter, r *http.Request) {
	db := h.engine.DB()

	payloads, _ := db.ListPayloads()
	nodes, _ := db.ListLocationNodes()
	coreNodes := h.engine.CoreNodes()
	coreNodeNames := make([]string, 0, len(coreNodes))
	for name := range coreNodes {
		coreNodeNames = append(coreNodeNames, name)
	}
	anomalies, rpMap := loadAnomalyData(h)

	data := map[string]interface{}{
		"Page":              "manual-order",
		"Payloads":          payloads,
		"Nodes":             nodes,
		"CoreNodes":         coreNodeNames,
		"Anomalies":         anomalies,
		"ReportingPointMap": rpMap,
	}

	h.renderTemplate(w, "manual-order.html", data)
}
