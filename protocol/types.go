package protocol

// Message type constants for the unified protocol.
const (
	// Generic data channel
	TypeData = "data"

	// Edge -> Core (published on orders topic)
	TypeOrderRequest       = "order.request"
	TypeOrderCancel        = "order.cancel"
	TypeOrderReceipt       = "order.receipt"
	TypeOrderRedirect      = "order.redirect"
	TypeOrderStorageWaybill = "order.storage_waybill"

	// Core -> Edge (published on dispatch topic)
	TypeOrderAck        = "order.ack"
	TypeOrderWaybill    = "order.waybill"
	TypeOrderUpdate     = "order.update"
	TypeOrderDelivered  = "order.delivered"
	TypeOrderError      = "order.error"
	TypeOrderCancelled  = "order.cancelled"
)

// Data channel subject constants.
const (
	SubjectEdgeRegister    = "edge.register"
	SubjectEdgeRegistered  = "edge.registered"
	SubjectEdgeHeartbeat   = "edge.heartbeat"
	SubjectEdgeHeartbeatAck = "edge.heartbeat_ack"
)

// Roles for Address.Role.
const (
	RoleEdge = "edge"
	RoleCore = "core"
)

// Protocol version.
const Version = 1
