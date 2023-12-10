package service_test

import (
	"testing"

	"github.com/VanLavr/tg-bot/internal/bot/service"
)

func TestValidateIP(t *testing.T) {
	ip := "192.168.0.1/24"
	botService := service.New()
	result := botService.ValidateIP(ip)
	expected := true

	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	ip = "192.3.0.1/adfs"
	result = botService.ValidateIP(ip)
	expected = false

	if result != expected {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
