package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Proto"
	"math"
)

// DistancePawn pawn之间距离
func DistancePawn(self, target *Pawn) float32 {
	return Distance(self.GetPos(), target.GetPos())
}

// Distance 两点间距离
func Distance(pos1, pos2 linemath.Vector2) float32 {
	return pos1.Sub(pos2).Len()
}

// CalcAngle 计算角度   pos1点 到 pos2点形成的向量角度
func CalcAngle(pos1, pos2 linemath.Vector2) float64 {
	var dv linemath.Vector2

	// 转换为方位向量
	dv.X = pos2.X - pos1.X
	dv.Y = pos2.Y - pos1.Y

	// 朝向目标角度
	angle := math.Atan2(float64(dv.Y), float64(dv.X))
	angle = math.Mod(angle, math.Pi*2)
	if angle < 0 {
		angle = 2*math.Pi + angle
	}

	return angle
}

//AddAngle 角度相加
func AddAngle(angle, addAngle float32) float32 {
	newAngle := math.Mod(float64(angle+addAngle), math.Pi*2)
	if newAngle < 0 {
		newAngle = 2*math.Pi + newAngle
	}

	return float32(newAngle)
}

// BetweenAngle 检测是否在指定角度范围之间
func BetweenAngle(angle, low, high float32) bool {
	normalFun := func(angle float32) float32 {
		angle = float32(math.Mod(float64(angle), 2*math.Pi))
		if angle < 0 {
			angle += 2 * math.Pi
		}
		return angle
	}

	angle = normalFun(angle)
	low = normalFun(low)
	high = normalFun(high)

	angle = normalFun(angle - low)
	high = normalFun(high - low)
	low = 0

	if low <= high {
		return angle >= low && angle <= high
	} else {
		return angle >= high && angle <= low
	}
}

// FanAndCircleOverlap 检测扇形与圆形是否重叠
func FanAndCircleOverlap(circlePos linemath.Vector2, circleRadius float32, fanPos linemath.Vector2, fanRadius, fanOrientAngle, fanIncAngle float32) bool {
	centerVec := circlePos.Sub(fanPos)
	centerDis := centerVec.Len()

	// 1.检测圆心距离
	if centerDis > circleRadius+fanRadius {
		return false
	}

	// 2.检测扇形圆心是否在圆内
	if centerDis <= circleRadius {
		return true
	}

	lowAngle := fanOrientAngle - fanIncAngle
	upperAngle := fanOrientAngle + fanIncAngle

	// 3.扇形上边界在圆内
	upperPos := linemath.Vector2{
		X: float32(math.Cos(float64(upperAngle))),
		Y: float32(math.Sin(float64(upperAngle))),
	}.Mul(fanRadius).Add(fanPos)

	if upperPos.Sub(circlePos).Len() <= circleRadius {
		return true
	}

	// 4.扇形下边界在圆内
	lowPos := linemath.Vector2{
		X: float32(math.Cos(float64(lowAngle))),
		Y: float32(math.Sin(float64(lowAngle))),
	}.Mul(fanRadius).Add(fanPos)

	if lowPos.Sub(circlePos).Len() <= circleRadius {
		return true
	}

	// 5.目标圆心处于扇形角度区域范围内
	angle := float32(CalcAngle(fanPos, circlePos))
	if BetweenAngle(angle, lowAngle, upperAngle) {
		return true
	}

	boundaryAngle := float32(math.Asin(float64(circleRadius / centerDis)))

	// 6.扇形上边界与圆形相交
	maxAngle := upperAngle + boundaryAngle
	if centerVec.Dot(upperPos.Sub(fanPos)) >= 0 && circlePos.Sub(upperPos).Dot(fanPos.Sub(upperPos)) >= 0 && BetweenAngle(angle, lowAngle, maxAngle) {
		return true
	}

	// 7.扇形下边界与圆形相交
	minAngle := lowAngle - boundaryAngle
	if centerVec.Dot(lowPos.Sub(fanPos)) >= 0 && circlePos.Sub(lowPos).Dot(fanPos.Sub(lowPos)) >= 0 && BetweenAngle(angle, minAngle, upperAngle) {
		return true
	}

	return false
}

// TransPos 移动位置
func TransPos(pos, offset linemath.Vector2, angle float32) linemath.Vector2 {
	pos.AddS(linemath.Vector2{
		X: float32(math.Cos(float64(angle))),
		Y: float32(math.Sin(float64(angle))),
	}.Mul(offset.Y).Add(linemath.Vector2{
		X: float32(math.Cos(float64(angle) + linemath.PI_HALF)),
		Y: float32(math.Sin(float64(angle) + linemath.PI_HALF)),
	}.Mul(offset.X)))

	return pos
}

// GetEnemyCamp 获取敌方阵营
func GetEnemyCamp(camp Proto.Camp_Enum) Proto.Camp_Enum {
	if Proto.Camp_Blue == camp {
		return Proto.Camp_Red
	} else {
		return Proto.Camp_Blue
	}
}

// Max 返回最大值
func Max(a, b float64) float64 {
	return math.Max(a, b)
}

// Min 返回最小值
func Min(a, b float64) float64 {
	return math.Min(a, b)
}

const FloatZero float64 = 1e-6

// IsZero 浮点数0值判断
func IsZero(num float64) bool {
	return math.Abs(num) <= FloatZero
}

// FloatEqual 判断相等
func FloatEqual(a, b float64) bool {
	if a > b {
		return a-b < FloatZero
	} else {
		return b-a < FloatZero
	}
}

// FloatLessEqual 判断小于等于
func FloatLessEqual(a, b float64) bool {
	if a < b {
		return true
	} else {
		return a-b < FloatZero
	}
}

// FloatGreaterEqual 判断大于等于
func FloatGreaterEqual(a, b float64) bool {
	if a > b {
		return true
	} else {
		return b-a < FloatZero
	}
}

// Ceil 取上整
func Ceil(value float64) float64 {
	if FloatEqual(value, 0) {
		return 0
	}

	return math.Ceil(value)
}

// Floor 取下整
func Floor(value float64) float64 {
	if FloatEqual(value, 0) {
		return 0
	}

	return math.Floor(value)
}

// Vector2Equal 向量判断相等
func Vector2Equal(a, b linemath.Vector2) bool {
	return FloatEqual(float64(a.X), float64(b.X)) && FloatEqual(float64(a.Y), float64(b.Y))
}

// Bits 位运算
type Bits uint64

func (bits *Bits) Reset() *Bits {
	*bits = 0
	return bits
}

func (bits *Bits) TurnOn(bit int32) *Bits {
	if bit >= 0 && bit < 64 {
		*bits |= 1 << bit
	}
	return bits
}

func (bits *Bits) TurnOff(bit int32) *Bits {
	if bit >= 0 && bit < 64 {
		*bits &= ^(1 << bit)
	}
	return bits
}

func (bits *Bits) Test(bit int32) bool {
	if bit >= 0 && bit < 64 {
		return (*bits)&(1<<bit) != 0
	}
	return false
}

func (bits *Bits) Any(other Bits) bool {
	return (uint64(*bits) & uint64(other)) != 0
}

func (bits *Bits) All(other Bits) bool {
	return (uint64(*bits) & uint64(other)) == uint64(other)
}

// BitsStick 拼接Bits
func BitsStick(bit ...int32) (bits Bits) {
	for _, v := range bit {
		bits.TurnOn(v)
	}
	return
}

// IntAbs 整数绝对值
func IntAbs(value int64) int64 {
	if value < 0 {
		return -value
	}

	return value
}
