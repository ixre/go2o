package impl

import "go2o/core/service/proto"

var _ proto.TemplateServiceServer = new(templateServiceImpl)

type templateServiceImpl struct {
}
