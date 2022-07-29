package internal

import (
	"Cinder/Base/linemath"
	"github.com/ByteArena/box2d"
	"math"
)

type _Terrain struct {
	scene              *Scene
	worldBoundaryShape *box2d.B2ChainShape //场景边界形状
	worldPolygonShape  *box2d.B2PolygonShape
	shapeList          map[uint32]*_PawnShape //key PawnUID
}

type _PawnShape struct {
	*box2d.B2CircleShape
	owner *Pawn
}

func (b *_Terrain) init(scene *Scene) {
	b.shapeList = make(map[uint32]*_PawnShape)
	b.buildWorldBoundary(scene.Info.BoundaryPoints)
}

//buildWorldBoundary 构建世界边界
func (b *_Terrain) buildWorldBoundary(boundaryPoints []linemath.Vector2) {
	chainShape := box2d.MakeB2ChainShape()
	var vertices []box2d.B2Vec2
	for _, v := range boundaryPoints {
		vertices = append(vertices, box2d.B2Vec2{
			X: float64(v.X),
			Y: float64(v.Y),
		})
	}

	chainShape.CreateLoop(vertices, len(vertices))
	b.worldBoundaryShape = &chainShape

	polygonShape := box2d.MakeB2PolygonShape()
	polygonShape.Set(vertices, len(vertices))
	b.worldPolygonShape = &polygonShape
}

//createPawnShape 创建pawn形状 用于碰撞检测
func (b *_Terrain) createPawnShape(x, y, radius float64, pawn *Pawn) {
	circleShape := box2d.MakeB2CircleShape()
	circleShape.M_p.Set(x, y)
	circleShape.SetRadius(radius)
	shape := &_PawnShape{
		B2CircleShape: &circleShape,
		owner:         pawn,
	}

	b.shapeList[pawn.UID] = shape
}

//destroyPawnShape 销毁pawn形状
func (b *_Terrain) destroyPawnShape(uid uint32) {
	delete(b.shapeList, uid)
}

//setPawnShapeRadius 设置pawn形状 圆形半径
func (b *_Terrain) setPawnShapeRadius(uid uint32, radius float64) {
	shape, find := b.shapeList[uid]
	if !find {
		return
	}
	shape.B2CircleShape.SetRadius(radius)
}

//setPawnShapePos 设置pawn形状 位置
func (b *_Terrain) setPawnShapePos(uid uint32, x, y float64) {
	shape, ok := b.shapeList[uid]
	if !ok {
		return
	}

	shape.B2CircleShape.M_p.Set(x, y)
}

//overlapCircleShape 圆形碰撞检测
func (b *_Terrain) overlapCircleShape(x, y, radius float64) (targets []*Pawn) {
	shapeA := box2d.MakeB2CircleShape()
	shapeA.M_p.Set(x, y)
	shapeA.SetRadius(radius)
	xfA := box2d.MakeB2Transform()
	xfA.SetIdentity()
	xfB := box2d.MakeB2Transform()
	xfB.SetIdentity()

	for _, val := range b.shapeList {
		manifold := box2d.NewB2Manifold()
		box2d.B2CollideCircles(manifold, &shapeA, xfA, val.B2CircleShape, xfB)
		if manifold.PointCount > 0 {
			targets = append(targets, val.owner)
		}
	}

	return
}

//overlapRectangleShape 矩形碰撞检测
func (b *_Terrain) overlapRectangleShape(x, y, halfX, halfY, angle float64) (targets []*Pawn) {
	angle += math.Pi / 2

	shapeA := box2d.MakeB2PolygonShape()
	shapeA.SetAsBoxFromCenterAndAngle(halfX, halfY, box2d.MakeB2Vec2(x, y), angle)
	xfA := box2d.MakeB2Transform()
	xfA.SetIdentity()
	xfB := box2d.MakeB2Transform()
	xfB.SetIdentity()

	for _, val := range b.shapeList {
		manifold := box2d.NewB2Manifold()
		box2d.B2CollidePolygonAndCircle(manifold, &shapeA, xfA, val.B2CircleShape, xfB)
		if manifold.PointCount > 0 {
			targets = append(targets, val.owner)
		}
	}

	return
}

