package backend

import (
	m "github.com/murphy214/mercantile"
	"github.com/paulmach/go.geojson"
	"reflect"
	
)

type Operator string

// logic operators
var AndOperator Operator = "all"
var OrOperator Operator = "any"
var XorOperator Operator = "xor"

// logical operators singular
var GreaterThan Operator = ">"
var GreaterThanOrEqual Operator = ">="
var LessThan Operator = "<"
var LessThanOrEqual Operator = "<="
var Equal Operator = "=="
var NotEqual Operator = "!="

// other operators
var Any Operator = "in"   // operator for if this key exists
var None Operator = "!in" // operator for if this doesn't exist

// special operators
var Intersects Operator = "intersects"

// special keys
var GeometryType = "$type"
var ID = "ID"
var Area = "AREA"

// filter types
type FeatureFilter struct {
	Operator    Operator
	Filters     []*FeatureFilter
	Map 		map[string]string
	Bounds 		m.Extrema
	Key         string
	FloatValue  float64
	StringValue string
	StringBool  bool
	BoolValue   bool
	BoolBool    bool
	InBool 		bool
	BoundsBool 	bool
	Area        float64
}


// structure for finding overlapping values
func Overlapping_1D(box1min float64,box1max float64,box2min float64,box2max float64) bool {
	if box1max >= box2min && box2max >= box1min {
		return true
	} else {
		return false
	}
	return false
}


// returns a boolval for whether or not the bb intersects
func IntersectBB(bdsref m.Extrema,bds m.Extrema) bool {
	if Overlapping_1D(bdsref.W,bdsref.E,bds.W,bds.E) && Overlapping_1D(bdsref.S,bdsref.N,bds.S,bds.N) {
		return true
	} else {
		return false
	}

	return false
}

// returns a bool for within filter
func Within(bdsref,bds m.Extrema) bool {
	return bdsref.W <= bds.W && bdsref.S <= bds.S && bdsref.E >= bds.E && bdsref.N >= bds.N 
}

// returns whether a box intersects
func IntersectsAll(bds,bdsref m.Extrema) bool {
	return Within(bdsref,bds) || Within(bds,bdsref) || IntersectBB(bdsref, bds)
}



//
func (filter *FeatureFilter) Filter(feature *geojson.Feature) bool {
	// going through each filter
	if len(filter.Filters) > 0 {
		// logic for dealing with an operator
		boolval := true
		if filter.Operator != AndOperator {
			boolval = false
		}
		count := 0
		for _, f := range filter.Filters {
			tempbool := f.Filter(feature)
			switch filter.Operator {
			case AndOperator:
				boolval = boolval && tempbool
				if !boolval {
					return boolval
				}
			case OrOperator:
				boolval = boolval || tempbool
			case XorOperator:
				boolval = boolval || tempbool
				if tempbool {
					count++
				}

			}
		}
		if filter.Operator == XorOperator && count != 1 && boolval {
			return false
		}
		return boolval
	} else {

		// logic for dealing with a single filter
		if filter.Area > 0 {
			// TO DO IMPLEMENT THIS RIGHT
			//area := GetArea(feature.Geometry)
			//return area > filter.Area
			return true
		}

		switch filter.Key {
		case GeometryType:
			switch filter.Operator {
			case Equal:
				return string(feature.Geometry.Type) == filter.StringValue
			case NotEqual:
				return string(feature.Geometry.Type) != filter.StringValue
			}
		case ID:
			var myval float64
			v := reflect.ValueOf(feature.ID)
			vv := v.Kind()

			switch vv {
			case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
				myval = float64(v.Uint())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				myval = float64(v.Int())
			case reflect.Float32, reflect.Float64:
				myval = v.Float()
			}

			switch filter.Operator {
			case Equal:
				return myval == filter.FloatValue
			case NotEqual:
				return myval != filter.FloatValue
			}
		case Area:
			// area case
		default:
			// logic for in bool filter
			if filter.InBool {
				val, boolval := feature.Properties[filter.Key]
				valstring,boolval2 := val.(string)
				if boolval && boolval2 {
					_,boolval := filter.Map[valstring]
					if filter.Operator == Any {
						if boolval {
							return true
						} else {
							return false
						}
					} else if filter.Operator == None {
						if boolval {
							return false
						} else {
							return true
						}
					}
				}
			}

			// filters about the bounds
			if filter.BoundsBool && len(feature.BoundingBox) == 4 {
				w,s,e,n := feature.BoundingBox[0],feature.BoundingBox[1],feature.BoundingBox[2],feature.BoundingBox[3]
				return IntersectsAll(m.Extrema{N:n,S:s,E:e,W:w}, filter.Bounds)
			}


			if filter.BoolBool {
				val, boolval := feature.Properties[filter.Key]
				if boolval {
					newval, boolval := val.(bool)
					if boolval {
						return newval == filter.BoolValue
					}
				}
				return false
			}

			// default behavior
			val, boolval := feature.Properties[filter.Key]
			if !boolval {
				if filter.Operator == None {
					return true
				} else {
					return false
				}
			} else if boolval && filter.Operator == Any {
				return true
			}

			if filter.StringBool {
				testvalue, boolval := val.(string)
				if !boolval {
					return false
				}
				switch filter.Operator {
				case Equal:
					return filter.StringValue == testvalue
				case NotEqual:
					return filter.StringValue != testvalue
				default:
					return false
				}
			} else {
				var myval float64
				v := reflect.ValueOf(val)
				vv := v.Kind()

				switch vv {
				case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
					myval = float64(v.Uint())
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					myval = float64(v.Int())
				case reflect.Float32, reflect.Float64:
					myval = v.Float()
				}

				switch filter.Operator {
				case Equal:
					return filter.FloatValue == myval
				case NotEqual:
					return filter.FloatValue != myval
				case GreaterThan:
					return filter.FloatValue < myval
				case GreaterThanOrEqual:
					return filter.FloatValue <= myval
				case LessThan:
					return filter.FloatValue > myval
				case LessThanOrEqual:
					return filter.FloatValue >= myval
				default:
					return false
				}
			}
		}
	}
	return false
}

