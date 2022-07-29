package sceneconfigdata

import "Cinder/Base/linemath"

//ConvertVector3 y-up转z-up，因为物理和linemath基于z-up
func ConvertVector3(loc linemath.Vector3) linemath.Vector3 {
	return linemath.Vector3{X: loc.Z, Y: loc.X, Z: loc.Y}
}

//UnconvertVector3 z-up转y-up，因为寻路基于y-up
func UnconvertVector3(loc linemath.Vector3) linemath.Vector3 {
	return linemath.Vector3{X: loc.Y, Y: loc.Z, Z: loc.X}
}
