package backend

import (
	"github.com/paulmach/go.geojson"
	"reflect"
	"fmt"
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
var Equal Operator = "="
var NotEqual Operator = "!="

// other operators
var Any Operator = "any"   // operator for if this key exists
var None Operator = "none" // operator for if this doesn't exist

// special operators
var Within Operator = "within"
var Intersects Operator = "intersects"

// special keys
var GeometryType = "GeometryType"
var ID = "ID"
var Area = "AREA"

// filter types
type FeatureFilter struct {
	Operator    Operator
	Filters     []*FeatureFilter
	Key         string
	FloatValue  float64
	StringValue string
	StringBool  bool
	BoolValue   bool
	BoolBool    bool
	Area        float64
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


func EncodeFromInterface(op Operator,key string,value interface{}) *FeatureFilter {
	total := &FeatureFilter{}
	mult,boolval := value.([]interface{})
	if boolval {
		tempop := Operator(mult[0].(string))
		for _,i := range mult {
			for tempop == AndOperator || tempop == OrOperator {

				total.Filters = append(total.Filters,EncodeFromInterface(tempop,"",newmult))
			}
		}
		fmt.Println(mult,tempop)
		if tempop == AndOperator || tempop == OrOperator {

			newmult := mult[1]
			total.Filters = append(total.Filters,EncodeFromInterface(tempop,"",newmult))
			//for _,i := range newmult {
			//	total.Filters = append(total.Filters,EncodeFromInterface(tempop,"",i))
			//}		
		}
		total.Operator = tempop
	
	} else {
		total = ParseFeatureFilter(op,key,value)
	}

	return total

}
