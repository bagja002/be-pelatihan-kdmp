package service

// ResolveScope menentukan cakupan query peserta dari role & satdik user.
//   - super_admin        → all=true (tanpa filter)
//   - admin dgn satdik    → all=false, satdikID = satdiknya
//   - admin tanpa satdik   → all=false, satdikID=0 (tak melihat apa pun)
func ResolveScope(role string, userSatdik *uint) (all bool, satdikID uint) {
	if role == "super_admin" {
		return true, 0
	}
	if userSatdik != nil {
		return false, *userSatdik
	}
	return false, 0
}
