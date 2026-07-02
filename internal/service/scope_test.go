package service

import "testing"

func TestResolveScope(t *testing.T) {
	sid := uint(3)
	cases := []struct {
		name      string
		role      string
		userSat   *uint
		wantAll   bool
		wantSatID uint
	}{
		{"super admin lihat semua", "super_admin", nil, true, 0},
		{"admin ter-scope ke satdiknya", "admin", &sid, false, 3},
		{"admin tanpa satdik tak lihat apa pun", "admin", nil, false, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			all, satID := ResolveScope(c.role, c.userSat)
			if all != c.wantAll || satID != c.wantSatID {
				t.Fatalf("ResolveScope(%q,%v)=(%v,%d), want (%v,%d)",
					c.role, c.userSat, all, satID, c.wantAll, c.wantSatID)
			}
		})
	}
}
