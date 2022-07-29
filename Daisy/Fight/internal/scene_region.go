package internal

import (
	"Cinder/Base/linemath"
	"Daisy/Proto"
	"fmt"
	"github.com/ByteArena/box2d"
	"math"
)

// _RegionFlag 区域标记
type _RegionFlag uint32

const (
	RegionFlag_AoeWarning _RegionFlag = iota // AOE伤害警告区域
)

// _RegionShape 区域形状
type _RegionShape struct {
	Type     Proto.AttackShapeType_Enum // 形状类型
	Extend   linemath.Vector2           // 矩形区域全长全高
	Radius   float32                    // 圆形或扇形区域半径
	FanAngle float32                    // 扇形区域夹角
}

// _RegionInfo 区域信息
type _RegionInfo struct {
	Creator *Pawn            // 创造者
	Flag    _RegionFlag      // 区域标记
	Shape   _RegionShape     // 区域形状
	Pos     linemath.Vector2 // 位置
	Angle   float32          // 角度
}

// _Region 区域
type _Region struct {
	UID     uint32                 // uid
	Info    *_RegionInfo           // 区域信息
	b2Shape box2d.B2ShapeInterface // 物理图形
}

// _RegionMgr 区域管理器
type _RegionMgr struct {
	scene     *Scene              // 场景
	regionMap map[uint32]*_Region // 区域表
}

// init 初始化
func (regionMgr *_RegionMgr) init(scene *Scene) {
	regionMgr.scene = scene
	regionMgr.regionMap = make(map[uint32]*_Region)
}

// queryPosInRegionList 查询位置所在区域列表
func (regionMgr *_RegionMgr) queryPosInRegionList(pos linemath.Vector2, radius float32, camp Proto.Camp_Enum) []*_Region {
	tShape := box2d.MakeB2CircleShape()
	tShape.SetRadius(float64(radius))
	tShape.M_p.Set(float64(pos.X), float64(pos.Y))

	xf := box2d.MakeB2Transform()
	xf.SetIdentity()

	var regionList []*_Region

	for _, region := range regionMgr.regionMap {
		manifold := box2d.NewB2Manifold()

		switch region.b2Shape.GetType() {
		case box2d.B2Shape_Type.E_circle:
			box2d.B2CollideCircles(manifold, &tShape, xf, region.b2Shape.(*box2d.B2CircleShape), xf)
		case box2d.B2Shape_Type.E_polygon:
			box2d.B2CollidePolygonAndCircle(manifold, region.b2Shape.(*box2d.B2PolygonShape), xf, &tShape, xf)
		}

		if manifold.PointCount > 0 {
			if region.Info.Shape.Type == Proto.AttackShapeType_Fan {
				if !FanAndCircleOverlap(pos, radius, region.Info.Pos, region.Info.Shape.Radius, region.Info.Angle, region.Info.Shape.FanAngle) {
					continue
				}
			}

			regionList = append(regionList, region)
		}
	}

	return regionList
}

// createRegion 创建区域
func (regionMgr *_RegionMgr) createRegion(info *_RegionInfo) (region *_Region, err error) {
	shape := createRegionB2Shape(info)
	if shape == nil {
		return nil, fmt.Errorf("create b2Shape failed, %+v", *info)
	}

	region = &_Region{
		UID:     regionMgr.scene.generateUID(),
		Info:    info,
		b2Shape: shape,
	}

	regionMgr.regionMap[region.UID] = region

	return region, nil
}

// destroyRegion 销毁区域
func (regionMgr *_RegionMgr) destroyRegion(uid uint32) {
	delete(regionMgr.regionMap, uid)
}

// createRegionB2Shape 创建物理图形
func createRegionB2Shape(regionInfo *_RegionInfo) box2d.B2ShapeInterface {
	switch regionInfo.Shape.Type {
	case Proto.AttackShapeType_Rect:
		halfX := float64(regionInfo.Shape.Extend.X / 2)
		halfY := float64(regionInfo.Shape.Extend.Y / 2)
		angle := float64(regionInfo.Angle + math.Pi/2)

		b2Shape := box2d.MakeB2PolygonShape()
		b2Shape.SetAsBoxFromCenterAndAngle(halfX, halfY, box2d.MakeB2Vec2(float64(regionInfo.Pos.X), float64(regionInfo.Pos.Y)), angle)
		return &b2Shape

	case Proto.AttackShapeType_Circle:
		fallthrough

	case Proto.AttackShapeType_Fan:
		b2Shape := box2d.MakeB2CircleShape()
		b2Shape.SetRadius(float64(regionInfo.Shape.Radius))
		b2Shape.M_p.Set(float64(regionInfo.Pos.X), float64(regionInfo.Pos.Y))

		return &b2Shape
	}

	return nil
}
