// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v4.25.3
// source: proto/spiceroute.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ProfileService_UpsertPreference_FullMethodName = "/spiceroute.v1.ProfileService/UpsertPreference"
	ProfileService_GetPreference_FullMethodName    = "/spiceroute.v1.ProfileService/GetPreference"
)

// ProfileServiceClient is the client API for ProfileService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProfileServiceClient interface {
	UpsertPreference(ctx context.Context, in *Preference, opts ...grpc.CallOption) (*Preference, error)
	GetPreference(ctx context.Context, in *Preference, opts ...grpc.CallOption) (*Preference, error)
}

type profileServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProfileServiceClient(cc grpc.ClientConnInterface) ProfileServiceClient {
	return &profileServiceClient{cc}
}

func (c *profileServiceClient) UpsertPreference(ctx context.Context, in *Preference, opts ...grpc.CallOption) (*Preference, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Preference)
	err := c.cc.Invoke(ctx, ProfileService_UpsertPreference_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *profileServiceClient) GetPreference(ctx context.Context, in *Preference, opts ...grpc.CallOption) (*Preference, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Preference)
	err := c.cc.Invoke(ctx, ProfileService_GetPreference_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProfileServiceServer is the server API for ProfileService service.
// All implementations must embed UnimplementedProfileServiceServer
// for forward compatibility.
type ProfileServiceServer interface {
	UpsertPreference(context.Context, *Preference) (*Preference, error)
	GetPreference(context.Context, *Preference) (*Preference, error)
	mustEmbedUnimplementedProfileServiceServer()
}

// UnimplementedProfileServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProfileServiceServer struct{}

func (UnimplementedProfileServiceServer) UpsertPreference(context.Context, *Preference) (*Preference, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertPreference not implemented")
}
func (UnimplementedProfileServiceServer) GetPreference(context.Context, *Preference) (*Preference, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPreference not implemented")
}
func (UnimplementedProfileServiceServer) mustEmbedUnimplementedProfileServiceServer() {}
func (UnimplementedProfileServiceServer) testEmbeddedByValue()                        {}

// UnsafeProfileServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProfileServiceServer will
// result in compilation errors.
type UnsafeProfileServiceServer interface {
	mustEmbedUnimplementedProfileServiceServer()
}

func RegisterProfileServiceServer(s grpc.ServiceRegistrar, srv ProfileServiceServer) {
	// If the following call pancis, it indicates UnimplementedProfileServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProfileService_ServiceDesc, srv)
}

func _ProfileService_UpsertPreference_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Preference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).UpsertPreference(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_UpsertPreference_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).UpsertPreference(ctx, req.(*Preference))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProfileService_GetPreference_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Preference)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProfileServiceServer).GetPreference(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProfileService_GetPreference_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProfileServiceServer).GetPreference(ctx, req.(*Preference))
	}
	return interceptor(ctx, in, info, handler)
}

// ProfileService_ServiceDesc is the grpc.ServiceDesc for ProfileService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProfileService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spiceroute.v1.ProfileService",
	HandlerType: (*ProfileServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertPreference",
			Handler:    _ProfileService_UpsertPreference_Handler,
		},
		{
			MethodName: "GetPreference",
			Handler:    _ProfileService_GetPreference_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/spiceroute.proto",
}

const (
	PlannerService_GeneratePlan_FullMethodName = "/spiceroute.v1.PlannerService/GeneratePlan"
)

// PlannerServiceClient is the client API for PlannerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlannerServiceClient interface {
	GeneratePlan(ctx context.Context, in *PlanRequest, opts ...grpc.CallOption) (*PlanResponse, error)
}

type plannerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlannerServiceClient(cc grpc.ClientConnInterface) PlannerServiceClient {
	return &plannerServiceClient{cc}
}