func IsBase(value interface{}) bool {
	myval,boolval := value.([]interface{}) 
	if !boolval {
		return true
	}
	if len(myval) > 2 {
		return !boolval
	}
	return false
}

func ParseFeatureFilter(op Operator,key string,singeinterface interface{}) *FeatureFilter {
	v := reflect.ValueOf(singeinterface)
	vv := v.Kind()
	var myval float64
	switch vv {
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		myval = float64(v.Uint())
		return &FeatureFilter{
			Operator:op,
			Key:key,
			FloatValue:myval,

		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		myval = float64(v.Int())
		return &FeatureFilter{
			Operator:op,
			Key:key,
			FloatValue:myval,

		}
	case reflect.Float32, reflect.Float64:
		myval = v.Float()
		return &FeatureFilter{
			Operator:op,
			Key:key,
			FloatValue:myval,

		}
	case reflect.String:
		return &FeatureFilter{
			Operator:op,
			Key:key,
			StringValue:v.String(),
			StringBool:true,
		}
	}
	return &FeatureFilter{}

}


// parses a simple feature filter
func ParseSimple(vals []interface{}) *FeatureFilter {
	opp,boolval := vals[0].(string)
	op := Operator(opp)
	if boolval {
		switch op {
		case Equal,NotEqual,GreaterThan,GreaterThanOrEqual,LessThan,LessThanOrEqual:
			return ParseFeatureFilter(op, vals[1].(string), vals[2])		
		case Any,None:
			key := vals[1].(string)
			newmap := map[string]string{}
			for _,myval := range vals[2:] {
				myvals,boolval := myval.(string)
				if boolval {
					newmap[myvals] = ""
				}
			}
			return &FeatureFilter{
				Operator:op,
				Map:newmap,
				InBool:true,
				Key:key,
			}
		case Intersects:
			key := vals[1].(string)
			bdsarr2 := vals[2].([]interface{})
			bdsarr := make([]float64,len(bdsarr2))
			for pos,val := range bdsarr2 {
				vall,boolval := val.(float64)
				if boolval {
					bdsarr[pos] = vall
				}
			}

			var ext m.Extrema
			if len(bdsarr) == 4 {
				w,s,e,n := bdsarr[0],bdsarr[1],bdsarr[2],bdsarr[3]
				ext = m.Extrema{W:w,S:s,E:e,N:n}
				return &FeatureFilter{
					Operator:op,
					Bounds:ext,
					BoundsBool:true,
					Key:key,
				}
			}
		} 
	}
	return &FeatureFilter{}
}

// given a blank interface containing a mapbox-gl js feature filter returns a feature filter for the backend
func ParseAll(vals []interface{}) *FeatureFilter {
	opp,boolval := vals[0].(string)
	op := Operator(opp)
	if boolval {
		switch op {
		case AndOperator,OrOperator:
			myfilter := &FeatureFilter{Operator:op}
			for _,val := range vals[1:] {
				valss,boolval := val.([]interface{})
				if boolval {
					myfilter.Filters = append(myfilter.Filters,ParseSimple(valss))
				}
			}
			return myfilter
		default:
			if boolval {
				return ParseSimple(vals)
			}
		}
	}
	return &FeatureFilter{}
}