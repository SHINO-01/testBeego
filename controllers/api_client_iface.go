// controllers/api_client_iface.go

package controllers

// APIClient is the interface that wraps the method to call The Cat API (or any external API).
type APIClient interface {
    MakeAPIRequest(endpoint, method string, payload interface{}) chan []byte
}
