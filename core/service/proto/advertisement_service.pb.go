// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.27.0
// source: advertisement_service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_advertisement_service_proto protoreflect.FileDescriptor

var file_advertisement_service_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x67,
	0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x2f, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x5f, 0x64, 0x74, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x99, 0x04, 0x0a,
	0x14, 0x41, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x27, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x12, 0x06, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x10, 0x2e, 0x41, 0x64, 0x47,
	0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x2c,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x2e,
	0x41, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1a, 0x0c, 0x2e, 0x53,
	0x41, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x0e,
	0x53, 0x61, 0x76, 0x65, 0x41, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0c,
	0x2e, 0x53, 0x41, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x09, 0x2e, 0x54,
	0x78, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2c, 0x0a, 0x10, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x41, 0x64, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0d, 0x2e, 0x41, 0x64,
	0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x1a, 0x09, 0x2e, 0x54, 0x78, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2f, 0x0a, 0x0c, 0x50, 0x75, 0x74, 0x44, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x41, 0x64, 0x12, 0x14, 0x2e, 0x53, 0x65, 0x74, 0x44, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x41, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x54, 0x78,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x51, 0x75, 0x65, 0x72, 0x79, 0x41,
	0x64, 0x12, 0x0f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x41, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x10, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x41, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x5b, 0x0a, 0x16, 0x51, 0x75, 0x65, 0x72, 0x79, 0x41,
	0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x1e, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x41, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x41, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x29, 0x0a, 0x09, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x64,
	0x12, 0x11, 0x2e, 0x53, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x41, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e, 0x54, 0x78, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x28,
	0x0a, 0x10, 0x47, 0x65, 0x74, 0x41, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x0c, 0x2e, 0x41, 0x64, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x04, 0x2e, 0x53, 0x41, 0x64, 0x22, 0x00, 0x12, 0x19, 0x0a, 0x06, 0x53, 0x61, 0x76, 0x65,
	0x41, 0x64, 0x12, 0x04, 0x2e, 0x53, 0x41, 0x64, 0x1a, 0x09, 0x2e, 0x54, 0x78, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x12, 0x23, 0x0a, 0x08, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x41, 0x64, 0x12,
	0x0c, 0x2e, 0x41, 0x64, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x09, 0x2e,
	0x54, 0x78, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x42, 0x1f, 0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x67, 0x6f, 0x32, 0x6f, 0x2e, 0x72, 0x70, 0x63, 0x5a,
	0x08, 0x2e, 0x2f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var file_advertisement_service_proto_goTypes = []interface{}{
	(*Empty)(nil),                          // 0: Empty
	(*AdPositionId)(nil),                   // 1: AdPositionId
	(*SAdPosition)(nil),                    // 2: SAdPosition
	(*SetDefaultAdRequest)(nil),            // 3: SetDefaultAdRequest
	(*QueryAdRequest)(nil),                 // 4: QueryAdRequest
	(*QueryAdvertisementDataRequest)(nil),  // 5: QueryAdvertisementDataRequest
	(*SetUserAdRequest)(nil),               // 6: SetUserAdRequest
	(*AdIdRequest)(nil),                    // 7: AdIdRequest
	(*SAd)(nil),                            // 8: SAd
	(*AdGroupResponse)(nil),                // 9: AdGroupResponse
	(*TxResult)(nil),                       // 10: TxResult
	(*QueryAdResponse)(nil),                // 11: QueryAdResponse
	(*QueryAdvertisementDataResponse)(nil), // 12: QueryAdvertisementDataResponse
}
var file_advertisement_service_proto_depIdxs = []int32{
	0,  // 0: AdvertisementService.GetGroups:input_type -> Empty
	1,  // 1: AdvertisementService.GetPosition:input_type -> AdPositionId
	2,  // 2: AdvertisementService.SaveAdPosition:input_type -> SAdPosition
	1,  // 3: AdvertisementService.DeleteAdPosition:input_type -> AdPositionId
	3,  // 4: AdvertisementService.PutDefaultAd:input_type -> SetDefaultAdRequest
	4,  // 5: AdvertisementService.QueryAd:input_type -> QueryAdRequest
	5,  // 6: AdvertisementService.QueryAdvertisementData:input_type -> QueryAdvertisementDataRequest
	6,  // 7: AdvertisementService.SetUserAd:input_type -> SetUserAdRequest
	7,  // 8: AdvertisementService.GetAdvertisement:input_type -> AdIdRequest
	8,  // 9: AdvertisementService.SaveAd:input_type -> SAd
	7,  // 10: AdvertisementService.DeleteAd:input_type -> AdIdRequest
	9,  // 11: AdvertisementService.GetGroups:output_type -> AdGroupResponse
	2,  // 12: AdvertisementService.GetPosition:output_type -> SAdPosition
	10, // 13: AdvertisementService.SaveAdPosition:output_type -> TxResult
	10, // 14: AdvertisementService.DeleteAdPosition:output_type -> TxResult
	10, // 15: AdvertisementService.PutDefaultAd:output_type -> TxResult
	11, // 16: AdvertisementService.QueryAd:output_type -> QueryAdResponse
	12, // 17: AdvertisementService.QueryAdvertisementData:output_type -> QueryAdvertisementDataResponse
	10, // 18: AdvertisementService.SetUserAd:output_type -> TxResult
	8,  // 19: AdvertisementService.GetAdvertisement:output_type -> SAd
	10, // 20: AdvertisementService.SaveAd:output_type -> TxResult
	10, // 21: AdvertisementService.DeleteAd:output_type -> TxResult
	11, // [11:22] is the sub-list for method output_type
	0,  // [0:11] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_advertisement_service_proto_init() }
func file_advertisement_service_proto_init() {
	if File_advertisement_service_proto != nil {
		return
	}
	file_global_proto_init()
	file_message_advertisement_dto_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_advertisement_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_advertisement_service_proto_goTypes,
		DependencyIndexes: file_advertisement_service_proto_depIdxs,
	}.Build()
	File_advertisement_service_proto = out.File
	file_advertisement_service_proto_rawDesc = nil
	file_advertisement_service_proto_goTypes = nil
	file_advertisement_service_proto_depIdxs = nil
}
