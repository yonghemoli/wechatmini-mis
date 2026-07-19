package mis

import (
	"testing"

	"yonghemolimis/src/dao/db"
)

func TestCompanyDecorationValuePreservesSavedFields(t *testing.T) {
	company := companyDecorationValue(`{"logoUrl":"https://example.com/logo.png","name":"永和护理","address":"南宁","introduction":"护理服务","contactPhone":"19994740191","serviceGuarantees":[{"icon":"verified","title":"身份核验","sub":"资料核验"}]}`)
	if got := company["contactPhone"]; got != "19994740191" {
		t.Fatalf("contactPhone lost after JSON decode: %#v", got)
	}
	if got := company["name"]; got != "永和护理" {
		t.Fatalf("company name lost after JSON decode: %#v", got)
	}
	items, ok := company["serviceGuarantees"].([]db.ServiceGuarantee)
	if !ok || len(items) != 1 || items[0].Title != "身份核验" {
		t.Fatalf("service guarantees were not normalized: %#v", company["serviceGuarantees"])
	}
}

func TestCompanyDecorationValueSupportsLegacyGuarantees(t *testing.T) {
	company := companyDecorationValue(`{"contactPhone":"19994740191","serviceGuarantees":["身份核验"]}`)
	items := company["serviceGuarantees"].([]db.ServiceGuarantee)
	if len(items) != 1 || items[0].Title != "身份核验" {
		t.Fatalf("legacy guarantees were not converted: %#v", items)
	}
}
