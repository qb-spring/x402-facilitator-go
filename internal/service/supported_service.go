package service

import (
	"x402-facilitator-go/internal/config"
	"x402-facilitator-go/internal/models"
)

// SupportedService provides information about supported schemes and networks
type SupportedService struct {
	NetworkInfos []config.NetworkInfo
}

// NewSupportedService creates a new SupportedService
func NewSupportedService(networkInfos []config.NetworkInfo) *SupportedService {
	return &SupportedService{
		NetworkInfos: networkInfos,
	}
}

// Supported returns the supported payment schemes and networks
func (s *SupportedService) Supported() *models.SupportedResponse {
	kinds := make([]models.SupportedKind, 0, len(s.NetworkInfos))

	for _, networkInfo := range s.NetworkInfos {
		kinds = append(kinds, models.SupportedKind{
			X402Version: networkInfo.X402Version,
			Scheme:      networkInfo.Scheme,
			Network:     networkInfo.Name,
		})
	}

	return &models.SupportedResponse{
		Kinds: kinds,
	}
}