func (c *plannerServiceClient) GeneratePlan(ctx context.Context, in *PlanRequest, opts ...grpc.CallOption) (*PlanResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PlanResponse)
	err := c.cc.Invoke(ctx, PlannerService_GeneratePlan_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlannerServiceServer is the server API for PlannerService service.
// All implementations must embed UnimplementedPlannerServiceServer
// for forward compatibility.
type PlannerServiceServer interface {
	GeneratePlan(context.Context, *PlanRequest) (*PlanResponse, error)
	mustEmbedUnimplementedPlannerServiceServer()
}

// UnimplementedPlannerServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPlannerServiceServer struct{}

func (UnimplementedPlannerServiceServer) GeneratePlan(context.Context, *PlanRequest) (*PlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GeneratePlan not implemented")
}
func (UnimplementedPlannerServiceServer) mustEmbedUnimplementedPlannerServiceServer() {}
func (UnimplementedPlannerServiceServer) testEmbeddedByValue()                        {}

// UnsafePlannerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlannerServiceServer will
// result in compilation errors.
type UnsafePlannerServiceServer interface {
	mustEmbedUnimplementedPlannerServiceServer()
}

func RegisterPlannerServiceServer(s grpc.ServiceRegistrar, srv PlannerServiceServer) {
	// If the following call pancis, it indicates UnimplementedPlannerServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PlannerService_ServiceDesc, srv)
}

func _PlannerService_GeneratePlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlannerServiceServer).GeneratePlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PlannerService_GeneratePlan_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlannerServiceServer).GeneratePlan(ctx, req.(*PlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PlannerService_ServiceDesc is the grpc.ServiceDesc for PlannerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlannerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spiceroute.v1.PlannerService",
	HandlerType: (*PlannerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GeneratePlan",
			Handler:    _PlannerService_GeneratePlan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/spiceroute.proto",
}

const (
	RecipeService_CreateRecipe_FullMethodName = "/spiceroute.v1.RecipeService/CreateRecipe"
	RecipeService_ListRecipes_FullMethodName  = "/spiceroute.v1.RecipeService/ListRecipes"
)

// RecipeServiceClient is the client API for RecipeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecipeServiceClient interface {
	CreateRecipe(ctx context.Context, in *Recipe, opts ...grpc.CallOption) (*RecipeID, error)
	ListRecipes(ctx context.Context, in *RecipeQuery, opts ...grpc.CallOption) (*RecipeList, error)
}

type recipeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecipeServiceClient(cc grpc.ClientConnInterface) RecipeServiceClient {
	return &recipeServiceClient{cc}
}

func (c *recipeServiceClient) CreateRecipe(ctx context.Context, in *Recipe, opts ...grpc.CallOption) (*RecipeID, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecipeID)
	err := c.cc.Invoke(ctx, RecipeService_CreateRecipe_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recipeServiceClient) ListRecipes(ctx context.Context, in *RecipeQuery, opts ...grpc.CallOption) (*RecipeList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RecipeList)
	err := c.cc.Invoke(ctx, RecipeService_ListRecipes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecipeServiceServer is the server API for RecipeService service.
// All implementations must embed UnimplementedRecipeServiceServer
// for forward compatibility.
type RecipeServiceServer interface {
	CreateRecipe(context.Context, *Recipe) (*RecipeID, error)
	ListRecipes(context.Context, *RecipeQuery) (*RecipeList, error)
	mustEmbedUnimplementedRecipeServiceServer()
}

// UnimplementedRecipeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRecipeServiceServer struct{}

func (UnimplementedRecipeServiceServer) CreateRecipe(context.Context, *Recipe) (*RecipeID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRecipe not implemented")
}
func (UnimplementedRecipeServiceServer) ListRecipes(context.Context, *RecipeQuery) (*RecipeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListRecipes not implemented")
}
func (UnimplementedRecipeServiceServer) mustEmbedUnimplementedRecipeServiceServer() {}
func (UnimplementedRecipeServiceServer) testEmbeddedByValue()                       {}

// UnsafeRecipeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecipeServiceServer will
// result in compilation errors.
type UnsafeRecipeServiceServer interface {
	mustEmbedUnimplementedRecipeServiceServer()
}

func RegisterRecipeServiceServer(s grpc.ServiceRegistrar, srv RecipeServiceServer) {
	// If the following call pancis, it indicates UnimplementedRecipeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RecipeService_ServiceDesc, srv)
}

func _RecipeService_CreateRecipe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Recipe)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).CreateRecipe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_CreateRecipe_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).CreateRecipe(ctx, req.(*Recipe))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecipeService_ListRecipes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecipeQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecipeServiceServer).ListRecipes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecipeService_ListRecipes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecipeServiceServer).ListRecipes(ctx, req.(*RecipeQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// RecipeService_ServiceDesc is the grpc.ServiceDesc for RecipeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecipeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spiceroute.v1.RecipeService",
	HandlerType: (*RecipeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRecipe",
			Handler:    _RecipeService_CreateRecipe_Handler,
		},
		{
			MethodName: "ListRecipes",
			Handler:    _RecipeService_ListRecipes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/spiceroute.proto",
}

const (
	FeedbackService_SubmitFeedback_FullMethodName = "/spiceroute.v1.FeedbackService/SubmitFeedback"
)

// FeedbackServiceClient is the client API for FeedbackService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FeedbackServiceClient interface {
	SubmitFeedback(ctx context.Context, in *FeedbackBatch, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type feedbackServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFeedbackServiceClient(cc grpc.ClientConnInterface) FeedbackServiceClient {
	return &feedbackServiceClient{cc}
}

func (c *feedbackServiceClient) SubmitFeedback(ctx context.Context, in *FeedbackBatch, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, FeedbackService_SubmitFeedback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FeedbackServiceServer is the server API for FeedbackService service.
// All implementations must embed UnimplementedFeedbackServiceServer
// for forward compatibility.
type FeedbackServiceServer interface {
	SubmitFeedback(context.Context, *FeedbackBatch) (*emptypb.Empty, error)
	mustEmbedUnimplementedFeedbackServiceServer()
}

// UnimplementedFeedbackServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFeedbackServiceServer struct{}

func (UnimplementedFeedbackServiceServer) SubmitFeedback(context.Context, *FeedbackBatch) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitFeedback not implemented")
}
func (UnimplementedFeedbackServiceServer) mustEmbedUnimplementedFeedbackServiceServer() {}
func (UnimplementedFeedbackServiceServer) testEmbeddedByValue()                         {}

// UnsafeFeedbackServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FeedbackServiceServer will
// result in compilation errors.
type UnsafeFeedbackServiceServer interface {
	mustEmbedUnimplementedFeedbackServiceServer()
}

func RegisterFeedbackServiceServer(s grpc.ServiceRegistrar, srv FeedbackServiceServer) {
	// If the following call pancis, it indicates UnimplementedFeedbackServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FeedbackService_ServiceDesc, srv)
}

func _FeedbackService_SubmitFeedback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FeedbackBatch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FeedbackServiceServer).SubmitFeedback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: FeedbackService_SubmitFeedback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FeedbackServiceServer).SubmitFeedback(ctx, req.(*FeedbackBatch))
	}
	return interceptor(ctx, in, info, handler)
}

// FeedbackService_ServiceDesc is the grpc.ServiceDesc for FeedbackService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FeedbackService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "spiceroute.v1.FeedbackService",
	HandlerType: (*FeedbackServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubmitFeedback",
			Handler:    _FeedbackService_SubmitFeedback_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/spiceroute.proto",
}