//overlapSectorShape 扇形碰撞检测 扇形原点、半径、半圆心角、扇形朝向
func (b *_Terrain) overlapSectorShape(originPos linemath.Vector2, radius, hCentralAngle, orientationAngle float32) (targets []*Pawn) {
	targetList := b.overlapCircleShape(float64(originPos.X), float64(originPos.Y), float64(radius))
	//半圆心角超过PI 相当于圆形碰撞检测
	if hCentralAngle >= math.Pi {
		targets = targetList
		return
	}

	//判断目标是否处于扇形夹角中间
	for _, pawn := range targetList {
		if FanAndCircleOverlap(pawn.GetPos(), pawn.Attr.CollisionRadius, originPos, radius, orientationAngle, hCentralAngle) {
			targets = append(targets, pawn)
		}
	}

	return
}

//sweepCircleShape 圆形滑动检测
func (b *_Terrain) sweepCircleShape(x, y, radius float64, direction linemath.Vector2, distance float32) (targets []*Pawn) {
	//map用于去重
	allResult := map[uint32]*Pawn{}

	//几何图形散点数量
	shapeCount := int(distance / float32(radius))

	//初始点以及中间散点检测
	for index := 0; index < shapeCount; index++ {
		tempX := float64(direction.Normalized().Mul(float32(radius)*float32(index)).X) + x
		tempY := float64(direction.Normalized().Mul(float32(radius)*float32(index)).Y) + y
		tempResult := b.overlapCircleShape(tempX, tempY, radius)
		for _, val := range tempResult {
			allResult[val.UID] = val
		}
	}

	//滑动终点检测
	tempX := float64(direction.Normalized().Mul(distance).X) + x
	tempY := float64(direction.Normalized().Mul(distance).Y) + y
	tempResult := b.overlapCircleShape(tempX, tempY, radius)
	for _, val := range tempResult {
		allResult[val.UID] = val
	}

	for _, val := range allResult {
		targets = append(targets, val)
	}

	return
}

//sweepRectangleShape 矩形滑动检测
func (b *_Terrain) sweepRectangleShape(x, y, halfX, halfY, angle float64, direction linemath.Vector2, distance float32) (targets []*Pawn) {
	//map用于去重
	allResult := map[uint32]*Pawn{}

	//基础单位距离
	baseDis := float32(halfX)
	if baseDis > float32(halfY) {
		baseDis = float32(halfY)
	}

	//几何图形散点数量
	shapeCount := int(distance / baseDis)

	//初始点以及中间散点检测
	for index := 0; index < shapeCount; index++ {
		tempX := float64(direction.Normalized().Mul(baseDis*float32(index)).X) + x
		tempY := float64(direction.Normalized().Mul(baseDis*float32(index)).Y) + y
		tempResult := b.overlapRectangleShape(tempX, tempY, halfX, halfY, angle)
		for _, val := range tempResult {
			allResult[val.UID] = val
		}
	}

	//滑动终点检测
	tempX := float64(direction.Normalized().Mul(distance).X) + x
	tempY := float64(direction.Normalized().Mul(distance).Y) + y
	tempResult := b.overlapRectangleShape(tempX, tempY, halfX, halfY, angle)
	for _, val := range tempResult {
		allResult[val.UID] = val
	}

	for _, val := range allResult {
		targets = append(targets, val)
	}

	return
}

//checkCircleShapeOverlapWorldBoundary 检查圆形是否触碰边界
func (b *_Terrain) checkCircleShapeOverlapWorldBoundary(x, y, radius float64) bool {
	shapeA := box2d.MakeB2CircleShape()
	shapeA.M_p.Set(x, y)
	shapeA.SetRadius(radius)
	xfA := box2d.MakeB2Transform()
	xfA.SetIdentity()

	for i := 0; i < b.worldBoundaryShape.GetChildCount(); i++ {

		manifold := box2d.NewB2Manifold()
		edge := box2d.MakeB2EdgeShape()
		b.worldBoundaryShape.GetChildEdge(&edge, i)
		xfB := box2d.MakeB2Transform()
		xfB.SetIdentity()

		box2d.B2CollideEdgeAndCircle(manifold, &edge, xfB, &shapeA, xfA)
		if manifold.PointCount > 0 {
			return true
		}

	}

	return false
}

// checkPointInWorldBoundary 检查点是否在战场区域内
func (b *_Terrain) checkPointInWorldBoundary(pos linemath.Vector2) bool {
	xf := box2d.MakeB2Transform()
	xf.SetIdentity()

	return b.worldPolygonShape.TestPoint(xf, box2d.B2Vec2{X: float64(pos.X), Y: float64(pos.Y)})
}
