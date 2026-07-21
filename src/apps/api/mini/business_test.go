package miniapi

import (
	"encoding/json"
	"testing"

	"yonghemolimis/src/dao/db"
)

func TestExperienceMatches(t *testing.T) {
	tests := []struct {
		years int
		value string
		want  bool
	}{
		{0, "LESS_THAN_1_YEAR", true}, {2, "YEAR_1_TO_3", true}, {4, "YEAR_3_TO_5", true},
		{8, "YEAR_5_TO_10", true}, {11, "MORE_THAN_10_YEARS", true}, {2, "YEAR_5_TO_10", false},
	}
	for _, tt := range tests {
		if got := experienceMatches(tt.years, tt.value); got != tt.want {
			t.Fatalf("experienceMatches(%d, %q)=%v, want %v", tt.years, tt.value, got, tt.want)
		}
	}
}

func TestMaskPhone(t *testing.T) {
	if got := maskPhone("13800138000"); got != "138****8000" {
		t.Fatalf("maskPhone=%q", got)
	}
}

func TestRandomNumericCode(t *testing.T) {
	for i := 0; i < 20; i++ {
		code := randomNumericCode()
		if len(code) != 6 || !numericCodePattern.MatchString(code) {
			t.Fatalf("invalid code shape: %q", code)
		}
	}
}

func TestCaregiverDetailViewDoesNotExposePhoneOrHiddenFields(t *testing.T) {
	row := db.CaregiverDO{
		ID: "CG1", Name: "张阿姨", ContactPhone: "13800138000", Introduction: "联系 13900139000",
		PersonalInfo:  `{"heightCm":165,"contactPhone":"13800138000","phone":"13900139000"}`,
		DisplayFields: `{"introduction":false,"personalInfo":true}`,
	}
	view := caregiverDetailView(row)
	raw, err := json.Marshal(view)
	if err != nil {
		t.Fatal(err)
	}
	if db.ContainsMainlandPhone(string(raw)) {
		t.Fatalf("public view exposes phone: %s", raw)
	}
	if _, exists := view["introduction"]; exists {
		t.Fatalf("hidden field was emitted: %#v", view)
	}
}
