package db

import "testing"

func TestSanitizeCaregiverPersonalInfoRemovesPhone(t *testing.T) {
	got := SanitizeCaregiverPersonalInfo(map[string]interface{}{
		"heightCm":     165,
		"languages":    []interface{}{"普通话", "13800138000"},
		"contactPhone": "13800138000",
		"phone":        "13900139000",
	})
	if got["heightCm"] != 165 {
		t.Fatalf("height was not preserved: %#v", got)
	}
	if _, exists := got["contactPhone"]; exists {
		t.Fatalf("contactPhone leaked: %#v", got)
	}
	if _, exists := got["phone"]; exists {
		t.Fatalf("phone leaked: %#v", got)
	}
	if _, exists := got["languages"]; exists {
		t.Fatalf("phone-containing public text was retained: %#v", got)
	}
}

func TestContainsPhoneInPublicContent(t *testing.T) {
	if !ContainsPhoneInPublicContent(map[string]interface{}{"description": "联系 13800138000"}) {
		t.Fatal("phone number was not detected")
	}
	if ContainsPhoneInPublicContent(map[string]interface{}{"imageUrls": []interface{}{"https://example.com/13800138000.jpg"}}) {
		t.Fatal("media URL should not be treated as public text")
	}
}

func TestBusinessStatusTransitions(t *testing.T) {
	if !validDemandTransition("PENDING_CONTACT", "CONTACTED") || validDemandTransition("PENDING_CONTACT", "MATCHING") {
		t.Fatal("demand transition rules are invalid")
	}
	if !validResumeTransition("VERIFYING", "APPROVED") || validResumeTransition("PENDING_CONTACT", "APPROVED") {
		t.Fatal("resume transition rules are invalid")
	}
}

func TestNormalizeHomeBannersSupportsCurrentAndLegacyConfig(t *testing.T) {
	current := NormalizeHomeBanners(map[string]interface{}{"items": []interface{}{map[string]interface{}{
		"id": "banner_1", "imageUrl": "https://example.com/banner.jpg", "title": "专业护理", "actionType": "DEMAND", "actionValue": "nursing", "sort": float64(10),
	}}})
	if len(current) != 1 || current[0].Title != "专业护理" || current[0].ActionValue != "nursing" {
		t.Fatalf("current banner config was not preserved: %#v", current)
	}
	legacy := NormalizeHomeBanners(map[string]interface{}{"urls": []interface{}{"https://example.com/legacy.jpg"}})
	if len(legacy) != 1 || legacy[0].ImageURL != "https://example.com/legacy.jpg" || legacy[0].ActionType != "NONE" {
		t.Fatalf("legacy URL config was not converted: %#v", legacy)
	}
}
