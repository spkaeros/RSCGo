package script

import (
	"github.com/d5/tengo/objects"
	"github.com/spkaeros/rscgo/pkg/server/clients"
	"github.com/spkaeros/rscgo/pkg/server/world"
	"time"
)

var scriptAttributes = map[string]objects.Object{
	"replaceObjectAt": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 3 {
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			object := world.GetObject(x, y)
			if object == nil {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "GetObject(x,y)",
					Expected: "An object",
					Found:    "Nothing",
				}
			}
			id, ok := objects.ToInt(args[2])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "objectID",
					Expected: "int",
					Found:    args[2].TypeName(),
				}
			}
			world.ReplaceObject(object, id)
			return objects.UndefinedValue, nil
		},
	},
	"replaceObject": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 2 {
				return nil, objects.ErrWrongNumArguments
			}
			object, ok := args[0].(*world.Object)
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "object",
					Expected: "*world.Object",
					Found:    args[0].TypeName(),
				}
			}
			id, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "objectID",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			return world.ReplaceObject(object, id), nil
		},
	},
	"getObject": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 2 {
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			object := world.GetObject(x, y)
			if object == nil {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "GetObject(x,y)",
					Expected: "An object",
					Found:    "Nothing",
				}
			}
			return object, nil
		},
	},
	"addObject": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 3 {
				return nil, objects.ErrWrongNumArguments
			}
			id, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "id",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			x, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[2])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[2].TypeName(),
				}
			}
			if object := world.GetObject(x, y); object != nil {
				world.ReplaceObject(object, id)
			} else {
				world.AddObject(world.NewObject(id, 0, x, y, false))
			}
			return objects.UndefinedValue, nil
		},
	},
	"removeObject": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 2 {
				return nil, objects.ErrWrongNumArguments
			}
			x, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "x",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			y, ok := objects.ToInt(args[1])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "y",
					Expected: "int",
					Found:    args[1].TypeName(),
				}
			}
			if object := world.GetObject(x, y); object == nil {
				world.RemoveObject(object)
			}
			return objects.UndefinedValue, nil
		},
	},
	"getPlayer": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 1 {
				return nil, objects.ErrWrongNumArguments
			}
			index, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "index",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			client, ok := clients.FromIndex(index)
			if !ok {
				return nil, objects.ErrIndexOutOfBounds
			}
			return client, nil
		},
	},
	"sleep": &objects.UserFunction{
		Value: func(args ...objects.Object) (ret objects.Object, err error) {
			if len(args) < 1 {
				return nil, objects.ErrWrongNumArguments
			}
			duration, ok := objects.ToInt(args[0])
			if !ok {
				return nil, objects.ErrInvalidArgumentType{
					Name:     "duration",
					Expected: "int",
					Found:    args[0].TypeName(),
				}
			}
			time.Sleep(time.Duration(duration) * time.Millisecond)
			return objects.UndefinedValue, nil
		},
	},
}

func NewWorldModule() *objects.BuiltinModule {
	return &objects.BuiltinModule{Attrs: scriptAttributes}
}
